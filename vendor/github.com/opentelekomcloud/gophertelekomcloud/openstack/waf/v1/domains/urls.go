package domains

import "github.com/opentelekomcloud/gophertelekomcloud"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instance")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instance", id)
}
