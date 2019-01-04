package gsuite

import (
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataUserAttributes() *schema.Resource {
	return &schema.Resource{
		Read: dataUserAttributesRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"string": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"bool": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"strings": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataUserAttributesRead(d *schema.ResourceData, meta interface{}) error {
	customAttributes := map[string]interface{}{}

	if statements, ok := d.GetOk("string"); ok {
		for _, statement := range statements.(*schema.Set).List() {
			stmt := statement.(map[string]interface{})
			customAttributes[stmt["name"].(string)] = stmt["value"].(string)
		}
	}

	if statements, ok := d.GetOk("bool"); ok {
		for _, statement := range statements.(*schema.Set).List() {
			stmt := statement.(map[string]interface{})
			customAttributes[stmt["name"].(string)] = stmt["value"].(bool)
		}
	}

	if statements, ok := d.GetOk("strings"); ok {
		for _, statement := range statements.(*schema.Set).List() {
			stmt := statement.(map[string]interface{})
			var values []interface{}
			for _, value := range stmt["value"].([]interface{}) {
				values = append(values, &struct {
					Value string `json:"value"`
				}{Value: value.(string)})
			}
			customAttributes[stmt["name"].(string)] = values
		}
	}

	out, err := json.Marshal(customAttributes)
	if err != nil {
		return err
	}

	outString := string(out)
	d.SetId(strconv.Itoa(hashcode.String(outString)))
	d.Set("json", outString)

	return nil
}
