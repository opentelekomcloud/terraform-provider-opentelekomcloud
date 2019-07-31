package volumes

import "github.com/huaweicloud/golangsdk"

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("os-vendor-volumes", id)
}
