package vpcs

import "github.com/gophercloud/gophercloud"

const resourcePath = "vpcs"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}
