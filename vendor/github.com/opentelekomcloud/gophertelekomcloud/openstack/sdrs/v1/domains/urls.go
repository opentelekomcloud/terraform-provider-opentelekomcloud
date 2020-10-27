package domains

import "github.com/opentelekomcloud/gophertelekomcloud"

func getURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("active-domains")
}
