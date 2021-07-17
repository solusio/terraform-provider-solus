package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/solusio/terraform-provider-solus/solus"
	"log"
)

func main() {
	var debugMode bool

	flag.BoolVar(
		&debugMode,
		"debug",
		false,
		"set to true to run the provider with support for debuggers like delve",
	)
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: solus.Provider,
	}

	if !debugMode {
		plugin.Serve(opts)
		return
	}

	err := plugin.Debug(
		context.Background(),
		"github.com/solusio/terraform-provider-solus",
		opts,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}
