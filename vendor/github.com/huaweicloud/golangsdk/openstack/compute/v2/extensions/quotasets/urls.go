package quotasets

import "github.com/huaweicloud/golangsdk"

const resourcePath = "os-quota-sets"

func resourceURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func getURL(c *golangsdk.ServiceClient, tenantID string) string {
	return c.ServiceURL(resourcePath, tenantID)
}

func getDetailURL(c *golangsdk.ServiceClient, tenantID string) string {
	return c.ServiceURL(resourcePath, tenantID, "detail")
}

func updateURL(c *golangsdk.ServiceClient, tenantID string) string {
	return getURL(c, tenantID)
}

func deleteURL(c *golangsdk.ServiceClient, tenantID string) string {
	return getURL(c, tenantID)
}
