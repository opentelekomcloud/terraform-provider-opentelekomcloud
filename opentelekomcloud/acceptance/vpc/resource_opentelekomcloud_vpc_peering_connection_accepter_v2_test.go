package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcPeeringConnectionAcceptorV2_basic(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, multipleRouters(2))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcPeeringConnectionAcceptorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionAcceptorV2Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_peering_connection_accepter_v2.peer", "status", "ACTIVE",
					),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_peering_connection_accepter_v2.peer", "accept", "true",
					),
					resource.TestCheckResourceAttrSet(
						"opentelekomcloud_vpc_peering_connection_accepter_v2.peer", "vpc_tenant_id",
					),
					resource.TestCheckResourceAttrSet(
						"opentelekomcloud_vpc_peering_connection_accepter_v2.peer", "created_at",
					),
					resource.TestCheckResourceAttrSet(
						"opentelekomcloud_vpc_peering_connection_accepter_v2.peer", "updated_at",
					),
				),
			},
			{
				ResourceName:      "opentelekomcloud_vpc_peering_connection_accepter_v2.peer",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckVpcPeeringConnectionAcceptorDestroy(_ *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

const testAccVpcPeeringConnectionAcceptorV2Basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "otc_vpc_pa1"
  cidr = "192.168.0.0/16"
}
resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "otc_vpc_pa2"
  cidr = "192.168.0.0/16"
}
resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name        = "opentelekomcloud"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}
resource "opentelekomcloud_vpc_peering_connection_accepter_v2" "peer" {
  vpc_peering_connection_id = opentelekomcloud_vpc_peering_connection_v2.peering_1.id
  accept                    = true
}
`
