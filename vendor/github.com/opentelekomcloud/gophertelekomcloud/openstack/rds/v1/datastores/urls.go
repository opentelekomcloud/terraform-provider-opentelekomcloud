package datastores

import "github.com/opentelekomcloud/gophertelekomcloud"

func listURL(c *golangsdk.ServiceClient, dataStoreName string) string {
	return c.ServiceURL("datastores", dataStoreName, "versions")
}
