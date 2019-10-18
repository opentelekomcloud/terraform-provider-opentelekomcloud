package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"testing"

	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/shares"
)

func TestAccVBSBackupShareV2_basic(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccVBSBackupShareCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVBSBackupShareV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupShareV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupShareV2Exists("opentelekomcloud_vbs_backup_share_v2.share", &share),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_share_v2.share", "to_project_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccVBSBackupShareV2_timeout(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccVBSBackupShareCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVBSBackupShareV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupShareV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupShareV2Exists("opentelekomcloud_vbs_backup_share_v2.share", &share),
				),
			},
		},
	})
}

func testAccVBSBackupShareV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vbsClient, err := config.vbsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating opentelekomcloud vbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vbs_backup_share_v2" {
			continue
		}

		_, err := shares.List(vbsClient, shares.ListOpts{BackupID: rs.Primary.ID})
		if err != nil {
			return fmt.Errorf("Backup share still exists")
		}
	}

	return nil
}

func testAccVBSBackupShareV2Exists(n string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vbsClient, err := config.vbsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating opentelekomcloud vbs client: %s", err)
		}

		shareList, err := shares.List(vbsClient, shares.ListOpts{BackupID: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := shareList[0]

		*share = found

		return nil
	}
}

var testAccVBSBackupShareV2_basic = fmt.Sprintf(`
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_123"
  description = "first test volume"
  size = 40
  cascade = true
}

resource "opentelekomcloud_vbs_backup_v2" "backup_1" {
  volume_id = "${opentelekomcloud_blockstorage_volume_v2.volume_1.id}"
  name = "vbs-backup"
  description = "Backup_Demo"
}

resource "opentelekomcloud_vbs_backup_share_v2" "share" {
  backup_id ="${opentelekomcloud_vbs_backup_v2.backup_1.id}"
  to_project_ids = ["%s"]
}
`, OS_TO_TENANT_ID)

var testAccVBSBackupShareV2_timeout = fmt.Sprintf(`
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_123"
  description = "first test volume"
  size = 40
  cascade = true
}

resource "opentelekomcloud_vbs_backup_v2" "backup_1" {
  volume_id = "${opentelekomcloud_blockstorage_volume_v2.volume_1.id}"
  name = "vbs-backup"
  description = "Backup_Demo"
}

resource "opentelekomcloud_vbs_backup_share_v2" "share" {
  backup_id ="${opentelekomcloud_vbs_backup_v2.backup_1.id}"
  to_project_ids = ["%s"]

timeouts {
    create = "5m"
    delete = "5m"
  }

}
`, OS_TO_TENANT_ID)
