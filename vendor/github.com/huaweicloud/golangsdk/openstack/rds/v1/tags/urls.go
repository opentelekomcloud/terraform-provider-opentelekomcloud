package tags

import "github.com/huaweicloud/golangsdk"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(id, "tags")
}
