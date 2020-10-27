package tenantnetworks

import "github.com/opentelekomcloud/gophertelekomcloud"

const resourcePath = "os-tenant-networks"

func resourceURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *golangsdk.ServiceClient) string {
	return resourceURL(c)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}
