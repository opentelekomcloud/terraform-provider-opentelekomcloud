package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const wafdCCRuleName = "opentelekomcloud_waf_dedicated_cc_rule_v1.rule_1"

func TestAccWafDedicatedCcAttackProtectionRuleV1_basic(t *testing.T) {
	var rule rules.CcRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedCcRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedCcAttackProtectionRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedCcRuleV1Exists(wafdCCRuleName, &rule),
					resource.TestCheckResourceAttr(wafdCCRuleName, "url", "/abc1"),
				),
			},
			{
				ResourceName:      wafdCCRuleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: dedicatedRuleImportStateIDFunc(wafdCCRuleName, wafdPolicyResourceName),
			},
		},
	})
}

func testAccCheckWafDedicatedCcRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_cc_rule_v1" {
			continue
		}

		_, err := rules.GetCc(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedCcRuleV1Exists(n string, rule *rules.CcRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
		if err != nil {
			return err
		}

		found, err := rules.GetCc(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("waf dedicated rule not found")
		}

		*rule = *found

		return nil
	}
}

const testAccWafDedicatedCcAttackProtectionRuleV1Basic = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_cc"
}

resource "opentelekomcloud_waf_dedicated_cc_rule_v1" "rule_1" {
  policy_id    = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  mode         = 0
  url          = "/abc1"
  limit_num    = 10
  limit_period = 60
  lock_time    = 10
  tag_type     = "cookie"
  tag_index    = "sessionid"

  action {
    category     = "block"
    content_type = "application/json"
    content      = "{\"error\":\"forbidden\"}"
  }
}
`
