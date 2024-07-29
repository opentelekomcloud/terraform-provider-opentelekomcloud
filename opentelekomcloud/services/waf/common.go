package waf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
)

const (
	ClientError                  = "error creating OpenTelekomCloud WAF client: %w"
	errCreationV1DedicatedClient = "error creating OpenTelekomCloud WAF dedicated client: %w"
	keyClientV1                  = "wafd-v1-client"
	// runStatusCreating the instance is creating.
	runStatusCreating = 0
	// runStatusRunning the instance has been created.
	runStatusRunning = 1
	// runStatusDeleting the instance deleting.
	runStatusDeleting = 2
	// runStatusDeleting the instance has be deleted.
	runStatusDeleted = 3
	// defaultCount the number of instances created.
	defaultCount = 1
	// Billing mode, payPerUseMode: pay pre use mode
	payPerUseMode = 30
	// ProtectionActionBlock block the request
	ProtectionActionBlock = 0
	// ProtectionActionAllow allow the request
	ProtectionActionAllow = 1
	// ProtectionActionLog log the request only
	ProtectionActionLog = 2
)

func wafRuleImporter() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: common.ImportByPath("policy_id", "id"),
	}
}

type ExtendOptions struct {
	DeepDecode            *bool `json:"deep_decode,omitempty"`
	CheckAllHeaders       *bool `json:"check_all_headers,omitempty"`
	ShiroRememberMeEnable *bool `json:"shiro_rememberMe_enable,omitempty"`
}
