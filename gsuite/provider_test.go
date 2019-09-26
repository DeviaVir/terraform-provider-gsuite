package gsuite

import (
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

var credsEnvVars = []string{
	"GOOGLE_CREDENTIALS",
	"GOOGLE_CLOUD_KEYFILE_JSON",
	"GCLOUD_KEYFILE_JSON",
	"GOOGLE_APPLICATION_CREDENTIALS",
}

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"gsuite": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func TestProvider_loadCredentialsFromFile(t *testing.T) {
	ws, es := validateCredentials(testFakeCredentialsPath, "")
	if len(ws) != 0 {
		t.Errorf("Expected %d warnings, got %v", len(ws), ws)
	}
	if len(es) != 0 {
		t.Errorf("Expected %d errors, got %v", len(es), es)
	}
}

func TestProvider_loadCredentialsFromJSON(t *testing.T) {
	contents, err := ioutil.ReadFile(testFakeCredentialsPath)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	ws, es := validateCredentials(string(contents), "")
	if len(ws) != 0 {
		t.Errorf("Expected %d warnings, got %v", len(ws), ws)
	}
	if len(es) != 0 {
		t.Errorf("Expected %d errors, got %v", len(es), es)
	}
}

func TestConfigOauthScopes(t *testing.T) {

	scopes := oauthScopesFromConfigOrDefault(&schema.Set{})

	if len(scopes) != len(defaultOauthScopes) {
		t.Fatalf("error: default oauth scopes not being set")
	}

	s := schema.NewSet(
		schema.HashString,
		[]interface{}{
			"https://www.googleapis.com/auth/admin.directory.group"})

	scopes = oauthScopesFromConfigOrDefault(s)

	if len(scopes) != len(convertStringSet(s)) {
		t.Fatalf("error: oauth scopes not being set")
	}
}
