package stacktemplates

import "github.com/opentelekomcloud/gophertelekomcloud"

func getURL(c *golangsdk.ServiceClient, stackName, stackID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "template")
}
