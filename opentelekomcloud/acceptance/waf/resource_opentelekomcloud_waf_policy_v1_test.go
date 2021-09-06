package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccWafPolicyV1_basic(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPolicyV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPolicyV1Exists("opentelekomcloud_waf_policy_v1.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_policy_v1.policy_1", "name", "policy_1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_policy_v1.policy_1", "level", "2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_policy_v1.policy_1", "full_detection", "false"),
				),
			},
			{
				Config: testAccWafPolicyV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPolicyV1Exists("opentelekomcloud_waf_policy_v1.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_policy_v1.policy_1", "name", "policy_updated"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_policy_v1.policy_1", "level", "1"),
				),
			},
		},
	})
}

func testAccCheckWafPolicyV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	wafClient, err := config.WafV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_policy_v1" {
			continue
		}

		_, err := policies.Get(wafClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("waf policy still exists")
		}
	}

	return nil
}

func testAccCheckWafPolicyV1Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
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

		found, err := policies.Get(wafClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("waf policy not found")
		}

		*policy = *found

		return nil
	}
}

const testAccWafPolicyV1_basic = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
	options {
		webattack = true
		crawler = true
	}
	full_detection = false
}
`

const testAccWafPolicyV1_update = `
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_updated"
	level = 1
	action {
		category = "block"
	}
}
`
