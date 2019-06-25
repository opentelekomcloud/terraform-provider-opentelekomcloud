package whiteblackip_rules

import "github.com/huaweicloud/golangsdk"

func rootURL(c *golangsdk.ServiceClient, policy_id string) string {
	return c.ServiceURL("policy", policy_id, "whiteblackip")
}

func resourceURL(c *golangsdk.ServiceClient, policy_id, id string) string {
	return c.ServiceURL("policy", policy_id, "whiteblackip", id)
}
