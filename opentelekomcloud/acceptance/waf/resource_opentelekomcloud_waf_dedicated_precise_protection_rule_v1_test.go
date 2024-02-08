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

const wafdPreciseProtectionRuleName = "opentelekomcloud_waf_dedicated_precise_protection_rule_v1.rule_1"

func TestAccWafDedicatedPreciseProtectionRuleV1_basic(t *testing.T) {
	var rule rules.CustomRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedPreciseProtectionRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedPreciseProtectionRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedPreciseProtectionRuleV1Exists(wafdPreciseProtectionRuleName, &rule),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "time", "false"),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "description", "desc"),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "priority", "50"),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "action.0.category", "block"),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "conditions.0.category", "url"),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "conditions.0.logic_operation", "contain"),
				),
			},
			{
				ResourceName:      wafdPreciseProtectionRuleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: dedicatedRuleImportStateIDFunc(wafdPreciseProtectionRuleName, wafdPolicyResourceName),
			},
		},
	})
}

func TestAccWafDedicatedPreciseProtectionRuleV1_Issue2434(t *testing.T) {
	var rule rules.CustomRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedPreciseProtectionRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedPreciseProtectionRuleV1Issue2434,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedPreciseProtectionRuleV1Exists(wafdPreciseProtectionRuleName, &rule),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "time", "false"),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "priority", "50"),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "action.0.category", "pass"),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "conditions.0.category", "url"),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "conditions.0.logic_operation", "contain"),
				),
			},
		},
	})
}

func testAccCheckWafDedicatedPreciseProtectionRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_precise_protection_rule_v1" {
			continue
		}

		_, err := rules.GetCustom(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedPreciseProtectionRuleV1Exists(n string, rule *rules.CustomRule) resource.TestCheckFunc {
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

		found, err := rules.GetCustom(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
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

const testAccWafDedicatedPreciseProtectionRuleV1Basic = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_pp"
}

resource "opentelekomcloud_waf_dedicated_precise_protection_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  time        = false
  description = "desc"
  priority    = 50

  conditions {
    category        = "url"
    contents        = ["test"]
    logic_operation = "contain"
  }
  action {
    category = "block"
  }
}
`

const testAccWafDedicatedPreciseProtectionRuleV1Issue2434 = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_pp"
}

resource "opentelekomcloud_waf_dedicated_precise_protection_rule_v1" "rule_1" {
  policy_id = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  priority  = 50
  time      = false
  action {
    category = "pass"
  }
  conditions {
    category = "url"
    contents = [
      "/xxx-yyyy-zzz.html"
    ]
    logic_operation = "contain"
  }
  conditions {
    category = "params"
    contents = [
      "32"
    ]
    index           = "u"
    logic_operation = "len_equal"
  }
  conditions {
    category = "params"
    contents = [
      "3"
    ]
    index           = "t"
    logic_operation = "len_greater"
  }
  conditions {
    category = "params"
    contents = [
      "13"
    ]
    index           = "t"
    logic_operation = "len_less"
  }
}
`
