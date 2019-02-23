package gsuite

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

// Provider returns the actual provider instance.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CREDENTIALS",
					"GOOGLE_CLOUD_KEYFILE_JSON",
					"GCLOUD_KEYFILE_JSON",
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
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gsuite_user_attributes": dataUserAttributes(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gsuite_group":         resourceGroup(),
			"gsuite_user":          resourceUser(),
			"gsuite_user_schema":   resourceUserSchema(),
			"gsuite_group_member":  resourceGroupMember(),
			"gsuite_group_members": resourceGroupMembers(),
			"gsuite_domain":        resourceDomain(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func oauthScopesFromConfigOrDefault(oauthScopesSet *schema.Set) []string {
	oauthScopes := convertStringSet(oauthScopesSet)
	if len(oauthScopes) == 0 {
		log.Printf("[INFO] No Oauth Scopes provided. Using default oauth scopes.")
		oauthScopes = defaultOauthScopes
	}
	return oauthScopes
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var impersonatedUserEmail string
	var customerId string

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
		customerId = v.(string)
	} else {
		log.Printf("[INFO] No Customer ID provided. Using my_customer.")
		customerId = "my_customer"
	}

	oauthScopes := oauthScopesFromConfigOrDefault(d.Get("oauth_scopes").(*schema.Set))
	config := Config{
		Credentials:           credentials,
		ImpersonatedUserEmail: impersonatedUserEmail,
		OauthScopes:           oauthScopes,
		CustomerId:            customerId,
	}

	if err := config.loadAndValidate(); err != nil {
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
