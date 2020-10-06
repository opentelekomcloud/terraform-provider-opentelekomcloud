package softwareconfig

import "github.com/opentelekomcloud/gophertelekomcloud"

const resourcePath = "software_configs"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}
