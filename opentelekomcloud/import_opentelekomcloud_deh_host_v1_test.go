package opentelekomcloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccOTCDedicatedHostV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_deh_host_v1.deh1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCDeHV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDeHV1_basic,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
