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

const wafdGeoIpRuleName = "opentelekomcloud_waf_dedicated_geo_ip_rule_v1.rule_1"

func TestAccWafDedicatedGeoIpRuleV1_basic(t *testing.T) {
	var rule rules.GeoIpRule

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedGeoIpRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedGeoIpRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedGeoIPRuleV1Exists(wafdGeoIpRuleName, &rule),
					resource.TestCheckResourceAttr(wafdGeoIpRuleName, "region_code", "BR"),
					resource.TestCheckResourceAttr(wafdGeoIpRuleName, "action", "0"),
					resource.TestCheckResourceAttr(wafdGeoIpRuleName, "name", "test"),
					resource.TestCheckResourceAttr(wafdGeoIpRuleName, "description", "test description"),
				),
			},
			{
				Config: testAccWafDedicatedGeoIpRuleV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedGeoIPRuleV1Exists(wafdGeoIpRuleName, &rule),
					resource.TestCheckResourceAttr(wafdGeoIpRuleName, "region_code", "DE"),
					resource.TestCheckResourceAttr(wafdGeoIpRuleName, "action", "1"),
					resource.TestCheckResourceAttr(wafdGeoIpRuleName, "name", "testUpdate"),
					resource.TestCheckResourceAttr(wafdGeoIpRuleName, "description", "test description updated"),
				),
			},
			{
				ResourceName:      wafdGeoIpRuleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: dedicatedRuleImportStateIDFunc(wafdGeoIpRuleName, wafdPolicyResourceName),
			},
		},
	})
}

func testAccCheckWafDedicatedGeoIpRuleV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_geo_ip_rule_v1" {
			continue
		}

		_, err := rules.GetGeoIp(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated rule still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedGeoIPRuleV1Exists(n string, rule *rules.GeoIpRule) resource.TestCheckFunc {
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

		found, err := rules.GetGeoIp(client, rs.Primary.Attributes["policy_id"], rs.Primary.ID)
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

const testAccWafDedicatedGeoIpRuleV1Basic = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_gi"
}

resource "opentelekomcloud_waf_dedicated_geo_ip_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  region_code = "BR"
  action      = 0
  name        = "test"
  description = "test description"
}
`

const testAccWafDedicatedGeoIpRuleV1Update = `
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_gi"
}

resource "opentelekomcloud_waf_dedicated_geo_ip_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  region_code = "DE"
  action      = 1
  name        = "testUpdate"
  description = "test description updated"
}
`
