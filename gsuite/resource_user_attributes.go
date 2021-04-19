package gsuite

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	directory "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

func resourceUserAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributesCreate,
		Read:   resourceUserAttributesRead,
		Update: resourceUserAttributesUpdate,
		Delete: resourceUserAttributesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserAttributesImporter,
		},

		Schema: map[string]*schema.Schema{
			"primary_email": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},

			"custom_schema": {
				Type:     schema.TypeList,
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
		},
	}
}

func resourceUserAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	user := &directory.User{}

	if v, ok := d.GetOk("primary_email"); ok {
		log.Printf("[DEBUG] Setting %s: %s", "primary_email", v.(string))
		user.PrimaryEmail = strings.ToLower(v.(string))
	}

	customSchemas := map[string]googleapi.RawMessage{}
	for i := 0; i < d.Get("custom_schema.#").(int); i++ {
		entry := d.Get(fmt.Sprintf("custom_schema.%d", i)).(map[string]interface{})
		customSchemas[entry["name"].(string)] = []byte(entry["value"].(string))

		log.Printf("[DEBUG] setting entry %v", entry)
	}
	if len(customSchemas) > 0 {
		user.CustomSchemas = customSchemas
	}

	var err error
	var updatedUser *directory.User
	err = retry(func() error {
		updatedUser, err = config.directory.Users.Patch(user.PrimaryEmail, user).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return fmt.Errorf("[ERROR] Error creating user fields: %s", err)
	}

	d.SetId(updatedUser.Id)
	log.Printf("[INFO] Created user fields: %s", updatedUser.PrimaryEmail)
	return resourceUserAttributesRead(d, meta)
}

func resourceUserAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	/*
		Defined in terms of delete + create, because resourceUserAttributesDelete will delete the information from a user by its stored id - even if we're actually changing the primary_email the gsuite_user_attributes resource points to - and resourceUserAttributesCreate will create the attributes for the current primary_email
	*/
	err := resourceUserAttributesDelete(d, meta)
	if err != nil {
		return err
	}

	err = resourceUserAttributesCreate(d, meta)
	if err != nil {
		return err
	}

	return resourceUserAttributesRead(d, meta)
}

func resourceUserAttributesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var user *directory.User
	var err error
	err = retry(func() error {
		user, err = config.directory.Users.Get(d.Id()).Projection("full").Do()
		if user != nil && user.Name == nil {
			return errors.New("Eventual consistency. Please try again")
		}
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("User %q", d.Id()))
	}

	d.SetId(user.Id)
	d.Set("primary_email", user.PrimaryEmail)

	err, flattenedCustomSchema := flattenCustomSchema(user.CustomSchemas)
	if err != nil {
		return err
	}

	if err = d.Set("custom_schema", flattenedCustomSchema); err != nil {
		return fmt.Errorf("Error setting custom_schema in state: %s", err.Error())
	}

	return nil
}

func resourceUserAttributesDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	user := &directory.User{}

	customSchemas := map[string]googleapi.RawMessage{}
	for i := 0; i < d.Get("custom_schema.#").(int); i++ {
		entry := d.Get(fmt.Sprintf("custom_schema.%d", i)).(map[string]interface{})

		customAttributes := map[string]interface{}{}
		err := json.Unmarshal([]byte(entry["value"].(string)), &customAttributes)

		if err != nil {
			return fmt.Errorf("Error unmarshalling custom attributes in resource: %s", err)
		}

		schemaBody := "{"
		for field := range customAttributes {
			schemaBody = schemaBody + fmt.Sprintf("\n  \"%s\": null,", field)
		}
		// remove the training comma
		schemaBody = schemaBody[:len(schemaBody)-1] + "\n}"

		customSchemas[entry["name"].(string)] = []byte(schemaBody)
	}
	user.CustomSchemas = customSchemas

	var err error
	err = retry(func() error {
		_, err = config.directory.Users.Update(d.Id(), user).Do()
		return err
	}, config.TimeoutMinutes)
	if err != nil {
		return fmt.Errorf("Error deleting user fields: %s", err)
	}

	d.SetId("")
	return nil
}

// Allow importing using any key (id, email, alias)
func resourceUserAttributesImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	id, err := config.directory.Users.Get(d.Id()).Projection("full").Do()

	if err != nil {
		return nil, fmt.Errorf("Error fetching user fields. Make sure the user fields exist: %s ", err)
	}

	d.SetId(id.Id)

	err, flattenedCustomSchema := flattenCustomSchema(id.CustomSchemas)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	d.Set("custom_schema", flattenedCustomSchema)

	return []*schema.ResourceData{d}, nil
}
