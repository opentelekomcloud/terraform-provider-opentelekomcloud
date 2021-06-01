package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOTCRtsSoftwareDeploymentV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_rts_software_deployment_v1.deployment_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckOTCRtsSoftwareDeploymentV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareDeploymentV1_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
