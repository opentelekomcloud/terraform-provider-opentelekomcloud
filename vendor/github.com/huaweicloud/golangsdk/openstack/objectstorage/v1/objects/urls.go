package objects

import (
	"github.com/huaweicloud/golangsdk"
)

func listURL(c *golangsdk.ServiceClient, container string) string {
	return c.ServiceURL(container)
}

func copyURL(c *golangsdk.ServiceClient, container, object string) string {
	return c.ServiceURL(container, object)
}

func createURL(c *golangsdk.ServiceClient, container, object string) string {
	return copyURL(c, container, object)
}

func getURL(c *golangsdk.ServiceClient, container, object string) string {
	return copyURL(c, container, object)
}

func deleteURL(c *golangsdk.ServiceClient, container, object string) string {
	return copyURL(c, container, object)
}

func downloadURL(c *golangsdk.ServiceClient, container, object string) string {
	return copyURL(c, container, object)
}

func updateURL(c *golangsdk.ServiceClient, container, object string) string {
	return copyURL(c, container, object)
}
