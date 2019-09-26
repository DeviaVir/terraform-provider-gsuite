package gsuite

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	directory "google.golang.org/api/admin/directory/v1"
)

// myCustomerID is a stand-in for the `customerId` field that's a required
// argument to several API requests. This save us from having to query an
// external resource and/or require it as an argument on the provider itself.
const myCustomerID = "my_customer"

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
			"schema_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"field_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{
									"BOOL", "DATE", "DOUBLE", "EMAIL",
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

func resourceUserSchemaCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	userSchema := &directory.Schema{}
	if v, ok := d.GetOk("schema_name"); ok {
		value := v.(string)
		log.Printf("[DEBUG] Setting %s: %s", "schema_name", value)
		userSchema.SchemaName = value
	}

	if v, ok := d.GetOk("display_name"); ok {
		value := v.(string)
		log.Printf("[DEBUG] Setting %s: %s", "display_name", value)
		userSchema.DisplayName = value
	}

	if userSchema.DisplayName == "" {
		userSchema.DisplayName = userSchema.SchemaName
	}

	fields, err := getUserSchemaFieldSpecs(d)
	if err != nil {
		return err
	}
	userSchema.Fields = fields
	var created *directory.Schema

	err = retry(func() error {
		created, err = config.directory.Schemas.Insert(myCustomerID, userSchema).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("[ERROR] Error creating user schema: %s", err)
	}

	d.SetId(created.SchemaId)
	log.Printf("[INFO] Created user schema: %s", created.SchemaName)
	return resourceUserSchemaRead(d, meta)
}

func resourceUserSchemaRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var (
		read *directory.Schema
		err  error
	)
	err = retry(func() error {
		read, err = config.directory.Schemas.Get(myCustomerID, d.Id()).Do()
		return err
	})
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Schema %q", d.Get("schema_name").(string)))
	}

	d.SetId(read.SchemaId)
	d.Set("schema_id", read.SchemaId)
	d.Set("schema_name", read.SchemaName)
	d.Set("display_name", read.DisplayName)
	d.Set("field", read.Fields)

	return nil
}

func resourceUserSchemaUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userSchema, err := config.directory.Schemas.Get(myCustomerID, d.Id()).Do()
	if err != nil {
		return err
	}

	if d.HasChange("schema_name") {
		if v, ok := d.GetOk("schema_name"); ok {
			value := v.(string)
			log.Printf("[DEBUG] Updating schema %s: %s -> %s", "schema_name", value, d.Get("schema_name"))
			userSchema.SchemaName = value
		}
	}

	if d.HasChange("display_name") {
		if v, ok := d.GetOk("display_name"); ok {
			value := v.(string)
			log.Printf("[DEBUG] Updating schema %s: %s -> %s", "display_name", value, d.Get("display_name"))
			userSchema.DisplayName = value
		}
	}

	if d.HasChange("field") {
		specs, err := getUserSchemaFieldSpecs(d)
		if err != nil {
			return err
		}
		userSchema.Fields = specs
	}

	var updated *directory.Schema

	err = retry(func() error {
		updated, err = config.directory.Schemas.Update(myCustomerID, d.Id(), userSchema).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("[ERROR] Error updating user schema: %s", err)
	}

	log.Printf("[INFO] Updated schema: %s", updated.DisplayName)
	return resourceUserSchemaRead(d, meta)
}

func resourceUserSchemaDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	return retry(func() error {
		return config.directory.Schemas.Delete(myCustomerID, d.Id()).Do()
	})
}

func resourceUserSchemaImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	imported, err := config.directory.Schemas.Get(myCustomerID, d.Id()).Do()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error fetching schema. Make sure the schema exists: %s ", err)
	}

	d.SetId(imported.SchemaId)
	d.Set("schema_id", imported.SchemaId)
	d.Set("schema_name", imported.SchemaName)
	d.Set("display_name", imported.DisplayName)
	d.Set("field", imported.Fields)

	return []*schema.ResourceData{d}, nil
}

func getUserSchemaFieldSpecs(d *schema.ResourceData) ([]*directory.SchemaFieldSpec, error) {
	var specs []*directory.SchemaFieldSpec
	for i := 0; i < d.Get("field.#").(int); i++ {
		key := fmt.Sprintf("field.%d", i)
		fields := d.Get(key).(map[string]interface{})
		indexed := fields["indexed"].(bool)
		spec := &directory.SchemaFieldSpec{
			FieldName:      fields["field_name"].(string),
			FieldType:      fields["field_type"].(string),
			MultiValued:    fields["multi_valued"].(bool),
			ReadAccessType: fields["read_access_type"].(string),
			Indexed:        &indexed,
		}

		if fields["display_name"] != "" {
			spec.DisplayName = fields["display_name"].(string)
		}

		if spec.DisplayName == "" {
			spec.DisplayName = spec.FieldName
		}

		if values, ok := fields["range"].(map[string]interface{}); ok {
			var (
				minValue float64
				maxValue float64
				err      error
			)
			switch spec.FieldType {
			case "DOUBLE":
				if v, ok := values["min_value"]; ok {
					minValue, err = strconv.ParseFloat(v.(string), 64)
					if err != nil {
						return nil, err
					}
				}

				if v, ok := values["max_value"]; ok {
					maxValue, err = strconv.ParseFloat(v.(string), 64)
					if err != nil {
						return nil, err
					}
				}

			case "INT64":
				if v, ok := values["min_value"]; ok {
					integer, err := strconv.Atoi(v.(string))
					if err != nil {
						return nil, err
					}
					minValue = float64(integer)
				}
				if v, ok := values["max_value"]; ok {
					integer, err := strconv.Atoi(v.(string))
					if err != nil {
						return nil, err
					}
					maxValue = float64(integer)
				}
			}
			if minValue > 0 || maxValue > 0 {
				spec.NumericIndexingSpec = &directory.SchemaFieldSpecNumericIndexingSpec{
					MinValue: minValue,
					MaxValue: maxValue,
				}
			}
		}

		specs = append(specs, spec)
	}

	return specs, nil
}
