package aggregates

import "github.com/huaweicloud/golangsdk"

func aggregatesListURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-aggregates")
}

func aggregatesCreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-aggregates")
}

func aggregatesDeleteURL(c *golangsdk.ServiceClient, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID)
}

func aggregatesGetURL(c *golangsdk.ServiceClient, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID)
}

func aggregatesUpdateURL(c *golangsdk.ServiceClient, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID)
}

func aggregatesAddHostURL(c *golangsdk.ServiceClient, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID, "action")
}

func aggregatesRemoveHostURL(c *golangsdk.ServiceClient, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID, "action")
}

func aggregatesSetMetadataURL(c *golangsdk.ServiceClient, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID, "action")
}
