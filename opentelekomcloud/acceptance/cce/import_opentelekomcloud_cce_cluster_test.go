package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCEClusterV3_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3Basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cluster_version",
				},
			},
		},
	})
}
