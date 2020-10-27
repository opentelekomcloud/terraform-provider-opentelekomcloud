package stackresources

import "github.com/opentelekomcloud/gophertelekomcloud"

func listURL(c *golangsdk.ServiceClient, stackName string) string {
	return c.ServiceURL("stacks", stackName, "resources")
}
