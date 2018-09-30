package hosts

import "github.com/huaweicloud/golangsdk"

const resourcePath = "dedicated-hosts"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}
func listServerURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "servers")
}
