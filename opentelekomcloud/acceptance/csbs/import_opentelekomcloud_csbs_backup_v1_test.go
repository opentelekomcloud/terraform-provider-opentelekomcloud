package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCSBSBackupV1_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
			common.TestAccImagePreCheck(t)
			common.TestAccAzPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1Basic,
			},
			{
				ResourceName:      resourceBackupName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}
