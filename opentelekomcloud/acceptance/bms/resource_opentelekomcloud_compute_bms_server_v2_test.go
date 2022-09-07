package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/bms/v2/servers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	ecs "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/ecs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceServerName = "opentelekomcloud_compute_bms_server_v2.instance_1"

func TestAccComputeV2BmsInstance_basic(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2BmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2BmsInstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2BmsInstanceExists(resourceServerName, &instance),
					resource.TestCheckResourceAttr(resourceServerName, "availability_zone", env.OS_AVAILABILITY_ZONE),
				),
			},
			{
				Config: testAccComputeV2BmsInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2BmsInstanceExists(resourceServerName, &instance),
					resource.TestCheckResourceAttr(resourceServerName, "name", "instance_2"),
				),
			},
		},
	})
}

func TestAccComputeV2BmsInstance_bootFromVolumeImage(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2BmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2BmsInstanceBootFromVolumeImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2BmsInstanceExists(resourceServerName, &instance),
					resource.TestCheckResourceAttr(resourceServerName, "name", "instance_1"),
				),
			},
		},
	})
}

func TestAccComputeV2BmsInstance_timeout(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      ecs.TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2BmsInstanceTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2BmsInstanceExists(resourceServerName, &instance),
				),
			},
		},
	})
}

func testAccCheckComputeV2BmsInstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ComputeV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_bms_server_v2" {
			continue
		}

		server, err := servers.Get(client, rs.Primary.ID)
		if err == nil {
			if server.Status != "SOFT_DELETED" {
				return fmt.Errorf("instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeV2BmsInstanceExists(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ComputeV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
		}

		found, err := servers.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("bms Instance not found")
		}

		*instance = *found

		return nil
	}
}

var testAccComputeV2BmsInstanceBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_bms_server_v2" "instance_1" {
  name              = "instance_1"
  flavor_id         = "physical.o2.medium"
  flavor_name       = "physical.o2.medium"
  security_groups   = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccComputeV2BmsInstanceUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_bms_server_v2" "instance_1" {
  name              = "instance_2"
  flavor_id         = "physical.o2.medium"
  flavor_name       = "physical.o2.medium"
  security_groups   = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccComputeV2BmsInstanceTimeout = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_bms_server_v2" "instance_1" {
  name              = "instance_1"
  flavor_id         = "physical.o2.medium"
  flavor_name       = "physical.o2.medium"
  security_groups   = ["default"]
  availability_zone = "%s"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  timeouts {
    create = "20m"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccComputeV2BmsInstanceBootFromVolumeImage = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_bms_server_v2" "instance_1" {
  name              = "instance_1"
  flavor_id         = "physical.h2.large"
  flavor_name       = "physical.h2.large"
  security_groups   = ["default"]
  availability_zone = "%s"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  block_device {
    uuid                  = "d50b4060-92cc-4d38-ae88-bd91bc3df00f"
    source_type           = "image"
    volume_type           = "SATA"
    volume_size           = 100
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
    device_name           = "/dev/sda"
  }
  timeouts {
    create = "20m"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
