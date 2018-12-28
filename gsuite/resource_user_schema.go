package gsuite

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	directory "google.golang.org/api/admin/directory/v1"
)

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

	var fieldEntries []map[string]interface{}
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
		fieldEntries = append(fieldEntries, map[string]interface{}{
			"field_name":       spec.FieldName,
			"field_type":       spec.FieldType,
			"multi_valued":     spec.MultiValued,
			"read_access_type": spec.ReadAccessType,
			"indexed":          *spec.Indexed,
		})
	}

	var (
		created *directory.Schema
		err     error
	)

	err = retry(func() error {
		created, err = config.directory.Schemas.Insert(config.CustomerId, userSchema).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error creating user schema: %s", err)
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
		read, err = config.directory.Schemas.Get(config.CustomerId, d.Id()).Do()
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

// TODO: resourceUserSchemaUpdate
func resourceUserSchemaUpdate(d *schema.ResourceData, meta interface{}) error {
	panic("update")
	return nil
}

// TODO: resourceUserSchemaDelete
func resourceUserSchemaDelete(d *schema.ResourceData, meta interface{}) error {
	panic("delete")
	return nil
}

// TODO: resourceUserSchemaImporter
func resourceUserSchemaImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	panic("import")
	return nil, nil
}
