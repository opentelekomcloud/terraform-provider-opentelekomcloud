package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/waf/v1/whiteblackip_rules"
)

func TestAccWafWhiteBlackIpRuleV1_basic(t *testing.T) {
	var rule whiteblackip_rules.WhiteBlackIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafWhiteBlackIpRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafWhiteBlackIpRuleV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafWhiteBlackIpRuleV1Exists("opentelekomcloud_waf_whiteblackip_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_whiteblackip_rule_v1.rule_1", "addr", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_whiteblackip_rule_v1.rule_1", "white", "0"),
				),
			},
			{
				Config: testAccWafWhiteBlackIpRuleV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafWhiteBlackIpRuleV1Exists("opentelekomcloud_waf_whiteblackip_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_whiteblackip_rule_v1.rule_1", "addr", "192.168.0.125"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_whiteblackip_rule_v1.rule_1", "white", "1"),
				),
			},
		},
	})
}

func testAccCheckWafWhiteBlackIpRuleV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.wafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_whiteblackip_rule_v1" {
			continue
		}

		_, err := whiteblackip_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafWhiteBlackIpRuleV1Exists(n string, rule *whiteblackip_rules.WhiteBlackIP) resource.TestCheckFunc {
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

		found, err := whiteblackip_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Waf whiteblackip rule not found")
		}

		*rule = *found

		return nil
	}
}

const testAccWafWhiteBlackIpRuleV1_basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_updated"
}

resource "opentelekomcloud_waf_whiteblackip_rule_v1" "rule_1" {
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	addr = "192.168.0.0/24"
}
`

const testAccWafWhiteBlackIpRuleV1_update = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_updated"
}

resource "opentelekomcloud_waf_whiteblackip_rule_v1" "rule_1" {
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	addr = "192.168.0.125"
	white = 1
}
`
