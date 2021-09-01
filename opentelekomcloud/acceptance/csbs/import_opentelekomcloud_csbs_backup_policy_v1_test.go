package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCSBSBackupPolicyV1_importWeekMonth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1WeekMonth,
			},
			{
				ResourceName:      resourcePolicyName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"scheduled_operation",
				},
			},
		},
	})
}
