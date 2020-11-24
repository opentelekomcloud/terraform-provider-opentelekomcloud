package flavors

import "github.com/opentelekomcloud/gophertelekomcloud"

func baseURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("flavors")
}
