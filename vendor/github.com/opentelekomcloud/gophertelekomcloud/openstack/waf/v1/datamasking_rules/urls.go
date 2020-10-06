package datamasking_rules

import "github.com/opentelekomcloud/gophertelekomcloud"

func rootURL(c *golangsdk.ServiceClient, policy_id string) string {
	return c.ServiceURL("policy", policy_id, "privacy")
}

func resourceURL(c *golangsdk.ServiceClient, policy_id, id string) string {
	return c.ServiceURL("policy", policy_id, "privacy", id)
}
