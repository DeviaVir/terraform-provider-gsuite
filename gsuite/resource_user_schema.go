package gsuite

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	directory "google.golang.org/api/admin/directory/v1"
)

const fieldKind = "admin#directory#schema#fieldspec"

func resourceUserSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserSchemaCreate,
		Read:   resourceUserSchemaRead,
		Update: resourceUserSchemaUpdate,
		Delete: resourceUserSchemaDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserSchemaImporter,
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"schema_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"field": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"field_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{
									"BOOL", "DATA", "DOUBLE", "EMAIL",
									"INT64", "PHONE", "STRING",
								},
								false,
							),
						},

						"multi_valued": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"read_access_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ADMINS_AND_SELF",
							ValidateFunc: validation.StringInSlice(
								[]string{"ADMINS_AND_SELF", "ALL_DOMAIN_USERS"},
								false,
							),
						},

						"indexed": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"range": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_value": {
										Type:     schema.TypeFloat,
										Optional: false,
									},

									"max_value": {
										Type:     schema.TypeFloat,
										Optional: false,
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

// TODO: resourceUserSchemaCreate
func resourceUserSchemaCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	userSchema := &directory.Schema{
		SchemaName: d.Get("schema_name").(string),
	}

	for i := 0; i < d.Get("field.#").(int); i++ {
		fields := d.Get(fmt.Sprintf("field.%d", i)).(map[string]interface{})

		indexed := fields["indexed"].(bool)
		spec := &directory.SchemaFieldSpec{
			FieldName:      fields["field_name"].(string),
			FieldType:      fields["field_type"].(string),
			MultiValued:    fields["multi_valued"].(bool),
			ReadAccessType: fields["read_access_type"].(string),
			Indexed:        &indexed,
			Kind:           fieldKind,
		}

		if values, ok := fields["range"].(map[string]interface{}); ok {
			var (
				minValue float64
				maxValue float64
				err      error
			)
			switch spec.FieldType {
			case "DOUBLE":
				minValue, err = strconv.ParseFloat(values["min_value"].(string), 64)
				if err != nil {
					return err
				}

				maxValue, err = strconv.ParseFloat(values["max_value"].(string), 64)
				if err != nil {
					return err
				}

			case "INT64":
				min, err := strconv.Atoi(values["min_value"].(string))
				if err != nil {
					return err
				}
				minValue = float64(min)

				max, err := strconv.Atoi(values["max_value"].(string))
				if err != nil {
					return err
				}
				maxValue = float64(max)
			}

			spec.NumericIndexingSpec = &directory.SchemaFieldSpecNumericIndexingSpec{
				MinValue: minValue,
				MaxValue: maxValue,
			}
		}

		userSchema.Fields = append(userSchema.Fields, spec)
	}

	var (
		createdUserSchema *directory.Schema
		err               error
	)

	err = retry(func() error {
		createdUserSchema, err = config.directory.Schemas.Insert(config.CustomerId, userSchema).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error creating user schema: %s", err)
	}

	log.Printf("[INFO] Created user schema: %s", createdUserSchema.SchemaName)
	return nil
}

// TODO: resourceUserSchemaRead
func resourceUserSchemaRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: resourceUserSchemaUpdate
func resourceUserSchemaUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: resourceUserSchemaDelete
func resourceUserSchemaDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// TODO: resourceUserSchemaImporter
func resourceUserSchemaImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
