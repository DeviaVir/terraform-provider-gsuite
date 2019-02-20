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
		Update: resourceDomainUpdate,

		Schema: map[string]*schema.Schema{

			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_aliases": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},

			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_primary": {
				Type:     schema.TypeBool,
				Optional: true,
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
	d.Set("domain_aliases", domain.DomainAliases)
	d.Set("etag", domain.Etag)
	d.Set("is_primary", domain.IsPrimary)
	d.Set("creation_time", domain.CreationTime)

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

func resourceDomainUpdate(d *schema.ResourceData, meta interface{}) error {
		// There is no update method in https://developers.google.com/admin-sdk/directory/v1/reference/domains,
		// therefore returning an error message to the user.
		return fmt.Errorf("There is no update method for gsuite_domain resource")
}
