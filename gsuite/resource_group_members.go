package gsuite

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	directory "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

func resourceGroupMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMembersCreate,
		Read:   resourceGroupMembersRead,
		Update: resourceGroupMembersUpdate,
		Delete: resourceGroupMembersDelete,

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"member": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: schemaMember,
				},
			},
		},
	}
}

func resourceGroupMembersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	members, err := getApiMembers(d.Id(), config)

	if err != nil {
		return err
	}

	d.Set("group", d.Id())
	d.Set("member", membersToCfg(members))
	return nil
}

func resourceGroupMembersCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Creating gsuite_group_members")
	gid, err := createOrUpdateGroupMembers(d, meta)

	if err != nil {
		return err
	}

	d.SetId(gid)
	return resourceGroupMembersRead(d, meta)
}

func resourceGroupMembersUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Updating gsuite_group_members")
	_, err := createOrUpdateGroupMembers(d, meta)

	if err != nil {
		return err
	}
	return resourceGroupMembersRead(d, meta)
}

func resourceGroupMembersDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Deleting gsuite_group_members")
	config := meta.(*Config)

	for _, raw_member := range d.Get("member").(*schema.Set).List() {
		member := raw_member.(map[string]interface{})
		deleteMember(member["email"].(string), d.Id(), config)
	}

	d.SetId("")
	return nil
}

func membersToCfg(members []*directory.Member) []map[string]interface{} {
	if members == nil {
		return nil
	}

	finalMembers := make([]map[string]interface{}, 0, len(members))

	for _, m := range members {
		finalMembers = append(finalMembers, map[string]interface{}{
			"email":  m.Email,
			"etag":   m.Etag,
			"kind":   m.Kind,
			"status": m.Status,
			"type":   m.Type,
			"role":   m.Role,
		})
	}

	return finalMembers
}

func resourceMembers(d *schema.ResourceData) (members []map[string]interface{}) {
	for _, raw_member := range d.Get("member").(*schema.Set).List() {
		member := raw_member.(map[string]interface{})
		members = append(members, member)
	}
	return members
}

func createOrUpdateGroupMembers(d *schema.ResourceData, meta interface{}) (string, error) {
	config := meta.(*Config)
	gid := d.Get("group").(string)

	// Get members from config
	cfgMembers := resourceMembers(d)

	// Get members from API
	apiMembers, err := getApiMembers(gid, config)
	if err != nil {
		return gid, fmt.Errorf("Error updating memberships: %v", err)
	}
	// This call removes any members that aren't defined in cfgMembers,
	// and adds all of those that are
	err = reconcileMembers(d, cfgMembers, membersToCfg(apiMembers), config, gid)
	if err != nil {
		return gid, fmt.Errorf("Error updating memberships: %v", err)
	}

	return gid, nil
}

// This function ensures that the members of a group exactly match that
// in a config by deleting any members that are returned by the API but not present
// in the config
func reconcileMembers(d *schema.ResourceData, cfgMembers, apiMembers []map[string]interface{}, config *Config, gid string) error {

	// Helper to convert slice to map
	m := func(vals []map[string]interface{}) map[string]map[string]interface{} {
		sm := make(map[string]map[string]interface{})
		for _, member := range vals {
			email := member["email"].(string)
			sm[email] = member
		}
		return sm
	}

	cfgMap := m(cfgMembers)
	apiMap := m(apiMembers)

	var cfgRole, apiRole string

	for k, apiMember := range apiMap {
		if cfgMember, ok := cfgMap[k]; !ok {
			// The member in the API is not in the config; disable it.
			err := deleteMember(k, gid, config)
			if err != nil {
				return err
			}
		} else {
			// The member exists in the config and the API
			// If role has changed update, otherwise do nothing
			cfgRole = cfgMember["role"].(string)
			apiRole = apiMember["role"].(string)
			if cfgRole != apiRole {
				groupMember := &directory.Member{
					Role: cfgRole,
				}
				var updatedGroupMember *directory.Member
				var err error
				err = retry(func() error {
					updatedGroupMember, err = config.directory.Members.Patch(
						d.Get("group").(string),
						cfgMember["email"].(string),
						groupMember).Do()
					return err
				})

				if err != nil {
					return fmt.Errorf("Error updating groupMember: %s", err)
				}

				log.Printf("[INFO] Updated groupMember: %s", updatedGroupMember.Email)
			}

			// Delete from cfgMap, we have already handled it
			delete(cfgMap, k)
		}
	}

	// Upsert memberships which are present in the config, but not in the api
	for email, _ := range cfgMap {
		err := upsertMember(email, gid, cfgMap[email]["role"].(string), config)
		if err != nil {
			return err
		}
	}
	return nil
}

