package gsuite

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// Provider returns the actual provider instance.
func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CREDENTIALS",
					"GOOGLE_CLOUD_KEYFILE_JSON",
					"GCLOUD_KEYFILE_JSON",
					"GOOGLE_APPLICATION_CREDENTIALS",
				}, nil),
				ValidateFunc: validateCredentials,
			},
			"impersonated_user_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth_scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"customer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeout_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1, // 1 + (n*2) roof 16 = 1+2+4+8+16 = 31 seconds, 1 min should be "normal" operations
			},
			"update_existing": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gsuite_group":           dataGroup(),
			"gsuite_group_settings":  dataGroupSettings(),
			"gsuite_user":            dataUser(),
			"gsuite_user_attributes": dataUserAttributes(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gsuite_domain":         resourceDomain(),
			"gsuite_group":          resourceGroup(),
			"gsuite_group_member":   resourceGroupMember(),
			"gsuite_group_members":  resourceGroupMembers(),
			"gsuite_group_settings": resourceGroupSettings(),
			"gsuite_user":           resourceUser(),
			"gsuite_user_fields":    resourceUserFields(),
			"gsuite_user_schema":    resourceUserSchema(),
		},
	}

	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return p
}

func oauthScopesFromConfigOrDefault(oauthScopesSet *schema.Set) []string {
	oauthScopes := convertStringSet(oauthScopesSet)
	if len(oauthScopes) == 0 {
		log.Printf("[INFO] No Oauth Scopes provided. Using default oauth scopes.")
		oauthScopes = defaultOauthScopes
	}
	return oauthScopes
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	var impersonatedUserEmail string
	var customerID string

	credentials := d.Get("credentials").(string)

	if v, ok := d.GetOk("impersonated_user_email"); ok {
		impersonatedUserEmail = v.(string)
	} else {
		if len(os.Getenv("IMPERSONATED_USER_EMAIL")) > 0 {
			impersonatedUserEmail = os.Getenv("IMPERSONATED_USER_EMAIL")
		}
	}

	// There shouldn't be the need to setup customer ID in the configuration,
	// but leaving the possibility to specify it explictly.
	// By default we use my_customer as customer ID, which means the API will use
	// the G Suite customer ID associated with the impersonating account.
	if v, ok := d.GetOk("customer_id"); ok {
		customerID = v.(string)
	} else {
		log.Printf("[INFO] No Customer ID provided. Using my_customer.")
		customerID = "my_customer"
	}

	timeoutMinutes := d.Get("timeout_minutes").(int)

	oauthScopes := oauthScopesFromConfigOrDefault(d.Get("oauth_scopes").(*schema.Set))

	updateExisting := true
	if v, ok := d.GetOk("update_existing"); ok {
		updateExisting = v.(bool)
	}

	config := Config{
		Credentials:           credentials,
		ImpersonatedUserEmail: impersonatedUserEmail,
		OauthScopes:           oauthScopes,
		CustomerId:            customerID,
		TimeoutMinutes:        timeoutMinutes,
		UpdateExisting:        updateExisting,
	}

	if err := config.loadAndValidate(terraformVersion); err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	return &config, nil
}

func validateCredentials(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}
	creds := v.(string)
	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(creds); err == nil {
		return
	}
	var account accountFile
	if err := json.Unmarshal([]byte(creds), &account); err != nil {
		errors = append(errors,
			fmt.Errorf("credentials are not valid JSON '%s': %s", creds, err))
	}

	return
}
