package gsuite

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	directory "google.golang.org/api/admin/directory/v1"
)

func schema() *schema.Resource {
	return &schema.Resource{
		Create: schemaCreate,
		Read:   schemaRead,
		Update: schemaUpdate,
		Delete: schemaDelete,
		Importer: &schema.ResourceImporter{
			State: schemaImporter,
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


// TODO: schemaCreate
func schemaCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: schemaRead
func schemaRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: schemaUpdate
func schemaUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: schemaDelete
func schemaDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: schemaImporter
func schemaImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}