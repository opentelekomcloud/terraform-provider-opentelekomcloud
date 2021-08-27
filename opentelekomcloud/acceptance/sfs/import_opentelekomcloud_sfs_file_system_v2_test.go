package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOTCSFSFileSystemV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_sfs_file_system_v2.sfs_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2Basic,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_level", "access_to", "access_type"},
			},
		},
	})
}
