package volumeactions

import "github.com/opentelekomcloud/gophertelekomcloud"

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("volumes", id, "action")
}
