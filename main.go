package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/norrland/terraform-provider-glesys/glesys"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: glesys.Provider})
}
