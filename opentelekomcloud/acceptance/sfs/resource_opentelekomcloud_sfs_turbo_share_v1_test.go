package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs_turbo/v1/shares"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccSFSTurboShareV1_basic(t *testing.T) {
	shareName := tools.RandomString("sfs-turbo-", 3)
	resourceName := "opentelekomcloud_sfs_turbo_share_v1.sfs-turbo"
	var turbo shares.Turbo

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSFSTurboShareV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboShareV1_basic(shareName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboShareV1Exists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", shareName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
				),
			},
			{
				Config: testAccSFSTurboShareV1_update(shareName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboShareV1Exists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "size", "600"),
				),
			},
		},
	})
}

func TestAccSFSTurboShareV1_withKMS(t *testing.T) {
	postfix := acctest.RandString(5)
	resourceName := "opentelekomcloud_sfs_turbo_share_v1.sfs-turbo"
	var turbo shares.Turbo

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSFSTurboShareV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboV1_crypt(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboShareV1Exists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", "sfs-turbo-"+postfix),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
					resource.TestCheckResourceAttrSet(resourceName, "crypt_key_id"),
				),
			},
		},
	})
}

func testAccCheckSFSTurboShareV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SfsTurboV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud SFSTurboV1 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_sfs_turbo_share_v1" {
			continue
		}

		_, err := shares.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("sfs turbo still exists")
		}
	}

	return nil
}

func testAccCheckSFSTurboShareV1Exists(n string, share *shares.Turbo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.SfsTurboV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud SFSTurboV1 client: %s", err)
		}

		found, err := shares.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("sfs turbo not found")
		}

		*share = *found
		return nil
	}
}

func testAccSFSTurboShareV1_basic(shareName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-sfs-turbo-acc"
}

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name        = "%s"
  size        = 500
  share_proto = "NFS"
  vpc_id      = "%s"
  subnet_id   = "%s"

  security_group_id = opentelekomcloud_networking_secgroup_v2.sg.id
  availability_zone = "%s"
}
`, shareName, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE)
}

func testAccSFSTurboShareV1_update(shareName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-sfs-turbo-acc"
}

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name        = "%s"
  size        = 600
  share_proto = "NFS"
  vpc_id      = "%s"
  subnet_id   = "%s"

  security_group_id = opentelekomcloud_networking_secgroup_v2.sg.id
  availability_zone = "%s"
}
`, shareName, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE)
}

func testAccSFSTurboV1_crypt(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-sfs-turbo-acc"
}

resource "opentelekomcloud_kms_key_v1" "key_1" {
  key_alias    = "kms-sfs-turbo-%[1]s"
  pending_days = "7"
}

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name        = "sfs-turbo-%[1]s"
  size        = 500
  share_proto = "NFS"
  vpc_id      = "%s"
  subnet_id   = "%s"

  security_group_id = opentelekomcloud_networking_secgroup_v2.sg.id
  availability_zone = "%s"
  crypt_key_id      = opentelekomcloud_kms_key_v1.key_1.id
}
`, postfix, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE)
}
