package waf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
)

const (
	ClientError = "error creating OpenTelekomCloud WAF client: %w"
)

func wafRuleImporter() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: common.ImportByPath("policy_id", "id"),
	}
}
