package gsuite

import (
	"context"
	"time"

	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
	"os"
)

var (
	// contextTimeout is the global context timeout for requests to complete.
	contextTimeout = 15 * time.Second
)

// Provider returns the actual provider instance.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CREDENTIALS",
					"GOOGLE_CLOUD_KEYFILE_JSON",
					"GCLOUD_KEYFILE_JSON",
				}, nil),
				ValidateFunc: validateCredentials,
			},
			"impersonated_user_email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth_scopes": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gsuite_group":         resourceGroup(),
			"gsuite_user":          resourceUser(),
			"gsuite_user_schema":   resourceUserSchema(),
			"gsuite_group_member":  resourceGroupMember(),
			"gsuite_group_members": resourceGroupMembers(),
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
	credentials := d.Get("credentials").(string)
	if v, ok := d.GetOk("impersonated_user_email"); ok {
		impersonatedUserEmail = v.(string)
	} else {
		if len(os.Getenv("IMPERSONATED_USER_EMAIL")) > 0 {
			impersonatedUserEmail = os.Getenv("IMPERSONATED_USER_EMAIL")
		}
	}
	oauthScopes := oauthScopesFromConfigOrDefault(d.Get("oauth_scopes").(*schema.Set))
	config := Config{
		Credentials:           credentials,
		ImpersonatedUserEmail: impersonatedUserEmail,
		OauthScopes:           oauthScopes,
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	return &config, nil
}

// contextWithTimeout creates a new context with the global context timeout.
func contextWithTimeout() (context.Context, func()) {
	return context.WithTimeout(context.Background(), contextTimeout)
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
