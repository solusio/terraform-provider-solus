package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/solusio/terraform-provider-solus/solus"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: solus.Provider,
	})
}
