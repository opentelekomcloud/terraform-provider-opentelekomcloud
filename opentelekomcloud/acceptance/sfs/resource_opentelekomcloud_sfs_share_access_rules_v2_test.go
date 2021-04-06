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

func TestAccSFSShareAccessRulesV2_basic(t *testing.T) {
	resourceName := "opentelekomcloud_sfs_share_access_rules_v2.sfs_rules"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckSFSShareAccessRulesV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareAccessRulesV2_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_rule.#", "2"),
				),
			},
			{
				Config: testAccSFSShareAccessRulesV2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_rule.#", "1"),
				),
			},
		},
	})
}

func testAccCheckSFSShareAccessRulesV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SfsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud SFSv2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_sfs_share_access_rules_v2" {
			continue
		}

		_, err := shares.ListAccessRights(client, rs.Primary.ID).ExtractAccessRights()
		if err == nil {
			return fmt.Errorf("share file rules still exists")
		}
	}

	return nil
}

var testAccSFSShareAccessRulesV2_basic = fmt.Sprintf(`
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
}

resource "opentelekomcloud_sfs_share_access_rules_v2" "sfs_rules" {
  share_id = opentelekomcloud_sfs_file_system_v2.sfs_1.id

  access_rule {
    access_to    = opentelekomcloud_vpc_v1.vpc_1.id
    access_type  = "cert"
    access_level = "rw"
  }

  access_rule {
    access_to    = opentelekomcloud_vpc_v1.vpc_2.id
    access_type  = "cert"
    access_level = "rw"
  }
}
`)

var testAccSFSShareAccessRulesV2_update = fmt.Sprintf(`
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
}

resource "opentelekomcloud_sfs_share_access_rules_v2" "sfs_rules" {
  share_id = opentelekomcloud_sfs_file_system_v2.sfs_1.id

  access_rule {
    access_to    = opentelekomcloud_vpc_v1.vpc_1.id
    access_type  = "cert"
    access_level = "rw"
  }
}
`)
