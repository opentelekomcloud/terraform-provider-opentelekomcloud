package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vbs/v2/backups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccVBSBackupV2_basic(t *testing.T) {
	var config backups.Backup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVBSBackupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVBSBackupV2Exists("opentelekomcloud_vbs_backup_v2.backup_1", &config),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_v2.backup_1", "name", "vbs-backup"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_v2.backup_1", "description", "Backup_Demo"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_v2.backup_1", "status", "available"),
				),
			},
		},
	})
}

func TestAccVBSBackupV2_timeout(t *testing.T) {
	var config backups.Backup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVBSBackupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVBSBackupV2Exists("opentelekomcloud_vbs_backup_v2.backup_1", &config),
				),
			},
		},
	})
}

func testAccCheckVBSBackupV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	vbsClient, err := config.VbsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud vbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vbs_backup_v2" {
			continue
		}

		_, err := backups.Get(vbsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VBS backup still exists")
		}
	}

	return nil
}

func testAccCheckVBSBackupV2Exists(n string, configs *backups.Backup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		vbsClient, err := config.VbsV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud vbs client: %s", err)
		}

		found, err := backups.Get(vbsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("VBS backup not found")
		}

		*configs = *found

		return nil
	}
}

const testAccVBSBackupV2_basic = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name        = "volume_123"
  description = "first test volume"
  size        = 40
  cascade     = true
}

resource "opentelekomcloud_vbs_backup_v2" "backup_1" {
  volume_id   = opentelekomcloud_blockstorage_volume_v2.volume_1.id
  name        = "vbs-backup"
  description = "Backup_Demo"
}
`

const testAccVBSBackupV2_timeout = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name        = "volume_123"
  description = "first test volume"
  size        = 40
  cascade     = true
}

resource "opentelekomcloud_vbs_backup_v2" "backup_1" {
  volume_id   = opentelekomcloud_blockstorage_volume_v2.volume_1.id
  name        = "vbs-backup"
  description = "Backup_Demo"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
