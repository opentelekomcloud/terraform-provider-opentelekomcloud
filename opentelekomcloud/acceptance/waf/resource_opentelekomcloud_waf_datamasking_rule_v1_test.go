package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/datamasking_rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceDMRuleName = "opentelekomcloud_waf_datamasking_rule_v1.rule_1"

func TestAccWafDataMaskingRuleV1_basic(t *testing.T) {
	var rule datamasking_rules.DataMasking

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDataMaskingRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDataMaskingRuleV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDataMaskingRuleV1Exists(resourceDMRuleName, &rule),
					resource.TestCheckResourceAttr(resourceDMRuleName, "url", "/login"),
					resource.TestCheckResourceAttr(resourceDMRuleName, "category", "params"),
					resource.TestCheckResourceAttr(resourceDMRuleName, "index", "password"),
				),
			},
			{
				Config: testAccWafDataMaskingRuleV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDataMaskingRuleV1Exists(resourceDMRuleName, &rule),
					resource.TestCheckResourceAttr(resourceDMRuleName, "url", "/login_new"),
					resource.TestCheckResourceAttr(resourceDMRuleName, "category", "params"),
					resource.TestCheckResourceAttr(resourceDMRuleName, "index", "password"),
				),
			},
		},
	})
}

func TestAccWafDataMaskingRuleV1_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDataMaskingRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDataMaskingRuleV1_basic,
			},
			stepWAFRuleImport(resourceDMRuleName),
		},
	})
}

func testAccCheckWafDataMaskingRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	wafClient, err := config.WafV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_data_masking_rule_v1" {
			continue
		}

		_, err := datamasking_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDataMaskingRuleV1Exists(n string, rule *datamasking_rules.DataMasking) resource.TestCheckFunc {
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

		found, err := datamasking_rules.Get(wafClient, rs.Primary.Attributes["policy_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("waf datamasking rule not found")
		}

		*rule = *found

		return nil
	}
}

const testAccWafDataMaskingRuleV1_basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
}

resource "opentelekomcloud_waf_datamasking_rule_v1" "rule_1" {
	policy_id = opentelekomcloud_waf_policy_v1.policy_1.id
	url = "/login"
	category = "params"
	index = "password"
}
`

const testAccWafDataMaskingRuleV1_update = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_updated"
}

resource "opentelekomcloud_waf_datamasking_rule_v1" "rule_1" {
	policy_id = opentelekomcloud_waf_policy_v1.policy_1.id
	url = "/login_new"
	category = "params"
	index = "password"
}
`