// Retrieve a group's members from the API
func getApiMembers(gid string, config *Config) ([]*directory.Member, error) {
	// Get members from the API
	groupMembers := make([]*directory.Member, 0)
	token := ""
	var membersResponse *directory.Members
	var err error
	for paginate := true; paginate; {

		err = retry(func() error {
			membersResponse, err = config.directory.Members.List(gid).PageToken(token).Do()
			return err
		})

		if err != nil {
			return groupMembers, err
		}
		for _, v := range membersResponse.Members {
			groupMembers = append(groupMembers, v)
		}
		token = membersResponse.NextPageToken
		paginate = token != ""
	}
	return groupMembers, nil
}

func upsertMember(email, gid, role string, config *Config) error {
	groupMember := &directory.Member{
		Role:  role,
		Email: email,
	}

	var isGroup bool
	var group *directory.Group
	var err error
	err = retry(func() error {
		group, err = config.directory.Groups.Get(email).Do()
		return err
	})
	isGroup = true
	if err != nil {
		isGroup = false
	}

	if isGroup == true {
		var currentMember *directory.Member
		var err error
		err = retry(func() error {
			currentMember, err = config.directory.Members.Get(gid, email).Do()
			return err
		})

		if err != nil {
			var createdGroupMember *directory.Member
			err = retry(func() error {
				createdGroupMember, err = config.directory.Members.Insert(gid, groupMember).Do()
				return err
			})
			if err != nil {
				return fmt.Errorf("Error creating groupMember: %s, %s", err, email)
			}
			log.Printf("[INFO] Created groupMember: %s", createdGroupMember.Email)
		} else {
			var updatedGroupMember *directory.Member
			err = retry(func() error {
				updatedGroupMember, err = config.directory.Members.Update(gid, email, groupMember).Do()
				return err
			})
			if err != nil {
				return fmt.Errorf("Error updating groupMember: %s, %s", err, email)
			}
			log.Printf("[INFO] Updated groupMember: %s", updatedGroupMember.Email)
		}
	}

	if isGroup == false {
		var hasMemberResponse *directory.MembersHasMember
		var err error
		err = retry(func() error {
			hasMemberResponse, err = config.directory.Members.HasMember(gid, email).Do()
			if err == nil {
				return err
			}

			// When a user does not exist, the API returns a 400 "memberKey, required"
			// Returning a friendly message
			if gerr, ok := err.(*googleapi.Error); ok && (gerr.Errors[0].Reason == "required" && gerr.Code == 400) {
			  return fmt.Errorf("Error adding groupMember %s. Please make sure the user exists beforehand.", email)
			}
			return err
		})
		if err != nil {
			return fmt.Errorf("Error checking hasmember: %s, %s", err, email)
		}

		if hasMemberResponse.IsMember == true {
			var updatedGroupMember *directory.Member
			err = retry(func() error {
				updatedGroupMember, err = config.directory.Members.Update(gid, email, groupMember).Do()
				return err
			})
			if err != nil {
				return fmt.Errorf("Error updating groupMember: %s, %s", err, email)
			}
			log.Printf("[INFO] Updated groupMember: %s", updatedGroupMember.Email)
		} else {
			var createdGroupMember *directory.Member
			err = retry(func() error {
				createdGroupMember, err = config.directory.Members.Insert(gid, groupMember).Do()
				return err
			})
			if err != nil {
				return fmt.Errorf("Error creating groupMember: %s, %s", err, email)
			}
			log.Printf("[INFO] Created groupMember: %s", createdGroupMember.Email)
		}
	}

	return nil
}

func deleteMember(email, gid string, config *Config) (err error) {
	err = retry(func() error {
		err = config.directory.Members.Delete(gid, email).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error deleting member: %s", err)
	}
	return nil
}
