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

const wafdAlarmMaskingRuleName = "opentelekomcloud_waf_dedicated_alarm_masking_rule_v1.rule_1"

func TestAccWafDedicatedAlarmMaskingRuleV1_basic(t *testing.T) {
	var rule rules.IgnoreRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedAlarmMaskingRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedAlarmMaskingRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedAlarmMaskingRuleV1Exists(wafdAlarmMaskingRuleName, &rule),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "domains.#", "1"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "domains.0", "www.example.com"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "rule", "all"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "description", "description"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "conditions.0.category", "url"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "advanced_settings.0.index", "header"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "advanced_settings.0.contents.0", "content-type"),
				),
			},
			{
				Config: testAccWafDedicatedAlarmMaskingRuleV1AdvEmptyContents,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedAlarmMaskingRuleV1Exists(wafdAlarmMaskingRuleName, &rule),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "domains.#", "1"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "domains.0", "www.example.com"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "rule", "all"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "description", "description"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "conditions.0.category", "url"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "advanced_settings.0.index", "header"),
				),
			},
			{
				Config: testAccWafDedicatedAlarmMaskingRuleV1AdvEmptyBodyContents,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedAlarmMaskingRuleV1Exists(wafdAlarmMaskingRuleName, &rule),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "domains.#", "1"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "domains.0", "www.example.com"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "rule", "all"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "description", "description"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "conditions.0.category", "url"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "advanced_settings.0.index", "body"),
				),
			},
			{
				Config: testAccWafDedicatedAlarmMaskingRuleV1AdvEmptyBodyCookie,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedAlarmMaskingRuleV1Exists(wafdAlarmMaskingRuleName, &rule),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "domains.#", "1"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "domains.0", "www.example.com"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "rule", "all"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "description", "description"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "conditions.0.category", "url"),
					resource.TestCheckResourceAttr(wafdAlarmMaskingRuleName, "advanced_settings.0.index", "cookie"),
				),
			},
			{
				ResourceName:            wafdAlarmMaskingRuleName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       dedicatedRuleImportStateIDFunc(wafdAlarmMaskingRuleName, wafdPolicyResourceName),
				ImportStateVerifyIgnore: []string{"advanced_settings.0.contents"},
			},
		},
	})
}

func testAccCheckWafDedicatedAlarmMaskingRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_alarm_masking_rule_v1" {
			continue
		}

		_, err := rules.GetIgnore(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedAlarmMaskingRuleV1Exists(n string, rule *rules.IgnoreRule) resource.TestCheckFunc {
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

		found, err := rules.GetIgnore(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
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

const testAccWafDedicatedAlarmMaskingRuleV1Basic = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_am"
}

resource "opentelekomcloud_waf_dedicated_alarm_masking_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  domains     = ["www.example.com"]
  rule        = "all"
  description = "description"

  conditions {
    category        = "url"
    contents        = ["/login"]
    logic_operation = "equal"
  }
  advanced_settings {
    index    = "header"
    contents = ["content-type"]
  }
}
`

const testAccWafDedicatedAlarmMaskingRuleV1AdvEmptyContents = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_am"
}

resource "opentelekomcloud_waf_dedicated_alarm_masking_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  domains     = ["www.example.com"]
  rule        = "all"
  description = "description"

  conditions {
    category        = "url"
    contents        = ["/login"]
    logic_operation = "equal"
  }
  advanced_settings {
    index = "header"
    contents = ["all"]
  }
}
`

const testAccWafDedicatedAlarmMaskingRuleV1AdvEmptyBodyContents = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_am"
}

resource "opentelekomcloud_waf_dedicated_alarm_masking_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  domains     = ["www.example.com"]
  rule        = "all"
  description = "description"

  conditions {
    category        = "url"
    contents        = ["/login"]
    logic_operation = "equal"
  }
  advanced_settings {
    index = "body"
  }
}
`

const testAccWafDedicatedAlarmMaskingRuleV1AdvEmptyBodyCookie = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_am"
}

resource "opentelekomcloud_waf_dedicated_alarm_masking_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  domains     = ["www.example.com"]
  rule        = "all"
  description = "description"

  conditions {
    category        = "url"
    contents        = ["/login"]
    logic_operation = "equal"
  }
  advanced_settings {
    index = "cookie"
    contents = []
  }
}
`
