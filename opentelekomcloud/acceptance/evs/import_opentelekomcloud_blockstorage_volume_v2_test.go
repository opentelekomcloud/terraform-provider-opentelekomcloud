package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccBlockStorageV2Volume_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2VolumeBasic,
			},
			{
				ResourceName:      resourceVolumeV2Name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cascade",
				},
			},
		},
	})
}
