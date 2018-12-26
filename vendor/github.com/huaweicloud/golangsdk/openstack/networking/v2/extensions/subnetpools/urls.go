package subnetpools

import "github.com/huaweicloud/golangsdk"

const resourcePath = "subnetpools"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *golangsdk.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func createURL(c *golangsdk.ServiceClient) string {
	return rootURL(c)
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return resourceURL(c, id)
}
