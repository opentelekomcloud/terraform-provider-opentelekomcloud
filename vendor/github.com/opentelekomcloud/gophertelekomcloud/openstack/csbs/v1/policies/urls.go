package policies

import "github.com/opentelekomcloud/gophertelekomcloud"

const rootPath = "policies"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, policyId string) string {
	return c.ServiceURL(rootPath, policyId)
}
