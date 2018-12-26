package containers

import "github.com/huaweicloud/golangsdk"

func listURL(c *golangsdk.ServiceClient) string {
	return c.Endpoint
}

func createURL(c *golangsdk.ServiceClient, container string) string {
	return c.ServiceURL(container)
}

func getURL(c *golangsdk.ServiceClient, container string) string {
	return createURL(c, container)
}

func deleteURL(c *golangsdk.ServiceClient, container string) string {
	return createURL(c, container)
}

func updateURL(c *golangsdk.ServiceClient, container string) string {
	return createURL(c, container)
}
