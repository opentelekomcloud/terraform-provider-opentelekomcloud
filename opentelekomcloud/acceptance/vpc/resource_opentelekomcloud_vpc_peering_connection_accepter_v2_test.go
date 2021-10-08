package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcPeeringConnectionAcceptorV2_basic(t *testing.T) {
	t.Parallel()
	th.AssertNoErr(t, quotas.Router.AcquireMultiple(2))
	defer quotas.Router.ReleaseMultiple(2)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcPeeringConnectionAcceptorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccVpcPeeringConnectionAcceptorV2Basic, // TODO: Research why normal scenario with peer tenant id is not working in acceptance tests
				ExpectError: regexp.MustCompile(`VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.`),
			},
		},
	})
}

func testAccCheckOTCVpcPeeringConnectionAcceptorDestroy(_ *terraform.State) error {
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
