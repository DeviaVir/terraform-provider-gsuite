package gsuite

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	directory "google.golang.org/api/admin/directory/v1"
)

var oauthScopes = []string{
	directory.AdminDirectoryCustomerScope,
	directory.AdminDirectoryGroupScope,
	directory.AdminDirectoryGroupMemberScope,
	directory.AdminDirectoryOrgunitScope,
	directory.AdminDirectoryUserScope,
	directory.AdminDirectoryUserAliasScope,
	directory.AdminDirectoryUserSecurityScope,
	directory.AdminDirectoryUserschemaScope,
}

// Config is the structure used to instantiate the GSuite provider.
type Config struct {
	directory *directory.Service
}

// loadAndValidate loads the application default credentials from the
// environment and creates a client for communicating with Google APIs.
func (c *Config) loadAndValidate() error {
	log.Printf("[INFO] authenticating with local client")
	client, err := google.DefaultClient(context.Background(), oauthScopes...)
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	// Use a custom user-agent string. This helps google with analytics and it's
	// just a nice thing to do.
	client.Transport = logging.NewTransport("Google", client.Transport)
	userAgent := fmt.Sprintf("(%s %s) Terraform/%s",
		runtime.GOOS, runtime.GOARCH, terraform.VersionString())

	// Create the directory service.
	directorySvc, err := directory.New(client)
	if err != nil {
		return nil
	}
	directorySvc.UserAgent = userAgent
	c.directory = directorySvc

	return nil
}
