package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccRTSSoftwareDeploymentV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_rts_software_deployment_v1.deployment_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccImagePreCheck(t)
			common.TestAccSubnetPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRTSSoftwareDeploymentV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareDeploymentV1Basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
