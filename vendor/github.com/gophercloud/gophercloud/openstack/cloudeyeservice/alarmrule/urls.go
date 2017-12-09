package alarmrule

import "github.com/gophercloud/gophercloud"

const (
	//rootPath = "V1.0/bf74229f30c0421fae270386a43315ee"
	rootPath     = "bf74229f30c0421fae270386a43315ee"
	resourcePath = "alarms"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}
