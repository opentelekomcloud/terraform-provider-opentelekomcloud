package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataServerName = "data.opentelekomcloud_deh_server_v1.servers"

func TestAccDedicatedHostServerV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedHostServerV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedHostServerV1DataSourceID(dataServerName),
					resource.TestCheckResourceAttr(dataServerName, "name", "ecs-instance-1"),
					resource.TestCheckResourceAttr(dataServerName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckDedicatedHostServerV1DataSourceID(n string) resource.TestCheckFunc {
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

var testAccDedicatedHostServerV1DataSourceBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_deh_host_v1" "deh1" {
  availability_zone = "%[2]s"
  auto_placement    = "on"
  host_type         = "s2-medium"
  name              = "deh-test-1"
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "ecs-instance-1"
  security_groups   = ["default"]
  flavor_name       = "s2.medium.1"
  availability_zone = "%[2]s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  scheduler_hints {
    tenancy = "dedicated"
    deh_id  = opentelekomcloud_deh_host_v1.deh1.id
  }
}

data "opentelekomcloud_deh_server_v1" "servers" {
  dedicated_host_id = opentelekomcloud_deh_host_v1.deh1.id
  server_id         = opentelekomcloud_compute_instance_v2.instance_1.id
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
