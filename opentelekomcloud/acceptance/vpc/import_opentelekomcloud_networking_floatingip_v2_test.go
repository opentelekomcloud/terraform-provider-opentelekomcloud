package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccNetworkingV2FloatingIP_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_networking_floatingip_v2.fip_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIP_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
