package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/csbs/v1/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	ecs "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/ecs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourcePolicyName = "opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1"

func TestAccCSBSBackupPolicyV1_basic(t *testing.T) {
	var policy policies.BackupPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookMany(t, policyInstanceQuotas())
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists(resourcePolicyName, &policy),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "backup-policy"),
					resource.TestCheckResourceAttr(resourcePolicyName, "status", "suspended"),
				),
			},
			{
				Config: testAccCSBSBackupPolicyV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists(resourcePolicyName, &policy),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "backup-policy-update"),
				),
			},
		},
	})
}

func TestAccCSBSBackupPolicyV1_importWeekMonth(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookMany(t, policyInstanceQuotas())
		},
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

func TestAccCSBSBackupPolicyV1_timeout(t *testing.T) {
	var policy policies.BackupPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookMany(t, policyInstanceQuotas())
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists(resourcePolicyName, &policy),
				),
			},
		},
	})
}

func TestAccCSBSBackupPolicyV1_weekMonth(t *testing.T) {
	var policy policies.BackupPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookMany(t, policyInstanceQuotas())
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1WeekMonth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists(resourcePolicyName, &policy),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupPolicyV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.CsbsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating CSBSv1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_csbs_backup_policy_v1" {
			continue
		}

		_, err := policies.Get(client, rs.Primary.ID)
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
		client, err := config.CsbsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CSBSv1 client: %w", err)
		}

		found, err := policies.Get(client, rs.Primary.ID)
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

var testAccCSBSBackupPolicyV1Basic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "%s"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
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
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OsFlavorID)

var testAccCSBSBackupPolicyV1Update = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "%s"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
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
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OsFlavorID)

var testAccCSBSBackupPolicyV1Timeout = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "%s"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
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
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OsFlavorID)

var testAccCSBSBackupPolicyV1WeekMonth = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "%s"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
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
  tags = {
    muh = "kuh"
  }
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OsFlavorID)

func policyInstanceQuotas() quotas.MultipleQuotas {
	qts := ecs.QuotasForFlavor(env.OsFlavorID)
	qts = append(qts,
		&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
		&quotas.ExpectedQuota{Q: quotas.Volume, Count: 1},
		&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4},
		&quotas.ExpectedQuota{Q: quotas.CBRPolicy, Count: 1}, // the quota is shared with CBR service
	)
	return qts
}
