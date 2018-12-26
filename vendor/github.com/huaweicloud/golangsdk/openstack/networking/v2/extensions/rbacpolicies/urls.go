package rbacpolicies

import "github.com/huaweicloud/golangsdk"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("rbac-policies", id)
}

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("rbac-policies")
}

func createURL(c *golangsdk.ServiceClient) string {
	return rootURL(c)
}

func listURL(c *golangsdk.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return resourceURL(c, id)
}
