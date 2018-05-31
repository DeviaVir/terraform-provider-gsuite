package gsuite

import (
	"io/ioutil"
	"testing"
)

const testFakeCredentialsPath = "./test-fixtures/fake_account.json"

func TestConfigLoadAndValidate_accountFilePath(t *testing.T) {
	config := Config{
		Credentials:           testFakeCredentialsPath,
		ImpersonatedUserEmail: "xxx@xxx.xom",
	}

	err := config.loadAndValidate()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestConfigLoadAndValidate_accountFileJSON(t *testing.T) {
	contents, err := ioutil.ReadFile(testFakeCredentialsPath)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	config := Config{
		Credentials:           string(contents),
		ImpersonatedUserEmail: "xxx@xxx.xom",
	}

	err = config.loadAndValidate()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestConfigLoadAndValidate_accountFileJSONInvalid(t *testing.T) {
	config := Config{
		Credentials: "{this is not json}",
	}

	if config.loadAndValidate() == nil {
		t.Fatalf("expected error, but got nil")
	}
}

func TestConfigLoadAndValidate_noImpersonatedEmail(t *testing.T) {
	// ImpersonatedUserEmail empty string when credentials set
	config := Config{
		Credentials:           testFakeCredentialsPath,
		ImpersonatedUserEmail: "",
	}

	err := config.loadAndValidate()
	if err == nil {
		t.Fatalf("error: %v", err)
	}
	if err.Error() != "required field missing: impersonated_user_email" {
		t.Fatalf("error: %v", err)
	}

	// ImpersonatedUserEmail not provided when credentials set
	config = Config{
		Credentials: testFakeCredentialsPath,
	}

	err = config.loadAndValidate()
	if err == nil {
		t.Fatalf("error: %v", err)
	}
	if err.Error() != "required field missing: impersonated_user_email" {
		t.Fatalf("error: %v", err)
	}
}
