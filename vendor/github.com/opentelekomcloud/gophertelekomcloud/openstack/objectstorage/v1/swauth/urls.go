package swauth

import "github.com/opentelekomcloud/gophertelekomcloud"

func getURL(c *golangsdk.ProviderClient) string {
	return c.IdentityBase + "auth/v1.0"
}
