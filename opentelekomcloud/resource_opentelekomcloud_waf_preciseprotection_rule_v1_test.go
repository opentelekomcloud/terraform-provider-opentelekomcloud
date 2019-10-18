package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/waf/v1/preciseprotection_rules"
)

func TestAccWafPreciseProtectionRuleV1_basic(t *testing.T) {
	var rule preciseprotection_rules.Precise

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafPreciseProtectionRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPreciseProtectionRuleV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPreciseProtectionRuleV1Exists("opentelekomcloud_waf_preciseprotection_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_preciseprotection_rule_v1.rule_1", "name", "rule_1"),
				),
			},
		},
	})
}

func testAccCheckWafPreciseProtectionRuleV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.wafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_preciseprotection_rule_v1" {
			continue
		}

		_, err := preciseprotection_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafPreciseProtectionRuleV1Exists(n string, rule *preciseprotection_rules.Precise) resource.TestCheckFunc {
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

		found, err := preciseprotection_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
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

const testAccWafPreciseProtectionRuleV1_basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
}

resource "opentelekomcloud_waf_preciseprotection_rule_v1" "rule_1" {
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	name = "rule_1"
	conditions {
		category = "url"
		contents = ["/login"]
		logic = 1
	}
	conditions {
		category = "ip"
		contents = ["192.168.1.1"]
		logic = 3
	}
	action_category = "block"
	priority = 10
}
`
