package main

import (
	"github.com/terraform-providers/terraform-provider-opentelekomcloud/opentelekomcloud"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opentelekomcloud.Provider})
}
