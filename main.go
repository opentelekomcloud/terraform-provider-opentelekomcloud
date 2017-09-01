package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/khdegraaf/terraform-provider-hwcloud/hwcloud" // TODO: Revert path when merge
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hwcloud.Provider})
}
