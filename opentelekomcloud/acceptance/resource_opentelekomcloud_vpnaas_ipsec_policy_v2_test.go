package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccVpnIPSecPolicyV2_basic(t *testing.T) {
	var policy ipsecpolicies.Policy
	resourceName := "opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(resourceName, &policy),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &policy.Name),
					resource.TestCheckResourceAttrPtr(resourceName, "description", &policy.Description),
					resource.TestCheckResourceAttrPtr(resourceName, "tenant_id", &policy.TenantID),
					resource.TestCheckResourceAttrPtr(resourceName, "pfs", &policy.PFS),
					resource.TestCheckResourceAttrPtr(resourceName, "transform_protocol", &policy.TransformProtocol),
					resource.TestCheckResourceAttrPtr(resourceName, "encapsulation_mode", &policy.EncapsulationMode),
					resource.TestCheckResourceAttrPtr(resourceName, "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr(resourceName, "encryption_algorithm", &policy.EncryptionAlgorithm),
				),
			},
			{
				Config: testAccIPSecPolicyV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(resourceName, &policy),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccVpnIPSecPolicyV2_withLifetime(t *testing.T) {
	var policy ipsecpolicies.Policy
	resourceName := "opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_withLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(resourceName, &policy),
					resource.TestCheckResourceAttrPtr(resourceName, "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr(resourceName, "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIPSecPolicyV2_withLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(resourceName, &policy),
				),
			},
		},
	})
}

func testAccCheckIPSecPolicyV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpnaas_ipsec_policy_v2" {
			continue
		}
		_, err = ipsecpolicies.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IPSec policy (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckIPSecPolicyV2Exists(n string, policy *ipsecpolicies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := testAccProvider.Meta().(*cfg.Config)
		client, err := config.NetworkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
		}

		found, err := ipsecpolicies.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*policy = *found

		return nil
	}
}

const testAccIPSecPolicyV2_basic = `
resource "opentelekomcloud_vpnaas_ipsec_policy_v2" "policy_1" { }
`

const testAccIPSecPolicyV2_update = `
resource "opentelekomcloud_vpnaas_ipsec_policy_v2" "policy_1" {
  name = "updatedname"
}
`

const testAccIPSecPolicyV2_withLifetime = `
resource "opentelekomcloud_vpnaas_ipsec_policy_v2" "policy_1" {
  auth_algorithm = "md5"
  pfs            = "group14"
  lifetime {
    units = "seconds"
    value = 1200
  }
}
`

const testAccIPSecPolicyV2_withLifetimeUpdate = `
resource "opentelekomcloud_vpnaas_ipsec_policy_v2" "policy_1" {
  auth_algorithm = "md5"
  pfs            = "group14"
  lifetime {
    units = "seconds"
    value = 1400
  }
}
`
