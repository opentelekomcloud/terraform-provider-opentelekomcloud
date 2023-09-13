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

const wafdWebTamperRuleName = "opentelekomcloud_waf_dedicated_web_tamper_rule_v1.rule_1"

func TestAccWafDedicatedWebTamperRuleV1_basic(t *testing.T) {
	var rule rules.AntiTamperRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedWebTamperRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedWebTamperRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedWebTamperRuleV1Exists(wafdWebTamperRuleName, &rule),
					resource.TestCheckResourceAttr(wafdWebTamperRuleName, "hostname", "www.domain.com"),
					resource.TestCheckResourceAttr(wafdWebTamperRuleName, "url", "/login"),
					resource.TestCheckResourceAttr(wafdWebTamperRuleName, "description", "test description"),
				),
			},
			{
				Config: testAccWafDedicatedWebTamperRuleV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedWebTamperRuleV1Exists(wafdWebTamperRuleName, &rule),
					resource.TestCheckResourceAttr(wafdWebTamperRuleName, "hostname", "www.domain.com"),
					resource.TestCheckResourceAttr(wafdWebTamperRuleName, "url", "/login"),
					resource.TestCheckResourceAttr(wafdWebTamperRuleName, "description", "test description"),
				),
			},
			{
				ResourceName:            wafdWebTamperRuleName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       dedicatedRuleImportStateIDFunc(wafdWebTamperRuleName, wafdPolicyResourceName),
				ImportStateVerifyIgnore: []string{"update_cache"},
			},
		},
	})
}

func testAccCheckWafDedicatedWebTamperRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_web_tamper_rule_v1" {
			continue
		}

		_, err := rules.GetAntiTamper(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedWebTamperRuleV1Exists(n string, rule *rules.AntiTamperRule) resource.TestCheckFunc {
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

		found, err := rules.GetAntiTamper(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
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

const testAccWafDedicatedWebTamperRuleV1Basic = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_at"
}

resource "opentelekomcloud_waf_dedicated_web_tamper_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  hostname    = "www.domain.com"
  url         = "/login"
  description = "test description"
}
`

const testAccWafDedicatedWebTamperRuleV1Update = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_at"
}

resource "opentelekomcloud_waf_dedicated_web_tamper_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  hostname    = "www.domain.com"
  url         = "/login"
  description = "test description"

  update_cache = true
}
`
