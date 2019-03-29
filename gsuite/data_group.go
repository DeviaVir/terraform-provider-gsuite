package gsuite

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	directory "google.golang.org/api/admin/directory/v1"
)

func dataGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataGroupRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},

			"aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
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

			"member": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: schemaGroupMembers,
				},
			},
		},
	}
}

func dataGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var group *directory.Group
	var err error
	err = retry(func() error {
		group, err = config.directory.Groups.Get(d.Get("email").(string)).Do()
		return err
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Group %q", d.Get("name").(string)))
	}

	members, err := getAPIMembers(d.Get("email").(string), config)

	d.SetId(group.Id)
	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("direct_members_count", group.DirectMembersCount)
	d.Set("admin_created", group.AdminCreated)
	d.Set("aliases", group.Aliases)
	d.Set("non_editable_aliases", group.NonEditableAliases)
	d.Set("member", membersToCfg(members))

	return nil
}
