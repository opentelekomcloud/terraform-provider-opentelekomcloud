package loadbalancer_elbs

import "github.com/gophercloud/gophercloud"

const (
	rootPath     = "elbaas"
	resourcePath = "loadbalancers"
	statusPath   = "statuses"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func statusRootURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, statusPath)
}
