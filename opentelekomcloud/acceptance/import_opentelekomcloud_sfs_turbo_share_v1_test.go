package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
)

func TestAccSFSTurboShareV1_importBasic(t *testing.T) {
	shareName := tools.RandomString("sfs-turbo-", 3)
	resourceName := "opentelekomcloud_sfs_turbo_share_v1.sfs-turbo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSTurboShareV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboShareV1_basic(shareName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
