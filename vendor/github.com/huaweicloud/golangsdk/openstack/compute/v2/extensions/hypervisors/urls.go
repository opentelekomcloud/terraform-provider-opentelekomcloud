package hypervisors

import "github.com/huaweicloud/golangsdk"

func hypervisorsListDetailURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-hypervisors", "detail")
}

func hypervisorsStatisticsURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-hypervisors", "statistics")
}

func hypervisorsGetURL(c *golangsdk.ServiceClient, hypervisorID string) string {
	return c.ServiceURL("os-hypervisors", hypervisorID)
}

func hypervisorsUptimeURL(c *golangsdk.ServiceClient, hypervisorID string) string {
	return c.ServiceURL("os-hypervisors", hypervisorID, "uptime")
}
