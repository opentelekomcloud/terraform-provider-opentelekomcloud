package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccCSBSBackupV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupV1DataSourceID("data.opentelekomcloud_csbs_backup_v1.csbs"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_csbs_backup_v1.csbs", "backup_name", "csbs-test"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_csbs_backup_v1.csbs", "resource_name", "instance_1"),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find backup data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("backup data source ID not set ")
		}

		return nil
	}
}

var testAccCSBSBackupV1DataSource_basic = fmt.Sprintf(`
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
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name   = "csbs-test"
  description   = "test-code"
  resource_id   = opentelekomcloud_compute_instance_v2.instance_1.id
  resource_type = "OS::Nova::Server"
}
data "opentelekomcloud_csbs_backup_v1" "csbs" {
  id = opentelekomcloud_csbs_backup_v1.csbs.id
}
`, env.OsImageID, env.OsAvailabilityZone, env.OsFlavorID, env.OsNetworkID)
