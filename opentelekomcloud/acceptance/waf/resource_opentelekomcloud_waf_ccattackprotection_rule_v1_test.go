package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/ccattackprotection_rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceCCRuleName = "opentelekomcloud_waf_ccattackprotection_rule_v1.rule_1"

func TestAccWafCcAttackProtectionRuleV1_basic(t *testing.T) {
	var rule ccattackprotection_rules.CcAttack

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafCcAttackProtectionRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCcAttackProtectionRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafCcAttackProtectionRuleV1Exists(resourceCCRuleName, &rule),
					resource.TestCheckResourceAttr(resourceCCRuleName, "url", "/abc1"),
				),
			},
		},
	})
}

func TestAccWafCcAttackProtectionRuleV1_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafCcAttackProtectionRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCcAttackProtectionRuleV1Basic,
			},
			stepWAFRuleImport(resourceCCRuleName),
		},
	})
}

func testAccCheckWafCcAttackProtectionRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	wafClient, err := config.WafV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_ccattackprotection_rule_v1" {
			continue
		}

		_, err := ccattackprotection_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafCcAttackProtectionRuleV1Exists(n string, rule *ccattackprotection_rules.CcAttack) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		wafClient, err := config.WafV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
		}

		found, err := ccattackprotection_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("waf rule not found")
		}

		*rule = *found

		return nil
	}
}

const testAccWafCcAttackProtectionRuleV1Basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_waf_ccattackprotection_rule_v1" "rule_1" {
  policy_id          = opentelekomcloud_waf_policy_v1.policy_1.id
  url                = "/abc1"
  limit_num          = 10
  limit_period       = 60
  lock_time          = 10
  tag_type           = "cookie"
  tag_index          = "sessionid"
  action_category    = "block"
  block_content_type = "application/json"
  block_content      = "{\"error\":\"forbidden\"}"
}
`
