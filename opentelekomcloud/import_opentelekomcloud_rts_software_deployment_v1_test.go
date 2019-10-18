package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOTCRtsSoftwareDeploymentV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_rts_software_deployment_v1.deployment_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
