package gsuite

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	directory "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

var schemaGroupMembersEmail = map[string]*schema.Schema{
	"email": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: false,
		StateFunc: func(val interface{}) string {
			return strings.ToLower(val.(string))
		},
		ValidateFunc: validateEmail,
	},
}

var schemaGroupMembers = mergeSchemas(schemaMember, schemaGroupMembersEmail)

func resourceGroupMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMembersCreate,
		Read:   resourceGroupMembersRead,
		Update: resourceGroupMembersUpdate,
		Delete: resourceGroupMembersDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupMembersImporter,
		},

		Schema: map[string]*schema.Schema{
			"group_email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"member": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: schemaGroupMembers,
				},
			},
		},
	}
}

func resourceGroupMembersRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Reading gsuite_group_members")
	config := meta.(*Config)

	groupEmail := d.Id()

	members, err := getAPIMembers(groupEmail, config)

	if err != nil {
		return err
	}

	d.Set("group_email", strings.ToLower(groupEmail))
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

	for _, rawMember := range d.Get("member").(*schema.Set).List() {
		member := rawMember.(map[string]interface{})
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
	for _, rawMember := range d.Get("member").(*schema.Set).List() {
		member := rawMember.(map[string]interface{})
		members = append(members, member)
	}
	return members
}

func createOrUpdateGroupMembers(d *schema.ResourceData, meta interface{}) (string, error) {
	config := meta.(*Config)
	groupEmail := strings.ToLower(d.Get("group_email").(string))

	// Get members from config
	cfgMembers := resourceMembers(d)

	// Get members from API
	apiMembers, err := getAPIMembers(groupEmail, config)
	if err != nil {
		return groupEmail, fmt.Errorf("[ERROR] Error updating memberships: %v", err)
	}
	// This call removes any members that aren't defined in cfgMembers,
	// and adds all of those that are
	err = reconcileMembers(d, cfgMembers, membersToCfg(apiMembers), config, groupEmail)
	if err != nil {
		return groupEmail, fmt.Errorf("[ERROR] Error updating memberships: %v", err)
	}

	return groupEmail, nil
}

// This function ensures that the members of a group exactly match that
// in a config by deleting any members that are returned by the API but not present
// in the config
func reconcileMembers(d *schema.ResourceData, cfgMembers, apiMembers []map[string]interface{}, config *Config, gid string) error {

	// Helper to convert slice to map
	m := func(vals []map[string]interface{}) map[string]map[string]interface{} {
		sm := make(map[string]map[string]interface{})
		for _, member := range vals {
			email := strings.ToLower(member["email"].(string))
			member["email"] = strings.ToLower(member["email"].(string))
			sm[email] = member
		}
		return sm
	}

	cfgMap := m(cfgMembers)
	log.Println("[DEBUG] Members in cfg: ", cfgMap)
	apiMap := m(apiMembers)
	log.Println("[DEBUG] Member in API: ", apiMap)

	var cfgRole, apiRole string

	for k, apiMember := range apiMap {
		if cfgMember, ok := cfgMap[k]; !ok {
			// The member in the API is not in the config; disable it.
			log.Printf("[DEBUG] Member in API not in config. Disabling it: %s", k)
			err := deleteMember(k, gid, config)
			if err != nil {
				return err
			}
		} else {
			// The member exists in the config and the API
			// If role has changed update, otherwise do nothing
			cfgRole = strings.ToUpper(cfgMember["role"].(string))
			apiRole = strings.ToUpper(apiMember["role"].(string))
			if cfgRole != apiRole {
				groupMember := &directory.Member{
					Role: cfgRole,
				}

				if cfgRole != "MEMBER" {
					return fmt.Errorf("[ERROR] Error updating groupMember (%s): nested groups should be role MEMBER", cfgMember["email"].(string))
				}

				var updatedGroupMember *directory.Member
				var err error
				err = retryNotFound(func() error {
					updatedGroupMember, err = config.directory.Members.Patch(
						strings.ToLower(d.Get("group_email").(string)),
						strings.ToLower(cfgMember["email"].(string)),
						groupMember).Do()
					return err
				}, config.TimeoutMinutes)

				if err != nil {
					return fmt.Errorf("[ERROR] Error updating groupMember: %s", err)
				}

				log.Printf("[INFO] Updated groupMember: %s", updatedGroupMember.Email)
			}

			// Delete from cfgMap, we have already handled it
			delete(cfgMap, k)
		}
	}

	// Upsert memberships which are present in the config, but not in the api
	for email := range cfgMap {
		err := upsertMember(email, gid, cfgMap[email]["role"].(string), config)
		if err != nil {
			return err
		}
	}
	return nil
}

