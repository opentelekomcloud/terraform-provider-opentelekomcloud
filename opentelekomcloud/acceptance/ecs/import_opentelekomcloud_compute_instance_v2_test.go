package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccComputeV2InstanceImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
			common.TestAccImagePreCheck(t)
			common.TestAccAzPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBasic,
			},
			{
				ResourceName:      resourceInstanceV2Name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
					"force_delete",
				},
			},
		},
	})
}

func TestAccComputeV2InstanceImportBootFromVolumeImage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
			common.TestAccImagePreCheck(t)
			common.TestAccAzPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBootFromVolumeImage,
			},
			{
				ResourceName:      resourceInstanceV2Name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
					"force_delete",
				},
			},
		},
	})
}
