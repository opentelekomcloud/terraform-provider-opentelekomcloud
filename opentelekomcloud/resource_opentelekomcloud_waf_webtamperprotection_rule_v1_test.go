package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/waf/v1/webtamperprotection_rules"
)

func TestAccWebTamperProtectionRuleV1_basic(t *testing.T) {
	var rule webtamperprotection_rules.WebTamper

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafWebTamperProtectionRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafWebTamperProtectionRuleV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebTamperProtectionRuleV1Exists("opentelekomcloud_waf_webtamperprotection_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_webtamperprotection_rule_v1.rule_1", "hostname", "www.abc.com"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_webtamperprotection_rule_v1.rule_1", "url", "/a"),
				),
			},
		},
	})
}

func testAccCheckWafWebTamperProtectionRuleV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.wafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_webtamperprotection_rule_v1" {
			continue
		}

		_, err := webtamperprotection_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWebTamperProtectionRuleV1Exists(n string, rule *webtamperprotection_rules.WebTamper) resource.TestCheckFunc {
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

		found, err := webtamperprotection_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Waf web tamper protection rule not found")
		}

		*rule = *found

		return nil
	}
}

const testAccWafWebTamperProtectionRuleV1_basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_updated"
}

resource "opentelekomcloud_waf_webtamperprotection_rule_v1" "rule_1" {
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	hostname = "www.abc.com"
	url = "/a"
}
`
