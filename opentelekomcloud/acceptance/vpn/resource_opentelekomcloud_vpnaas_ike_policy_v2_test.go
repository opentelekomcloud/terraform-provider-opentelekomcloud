package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccVpnIKEPolicyV2_basic(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIKEPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", "name", &policy.Name),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", "description", &policy.Description),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", "tenant_id", &policy.TenantID),
				),
			},
			{
				Config: testAccIKEPolicyV2_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccVpnIKEPolicyV2_withLifetime(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIKEPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2_withLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIKEPolicyV2_withLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", &policy),
				),
			},
		},
	})
}

func TestAccVpnIKEPolicyV2_withNewParams(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIKEPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2_withNewParams,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttr("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", "ike_version", "v2"),
					resource.TestCheckResourceAttr("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", "phase1_negotiation_mode", "aggressive"),
					resource.TestCheckResourceAttr("opentelekomcloud_vpnaas_ike_policy_v2.policy_1", "pfs", "group16"),
				),
			},
		},
	})
}

func testAccCheckIKEPolicyV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpnaas_ike_policy_v2" {
			continue
		}
		_, err = ikepolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IKE policy (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckIKEPolicyV2Exists(n string, policy *ikepolicies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
		}

		found, err := ikepolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*policy = *found

		return nil
	}
}

const testAccIKEPolicyV2_basic = `
resource "opentelekomcloud_vpnaas_ike_policy_v2" "policy_1" {}
`

const testAccIKEPolicyV2_Update = `
resource "opentelekomcloud_vpnaas_ike_policy_v2" "policy_1" {
  name = "updatedname"
}
`

const testAccIKEPolicyV2_withLifetime = `
resource "opentelekomcloud_vpnaas_ike_policy_v2" "policy_1" {
  auth_algorithm = "sha2-256"
  pfs            = "group14"
  lifetime {
    units = "seconds"
    value = 1200
  }
}
`

const testAccIKEPolicyV2_withLifetimeUpdate = `
resource "opentelekomcloud_vpnaas_ike_policy_v2" "policy_1" {
  auth_algorithm = "sha2-256"
  pfs            = "group14"
  lifetime {
    units = "seconds"
    value = 1400
  }
}
`

const testAccIKEPolicyV2_withNewParams = `
resource "opentelekomcloud_vpnaas_ike_policy_v2" "policy_1" {
  auth_algorithm = "sha2-256"
  pfs            = "group16"
  ike_version    = "v2"

  phase1_negotiation_mode = "aggressive"

  lifetime {
    units = "seconds"
    value = 120
  }
}
`
