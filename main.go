package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: opentelekomcloud.Provider}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/opentelekomcloud/opentelekomcloud", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
