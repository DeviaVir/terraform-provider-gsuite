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
	},
}

var schemaGroup = map[string]*schema.Schema{
	"group": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
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
		Role: d.Get("role").(string),
		Email: d.Get("email").(string),
	}

	var createdGroupMember *directory.Member
	var err error
	err = retry(func() error {
		createdGroupMember, err = config.directory.Members.Insert(group, groupMember).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error creating group member: %s", err)
	}

  d.SetId(createdGroupMember.Id)
	log.Printf("[INFO] Created group: %s", createdGroupMember.Email)
	return resourceGroupMemberRead(d, meta)
}

func resourceGroupMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	groupMember := &directory.Member{}
	nullFields := []string{}

	if d.HasChange("email") {
		log.Printf("[DEBUG] Updating groupMember email: %s", d.Get("email").(string))
		groupMember.Email = d.Get("email").(string)
	}

	if len(nullFields) > 0 {
		groupMember.NullFields = nullFields
	}

	var updatedGroupMember *directory.Member
	var err error
	err = retry(func() error {
		updatedGroupMember, err = config.directory.Members.Patch(d.Get("group").(string), d.Id(), groupMember).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error updating group member: %s", err)
	}

	log.Printf("[INFO] Updated groupMember: %s", updatedGroupMember.Email)
	return resourceGroupMemberRead(d, meta)
}

func resourceGroupMemberRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var groupMember *directory.Member
	var err error
	err = retry(func() error {
		groupMember, err = config.directory.Members.Get(d.Get("group").(string), d.Id()).Do()
		return err
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Group member %q", d.Get("name").(string)))
	}

  d.SetId(groupMember.Id)
	d.Set("email", groupMember.Email)
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
		err = config.directory.Members.Delete(d.Get("group").(string), d.Id()).Do()
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
		return nil, fmt.Errorf("Import via [group]:[member email] or [group]/[member email]!")
	}
	group, member := s[0], s[1]

	id, err := config.directory.Members.Get(group, member).Do()

	if err != nil {
		return nil, fmt.Errorf("Error fetching member. Make sure the member exists: %s ", err)
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
