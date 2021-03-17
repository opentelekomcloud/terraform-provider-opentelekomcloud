package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccFWFirewallV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_fw_firewall_group_v2.fw_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallGroupV2_basic_1,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
