package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/waf/v1/ccattackprotection_rules"
)

func TestAccWafCcAttackProtectionRuleV1_basic(t *testing.T) {
	var rule ccattackprotection_rules.CcAttack

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafCcAttackProtectionRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCcAttackProtectionRuleV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafCcAttackProtectionRuleV1Exists("opentelekomcloud_waf_ccattackprotection_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_ccattackprotection_rule_v1.rule_1", "url", "/abc1"),
				),
			},
		},
	})
}

func testAccCheckWafCcAttackProtectionRuleV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.wafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_ccattackprotection_rule_v1" {
			continue
		}

		_, err := ccattackprotection_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafCcAttackProtectionRuleV1Exists(n string, rule *ccattackprotection_rules.CcAttack) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		wafClient, err := config.wafV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
		}

		found, err := ccattackprotection_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Waf rule not found")
		}

		*rule = *found

		return nil
	}
}

const testAccWafCcAttackProtectionRuleV1_basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
}

resource "opentelekomcloud_waf_ccattackprotection_rule_v1" "rule_1" {
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	url = "/abc1"
	limit_num = 10
	limit_period = 60
	lock_time = 10
	tag_type = "cookie"
	tag_index = "sessionid"
	action_category = "block"
	block_content_type = "application/json"
	block_content = "{\"error\":\"forbidden\"}"
}
`
