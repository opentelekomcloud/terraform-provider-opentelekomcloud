package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/whiteblackip_rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceRuleName = "opentelekomcloud_waf_whiteblackip_rule_v1.rule_1"

func TestAccWafWhiteBlackIpRuleV1_basic(t *testing.T) {
	var rule whiteblackip_rules.WhiteBlackIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafWhiteBlackIpRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafWhiteBlackIpRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafWhiteBlackIpRuleV1Exists(resourceRuleName, &rule),
					resource.TestCheckResourceAttr(resourceRuleName, "addr", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resourceRuleName, "white", "0"),
				),
			},
			{
				Config: testAccWafWhiteBlackIpRuleV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafWhiteBlackIpRuleV1Exists(resourceRuleName, &rule),
					resource.TestCheckResourceAttr(resourceRuleName, "addr", "192.168.0.125"),
					resource.TestCheckResourceAttr(resourceRuleName, "white", "1"),
				),
			},
		},
	})
}

func testAccCheckWafWhiteBlackIpRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	wafClient, err := config.WafV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_whiteblackip_rule_v1" {
			continue
		}

		_, err := whiteblackip_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafWhiteBlackIpRuleV1Exists(n string, rule *whiteblackip_rules.WhiteBlackIP) resource.TestCheckFunc {
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

		found, err := whiteblackip_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("waf whiteblackip rule not found")
		}

		*rule = *found

		return nil
	}
}

const testAccWafWhiteBlackIpRuleV1Basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_updated"
}

resource "opentelekomcloud_waf_whiteblackip_rule_v1" "rule_1" {
  policy_id = opentelekomcloud_waf_policy_v1.policy_1.id
  addr      = "192.168.0.0/24"
}
`

const testAccWafWhiteBlackIpRuleV1Update = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_updated"
}

resource "opentelekomcloud_waf_whiteblackip_rule_v1" "rule_1" {
  policy_id = opentelekomcloud_waf_policy_v1.policy_1.id
  addr      = "192.168.0.125"
  white     = 1
}
`
