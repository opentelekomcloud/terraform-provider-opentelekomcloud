package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/csbs/v1/backup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccCSBSBackupV1_basic(t *testing.T) {
	var backups backup.Backup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("opentelekomcloud_csbs_backup_v1.csbs", &backups),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_v1.csbs", "backup_name", "csbs-test1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_v1.csbs", "resource_type", "OS::Nova::Server"),
				),
			},
		},
	})
}

func TestAccCSBSBackupV1_timeout(t *testing.T) {
	var backups backup.Backup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("opentelekomcloud_csbs_backup_v1.csbs", &backups),
				),
			},
		},
	})
}

func testAccCSBSBackupV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	backupClient, err := config.CsbsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating csbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_csbs_backup_v1" {
			continue
		}

		_, err := backup.Get(backupClient, rs.Primary.ID).ExtractBackup()
		if err == nil {
			return fmt.Errorf("Backup still exists")
		}
	}

	return nil
}

func testAccCSBSBackupV1Exists(n string, backups *backup.Backup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Backup not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		backupClient, err := config.CsbsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating csbs client: %s", err)
		}

		found, err := backup.Get(backupClient, rs.Primary.ID).ExtractBackup()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("backup not found")
		}

		*backups = *found

		return nil
	}
}

var testAccCSBSBackupV1_basic = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test1"
  description      = "test-code"
  resource_id = opentelekomcloud_compute_instance_v2.instance_1.id
  resource_type = "OS::Nova::Server"
}
`, env.OS_IMAGE_ID, env.OS_AVAILABILITY_ZONE, env.OS_FLAVOR_ID, env.OS_NETWORK_ID)

var testAccCSBSBackupV1_timeout = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test1"
  description      = "test-code"
  resource_id = opentelekomcloud_compute_instance_v2.instance_1.id
  resource_type = "OS::Nova::Server"
}
`, env.OS_IMAGE_ID, env.OS_AVAILABILITY_ZONE, env.OS_FLAVOR_ID, env.OS_NETWORK_ID)
