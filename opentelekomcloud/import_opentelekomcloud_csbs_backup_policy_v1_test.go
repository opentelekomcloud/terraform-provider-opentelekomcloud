package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCSBSBackupPolicyV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCSBSBackupPolicyV1_importWeekMonth(t *testing.T) {
	resourceName := "opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1_weekMonth,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"scheduled_operation.1795253542.day_backups",
					"scheduled_operation.1795253542.week_backups",
					"scheduled_operation.1795253542.month_backups",
					"scheduled_operation.1795253542.year_backups",
					"scheduled_operation.1795253542.timezone",
					"scheduled_operation.2772519308.day_backups",
					"scheduled_operation.2772519308.week_backups",
					"scheduled_operation.2772519308.month_backups",
					"scheduled_operation.2772519308.year_backups",
					"scheduled_operation.2772519308.timezone",
				},
			},
		},
	})
}
