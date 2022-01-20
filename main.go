package main

//go:generate terraform fmt  -recursive ./examples/

import (
	"context"
	"flag"
	"log"

	"github.com/solusio/terraform-provider-solus/internal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
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
		ProviderFunc: provider.New,
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
