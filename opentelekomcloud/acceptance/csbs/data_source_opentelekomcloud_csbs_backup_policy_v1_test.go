package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataPolicyName = "data.opentelekomcloud_csbs_backup_policy_v1.csbs_policy"

func TestAccCSBSBackupPolicyV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookMany(t, policyInstanceQuotas().X(2))
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1DataSourceID(dataPolicyName),
					resource.TestCheckResourceAttr(dataPolicyName, "name", "backup-policy"),
					resource.TestCheckResourceAttr(dataPolicyName, "status", "suspended"),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupPolicyV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find backup data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("backup data source ID not set")
		}

		return nil
	}
}

var testAccCSBSBackupPolicyV1DataSourceBasic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id
  security_groups   = ["default"]
  availability_zone = "%s"
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
  tags {
    key   = "kuh"
    value = "muh"
  }
}

data "opentelekomcloud_csbs_backup_policy_v1" "csbs_policy" {
  id = opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1.id
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
