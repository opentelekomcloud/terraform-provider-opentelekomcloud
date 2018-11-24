package defsecrules

import "github.com/huaweicloud/golangsdk"

const rulepath = "os-security-group-default-rules"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rulepath, id)
}

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rulepath)
}
