package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/csbs/v1/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccCSBSBackupPolicyV1_basic(t *testing.T) {
	var policy policies.BackupPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", "name", "backup-policy"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", "status", "suspended"),
				),
			},
			{
				Config: testAccCSBSBackupPolicyV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", "name", "backup-policy-update"),
				),
			},
		},
	})
}

func TestAccCSBSBackupPolicyV1_timeout(t *testing.T) {
	var policy policies.BackupPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", &policy),
				),
			},
		},
	})
}

func TestAccCSBSBackupPolicyV1_weekMonth(t *testing.T) {
	var policy policies.BackupPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1_weekMonth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", &policy),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupPolicyV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	policyClient, err := config.CsbsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating csbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_csbs_backup_policy_v1" {
			continue
		}

		_, err := policies.Get(policyClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("backup policy still exists")
		}
	}

	return nil
}

func testAccCheckCSBSBackupPolicyV1Exists(n string, policy *policies.BackupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		policyClient, err := config.CsbsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CSBS client: %s", err)
		}

		found, err := policies.Get(policyClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("backup policy not found")
		}

		*policy = *found

		return nil
	}
}

var testAccCSBSBackupPolicyV1_basic = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  image_id          = "%s"
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "%s"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_csbs_backup_policy_v1" "backup_policy_v1" {
  name = "backup-policy"

  resource {
    id   = opentelekomcloud_compute_instance_v2.instance_1.id
    type = "OS::Nova::Server"
    name = "resource4"
  }
  scheduled_operation {
    name            = "mybackup"
    enabled         = true
    operation_type  = "backup"
    max_backups     = "2"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }
}
`, env.OS_IMAGE_ID, env.OS_AVAILABILITY_ZONE, env.OS_FLAVOR_ID, env.OS_NETWORK_ID)

var testAccCSBSBackupPolicyV1_update = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  image_id          = "%s"
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "%s"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}

resource "opentelekomcloud_csbs_backup_policy_v1" "backup_policy_v1" {
  name = "backup-policy-update"

  resource {
    id   = opentelekomcloud_compute_instance_v2.instance_1.id
    type = "OS::Nova::Server"
    name = "resource4"
  }
  scheduled_operation {
    name            = "mybackup"
    enabled         = true
    operation_type  = "backup"
    max_backups     = "2"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }
}
`, env.OS_IMAGE_ID, env.OS_AVAILABILITY_ZONE, env.OS_FLAVOR_ID, env.OS_NETWORK_ID)

var testAccCSBSBackupPolicyV1_timeout = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  image_id          = "%s"
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "%s"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_csbs_backup_policy_v1" "backup_policy_v1" {
  name = "backup-policy"

  resource {
    id   = opentelekomcloud_compute_instance_v2.instance_1.id
    type = "OS::Nova::Server"
    name = "resource4"
  }

  scheduled_operation {
    name            = "mybackup"
    enabled         = true
    operation_type  = "backup"
    max_backups     = "2"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, env.OS_IMAGE_ID, env.OS_AVAILABILITY_ZONE, env.OS_FLAVOR_ID, env.OS_NETWORK_ID)

var testAccCSBSBackupPolicyV1_weekMonth = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  image_id          = "%s"
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "%s"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_csbs_backup_policy_v1" "backup_policy_v1" {
  name = "backup-policy"

  resource {
    id   = opentelekomcloud_compute_instance_v2.instance_1.id
    type = "OS::Nova::Server"
    name = "resource1"
  }
  scheduled_operation {
    name            = "mybackup"
    enabled         = true
    operation_type  = "backup"
    max_backups     = "6"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
    week_backups    = "4"
    month_backups   = "2"
    timezone        = "UTC+03:00"
  }
}
`, env.OS_IMAGE_ID, env.OS_AVAILABILITY_ZONE, env.OS_FLAVOR_ID, env.OS_NETWORK_ID)
