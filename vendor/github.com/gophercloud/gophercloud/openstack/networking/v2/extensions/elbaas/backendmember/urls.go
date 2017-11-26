package backendmember

import "github.com/gophercloud/gophercloud"

const (
	rootPath     = "elbaas"
	resourcePath = "listeners"
)

func addURL(c *gophercloud.ServiceClient, listener_id string) string {
	return c.ServiceURL(rootPath, resourcePath, listener_id, "members")
}

func removeURL(c *gophercloud.ServiceClient, listener_id string) string {
	return c.ServiceURL(rootPath, resourcePath, listener_id, "members", "action")
}

func resourceURL(c *gophercloud.ServiceClient, listener_id string, id string) string {
	return c.ServiceURL(rootPath, resourcePath, listener_id, "members?id=" + id)
}
