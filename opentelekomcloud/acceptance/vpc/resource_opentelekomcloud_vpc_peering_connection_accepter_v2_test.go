package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcPeeringConnectionAccepterV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcPeeringConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccVpcPeeringConnectionAccepterV2_basic, // TODO: Research why normal scenario with peer tenant id is not working in acceptance tests
				ExpectError: regexp.MustCompile(`VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.`),
			},
		},
	})
}

func testAccCheckOTCVpcPeeringConnectionAccepterDestroy(_ *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

const testAccVpcPeeringConnectionAccepterV2_basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "otc_vpc_1"
  cidr = "192.168.0.0/16"
}
resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "otc_vpc_2"
  cidr = "192.168.0.0/16"
}
resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
    name = "opentelekomcloud"
    vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
    peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
  }
resource "opentelekomcloud_vpc_peering_connection_accepter_v2" "peer" {
  vpc_peering_connection_id = opentelekomcloud_vpc_peering_connection_v2.peering_1.id
  accept = true

}
`
