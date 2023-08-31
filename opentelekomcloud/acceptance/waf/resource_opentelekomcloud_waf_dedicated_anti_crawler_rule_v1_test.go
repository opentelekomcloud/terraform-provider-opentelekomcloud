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

const wafdAntiCrawlerRuleName = "opentelekomcloud_waf_dedicated_anti_crawler_rule_v1.rule_1"

func TestAccWafDedicatedAntiCrawlerRuleV1_basic(t *testing.T) {
	var rule rules.AntiCrawlerRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedAntiCrawlerRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedAntiCrawlerRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedAntiCrawlerRuleV1Exists(wafdAntiCrawlerRuleName, &rule),
					resource.TestCheckResourceAttr(wafdAntiCrawlerRuleName, "url", "/patent/id"),
					resource.TestCheckResourceAttr(wafdAntiCrawlerRuleName, "logic", "3"),
					resource.TestCheckResourceAttr(wafdAntiCrawlerRuleName, "protection_mode", "anticrawler_except_url"),
					resource.TestCheckResourceAttr(wafdAntiCrawlerRuleName, "name", "anticrawler_1"),
				),
			},
			{
				Config: testAccWafDedicatedAntiCrawlerRuleV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedAntiCrawlerRuleV1Exists(wafdAntiCrawlerRuleName, &rule),
					resource.TestCheckResourceAttr(wafdAntiCrawlerRuleName, "url", "/patent/id/update"),
					resource.TestCheckResourceAttr(wafdAntiCrawlerRuleName, "logic", "4"),
					resource.TestCheckResourceAttr(wafdAntiCrawlerRuleName, "name", "anticrawler_1_update"),
				),
			},
			{
				ResourceName:      wafdAntiCrawlerRuleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: dedicatedRuleImportStateIDFunc(wafdAntiCrawlerRuleName, wafdPolicyResourceName),
			},
		},
	})
}

func testAccCheckWafDedicatedAntiCrawlerRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_anti_crawler_rule_v1" {
			continue
		}

		_, err := rules.GetAntiCrawler(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedAntiCrawlerRuleV1Exists(n string, rule *rules.AntiCrawlerRule) resource.TestCheckFunc {
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

		found, err := rules.GetAntiCrawler(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
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

const testAccWafDedicatedAntiCrawlerRuleV1Basic = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_anticrawler"
}

resource "opentelekomcloud_waf_dedicated_anti_crawler_rule_v1" "rule_1" {
  policy_id       = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  name            = "anticrawler_1"
  url             = "/patent/id"
  logic           = 3
  protection_mode = "anticrawler_except_url"
}
`

const testAccWafDedicatedAntiCrawlerRuleV1Update = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_anticrawler"
}

resource "opentelekomcloud_waf_dedicated_anti_crawler_rule_v1" "rule_1" {
  policy_id       = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  name            = "anticrawler_1_update"
  url             = "/patent/id/update"
  logic           = 4
  protection_mode = "anticrawler_except_url"
}
`
