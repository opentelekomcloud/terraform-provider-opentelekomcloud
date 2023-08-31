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

const wafdDataMaskingRuleName = "opentelekomcloud_waf_dedicated_data_masking_rule_v1.rule_1"

func TestAccWafDedicatedDataMaskingRuleV1_basic(t *testing.T) {
	var rule rules.PrivacyRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedDataMaskingRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDataMaskingRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDataMaskingRuleV1Exists(wafdDataMaskingRuleName, &rule),
					resource.TestCheckResourceAttr(wafdDataMaskingRuleName, "name", "data_masking_1"),
					resource.TestCheckResourceAttr(wafdDataMaskingRuleName, "url", "/login"),
					resource.TestCheckResourceAttr(wafdDataMaskingRuleName, "category", "params"),
					resource.TestCheckResourceAttr(wafdDataMaskingRuleName, "description", "test description"),
				),
			},
			{
				Config: testAccWafDedicatedDataMaskingRuleV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDataMaskingRuleV1Exists(wafdDataMaskingRuleName, &rule),
					resource.TestCheckResourceAttr(wafdDataMaskingRuleName, "name", "data_masking_1_update"),
					resource.TestCheckResourceAttr(wafdDataMaskingRuleName, "url", "/login/update"),
					resource.TestCheckResourceAttr(wafdDataMaskingRuleName, "category", "header"),
					resource.TestCheckResourceAttr(wafdDataMaskingRuleName, "description", "test description update"),
				),
			},
			{
				ResourceName:      wafdDataMaskingRuleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: dedicatedRuleImportStateIDFunc(wafdDataMaskingRuleName, wafdPolicyResourceName),
			},
		},
	})
}

func testAccCheckWafDedicatedDataMaskingRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_data_masking_rule_v1" {
			continue
		}

		_, err := rules.GetPrivacy(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedDataMaskingRuleV1Exists(n string, rule *rules.PrivacyRule) resource.TestCheckFunc {
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

		found, err := rules.GetPrivacy(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
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

const testAccWafDedicatedDataMaskingRuleV1Basic = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_dm"
}

resource "opentelekomcloud_waf_dedicated_data_masking_rule_v1" "rule_1" {
  policy_id    = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  name        = "data_masking_1"
  url         = "/login"
  category    = "params"
  description = "test description"
}
`

const testAccWafDedicatedDataMaskingRuleV1Update = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_dm"
}

resource "opentelekomcloud_waf_dedicated_data_masking_rule_v1" "rule_1" {
  policy_id    = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  name        = "data_masking_1_update"
  url         = "/login/update"
  category    = "header"
  description = "test description update"
}
`
