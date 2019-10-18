package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/waf/v1/falsealarmmasking_rules"
)

func TestAccWafFalseAlarmMaskingRuleV1_basic(t *testing.T) {
	var rule falsealarmmasking_rules.AlarmMasking

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafFalseAlarmMaskingRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafFalseAlarmMaskingRuleV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafFalseAlarmMaskingRuleV1Exists("opentelekomcloud_waf_falsealarmmasking_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_falsealarmmasking_rule_v1.rule_1", "url", "/a"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_falsealarmmasking_rule_v1.rule_1", "rule", "100001"),
				),
			},
		},
	})
}

func testAccCheckWafFalseAlarmMaskingRuleV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.wafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_falsealarmmasking_rule_v1" {
			continue
		}

		rules, err := falsealarmmasking_rules.List(wafClient, rs.Primary.Attributes["policy_id"]).Extract()
		if err != nil {
			return nil
		}
		for _, r := range rules {
			if r.Id == rs.Primary.ID {
				return fmt.Errorf("Waf rule still exists")
			}
		}
	}

	return nil
}

func testAccCheckWafFalseAlarmMaskingRuleV1Exists(n string, rule *falsealarmmasking_rules.AlarmMasking) resource.TestCheckFunc {
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

		rules, err := falsealarmmasking_rules.List(wafClient, rs.Primary.Attributes["policy_id"]).Extract()
		if err != nil {
			return err
		}
		for _, r := range rules {
			if r.Id == rs.Primary.ID {
				*rule = r
				return nil
			}
		}

		return fmt.Errorf("Waf falsealarmmasking rule not found")
	}
}

const testAccWafFalseAlarmMaskingRuleV1_basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
}

resource "opentelekomcloud_waf_falsealarmmasking_rule_v1" "rule_1" {
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	url = "/a"
	rule = "100001"
}
`
