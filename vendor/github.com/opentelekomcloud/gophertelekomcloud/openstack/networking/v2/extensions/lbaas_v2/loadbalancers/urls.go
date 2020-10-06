package loadbalancers

import "github.com/opentelekomcloud/gophertelekomcloud"

const (
	rootPath     = "lbaas"
	resourcePath = "loadbalancers"
	statusPath   = "statuses"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func statusRootURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, statusPath)
}
