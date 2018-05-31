package main

import (
	"github.com/DeviaVir/terraform-provider-gsuite/gsuite"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return gsuite.Provider()
		},
	})
}
