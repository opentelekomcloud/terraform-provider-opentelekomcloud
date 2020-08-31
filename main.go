package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opentelekomcloud.Provider})
}
