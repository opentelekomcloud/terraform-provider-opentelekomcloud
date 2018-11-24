package pools

import "github.com/huaweicloud/golangsdk"

const (
	rootPath     = "lb"
	resourcePath = "pools"
	monitorPath  = "health_monitors"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func associateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, monitorPath)
}

func disassociateURL(c *golangsdk.ServiceClient, poolID, monitorID string) string {
	return c.ServiceURL(rootPath, resourcePath, poolID, monitorPath, monitorID)
}
