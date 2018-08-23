package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOTCBMSServerV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckRequiredEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCBMSServerV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSServerV2DataSourceID("data.opentelekomcloud_compute_bms_server_v2.server1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_compute_bms_server_v2.server1", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckBMSServerV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find servers data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("server data source ID not set ")
		}

		return nil
	}
}

var testAccOTCBMSServerV2DataSource_basic = fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_otc123"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "sub_otc123"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  availability_zone = "%s"
}
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "BMSinstance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "physical.o2.medium"
  metadata {
    foo = "bar"
  }
  network {
    uuid = "${opentelekomcloud_vpc_subnet_v1.subnet_1.id}"
  }
}

data "opentelekomcloud_compute_bms_server_v2" "server1" {
  id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
}
`, OS_AVAILABILITY_ZONE, OS_IMAGE_ID, OS_AVAILABILITY_ZONE)
