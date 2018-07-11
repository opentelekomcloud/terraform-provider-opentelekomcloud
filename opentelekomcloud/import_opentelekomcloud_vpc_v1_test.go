package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// PASS
func TestAccOTCVpcV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_vpc_v1.vpc_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCVpcV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcV1_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
