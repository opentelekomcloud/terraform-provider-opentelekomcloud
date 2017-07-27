package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/khdegraaf/terraform-provider-openstack/openstack" // TODO: Revert path when merge
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: openstack.Provider})
}
