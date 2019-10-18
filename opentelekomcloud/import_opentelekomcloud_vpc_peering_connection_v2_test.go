package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOTCVpcPeeringConnectionV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_vpc_peering_connection_v2.peering_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOTCVpcPeeringConnectionV2_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
