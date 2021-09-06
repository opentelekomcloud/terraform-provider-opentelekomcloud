package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccEcsV1Instance_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
			common.TestAccImagePreCheck(t)
			common.TestAccAzPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1InstanceBasic,
			},

			{
				ResourceName:      resourceInstanceV1Name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}
