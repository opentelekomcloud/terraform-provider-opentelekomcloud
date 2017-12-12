package alarmrule

import "github.com/gophercloud/gophercloud"

const (
	rootPath = "alarms"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, rootPath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, rootPath, id)
}

func actionURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, rootPath, id, "action")
}
