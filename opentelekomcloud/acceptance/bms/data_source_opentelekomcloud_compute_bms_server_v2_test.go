package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataServerName = "data.opentelekomcloud_compute_bms_server_v2.server1"

func TestAccBMSServerV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSServerV2DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSServerV2DataSourceID(dataServerName),
					resource.TestCheckResourceAttr(dataServerName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckBMSServerV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find servers data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("server data source ID not set ")
		}

		return nil
	}
}

var testAccBMSServerV2DataSourceBasic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "BMSinstance_1"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "physical.o2.medium"
  flavor_name       = "physical.o2.medium"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

data "opentelekomcloud_compute_bms_server_v2" "server1" {
  id = opentelekomcloud_compute_instance_v2.instance_1.id
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
