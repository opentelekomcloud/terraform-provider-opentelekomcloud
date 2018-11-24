package services

import "github.com/huaweicloud/golangsdk"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-services")
}
