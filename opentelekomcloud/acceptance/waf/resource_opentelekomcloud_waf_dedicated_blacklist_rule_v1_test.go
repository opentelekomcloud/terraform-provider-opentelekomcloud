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

const wafdBlacklistRuleName = "opentelekomcloud_waf_dedicated_blacklist_rule_v1.rule_1"

func TestAccWafDedicatedBlacklistRuleV1_basic(t *testing.T) {
	var rule rules.BlacklistRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedBlacklistRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedBlacklistRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedBlacklistRuleV1Exists(wafdBlacklistRuleName, &rule),
					resource.TestCheckResourceAttr(wafdBlacklistRuleName, "name", "my_blacklist"),
					resource.TestCheckResourceAttr(wafdBlacklistRuleName, "ip_address", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(wafdBlacklistRuleName, "action", "0"),
					resource.TestCheckResourceAttr(wafdBlacklistRuleName, "description", "test description"),
				),
			},
			{
				ResourceName:      wafdBlacklistRuleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: dedicatedRuleImportStateIDFunc(wafdBlacklistRuleName, wafdPolicyResourceName),
			},
		},
	})
}

func testAccCheckWafDedicatedBlacklistRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_blacklist_rule_v1" {
			continue
		}

		_, err := rules.GetBlacklist(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedBlacklistRuleV1Exists(n string, rule *rules.BlacklistRule) resource.TestCheckFunc {
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

		found, err := rules.GetBlacklist(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
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

const testAccWafDedicatedBlacklistRuleV1Basic = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_black"
}

resource "opentelekomcloud_waf_dedicated_blacklist_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  name        = "my_blacklist"
  ip_address  = "192.168.1.0/24"
  action      = 0
  description = "test description"
}
`
