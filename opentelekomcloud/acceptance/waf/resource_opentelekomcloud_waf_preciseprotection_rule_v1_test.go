package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/preciseprotection_rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourcePPRuleName = "opentelekomcloud_waf_preciseprotection_rule_v1.rule_1"

func TestAccWafPreciseProtectionRuleV1_basic(t *testing.T) {
	var rule preciseprotection_rules.Precise

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafPreciseProtectionRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPreciseProtectionRuleV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPreciseProtectionRuleV1Exists(resourcePPRuleName, &rule),
					resource.TestCheckResourceAttr(resourcePPRuleName, "name", "rule_1"),
				),
			},
		},
	})
}

func TestAccWafPreciseProtectionRuleV1_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafPreciseProtectionRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPreciseProtectionRuleV1_basic,
			},
			stepWAFRuleImport(resourcePPRuleName),
		},
	})
}

func testAccCheckWafPreciseProtectionRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	wafClient, err := config.WafV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_preciseprotection_rule_v1" {
			continue
		}

		_, err := preciseprotection_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafPreciseProtectionRuleV1Exists(n string, rule *preciseprotection_rules.Precise) resource.TestCheckFunc {
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

		found, err := preciseprotection_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
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

const testAccWafPreciseProtectionRuleV1_basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
}

resource "opentelekomcloud_waf_preciseprotection_rule_v1" "rule_1" {
	policy_id = opentelekomcloud_waf_policy_v1.policy_1.id
	name = "rule_1"
	conditions {
		category = "path"
		contents = ["/login"]
		logic = "contain"
	}
	conditions {
		category = "ip"
		contents = ["192.168.1.1"]
		logic = "equal"
	}
	action_category = "block"
	priority = 10
}
`
