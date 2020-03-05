package gsuite

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	directory "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

// isDuplicateError returns true when the error is googleapi 409 "Entity already exists".
func isDuplicateError(err error) bool {
	gerr, ok := err.(*googleapi.Error)
	if !ok {
		return false
	}
	return gerr.Code == 409
}

// groupMatchesActual checks if a Group made for creation matches a currently existing (actual) group.
// If the returned error is non-nil, there is a nonrecoverable problem. Otherwise, true is
// returned when the groups match.  A group matches an actual group when the names, emails
// and aliases are the same. In this or similar cases there is no error and true is
// returned. If the names are different, it is not a match, with no error. If the two
// groups are incompatible, for example the names match but the emails are different, then
// a human-readable error is returned. The description is ignored.
func groupMatchesActual(group, actual *directory.Group) (bool, error) {
	if group.Name != actual.Name {
		return false, nil
	}
	if group.Email != actual.Email {
		return false, fmt.Errorf("Emails for group %s do not match: %s vs %s", group.Name, group.Email, actual.Email)
	}
	if len(group.Aliases) != len(actual.Aliases) {
		return false, fmt.Errorf("Aliases don't match for group %s: %d in group vs %d in actual", group.Name, len(group.Aliases), len(actual.Aliases))
	}
	if len(group.Aliases) > 0 {
		groupAliases := make([]string, len(group.Aliases))
		copy(groupAliases, group.Aliases)
		sort.Strings(groupAliases)
		actualAliases := make([]string, len(actual.Aliases))
		copy(actualAliases, actual.Aliases)
		sort.Strings(actualAliases)
		for i := 0; i < len(groupAliases); i++ {
			if groupAliases[i] != actualAliases[i] {
				return false, fmt.Errorf("Aliase mismatch for group %s: %s vs %s", groupAliases[i], actualAliases[i])
			}
		}
	}
	return true, nil
}

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupImporter,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},

			"aliases": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"direct_members_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"admin_created": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"non_editable_aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"ignore_duplicates": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,  // Will default to false.
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	group := &directory.Group{
		Email: strings.ToLower(d.Get("email").(string)),
	}

	if v, ok := d.GetOk("name"); ok {
		log.Printf("[DEBUG] Setting group name: %s", v.(string))
		group.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		log.Printf("[DEBUG] Setting group description: %s", v.(string))
		group.Description = v.(string)
	}

	ignoreDuplicates := false
	if v, ok := d.GetOk("ignore_duplicates"); ok {
		log.Printf("[DEBUG] Setting ignore duplicates: %t", v.(bool))
		ignoreDuplicates = v.(bool)
	}

	var createdGroup *directory.Group
	var err error
	err = retryPassDuplicate(func() error {
		createdGroup, err = config.directory.Groups.Insert(group).Do()
		return err
	}, config.TimeoutMinutes)

	isDuplicate := false
	if ignoreDuplicates && isDuplicateError(err) {
		domainSep := strings.LastIndex(group.Email, "@")
		if domainSep == -1 {
			return fmt.Errorf("[ERROR] Could not find domain in %s", group.Email)
		}
		domain := group.Email[domainSep + 1:]
		log.Printf("[DEBUG] listing groups to match duplicate for %s in %s", group.Name, domain)
		var groupsList *directory.Groups
		err = retry(func() error {
			groupsList, err = config.directory.Groups.List().Domain(domain).Query(fmt.Sprintf("name=%s", group.Name)).Do()
			return err
		}, config.TimeoutMinutes)
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing groups to match for duplicate group %s: %s", group.Name, err.Error())
		}
		createdGroup = nil
		log.Printf("[DEBUG] found %d groups to match to %s", len(groupsList.Groups), group.Name)
		for _, foundGroup := range groupsList.Groups {
			var match bool
			if match, err = groupMatchesActual(group, foundGroup); err != nil {
				return fmt.Errorf("[ERROR] Unresolvable duplicate group mismatch: %s", err)
			} else if match {
				createdGroup = foundGroup
				break
			}
		}
		if createdGroup == nil {
			return fmt.Errorf("[ERROR] No match found for duplicate group %s", group.Name)
		}
		log.Printf("[DEBUG] found matching duplicate for %s", group.Name)
		isDuplicate = true
	} else if err != nil {
		return fmt.Errorf("[ERROR] Error creating group: %s", err)
	}

	if !isDuplicate {
		// Handle group aliases
		aliasesCount := d.Get("aliases.#").(int)
		for i := 0; i < aliasesCount; i++ {
			cfgAlias := d.Get(fmt.Sprintf("aliases.%d", i)).(string)
			err = retry(func() error {
				alias := &directory.Alias{
					Alias: cfgAlias,
				}
				_, err = config.directory.Groups.Aliases.Insert(d.Id(), alias).Do()
				return err
			}, config.TimeoutMinutes)
		}

		if err != nil {
			return fmt.Errorf("[ERROR] Error creating group aliases: %s", err)
		}
	}

	// Try to read the group, retrying for 404's
	err = retryNotFound(func() error {
		group, err = config.directory.Groups.Get(createdGroup.Id).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return fmt.Errorf("[ERROR] Taking too long to create this group: %s", err)
	}

	d.SetId(createdGroup.Id)
	log.Printf("[INFO] Created group: %s", createdGroup.Email)
	return resourceGroupRead(d, meta)
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	group := &directory.Group{}
	nullFields := []string{}

	if d.HasChange("email") {
		log.Printf("[DEBUG] Updating group email: %s", d.Get("email").(string))
		group.Email = strings.ToLower(d.Get("email").(string))
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			log.Printf("[DEBUG] Updating group name: %s", v.(string))
			group.Name = v.(string)
		} else {
			log.Printf("[DEBUG] Removing group name")
			group.Name = ""
			nullFields = append(nullFields, "name")
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			log.Printf("[DEBUG] Updating group description: %s", v.(string))
			group.Description = v.(string)
		} else {
			log.Printf("[DEBUG] Removing group description")
			group.Description = ""
			nullFields = append(nullFields, "description")
		}
	}

	if len(nullFields) > 0 {
		group.NullFields = nullFields
	}

	var updatedGroup *directory.Group
	var err error
	err = retry(func() error {
		updatedGroup, err = config.directory.Groups.Patch(d.Id(), group).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return fmt.Errorf("[ERROR] Error updating group: %s", err)
	}

	// Handle group aliases
	var aliasesResponse *directory.Aliases
	err = retry(func() error {
		aliasesResponse, err = config.directory.Groups.Aliases.List(d.Id()).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return fmt.Errorf("[ERROR] Could not list group aliases: %s", err)
	}

	for _, v := range aliasesResponse.Aliases {
		c, ok := v.(map[string]interface{})
		if ok {
			alias := c["alias"].(string)
			log.Printf("[DEBUG] Removing alias: %s", alias)
			err = config.directory.Groups.Aliases.Delete(d.Id(), alias).Do()
		}
	}

	if err != nil {
		return fmt.Errorf("[ERROR] Error removing group aliases: %s", err)
	}

	aliasesCount := d.Get("aliases.#").(int)
	for i := 0; i < aliasesCount; i++ {
		cfgAlias := d.Get(fmt.Sprintf("aliases.%d", i)).(string)
		err = retry(func() error {
			alias := &directory.Alias{
				Alias: cfgAlias,
			}
			_, err = config.directory.Groups.Aliases.Insert(d.Id(), alias).Do()
			return err
		}, config.TimeoutMinutes)
	}

	if err != nil {
		return fmt.Errorf("[ERROR] Error creating group aliases: %s", err)
	}

	log.Printf("[INFO] Updated group: %s", updatedGroup.Email)
	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var group *directory.Group
	var err error
	err = retry(func() error {
		group, err = config.directory.Groups.Get(d.Id()).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Group %q", d.Get("name").(string)))
	}

	d.SetId(group.Id)
	d.Set("direct_members_count", group.DirectMembersCount)
	d.Set("admin_created", group.AdminCreated)
	d.Set("aliases", group.Aliases)
	d.Set("non_editable_aliases", group.NonEditableAliases)
	d.Set("description", group.Description)
	d.Set("name", group.Name)

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var err error
	err = retry(func() error {
		err = config.directory.Groups.Delete(d.Id()).Do()
		return err
	}, config.TimeoutMinutes)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting group: %s", err)
	}

	d.SetId("")
	return nil
}

// Allow importing using any key (id, email, alias)
func resourceGroupImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	id, err := config.directory.Groups.Get(d.Id()).Do()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error fetching group. Make sure the group exists: %s ", err)
	}

	d.SetId(id.Id)
	d.Set("email", id.Email)
	d.Set("description", id.Description)
	d.Set("name", id.Name)

	return []*schema.ResourceData{d}, nil
}
