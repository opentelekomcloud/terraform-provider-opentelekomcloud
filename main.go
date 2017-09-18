package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/khdegraaf/terraform-provider-huaweicloud/huaweicloud" // TODO: Revert path when merge
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: huaweicloud.Provider})
}
