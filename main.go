package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: opentelekomcloud.Provider,
		ProviderAddr: "registry.terraform.io/opentelekomcloud/opentelekomcloud",
	}

	if debugMode {
		opts.Debug = true
	}

	plugin.Serve(opts)
}
