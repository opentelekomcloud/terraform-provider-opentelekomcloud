package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccNatDnat_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
			common.TestAccAzPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatDnatBasic,
			},

			{
				ResourceName:      resourceDnatRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
