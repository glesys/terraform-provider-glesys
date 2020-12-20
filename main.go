package main

import (
	"github.com/glesys/terraform-provider-glesys/glesys"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: glesys.Provider})
}
