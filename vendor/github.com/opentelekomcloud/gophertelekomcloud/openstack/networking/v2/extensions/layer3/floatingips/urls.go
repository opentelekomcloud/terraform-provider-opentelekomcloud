package floatingips

import "github.com/opentelekomcloud/gophertelekomcloud"

const resourcePath = "floatingips"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}
