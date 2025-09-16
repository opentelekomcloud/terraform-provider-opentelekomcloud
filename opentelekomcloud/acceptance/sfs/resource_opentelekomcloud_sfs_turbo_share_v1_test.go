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

const sfsShareResourceName = "opentelekomcloud_sfs_turbo_share_v1.sfs-turbo"

func getSFSTurboShare(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.SfsTurboV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS Turbo Client: %s", err)
	}
	return shares.Get(client, state.Primary.ID)
}

func TestAccSFSTurboShareV1_basic(t *testing.T) {
	shareName := tools.RandomString("sfs-turbo-", 3)
	var turbo shares.Turbo
	rc := common.InitResourceCheck(
		sfsShareResourceName,
		&turbo,
		getSFSTurboShare,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboShareV1Basic(shareName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(sfsShareResourceName, "name", shareName),
					resource.TestCheckResourceAttr(sfsShareResourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(sfsShareResourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(sfsShareResourceName, "size", "500"),
				),
			},
			{
				Config: testAccSFSTurboShareV1Update(shareName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(sfsShareResourceName, "size", "600"),
				),
			},
			{
				ResourceName:      sfsShareResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
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

func TestAccSFSTurboShareV1_hpc(t *testing.T) {
	shareName := tools.RandomString("sfs-turbo-", 3)
	var turbo shares.Turbo
	rc := common.InitResourceCheck(
		sfsShareResourceName,
		&turbo,
		getSFSTurboShare,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboShareV1HPC(shareName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(sfsShareResourceName, "name", shareName),
					resource.TestCheckResourceAttr(sfsShareResourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(sfsShareResourceName, "expand_type", "hpc"),
					resource.TestCheckResourceAttr(sfsShareResourceName, "size", "3686"),
				),
			},
		},
	})
}

func testAccSFSTurboShareV1HPC(shareName string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name        = "%s"
  size        = 3686
  share_proto = "NFS"
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  expand_type = "hpc"
  hpc_bw      = "20M"

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone = "%s"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, shareName, env.OS_AVAILABILITY_ZONE)
}

func TestAccSFSTurboShareV1_withKMS(t *testing.T) {
	postfix := acctest.RandString(5)
	var turbo shares.Turbo
	rc := common.InitResourceCheck(
		sfsShareResourceName,
		&turbo,
		getSFSTurboShare,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboV1Crypt(postfix),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(sfsShareResourceName, "name", "sfs-turbo-"+postfix),
					resource.TestCheckResourceAttr(sfsShareResourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(sfsShareResourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(sfsShareResourceName, "size", "500"),
					resource.TestCheckResourceAttrSet(sfsShareResourceName, "crypt_key_id"),
				),
			},
		},
	})
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
