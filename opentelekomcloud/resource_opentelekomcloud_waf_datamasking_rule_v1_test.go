package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/waf/v1/datamasking_rules"
)

func TestAccWafDataMaskingRuleV1_basic(t *testing.T) {
	var rule datamasking_rules.DataMasking

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafDataMaskingRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDataMaskingRuleV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDataMaskingRuleV1Exists("opentelekomcloud_waf_datamasking_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_datamasking_rule_v1.rule_1", "url", "/login"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_datamasking_rule_v1.rule_1", "category", "params"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_datamasking_rule_v1.rule_1", "index", "password"),
				),
			},
			{
				Config: testAccWafDataMaskingRuleV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDataMaskingRuleV1Exists("opentelekomcloud_waf_datamasking_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_datamasking_rule_v1.rule_1", "url", "/login_new"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_datamasking_rule_v1.rule_1", "category", "params"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_datamasking_rule_v1.rule_1", "index", "password"),
				),
			},
		},
	})
}

func testAccCheckWafDataMaskingRuleV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.wafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_data_masking_rule_v1" {
			continue
		}

		_, err := datamasking_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDataMaskingRuleV1Exists(n string, rule *datamasking_rules.DataMasking) resource.TestCheckFunc {
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

		found, err := datamasking_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Waf datamasking rule not found")
		}

		*rule = *found

		return nil
	}
}

const testAccWafDataMaskingRuleV1_basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
}

resource "opentelekomcloud_waf_datamasking_rule_v1" "rule_1" {
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	url = "/login"
	category = "params"
	index = "password"
}
`

const testAccWafDataMaskingRuleV1_update = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_updated"
}

resource "opentelekomcloud_waf_datamasking_rule_v1" "rule_1" {
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	url = "/login_new"
	category = "params"
	index = "password"
}
`
