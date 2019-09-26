package gsuite

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// userAttrMappings is a mapping of top level keys in a gsuite_user_attributes
// data source to a struct defining type of of value expected as input and
// whether or not the value is a list. See also userAttrMapping.schema() which
// constructs the underlying *schema.Schema at runtime.
var userAttrMappings = map[string]*userAttrMapping{
	"string":   {schema.TypeString, false},
	"strings":  {schema.TypeString, true},
	"bool":     {schema.TypeBool, false},
	"bools":    {schema.TypeBool, true},
	"integer":  {schema.TypeInt, false},
	"integers": {schema.TypeInt, true},
	"double":   {schema.TypeFloat, false},
	"doubles":  {schema.TypeFloat, true},
	"date":     {schema.TypeString, false},
	"dates":    {schema.TypeString, true},
	"email":    {schema.TypeString, false},
	"emails":   {schema.TypeString, true},
	"phone":    {schema.TypeString, false},
	"phones":   {schema.TypeString, true},
}

type userAttrMapping struct {
	valueType schema.ValueType
	list      bool
}

func (s *userAttrMapping) schema() *schema.Schema {
	value := &schema.Schema{Required: true, Type: s.valueType}
	if s.list {
		value.Type = schema.TypeList
		value.Elem = &schema.Schema{Type: s.valueType}
	}
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": value,
			},
		},
	}
}

func dataUserAttributes() *schema.Resource {
	resource := &schema.Resource{
		Read: dataUserAttributesRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	for name, mapping := range userAttrMappings {
		resource.Schema[name] = mapping.schema()
	}
	return resource
}

type entry struct {
	Value interface{} `json:"value"`
}

// MarshalJSON converts the interface value to a reasonable string
// representation. Without this some types of values (basically anything not a string)
// ends up without quotes around it. While Google's directory SDK does allow you
// to pass in some types of values without quotes other types, like floats, require
// strings. Placing quotes around *all* types of values provides the most
// consistent behavior overall.
func (e *entry) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"value":"%v"}`, e.Value)), nil
}

func dataUserAttributesRead(d *schema.ResourceData, _ interface{}) error {
	customAttributes := map[string]interface{}{}

	for name, mapping := range userAttrMappings {
		if mapping.list {
			if statements, ok := d.GetOk(name); ok {
				for _, statement := range statements.(*schema.Set).List() {
					stmt := statement.(map[string]interface{})
					var values []*entry
					for _, value := range stmt["value"].([]interface{}) {
						values = append(values, &entry{Value: value})
					}
					customAttributes[stmt["name"].(string)] = values
				}
			}
			continue
		}

		if statements, ok := d.GetOk(name); ok {
			for _, statement := range statements.(*schema.Set).List() {
				stmt := statement.(map[string]interface{})
				customAttributes[stmt["name"].(string)] = stmt["value"]
			}
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
