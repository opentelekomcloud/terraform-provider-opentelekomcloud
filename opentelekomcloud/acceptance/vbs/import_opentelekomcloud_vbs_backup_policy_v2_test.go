package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVBSBackupPolicyV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_vbs_backup_policy_v2.vbs"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupPolicyV2_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
