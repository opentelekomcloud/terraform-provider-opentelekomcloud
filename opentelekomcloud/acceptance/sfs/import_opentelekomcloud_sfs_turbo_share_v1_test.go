package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccSFSTurboShareV1_importBasic(t *testing.T) {
	shareName := tools.RandomString("sfs-turbo-", 3)
	resourceName := "opentelekomcloud_sfs_turbo_share_v1.sfs-turbo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
			common.TestAccAzPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSFSTurboShareV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboShareV1Basic(shareName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
