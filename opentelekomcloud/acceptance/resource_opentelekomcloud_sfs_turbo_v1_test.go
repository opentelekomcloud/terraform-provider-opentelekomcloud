package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs_turbo/v1/shares"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccSFSTurboV1_basic(t *testing.T) {
	postfix := acctest.RandString(5)
	resourceName := "opentelekomcloud_sfs_turbo_v1.sfs-turbo"
	var turbo shares.Turbo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSTurboV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboV1_basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboV1Exists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", "sfs-turbo-"+postfix),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "enhanced", "false"),
					resource.TestCheckResourceAttr(resourceName, "size", "500.00"),
				),
			},
			{
				Config: testAccSFSTurboV1_update(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboV1Exists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "size", "600.00"),
				),
			},
		},
	})
}

func TestAccSFSTurboV1_withKMS(t *testing.T) {
	postfix := acctest.RandString(5)
	resourceName := "opentelekomcloud_sfs_turbo_v1.sfs-turbo"
	var turbo shares.Turbo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSTurboV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboV1_crypt(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboV1Exists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", "sfs-turbo-"+postfix),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "enhanced", "false"),
					resource.TestCheckResourceAttr(resourceName, "size", "500.00"),
					resource.TestCheckResourceAttrSet(resourceName, "crypt_key_id"),
				),
			},
		},
	})
}

func testAccCheckSFSTurboV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*cfg.Config)
	client, err := config.SfsTurboV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud SFSTurboV1 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_sfs_turbo_v1" {
			continue
		}

		_, err := shares.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SFS Turbo still exists")
		}
	}

	return nil
}

func testAccCheckSFSTurboV1Exists(n string, share *shares.Turbo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := testAccProvider.Meta().(*cfg.Config)
		client, err := config.SfsTurboV1Client(OS_REGION_NAME)
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

func testAccSFSTurboV1_basic(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-sfs-turbo-acc"
}

resource "opentelekomcloud_sfs_turbo" "sfs-turbo" {
  name        = "sfs-turbo-%s"
  size        = 500
  share_proto = "NFS"
  vpc_id      = "%s"
  subnet_id   = "%s"

  security_group_id = opentelekomcloud_networking_secgroup_v2.sg.id
  availability_zone = "%s"
}
`, postfix, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE)
}

func testAccSFSTurboV1_update(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-sfs-turbo-acc"
}

resource "opentelekomcloud_sfs_turbo" "sfs-turbo" {
  name        = "sfs-turbo-%s"
  size        = 600
  share_proto = "NFS"
  vpc_id      = "%s"
  subnet_id   = "%s"

  security_group_id = opentelekomcloud_networking_secgroup_v2.sg.id
  availability_zone = "%s"
}
`, postfix, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE)
}

func testAccSFSTurboV1_crypt(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-sfs-turbo-acc"
}

resource "opentelekomcloud_kms_key_v1" "key_1" {
  key_alias    = "kms-sfs-turbo-acc"
  pending_days = "7"
}

resource "opentelekomcloud_sfs_turbo_v1" "sfs-turbo" {
  name        = "sfs-turbo-%s"
  size        = 500
  share_proto = "NFS"
  vpc_id      = "%s"
  subnet_id   = "%s"

  security_group_id = opentelekomcloud_networking_secgroup_v2.sg.id
  availability_zone = "%s"
  crypt_key_id      = opentelekomcloud_kms_key_v1.key_1.id
}
`, postfix, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE)
}
