package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcPeeringConnectionV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOTCVpcPeeringConnectionV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringConnectionV2DataSourceID("data.opentelekomcloud_vpc_peering_connection_v2.by_id"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_peering_connection_v2.by_id", "name", "opentelekomcloud_peering"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_peering_connection_v2.by_id", "status", "ACTIVE"),
					testAccCheckVpcPeeringConnectionV2DataSourceID("data.opentelekomcloud_vpc_peering_connection_v2.by_vpc_id"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_peering_connection_v2.by_vpc_id", "name", "opentelekomcloud_peering"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_peering_connection_v2.by_vpc_id", "status", "ACTIVE"),
					testAccCheckVpcPeeringConnectionV2DataSourceID("data.opentelekomcloud_vpc_peering_connection_v2.by_peer_vpc_id"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_peering_connection_v2.by_peer_vpc_id", "name", "opentelekomcloud_peering"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_peering_connection_v2.by_peer_vpc_id", "status", "ACTIVE"),
					testAccCheckVpcPeeringConnectionV2DataSourceID("data.opentelekomcloud_vpc_peering_connection_v2.by_vpc_ids"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_peering_connection_v2.by_vpc_ids", "name", "opentelekomcloud_peering"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_peering_connection_v2.by_vpc_ids", "status", "ACTIVE"),
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
		name = "vpc_test"
		cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
		name = "vpc_test1"
        cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
		name = "opentelekomcloud_peering"
		vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
		peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
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
		vpc_id = opentelekomcloud_vpc_peering_connection_v2.peering_1.vpc_id
		peer_vpc_id = opentelekomcloud_vpc_peering_connection_v2.peering_1.peer_vpc_id
}
`
