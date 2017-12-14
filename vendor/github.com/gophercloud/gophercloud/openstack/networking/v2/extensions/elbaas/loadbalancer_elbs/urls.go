package loadbalancer_elbs

import "github.com/gophercloud/gophercloud"

const (
	rootPath     = "elbaas"
	resourcePath = "loadbalancers"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}
