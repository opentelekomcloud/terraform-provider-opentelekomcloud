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
				Config: testAccSFSTurboShareV1Basic(shareName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboShareV1Exists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", shareName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
				),
			},
			{
				Config: testAccSFSTurboShareV1Update(shareName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboShareV1Exists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "size", "600"),
				),
			},
		},
	})
}

func TestAccSFSTurboShareV1_enhanced(t *testing.T) {
	shareName := tools.RandomString("sfs-turbo-", 3)
	resourceName := "opentelekomcloud_sfs_turbo_share_v1.sfs-turbo"
	var turbo shares.Turbo

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSFSTurboShareV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboShareV1Enhanced(shareName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboShareV1Exists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", shareName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "PERFORMANCE"),
					resource.TestCheckResourceAttr(resourceName, "expand_type", "bandwidth"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
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
				Config: testAccSFSTurboV1Crypt(postfix),
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

func testAccSFSTurboShareV1Basic(shareName string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name        = "%s"
  size        = 500
  share_proto = "NFS"
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone = "%s"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, shareName, env.OS_AVAILABILITY_ZONE)
}

func testAccSFSTurboShareV1Update(shareName string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name        = "%s"
  size        = 600
  share_proto = "NFS"
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone = "%s"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, shareName, env.OS_AVAILABILITY_ZONE)
}

func testAccSFSTurboV1Crypt(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_kms_key_v1" "key_1" {
  key_alias    = "kms-sfs-turbo-%[3]s"
  pending_days = "7"
}

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name        = "sfs-turbo-%[3]s"
  size        = 500
  share_proto = "NFS"
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone = "%s"
  crypt_key_id      = opentelekomcloud_kms_key_v1.key_1.id
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccSFSTurboShareV1Enhanced(shareName string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name        = "%s"
  size        = 500
  share_proto = "NFS"
  share_type  = "PERFORMANCE"
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  enhanced    = true

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone = "%s"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, shareName, env.OS_AVAILABILITY_ZONE)
}
