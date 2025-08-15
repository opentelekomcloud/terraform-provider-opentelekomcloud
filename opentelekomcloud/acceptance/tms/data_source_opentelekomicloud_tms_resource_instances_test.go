package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccTmsResourceInstancesDS_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_tms_resource_tags_v1.tags.ds_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTagsDS_basic(tools.RandomString("tms-rt-acc-test-", 4)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTMSResourceInstancesV1DsId(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "resources.0.resource_type", "disk"),
				),
			},
		},
	})
}

func testAccCheckTMSResourceInstancesV1DsId(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TMS resource instances v1 data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("TMS resource instances v1 data source ID not set")
		}

		return nil
	}
}

func testAccResourceTagsDS_basic(volName string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_identity_project_v3" "project_1" {}

resource "opentelekomcloud_evs_volume_v3" "volume" {
  name              = "%s"
  description       = "tms test volume"
  availability_zone = "eu-de-01"
  volume_type       = "SSD"
  size              = 20

  lifecycle {
    ignore_changes = [
      tags
    ]
  }

}

resource "opentelekomcloud_tms_resource_tags_v1" "tags_1" {
  depends_on = ["opentelekomcloud_evs_volume_v3.volume"]
  project_id = data.opentelekomcloud_identity_project_v3.project_1.id

  resources {
    resource_type = "disk"
    resource_id   = opentelekomcloud_evs_volume_v3.volume.id
  }

  tags = {
    test = "test-tf-acc"
  }
}

data "opentelekomcloud_tms_resource_tags_v1" "tags_ds_1" {
  depends_on = ["opentelekomcloud_tms_resource_tags_v1.tags_1"]
  resources  = ["disk", "ecs"]
  tags {
    key    = "test"
    values = ["test-tf-acc"]
  }
  project_id = data.opentelekomcloud_identity_project_v3.project_1.id
}
`, volName)
}
