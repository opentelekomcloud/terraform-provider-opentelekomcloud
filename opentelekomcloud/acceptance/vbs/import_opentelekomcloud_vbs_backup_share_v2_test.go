package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVBSBackupShareV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_vbs_backup_share_v2.share"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccVBSBackupShareCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccVBSBackupShareV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupShareV2_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
