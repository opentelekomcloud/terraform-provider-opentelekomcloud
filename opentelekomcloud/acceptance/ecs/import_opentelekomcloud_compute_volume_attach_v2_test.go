package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccComputeV2VolumeAttach_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
			common.TestAccImagePreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2VolumeAttachBasic,
			},

			{
				ResourceName:      resourceVolumeAttach,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
