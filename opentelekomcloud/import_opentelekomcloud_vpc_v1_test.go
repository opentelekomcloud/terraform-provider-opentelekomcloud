package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// PASS
func TestAccNetworkingV1Vpc_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_vpc_v1.networking_vpc_v1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV1VpcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV1Vpc_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
