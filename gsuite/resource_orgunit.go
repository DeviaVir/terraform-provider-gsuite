package gsuite

import (
	"fmt"
	"log"
	"path"
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

			"orgunit_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
		},
	}
}

func resourceOrgUnitCreate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	orgunit := &directory.OrgUnit{}
	customerId := config.CustomerId
    orgUnitPath := ""

	if v, ok := d.GetOk("orgunit_path"); ok {
		log.Printf("[DEBUG] Creating %s: %s", "orgunit_path", v.(string))
		orgUnitPath = v.(string)
	}

	cumulativePath := ""
	pathComponents := strings.Split(strings.TrimPrefix(orgUnitPath, "/"), "/")
	for _, pathComponent := range pathComponents {
	    cumulativePath += "/" + pathComponent
	    parentPath, unitName := path.Split(cumulativePath)
		currentPath := "/" + unitName
	    if parentPath != "/" {
	        parentPath = strings.TrimSuffix(parentPath, "/")
			currentPath = parentPath + "/" + unitName
	    }

		orgunit.ParentOrgUnitPath = parentPath
		orgunit.Name = unitName

		var createdOrgUnit *directory.OrgUnit
		var err error

		// Try to query the organizational unit, lest creating an existing unit
		err = retry(func() error {
			createdOrgUnit, err = config.directory.Orgunits.Get(customerId,
				[]string{strings.TrimPrefix(strings.ToLower(currentPath), "/")},
			).Do()
			return err
		}, config.TimeoutMinutes)

		// If the organizational unit does not exist, try to create it
		if err != nil && strings.Contains(string(err.Error()), "not found") {
			err = retry(func() error {
				createdOrgUnit, err = config.directory.Orgunits.Insert(
					customerId,
					orgunit,
				).Do()
				return err
			}, config.TimeoutMinutes)
		} else if err != nil {
			return fmt.Errorf(
				"[ERROR] Error creating orgunit: %s",
				string(err.Error()),
			)
		}

        // Leaf organizational unit
		if currentPath == orgUnitPath {
			d.SetId(createdOrgUnit.OrgUnitPath)
			d.Set("orgunit_path", createdOrgUnit.OrgUnitPath)

			log.Printf("[INFO] Created orgunit: %s", orgUnitPath)
			return resourceOrgUnitRead(d, meta)
		}
	}
	return nil
}

func resourceOrgUnitRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	customerId := config.CustomerId
	var orgunit *directory.OrgUnit

	var orgunitPath string

	if v, ok := d.GetOk("orgunit_path"); ok {
		log.Printf("[DEBUG] Reading %s: %s", "orgunit_path", v.(string))
		orgunitPath = strings.TrimPrefix(strings.ToLower(v.(string)), "/")
	}

	var err error
	err = retry(func() error {
		orgunit, err = config.directory.Orgunits.Get(
			customerId,
			[]string{orgunitPath},
		).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf(
			"OrgUnit %q", d.Get("orgunit_path").(string),
		))
	}

	d.SetId(orgunit.OrgUnitPath)
	d.Set("orgunit_path", orgunit.OrgUnitPath)

	return nil
}

func resourceOrgUnitDelete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	customerId := config.CustomerId

	var orgunitPath string

	if v, ok := d.GetOk("orgunit_path"); ok {
		log.Printf("[DEBUG] Deleting %s: %s", "orgunit_path", v.(string))
		orgunitPath = strings.TrimPrefix(strings.ToLower(v.(string)), "/")
	}

	var err error
	err = retry(func() error {
		err = config.directory.Orgunits.Delete(
			customerId,
			[]string{orgunitPath},
		).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting orgunit: %s", err)
	}

	d.SetId("")

	return nil
}
