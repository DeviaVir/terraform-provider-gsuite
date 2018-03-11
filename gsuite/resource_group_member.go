package gsuite

import (
	"fmt"
	"log"

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

	createdGroupMember, err := config.directory.Members.Insert(group, groupMember).Do()
	if err != nil {
		return fmt.Errorf("Error creating groupMember: %s", err)
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

	updatedGroupMember, err := config.directory.Members.Patch(d.Get("group").(string), d.Id(), groupMember).Do()
	if err != nil {
		return fmt.Errorf("Error updating groupMember: %s", err)
	}

	log.Printf("[INFO] Updated groupMember: %s", updatedGroupMember.Email)
	return resourceGroupMemberRead(d, meta)
}

func resourceGroupMemberRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	groupMember, err := config.directory.Members.Get(d.Get("group").(string), d.Id()).Do()
	if err != nil {
		return err
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

	err := config.directory.Members.Delete(d.Get("group").(string), d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting group: %s", err)
	}

	d.SetId("")
	return nil
}
