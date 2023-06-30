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

const dataBackupName = "data.opentelekomcloud_csbs_backup_v1.csbs"

func TestAccCSBSBackupV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookMany(t, backupInstanceQuotas().X(2))
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupV1DataSourceID(dataBackupName),
					resource.TestCheckResourceAttr(dataBackupName, "backup_name", "csbs-test"),
					resource.TestCheckResourceAttr(dataBackupName, "resource_name", "instance_1"),
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

var testAccCSBSBackupV1DataSourceBasic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id
  security_groups   = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name   = "csbs-test"
  description   = "test-code"
  resource_id   = opentelekomcloud_compute_instance_v2.instance_1.id
  resource_type = "OS::Nova::Server"

  tags {
    key   = "kuh"
    value = "muh"
  }
}

data "opentelekomcloud_csbs_backup_v1" "csbs" {
  id = opentelekomcloud_csbs_backup_v1.csbs.id
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
