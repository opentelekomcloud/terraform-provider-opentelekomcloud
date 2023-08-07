package acceptance

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const wafdPolicyResourceName = "opentelekomcloud_waf_dedicated_policy_v1.policy_1"

func TestAccWafDedicatedPolicyV1_basic(t *testing.T) {
	var policy policies.Policy
	var policyName = fmt.Sprintf("wafd_policy_%s", acctest.RandString(5))
	log.Printf("[DEBUG] The opentelekomcloud Waf dedicated instance test running in '%s' region.", env.OS_REGION_NAME)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedPolicyV1_basic(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedPolicyV1Exists(wafdPolicyResourceName, &policy),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "name", policyName),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "level", "2"),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "full_detection", "false"),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "protection_mode", "log"),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "options.0.web_attack", "true"),
				),
			},
			{
				Config: testAccWafDedicatedPolicyV1_update(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedPolicyV1Exists(wafdPolicyResourceName, &policy),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "name", policyName+"-updated"),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "level", "3"),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "full_detection", "true"),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "protection_mode", "block"),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "options.0.web_attack", "false"),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "options.0.cc", "true"),
					resource.TestCheckResourceAttr(wafdPolicyResourceName, "options.0.web_shell", "true"),
				),
			},
			{
				ResourceName:      wafdPolicyResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckWafDedicatedPolicyV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_policy_v1" {
			continue
		}

		_, err := policies.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated policy (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckWafDedicatedPolicyV1Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
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

		found, err := policies.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return err
		}
		*policy = *found

		return nil
	}
}

func testAccWafDedicatedPolicyV1_basic(policyName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name            = "%s"
  protection_mode = "log"
  full_detection  = false
  level           = 2

  options {
    crawler    = true
    web_attack = true
  }
}
`, policyName)
}

func testAccWafDedicatedPolicyV1_update(policyName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name            = "%s-updated"
  level           = 3
  protection_mode = "block"
  full_detection  = true

  options {
    crawler    = false
    web_attack = false
    cc         = true
    web_shell  = true
  }
}
`, policyName)
}
