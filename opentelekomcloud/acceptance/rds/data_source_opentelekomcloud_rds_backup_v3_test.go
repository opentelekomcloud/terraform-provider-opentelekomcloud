package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/backups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/rds"
)

const dataSourceBackupName = "data.opentelekomcloud_rds_backup_v3.backup"

func TestAccDataSourceRDSV3Backup_basic(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3Basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
				),
			},
			{
				PreConfig: func() {
					forceRdsBackup(t, &rdsInstance.Id)
				},
				Config: testAccDataSourceRDSV3BackupBasic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(dataSourceBackupName, "databases.#", "0"),
					resource.TestCheckResourceAttr(dataSourceBackupName, "db_type", "postgresql"),
					resource.TestCheckResourceAttr(dataSourceBackupName, "db_version", "10"),
				),
			},
		},
	})
}

func testAccDataSourceRDSV3BackupBasic(postfix string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_rds_backup_v3" "backup" {
  instance_id = opentelekomcloud_rds_instance_v3.instance.id
  type        = "auto"
}
`, testAccRdsInstanceV3Basic(postfix))
}

func forceRdsBackup(t *testing.T, instanceID *string) {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.RdsV3Client(env.OS_REGION_NAME)
	th.AssertNoErr(t, err)

	// make sure RDS instance is not doing something else
	err = golangsdk.WaitFor(1200, func() (bool, error) {
		rdsInstance, err := rds.GetRdsInstance(client, *instanceID)
		if err != nil {
			return false, fmt.Errorf("error fetching RDS instance: %s", err)
		}
		if rdsInstance == nil {
			return false, fmt.Errorf("RDS instance %s is missing", err)
		}
		if rdsInstance.Status == "ACTIVE" {
			return true, nil
		}
		return false, nil
	})
	th.AssertNoErr(t, err)

	err = golangsdk.WaitFor(600, func() (bool, error) {
		bList, err := backups.List(client, backups.ListOpts{
			InstanceID: *instanceID,
		})
		if err != nil {
			return false, err
		}
		if len(bList) == 0 {
			return false, nil
		}
		backup := bList[0]
		if backup.Status == backups.StatusCompleted {
			return true, nil
		}
		if backup.Status == backups.StatusFailed {
			return false, fmt.Errorf("backup for instance %s failed", *instanceID)
		}
		return false, nil
	})
	th.AssertNoErr(t, err)
}
