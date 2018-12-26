package limits

import (
	"github.com/huaweicloud/golangsdk"
)

const resourcePath = "limits"

func getURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
