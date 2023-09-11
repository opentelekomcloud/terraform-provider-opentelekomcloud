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

const wafdKnownAttackRuleName = "opentelekomcloud_waf_dedicated_known_attack_source_rule_v1.rule_1"

func TestAccWafDedicatedKnownAttackRuleV1_basic(t *testing.T) {
	var rule rules.KnownAttackSourceRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedKnownAttackRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedKnownAttackRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedKnownAttackRuleV1Exists(wafdKnownAttackRuleName, &rule),
					resource.TestCheckResourceAttr(wafdKnownAttackRuleName, "block_time", "300"),
					resource.TestCheckResourceAttr(wafdKnownAttackRuleName, "category", "long_cookie_block"),
					resource.TestCheckResourceAttr(wafdKnownAttackRuleName, "description", "test description"),
				),
			},
			{
				Config: testAccWafDedicatedKnownAttackRuleV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedKnownAttackRuleV1Exists(wafdKnownAttackRuleName, &rule),
					resource.TestCheckResourceAttr(wafdKnownAttackRuleName, "block_time", "1200"),
					resource.TestCheckResourceAttr(wafdKnownAttackRuleName, "category", "long_cookie_block"),
					resource.TestCheckResourceAttr(wafdKnownAttackRuleName, "description", "test description update"),
				),
			},
			{
				ResourceName:      wafdKnownAttackRuleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: dedicatedRuleImportStateIDFunc(wafdKnownAttackRuleName, wafdPolicyResourceName),
			},
		},
	})
}

func testAccCheckWafDedicatedKnownAttackRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_known_attack_source_rule_v1" {
			continue
		}

		_, err := rules.GetKnownAttackSource(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedKnownAttackRuleV1Exists(n string, rule *rules.KnownAttackSourceRule) resource.TestCheckFunc {
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

		found, err := rules.GetKnownAttackSource(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
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

const testAccWafDedicatedKnownAttackRuleV1Basic = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_ka"
}

resource "opentelekomcloud_waf_dedicated_known_attack_source_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  block_time  = 300
  category    = "long_cookie_block"
  description = "test description"
}
`

const testAccWafDedicatedKnownAttackRuleV1Update = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_ka"
}

resource "opentelekomcloud_waf_dedicated_known_attack_source_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  block_time  = 1200
  category    = "long_cookie_block"
  description = "test description update"
}
`
