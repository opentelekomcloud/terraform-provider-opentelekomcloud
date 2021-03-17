package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOTCVpcSubnetV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_vpc_subnet_v1.subnet_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckOTCVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOTCVpcSubnetV1_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
