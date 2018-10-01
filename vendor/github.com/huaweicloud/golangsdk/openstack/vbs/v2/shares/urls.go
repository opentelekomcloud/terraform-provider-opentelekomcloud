package shares

import "github.com/huaweicloud/golangsdk"

const rootPath = "os-vendor-backup-sharing"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, "detail")
}
