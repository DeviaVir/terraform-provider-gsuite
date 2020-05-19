package gsuite

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	directory "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

func resourceUserFields() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserFieldsCreate,
		Read:   resourceUserFieldsRead,
		Update: resourceUserFieldsUpdate,
		Delete: resourceUserFieldsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserFieldsImporter,
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

func resourceUserFieldsCreate(d *schema.ResourceData, meta interface{}) error {
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
	return resourceUserFieldsRead(d, meta)
}

func resourceUserFieldsUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	user := &directory.User{}

	if d.HasChange("primary_email") {
		if v, ok := d.GetOk("primary_email"); ok {
			log.Printf("[DEBUG] Updating user primary_email: %s", d.Get("primary_email").(string))
			user.PrimaryEmail = v.(string)
		}
	}

	if d.HasChange("custom_schema") {
		customSchemas := map[string]googleapi.RawMessage{}
		for i := 0; i < d.Get("custom_schema.#").(int); i++ {
			entry := d.Get(fmt.Sprintf("custom_schema.%d", i)).(map[string]interface{})
			customSchemas[entry["name"].(string)] = []byte(entry["value"].(string))
		}
		user.CustomSchemas = customSchemas
	}

	var updatedUser *directory.User
	var err error
	err = retry(func() error {
		updatedUser, err = config.directory.Users.Update(d.Id(), user).Do()
		if e, ok := err.(*googleapi.Error); ok {
			return errors.Wrap(e, e.Body)
		}
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		log.Printf("[WARN] Please note, a persistent 503 backend error can mean you need to change your posix values to be unique.")
		return fmt.Errorf("[ERROR] Error updating user fields: %s", err)
	}

	log.Printf("[INFO] Updated user fields: %s", updatedUser.PrimaryEmail)
	return resourceUserFieldsRead(d, meta)
}

func resourceUserFieldsRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceUserFieldsDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	user := &directory.User{}
	user.CustomSchemas = map[string]googleapi.RawMessage{}

	var err error
	err = retry(func() error {
		_, err = config.directory.Users.Patch(d.Id(), user).Do()
		return err
	}, config.TimeoutMinutes)
	if err != nil {
		return fmt.Errorf("Error deleting user fields: %s", err)
	}

	d.SetId("")
	return nil
}

// Allow importing using any key (id, email, alias)
func resourceUserFieldsImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
