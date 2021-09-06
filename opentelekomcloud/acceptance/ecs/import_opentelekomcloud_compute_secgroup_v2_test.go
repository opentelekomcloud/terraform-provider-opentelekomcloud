package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccComputeV2SecGroup_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupBasicOrig,
			},
			{
				ResourceName:      resourceSecGroupName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
