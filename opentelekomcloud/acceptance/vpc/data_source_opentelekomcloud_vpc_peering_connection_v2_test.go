package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const (
	dataSourceVpcPeeringNameID        = "data.opentelekomcloud_vpc_peering_connection_v2.by_id"
	dataSourceVpcPeeringNameVpcID     = "data.opentelekomcloud_vpc_peering_connection_v2.by_vpc_id"
	dataSourceVpcPeeringNamePeerVpcID = "data.opentelekomcloud_vpc_peering_connection_v2.by_peer_vpc_id"
)

func TestAccVpcPeeringConnectionV2DataSource_basic(t *testing.T) {
	t.Parallel()
	th.AssertNoErr(t, quotas.Router.AcquireMultiple(3))
	defer quotas.Router.ReleaseMultiple(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOTCVpcPeeringConnectionV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringConnectionV2DataSourceID(dataSourceVpcPeeringNameID),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringNameID, "name", "opentelekomcloud_peering_ds"),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringNameID, "status", "ACTIVE"),
					testAccCheckVpcPeeringConnectionV2DataSourceID(dataSourceVpcPeeringNameVpcID),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringNameVpcID, "name", "opentelekomcloud_peering_ds"),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringNameVpcID, "status", "ACTIVE"),
					testAccCheckVpcPeeringConnectionV2DataSourceID(dataSourceVpcPeeringNamePeerVpcID),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringNameVpcID, "name", "opentelekomcloud_peering_ds"),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringNameVpcID, "status", "ACTIVE"),
					testAccCheckVpcPeeringConnectionV2DataSourceID(dataSourceVpcPeeringNameVpcID),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringNameVpcID, "name", "opentelekomcloud_peering_ds"),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringNameVpcID, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckVpcPeeringConnectionV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find vpc peering connection data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("vpc peering connection data source ID not set")
		}

		return nil
	}
}

const testAccDataSourceOTCVpcPeeringConnectionV2Config = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test_ds_peer"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test_ds_peer1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_3" {
  name = "vpc_test_ds_peer_other"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name        = "opentelekomcloud_peering_ds"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_other" {
  name        = "opentelekomcloud_peering_ds_other"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_2.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_3.id
}

data "opentelekomcloud_vpc_peering_connection_v2" "by_id" {
  id = opentelekomcloud_vpc_peering_connection_v2.peering_1.id
}

data "opentelekomcloud_vpc_peering_connection_v2" "by_vpc_id" {
  vpc_id = opentelekomcloud_vpc_peering_connection_v2.peering_1.vpc_id
}

data "opentelekomcloud_vpc_peering_connection_v2" "by_peer_vpc_id" {
  peer_vpc_id = opentelekomcloud_vpc_peering_connection_v2.peering_1.peer_vpc_id
}

data "opentelekomcloud_vpc_peering_connection_v2" "by_vpc_ids" {
  vpc_id      = opentelekomcloud_vpc_peering_connection_v2.peering_1.vpc_id
  peer_vpc_id = opentelekomcloud_vpc_peering_connection_v2.peering_1.peer_vpc_id
}
`
