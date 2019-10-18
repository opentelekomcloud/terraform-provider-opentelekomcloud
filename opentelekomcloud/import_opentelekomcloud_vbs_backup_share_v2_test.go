package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVBSBackupShareV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_vbs_backup_share_v2.share"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccVBSBackupShareCheck(t) },
		Providers:    testAccProviders,
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
