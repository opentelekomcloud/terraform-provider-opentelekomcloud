package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// PASS
func TestAccOTCVpcSubnetV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_vpc_subnet_v1.subnet_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCVpcSubnetV1_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
