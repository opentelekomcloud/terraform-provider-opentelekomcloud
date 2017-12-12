package tags

import (
	"github.com/gophercloud/gophercloud"
)

func createURL(c *gophercloud.ServiceClient, server_id string) string {
	return c.ServiceURL("servers", server_id, "tags")
}

func getURL(c *gophercloud.ServiceClient, server_id string) string {
	return c.ServiceURL("servers", server_id, "tags")
}

func deleteURL(c *gophercloud.ServiceClient, server_id string) string {
	return c.ServiceURL("servers", server_id, "tags")
}

