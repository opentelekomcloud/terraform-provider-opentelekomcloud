package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs/v2/shares"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccSFSShareAccessRuleV2_basic(t *testing.T) {
	resourceName := "opentelekomcloud_sfs_share_access_rule_v2.sfs_rules"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckSFSShareAccessRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareAccessRuleV2_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "", "sfs-test1"),
				),
			},
			{
				Config: testAccSFSShareAccessRuleV2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "sfs-test2"),
				),
			},
		},
	})
}

func testAccCheckSFSShareAccessRuleV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SfsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud SFSv2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_sfs_share_access_rule_v2" {
			continue
		}

		_, err := shares.ListAccessRights(client, rs.Primary.ID).ExtractAccessRights()
		if err == nil {
			return fmt.Errorf("share file rules still exists")
		}
	}

	return nil
}

var testAccSFSShareAccessRuleV2_basic = fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name   = "sfs_share_vpc_1"
  cidr   = "192.168.0.0/16"
}
resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name   = "sfs_share_vpc_2"
  cidr   = "192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
  share_proto       = "NFS"
  size              = 1
  name              = "sfs-test1"
  availability_zone = "eu-de-01"
  access_to         = "%s"
  access_type       = "cert"
  access_level      = "rw"
}
resource "opentelekomcloud_sfs_access_rule_v2" "sfs_rules" {
  share_id = opentelekomcloud_sfs_file_system_v2.sfs_1.id
  access_rules {
    access_to    = opentelekomcloud_vpc_v1.vpc_1.id
    access_type  = "cert"
    access_level = "rw"
  }
  access_rules {
    access_to    = opentelekomcloud_vpc_v1.vpc_2.id
    access_type  = "cert"
    access_level = "rw"
  }
}
`, env.OS_VPC_ID)

var testAccSFSShareAccessRuleV2_update = fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name   = "sfs_share_vpc_1"
  cidr   = "192.168.0.0/16"
}
resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name   = "sfs_share_vpc_2"
  cidr   = "192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
  share_proto       = "NFS"
  size              = 1
  name              = "sfs-test1"
  availability_zone = "eu-de-01"
  access_to         = "%s"
  access_type       = "cert"
  access_level      = "rw"
}
resource "opentelekomcloud_sfs_access_rule_v2" "sfs_rules" {
  share_id = opentelekomcloud_sfs_file_system_v2.sfs_1.id
  access_rules {
    access_to    = opentelekomcloud_vpc_v1.vpc_1.id
    access_type  = "cert"
    access_level = "rw"
  }
}
`, env.OS_VPC_ID)
