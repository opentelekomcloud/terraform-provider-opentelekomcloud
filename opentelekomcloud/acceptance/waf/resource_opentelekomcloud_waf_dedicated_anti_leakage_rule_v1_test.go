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

const wafdAntiLeakageRuleName = "opentelekomcloud_waf_dedicated_anti_leakage_rule_v1.rule_1"

func TestAccWafDedicatedAntiLeakageRuleV1_basic(t *testing.T) {
	var rule rules.AntiLeakageRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedAntiLeakageRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedAntiLeakageRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedAntiLeakageRuleV1Exists(wafdAntiLeakageRuleName, &rule),
					resource.TestCheckResourceAttr(wafdAntiLeakageRuleName, "url", "/attack"),
					resource.TestCheckResourceAttr(wafdAntiLeakageRuleName, "contents.#", "1"),
					resource.TestCheckResourceAttr(wafdAntiLeakageRuleName, "contents.0", "id_card"),
					resource.TestCheckResourceAttr(wafdAntiLeakageRuleName, "category", "sensitive"),
					resource.TestCheckResourceAttr(wafdAntiLeakageRuleName, "description", "test description"),
				),
			},
			{
				Config: testAccWafDedicatedAntiLeakageRuleV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedAntiLeakageRuleV1Exists(wafdAntiLeakageRuleName, &rule),
					resource.TestCheckResourceAttr(wafdAntiLeakageRuleName, "url", "/pass"),
					resource.TestCheckResourceAttr(wafdAntiLeakageRuleName, "contents.#", "1"),
					resource.TestCheckResourceAttr(wafdAntiLeakageRuleName, "contents.0", "id_card"),
					resource.TestCheckResourceAttr(wafdAntiLeakageRuleName, "category", "sensitive"),
					resource.TestCheckResourceAttr(wafdAntiLeakageRuleName, "description", "test description updated"),
				),
			},
			{
				ResourceName:      wafdAntiLeakageRuleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: dedicatedRuleImportStateIDFunc(wafdAntiLeakageRuleName, wafdPolicyResourceName),
			},
		},
	})
}

func testAccCheckWafDedicatedAntiLeakageRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_anti_leakage_rule_v1" {
			continue
		}

		_, err := rules.GetAntiLeakage(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedAntiLeakageRuleV1Exists(n string, rule *rules.AntiLeakageRule) resource.TestCheckFunc {
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

		found, err := rules.GetAntiLeakage(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
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

const testAccWafDedicatedAntiLeakageRuleV1Basic = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_al"
}

resource "opentelekomcloud_waf_dedicated_anti_leakage_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  url         = "/attack"
  category    = "sensitive"
  contents    = ["id_card"]
  description = "test description"
}
`

const testAccWafDedicatedAntiLeakageRuleV1Update = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_al"
}

resource "opentelekomcloud_waf_dedicated_anti_leakage_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  url         = "/pass"
  category    = "sensitive"
  contents    = ["id_card"]
  description = "test description updated"
}
`
