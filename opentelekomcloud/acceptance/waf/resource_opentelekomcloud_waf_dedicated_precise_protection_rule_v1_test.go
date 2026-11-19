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
	rc := common.InitResourceCheck(wafdPreciseProtectionRuleName, &rule, getWafDedicatedPreciseProtectionRuleFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedPreciseProtectionRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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
	rc := common.InitResourceCheck(wafdPreciseProtectionRuleName, &rule, getWafDedicatedPreciseProtectionRuleFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedPreciseProtectionRuleV1Issue2434,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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

func getWafDedicatedPreciseProtectionRuleFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, err
	}
	found, err := rules.GetCustom(client, state.Primary.Attributes["policy_id"], state.Primary.ID)
	if err != nil {
		return nil, err
	}
	if found.ID != state.Primary.ID {
		return nil, fmt.Errorf("unexpected WAF dedicated rule ID %q, want %q", found.ID, state.Primary.ID)
	}
	return found, nil
}

func TestAccWafDedicatedPreciseProtectionRuleV1_referenceTable(t *testing.T) {
	var rule rules.CustomRule
	var ruleID, referenceID string
	rc := common.InitResourceCheck(wafdPreciseProtectionRuleName, &rule, getWafDedicatedPreciseProtectionRuleFunc)
	name := common.RandomAccResourceName()
	initialConfig := testAccWafDedicatedPreciseProtectionRuleV1ReferenceTable(name, "/admin")
	updatedConfig := testAccWafDedicatedPreciseProtectionRuleV1ReferenceTable(name, "/private")
	referenceName := "opentelekomcloud_waf_dedicated_reference_table_v1.urls"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(wafdPreciseProtectionRuleName, "conditions.0.logic_operation", "contain_any"),
					resource.TestCheckResourceAttrPair(wafdPreciseProtectionRuleName, "conditions.0.value_list_id", referenceName, "id"),
					func(_ *terraform.State) error {
						if len(rule.Conditions) != 1 || rule.Conditions[0].ValueListId == "" {
							return fmt.Errorf("expected one reference-table condition, got %#v", rule.Conditions)
						}
						ruleID, referenceID = rule.ID, rule.Conditions[0].ValueListId
						return nil
					},
				),
			},
			{
				Config:             initialConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      wafdPreciseProtectionRuleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: dedicatedRuleImportStateIDFunc(wafdPreciseProtectionRuleName, wafdPolicyResourceName),
			},
			{
				Config:             initialConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(referenceName, "conditions.0", "/private"),
					resource.TestCheckResourceAttrPair(wafdPreciseProtectionRuleName, "conditions.0.value_list_id", referenceName, "id"),
					func(_ *terraform.State) error {
						if rule.ID != ruleID || len(rule.Conditions) != 1 || rule.Conditions[0].ValueListId != referenceID {
							return fmt.Errorf("updating the reference table replaced the rule or its reference: %#v", rule)
						}
						return nil
					},
				),
			},
			{
				Config:             updatedConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccWafDedicatedPreciseProtectionRuleV1ReferenceTable(name, path string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = %q
}

resource "opentelekomcloud_waf_dedicated_reference_table_v1" "urls" {
  name       = %q
  type       = "url"
  conditions = [%q]
}

resource "opentelekomcloud_waf_dedicated_precise_protection_rule_v1" "rule_1" {
  policy_id = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  priority  = 50
  time      = false

  action {
    category = "block"
  }

  conditions {
    category        = "url"
    logic_operation = "contain_any"
    value_list_id   = opentelekomcloud_waf_dedicated_reference_table_v1.urls.id
  }
}
`, name, name, path)
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
