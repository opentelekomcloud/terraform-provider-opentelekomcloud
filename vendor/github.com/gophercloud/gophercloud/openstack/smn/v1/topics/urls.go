package topics

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("topics")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("topics", id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("topics", id)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("topics", id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("topics")
}
