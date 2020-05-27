package gsuite

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	directory "google.golang.org/api/admin/directory/v1"
)

func resourceOrgUnit() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrgUnitCreate,
		Read:   resourceOrgUnitRead,
		Delete: resourceOrgUnitDelete,
		// There is no update method

		Schema: map[string]*schema.Schema{

			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"orgunit_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},

			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOrgUnitCreate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	orgunit := &directory.OrgUnit{}
	customerId := config.CustomerId

	if v, ok := d.GetOk("orgunit_name"); ok {
		log.Printf("[DEBUG] Creating %s: %s", "orgunit_name", v.(string))
		orgunit.Name = v.(string)
	}

	var createdOrgUnit *directory.OrgUnit

	var err error
	err = retry(func() error {
		createdOrgUnit, err = config.directory.Orgunits.Insert(customerId, orgunit).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return fmt.Errorf("[ERROR] Error creating orgunit: %s", err)
	}

	// There is no id as such for a OrgUnit resource, therefore we use
	// Name as unique identifier.
	d.SetId(createdOrgUnit.Name)
	d.Set("orgunit_name", orgunit.Name)

	log.Printf("[INFO] Created orgunit: %s", createdOrgUnit.Name)
	return resourceOrgUnitRead(d, meta)
}

func resourceOrgUnitRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	customerId := config.CustomerId
	var orgunit *directory.OrgUnit

	var orgunitPath []string

	if v, ok := d.GetOk("orgunit_path"); ok {
		log.Printf("[DEBUG] Reading %s: %s", "orgunit_path", v.(string))
		orgunitPath = v.([]string)
	}

	var err error
	err = retry(func() error {
		orgunit, err = config.directory.Orgunits.Get(customerId, orgunitPath).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("OrgUnit %q", d.Get("orgunit_path").(string)))
	}

	d.SetId(orgunit.Name)
	d.Set("orgunit_name", orgunit.Name)

	return nil
}

func resourceOrgUnitDelete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	customerId := config.CustomerId

	var orgunitPath []string

	if v, ok := d.GetOk("orgunit_path"); ok {
		log.Printf("[DEBUG] Deleting %s: %s", "orgunit_path", v.(string))
		orgunitPath = v.([]string)
	}

	var err error
	err = retry(func() error {
		err = config.directory.Orgunits.Delete(customerId, orgunitPath).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting orgunit: %s", err)
	}

	d.SetId("")

	return nil
}
