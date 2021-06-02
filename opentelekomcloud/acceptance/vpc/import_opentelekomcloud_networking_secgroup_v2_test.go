package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccNetworkingV2SecGroup_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_networking_secgroup_v2.secgroup_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroup_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
