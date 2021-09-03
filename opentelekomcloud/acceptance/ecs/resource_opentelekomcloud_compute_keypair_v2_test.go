package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/keypairs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceKeyPairName = "opentelekomcloud_compute_keypair_v2.kp_1"

func TestAccComputeV2Keypair_basic(t *testing.T) {
	var keypair keypairs.KeyPair

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2KeypairBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairExists(resourceKeyPairName, &keypair),
				),
			},
		},
	})
}

func TestAccComputeV2Keypair_shared(t *testing.T) {
	var keypair keypairs.KeyPair
	resourceName2 := "opentelekomcloud_compute_keypair_v2.kp_2"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2KeypairBasic,
			},
			{
				Config: testAccComputeV2KeypairShared,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairExists(resourceKeyPairName, &keypair),
					testAccCheckComputeV2KeypairExists(resourceName2, &keypair),
					resource.TestCheckResourceAttr(resourceName2, "shared", "true"),
					resource.TestCheckResourceAttr(resourceKeyPairName, "shared", "false"),
				),
			},
		},
	})
}

func TestAccComputeV2Keypair_private(t *testing.T) {
	var keypair keypairs.KeyPair

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2KeypairPrivate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairExists(resourceKeyPairName, &keypair),
					resource.TestCheckResourceAttrSet(resourceKeyPairName, "private_key"),
				),
			},
		},
	})
}

func testAccCheckComputeV2KeypairDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ComputeV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_keypair_v2" {
			continue
		}

		_, err := keypairs.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("keypair still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2KeypairExists(n string, kp *keypairs.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ComputeV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
		}

		found, err := keypairs.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("keypair not found")
		}

		*kp = *found

		return nil
	}
}

const testAccComputeV2KeypairBasic = `
resource "opentelekomcloud_compute_keypair_v2" "kp_1" {
  name       = "kp_1"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}
`

const testAccComputeV2KeypairShared = `
locals {
  public_name = "kp_1"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "opentelekomcloud_compute_keypair_v2" "kp_1" {
  name       = local.public_name
  public_key = local.public_key
}

resource "opentelekomcloud_compute_keypair_v2" "kp_2" {
  name       = local.public_name
  public_key = local.public_key

  depends_on = [opentelekomcloud_compute_keypair_v2.kp_1]
}
`

const testAccComputeV2KeypairPrivate = `
resource "opentelekomcloud_compute_keypair_v2" "kp_1" {
  name = "kp_1"
}
`
