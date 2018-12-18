package gsuite

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	directory "google.golang.org/api/admin/directory/v1"
)

var schemaMember = map[string]*schema.Schema{
	"etag": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},

	"kind": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},

	"status": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},

	"type": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},

	"role": &schema.Schema{
		Type:     schema.TypeString,
		Default:  "MEMBER",
		Optional: true,
	},

	"email": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		StateFunc: func(val interface{}) string {
			return strings.ToLower(val.(string))
		},
	},
}

var schemaGroup = map[string]*schema.Schema{
	"group": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		StateFunc: func(val interface{}) string {
			return strings.ToLower(val.(string))
		},
	},
}

var schemaMembership = mergeSchemas(schemaGroup, schemaMember)

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMemberCreate,
		Read:   resourceGroupMemberRead,
		Update: resourceGroupMemberUpdate,
		Delete: resourceGroupMemberDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupMemberImporter,
		},

		Schema: schemaMembership,
	}
}

func resourceGroupMemberCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	group := d.Get("group").(string)

	groupMember := &directory.Member{
		Role:  strings.ToUpper(d.Get("role").(string)),
		Email: strings.ToLower(d.Get("email").(string)),
	}

	var createdGroupMember *directory.Member
	var err error
	err = retry(func() error {
		createdGroupMember, err = config.directory.Members.Insert(group, groupMember).Do()
		return err
	})

	if err != nil {
		if !strings.Contains(err.Error(), "Member already exists") {
			return fmt.Errorf("error creating group member: %s", err)
		}
		log.Printf("[INFO] %s already part of this group. attempting to update", groupMember.Email)

		var existingGroupMembers *directory.Members
		err = retry(func() error {
			existingGroupMembers, err = config.directory.Members.List(group).Do()
			return err
		})
		if err != nil {
			return fmt.Errorf("error locating existing group members: %s", err)
		}
		var locatedGroupMember *directory.Member
		for _, existingGroupMember := range existingGroupMembers.Members {
			if existingGroupMember.Email == groupMember.Email {
				locatedGroupMember = existingGroupMember
				break
			}
		}
		if locatedGroupMember == nil {
			return fmt.Errorf("error locating existing group member %s", groupMember.Email)
		}
		log.Printf("[INFO] found existing group member %s", locatedGroupMember.Email)

		var err error
		err = retry(func() error {
			_, err = config.directory.Members.Patch(group, locatedGroupMember.Id, groupMember).Do()
			return err
		})

		if err != nil {
			return fmt.Errorf("error updating existing group member: %s", err)
		}
		log.Printf("[INFO] Updated group member: %s", groupMember.Email)
		d.SetId(locatedGroupMember.Id)
	} else {
		log.Printf("[INFO] Created group member: %s", createdGroupMember.Email)
		d.SetId(createdGroupMember.Id)
	}

	return resourceGroupMemberRead(d, meta)
}

func resourceGroupMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	groupMember := &directory.Member{}
	nullFields := []string{}

	if d.HasChange("email") {
		log.Printf("[DEBUG] Updating groupMember email (recreating member): %s", d.Get("email").(string))
		groupMember.Email = strings.ToLower(d.Get("email").(string))
	}

	if d.HasChange("role") {
		log.Printf("[DEBUG] Updating groupMember role: %s to %s", d.Get("email").(string), d.Get("role").(string))
		groupMember.Role = strings.ToUpper(d.Get("role").(string))
	}

	if len(nullFields) > 0 {
		groupMember.NullFields = nullFields
	}

	var updatedGroupMember *directory.Member
	var err error
	err = retry(func() error {
		updatedGroupMember, err = config.directory.Members.Patch(strings.ToLower(d.Get("group").(string)), d.Id(), groupMember).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("error updating group member: %s", err)
	}

	log.Printf("[INFO] Updated groupMember: %s", updatedGroupMember.Email)
	return resourceGroupMemberRead(d, meta)
}

func resourceGroupMemberRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var groupMember *directory.Member
	var err error
	err = retry(func() error {
		groupMember, err = config.directory.Members.Get(strings.ToLower(d.Get("group").(string)), d.Id()).Do()
		return err
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Group member %q", d.Get("email").(string)))
	}

	d.SetId(groupMember.Id)
	d.Set("role", strings.ToUpper(groupMember.Role))
	d.Set("email", strings.ToLower(groupMember.Email))
	d.Set("etag", groupMember.Etag)
	d.Set("kind", groupMember.Kind)
	d.Set("status", groupMember.Status)
	d.Set("type", groupMember.Type)

	return nil
}

func resourceGroupMemberDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var err error
	err = retry(func() error {
		err = config.directory.Members.Delete(strings.ToLower(d.Get("group").(string)), d.Id()).Do()
		return err
	})
	if err != nil {
		return fmt.Errorf("Error deleting group member: %s", err)
	}

	d.SetId("")
	return nil
}

// Allow importing using [group]{:,/}[email]
func resourceGroupMemberImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	s := strings.Split(d.Id(), ":")
	if len(s) < 2 {
		s = strings.Split(d.Id(), "/")
	}

	if len(s) < 2 {
		return nil, fmt.Errorf("import via [group]:[member email] or [group]/[member email]")
	}
	group, member := strings.ToLower(s[0]), strings.ToLower(s[1])

	id, err := config.directory.Members.Get(group, member).Do()

	if err != nil {
		return nil, fmt.Errorf("error fetching member, make sure the member exists: %s ", err)
	}

	d.SetId(id.Id)
	d.Set("group", group)
	d.Set("role", id.Role)
	d.Set("email", id.Email)
	d.Set("etag", id.Etag)
	d.Set("kind", id.Kind)
	d.Set("status", id.Status)
	d.Set("type", id.Type)

	return []*schema.ResourceData{d}, nil
}
