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
	t.Parallel()

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
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2KeypairSharedPre,
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
	t.Parallel()

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
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIALRzbIOR9HUYNwfKtII/et98eGXDJhf8YxHf9BtRdAU"
}
`

const testAccComputeV2KeypairSharedPre = `
locals {
  public_name = "kp_1-shared"
  public_key  = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINnnG9HsMplxW056UKoJWeiWYEMBZ0fKQoMOaPFRA5Zp"
}

resource "opentelekomcloud_compute_keypair_v2" "kp_1" {
  name       = local.public_name
  public_key = local.public_key
}
`

const testAccComputeV2KeypairShared = `
locals {
  public_name = "kp_1-shared"
  public_key  = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINnnG9HsMplxW056UKoJWeiWYEMBZ0fKQoMOaPFRA5Zp"
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
  name = "kp_1-private"
}
`
