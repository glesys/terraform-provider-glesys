package main

import (
	"github.com/glesys/terraform-provider-glesys/glesys"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return glesys.Provider()
		},
	})
}
