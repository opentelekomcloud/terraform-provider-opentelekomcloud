package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	ecs "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/ecs"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/csbs/v1/backup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceBackupName = "opentelekomcloud_csbs_backup_v1.csbs"

func TestAccCSBSBackupV1_basic(t *testing.T) {
	var backups backup.Backup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookMany(t, backupInstanceQuotas())
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists(resourceBackupName, &backups),
					resource.TestCheckResourceAttr(resourceBackupName, "backup_name", "csbs-test1"),
					resource.TestCheckResourceAttr(resourceBackupName, "resource_type", "OS::Nova::Server"),
				),
			},
		},
	})
}

func TestAccCSBSBackupV1_importBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookMany(t, backupInstanceQuotas())
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1Basic,
			},
			{
				ResourceName:      resourceBackupName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccCSBSBackupV1_timeout(t *testing.T) {
	var backups backup.Backup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookMany(t, backupInstanceQuotas())
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists(resourceBackupName, &backups),
				),
			},
		},
	})
}

func testAccCSBSBackupV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.CsbsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating CSBSv1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_csbs_backup_v1" {
			continue
		}

		_, err := backup.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("backup still exists")
		}
	}

	return nil
}

func testAccCSBSBackupV1Exists(n string, backups *backup.Backup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("backup not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CsbsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CSBSv1 client: %s", err)
		}

		found, err := backup.Get(client, rs.Primary.ID)
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

var testAccCSBSBackupV1Basic = fmt.Sprintf(`
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
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name   = "csbs-test1"
  description   = "test-code"
  resource_id   = opentelekomcloud_compute_instance_v2.instance_1.id
  resource_type = "OS::Nova::Server"

  tags = {
    muh = "kuh"
  }
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OsFlavorID)

var testAccCSBSBackupV1Timeout = fmt.Sprintf(`
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
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name   = "csbs-test1"
  description   = "test-code"
  resource_id   = opentelekomcloud_compute_instance_v2.instance_1.id
  resource_type = "OS::Nova::Server"
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OsFlavorID)

func backupInstanceQuotas() quotas.MultipleQuotas {
	qts := ecs.QuotasForFlavor(env.OsFlavorID)
	qts = append(qts,
		&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
		&quotas.ExpectedQuota{Q: quotas.Volume, Count: 1},
		&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4},
	)
	return qts
}
