package backups

import "github.com/opentelekomcloud/gophertelekomcloud"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "backups/policy")
}
