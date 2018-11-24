package ruletypes

import "github.com/huaweicloud/golangsdk"

func listRuleTypesURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("qos", "rule-types")
}