// Retrieve a group's members from the API
func getAPIMembers(groupEmail string, config *Config) ([]*directory.Member, error) {
	// Get members from the API
	groupMembers := make([]*directory.Member, 0)
	token := ""
	var membersResponse *directory.Members
	var err error
	for paginate := true; paginate; {

		err = retry(func() error {
			membersResponse, err = config.directory.Members.List(groupEmail).PageToken(token).Do()
			return err
		}, config.TimeoutMinutes)

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

func upsertMember(email, groupEmail, role string, config *Config) error {
	var err error
	groupMember := &directory.Member{
		Role:  strings.ToUpper(role),
		Email: strings.ToLower(email),
	}

	// Check if the email address belongs to a user, or to a group
	// we need to make sure, because we need to use different logic
	var isGroup = true
	err = retry(func() error {
		_, err := config.directory.Groups.Get(email).Do()
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			isGroup = false
			log.Printf("[DEBUG] Setting isGroup to false for %s after getting a 404", email)
			return nil
		}
		return err
	}, config.TimeoutMinutes)

	if isGroup == true {
		if role != "MEMBER" {
			return fmt.Errorf("[ERROR] Error creating groupMember (%s): nested groups should be role MEMBER", email)
		}

		var isGroupMember = true

		// Grab the group as a directory member of the current group
		err = retry(func() error {
			_, err := config.directory.Members.Get(groupEmail, email).Do()

			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
				isGroupMember = false
				log.Printf("[DEBUG] Setting isGroupMember to false for %s after getting a 404", email)
				return nil
			}

			return err
		}, config.TimeoutMinutes)

		// Based on the err return, either add as a new member, or update
		if isGroupMember == false {
			var createdGroupMember *directory.Member
			err = retry(func() error {
				createdGroupMember, err = config.directory.Members.Insert(groupEmail, groupMember).Do()
				return err
			}, config.TimeoutMinutes)
			if err != nil {
				return fmt.Errorf("[ERROR] Error creating groupMember: %s, %s", err, email)
			}
			log.Printf("[INFO] Created groupMember: %s", createdGroupMember.Email)
		} else {
			var updatedGroupMember *directory.Member
			err = retryNotFound(func() error {
				updatedGroupMember, err = config.directory.Members.Update(groupEmail, email, groupMember).Do()
				return err
			}, config.TimeoutMinutes)
			if err != nil {
				return fmt.Errorf("[ERROR] Error updating groupMember: %s, %s", err, email)
			}
			log.Printf("[INFO] Updated groupMember: %s", updatedGroupMember.Email)
		}
	}

	if isGroup == false {
		// Basically the same check as group, but using a more apt method "HasMember"
		// specifically meant for users
		var hasMemberResponse *directory.MembersHasMember
		var err error
		err = retry(func() error {
			hasMemberResponse, err = config.directory.Members.HasMember(groupEmail, email).Do()
			if err == nil {
				return err
			}

			// When a user does not exist, the API returns a 400 "memberKey, required"
			// Returning a friendly message
			if gerr, ok := err.(*googleapi.Error); ok && (gerr.Errors[0].Reason == "required" && gerr.Code == 400) {
				return fmt.Errorf("[ERROR] Error adding groupMember %s, please make sure the user exists beforehand", email)
			}
			return err
		}, config.TimeoutMinutes)
		if err != nil {
			return createGroupMember(groupMember, groupEmail, config)
		}

		if hasMemberResponse.IsMember == true {
			var updatedGroupMember *directory.Member
			err = retryNotFound(func() error {
				updatedGroupMember, err = config.directory.Members.Update(groupEmail, email, groupMember).Do()
				return err
			}, config.TimeoutMinutes)
			if err != nil {
				return fmt.Errorf("[ERROR] Error updating groupMember: %s, %s", err, email)
			}
			log.Printf("[INFO] Updated groupMember: %s", updatedGroupMember.Email)
		} else {
			return createGroupMember(groupMember, groupEmail, config)
		}
	}

	return nil
}

func createGroupMember(groupMember *directory.Member, groupEmail string, config *Config) (err error) {
	var createdGroupMember *directory.Member
	err = retry(func() error {
		createdGroupMember, err = config.directory.Members.Insert(groupEmail, groupMember).Do()
		return err
	}, config.TimeoutMinutes)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating groupMember: %s, %s", err, groupMember.Email)
	}
	log.Printf("[INFO] Created groupMember: %s", createdGroupMember.Email)

	return nil
}

func deleteMember(email, groupEmail string, config *Config) (err error) {
	err = retry(func() error {
		err = config.directory.Members.Delete(groupEmail, email).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting member: %s", err)
	}
	return nil
}

// Allow importing using any groupKey (id, email, alias)
func resourceGroupMembersImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Importing gsuite_group_members")
	config := meta.(*Config)

	var group *directory.Group
	var err error
	err = retry(func() error {
		group, err = config.directory.Groups.Get(strings.ToLower(d.Id())).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error fetching group. Make sure the group exists: %s ", err)
	}

	d.SetId(group.Email)

	return []*schema.ResourceData{d}, nil
}
