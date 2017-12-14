package tags

import (
	"github.com/gophercloud/gophercloud"
)

func createURL(c *gophercloud.ServiceClient, resource_type, resource_id string) string {
	return c.ServiceURL("os-vendor-tags", resource_type, resource_id)
}

func getURL(c *gophercloud.ServiceClient, resource_type, resource_id string) string {
	return c.ServiceURL("os-vendor-tags", resource_type, resource_id)
}

