package gsuite

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	directory "google.golang.org/api/admin/directory/v1"
)

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceSchemaCreate,
		Read:   resourceSchemaRead,
		Update: resourceSchemaUpdate,
		Delete: resourceSchemaDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSchemaImporter,
		},

		Schema: map[string]*schema.Schema{
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"schema_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"field": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"field_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"display_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"field_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"multi_valued": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},

						"read_access_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"indexed": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},

						"field": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{

								Schema: map[string]*schema.Schema{
									"min_value": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},

									"max_value": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// TODO: resourceSchemaCreate
func resourceSchemaCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: resourceSchemaRead
func resourceSchemaRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: resourceSchemaUpdate
func resourceSchemaUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: resourceSchemaDelete
func resourceSchemaDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: resourceSchemaImporter
func resourceSchemaImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
