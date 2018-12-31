package gsuite

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	directory "google.golang.org/api/admin/directory/v1"
)

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
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},

			"aliases": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"direct_members_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"admin_created": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"non_editable_aliases": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

	var createdGroup *directory.Group
	var err error
	err = retry(func() error {
		createdGroup, err = config.directory.Groups.Insert(group).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error creating group: %s", err)
	}

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
		})
	}

	if err != nil {
		return fmt.Errorf("Error creating group aliases: %s", err)
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
	})

	if err != nil {
		return fmt.Errorf("Error updating group: %s", err)
	}

	// Handle group aliases
	var aliasesResponse *directory.Aliases
	err = retry(func() error {
		aliasesResponse, err = config.directory.Groups.Aliases.List(d.Id()).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Could not list group aliases: %s", err)
	}

	for _, v := range aliasesResponse.Aliases {
		c, ok := v.(map[string]interface {})
		if ok {
			alias := c["alias"].(string)
			log.Printf("[DEBUG] Removing alias: %s", alias)
			err = config.directory.Groups.Aliases.Delete(d.Id(), alias).Do()
		}
	}

	if err != nil {
		return fmt.Errorf("Error removing group aliases: %s", err)
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
		})
	}

	if err != nil {
		return fmt.Errorf("Error creating group aliases: %s", err)
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
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Group %q", d.Get("name").(string)))
	}

	d.SetId(group.Id)
	d.Set("direct_members_count", group.DirectMembersCount)
	d.Set("admin_created", group.AdminCreated)
	d.Set("aliases", group.Aliases)
	d.Set("non_editable_aliases", group.NonEditableAliases)

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var err error
	err = retry(func() error {
		err = config.directory.Groups.Delete(d.Id()).Do()
		return err
	})
	if err != nil {
		return fmt.Errorf("Error deleting group: %s", err)
	}

	d.SetId("")
	return nil
}

// Allow importing using any key (id, email, alias)
func resourceGroupImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	id, err := config.directory.Groups.Get(d.Id()).Do()

	if err != nil {
		return nil, fmt.Errorf("Error fetching group. Make sure the group exists: %s ", err)
	}

	d.SetId(id.Id)
	d.Set("email", id.Email)
	d.Set("description", id.Description)
	d.Set("name", id.Name)

	return []*schema.ResourceData{d}, nil
}
