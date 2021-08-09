package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVBSBackupV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVBSBackupV2DataSourceID("data.opentelekomcloud_vbs_backup_v2.backups"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_vbs_backup_v2.backups", "name", "vbs-backup"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_vbs_backup_v2.backups", "description", "Backup_Demo"),
				),
			},
		},
	})
}

func testAccCheckVBSBackupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find backup data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("VBS backup data source ID not set ")
		}

		return nil
	}
}

var testAccVBSBackupV2DataSource_basic = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_123"
  description = "first test volume"
  size = 40
  cascade = true
}

resource "opentelekomcloud_vbs_backup_v2" "backup_1" {
  volume_id = opentelekomcloud_blockstorage_volume_v2.volume_1.id
  name = "vbs-backup"
  description = "Backup_Demo"
}

data "opentelekomcloud_vbs_backup_v2" "backups" {
  id = opentelekomcloud_vbs_backup_v2.backup_1.id
}
`
