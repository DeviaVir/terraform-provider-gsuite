package gsuite

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	directory "google.golang.org/api/admin/directory/v1"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Delete: resourceDomainDelete,
		// There is no update method

		Schema: map[string]*schema.Schema{

			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_name": {
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

func resourceDomainCreate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	domain := &directory.Domains{}
	customerId := config.CustomerId

	if v, ok := d.GetOk("domain_name"); ok {
		log.Printf("[DEBUG] Setting %s: %s", "domain_name", v.(string))
		domain.DomainName = strings.ToLower(v.(string))
	}

	var createdDomain *directory.Domains

	var err error
	err = retry(func() error {
		createdDomain, err = config.directory.Domains.Insert(customerId, domain).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error creating domain: %s", err)
	}

	// There is no id as such for a Domain resource, therefore we use
	// DomainName as unique identifier.
	d.SetId(createdDomain.DomainName)
	d.Set("domain_name", domain.DomainName)

	log.Printf("[INFO] Created domain: %s", createdDomain.DomainName)
	return resourceDomainRead(d, meta)
}

func resourceDomainRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	customerId := config.CustomerId
	var domain *directory.Domains

	var domainName string

	if v, ok := d.GetOk("domain_name"); ok {
		log.Printf("[DEBUG] Reading %s: %s", "domain_name", v.(string))
		domainName = v.(string)
	}

	var err error
	err = retry(func() error {
		domain, err = config.directory.Domains.Get(customerId, domainName).Do()
		return err
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Domain %q", d.Get("domain_name").(string)))
	}

	d.SetId(domain.DomainName)
	d.Set("domain_name", domain.DomainName)

	return nil
}

func resourceDomainDelete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	customerId := config.CustomerId

	var domainName string

	if v, ok := d.GetOk("domain_name"); ok {
		log.Printf("[DEBUG] Deleting %s: %s", "domain_name", v.(string))
		domainName = v.(string)
	}

	var err error
	err = retry(func() error {
		err = config.directory.Domains.Delete(customerId, domainName).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error deleting domain: %s", err)
	}

	d.SetId("")

	return nil
}
