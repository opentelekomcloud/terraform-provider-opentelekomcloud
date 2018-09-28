package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOTCDedicatedHostServerV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCDedicatedHostServerV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedHostServerV1DataSourceID("data.opentelekomcloud_deh_server_v1.servers"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_deh_server_v1.servers", "name", "ecs-instance-1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_deh_server_v1.servers", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckDedicatedHostServerV1DataSourceID(n string) resource.TestCheckFunc {
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

var testAccOTCDedicatedHostServerV1DataSource_basic = fmt.Sprintf(`
resource "opentelekomcloud_deh_host_v1" "deh1" {
	 availability_zone= "%s"
     auto_placement= "on"
     host_type= "s2"
	name = "deh-test-1"
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "ecs-instance-1"
  image_id = "%s"
  flavor_name = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  metadata {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
    scheduler_hints {
    tenancy = "dedicated"
    deh_id = "${opentelekomcloud_deh_host_v1.deh1.id}"
    }
}

data "opentelekomcloud_deh_server_v1" "servers" {
  dedicated_host_id  = "${opentelekomcloud_deh_host_v1.deh1.id}"
  server_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
}
`, OS_AVAILABILITY_ZONE, OS_IMAGE_ID, OS_FLAVOR_NAME, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
