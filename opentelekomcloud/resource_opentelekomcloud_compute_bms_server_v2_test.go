package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/bms/v2/servers"
)

func TestAccComputeV2BmsInstance_basic(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccBmsFlavorPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2BmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2BmsInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2BmsInstanceExists("opentelekomcloud_compute_bms_server_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_bms_server_v2.instance_1", "availability_zone", OS_AVAILABILITY_ZONE),
				),
			},
			{
				Config: testAccComputeV2BmsInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2BmsInstanceExists("opentelekomcloud_compute_bms_server_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_bms_server_v2.instance_1", "name", "instance_2"),
				),
			},
		},
	})
}

func TestAccComputeV2BmsInstance_bootFromVolumeImage(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccBmsFlavorPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2BmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2BmsInstance_bootFromVolumeImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2BmsInstanceExists("opentelekomcloud_compute_bms_server_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_bms_server_v2.instance_1", "name", "instance_1"),
				),
			},
		},
	})
}

func TestAccComputeV2BmsInstance_timeout(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccBmsFlavorPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2BmsInstance_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2BmsInstanceExists("opentelekomcloud_compute_bms_server_v2.instance_1", &instance),
				),
			},
		},
	})
}

func testAccCheckComputeV2BmsInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2HWClient(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_bms_server_v2" {
			continue
		}

		server, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "SOFT_DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeV2BmsInstanceExists(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2HWClient(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud compute client: %s", err)
		}

		found, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Bms Instance not found")
		}

		*instance = *found

		return nil
	}
}

var testAccComputeV2BmsInstance_basic = fmt.Sprintf(`
resource "opentelekomcloud_compute_bms_server_v2" "instance_1" {
  name = "instance_1"
  flavor_id = "physical.o2.medium"
  flavor_name = "physical.o2.medium"
  security_groups = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2BmsInstance_update = fmt.Sprintf(`
resource "opentelekomcloud_compute_bms_server_v2" "instance_1" {
  name = "instance_2"
  flavor_id = "physical.o2.medium"
  flavor_name = "physical.o2.medium"
  security_groups = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2BmsInstance_timeout = fmt.Sprintf(`
resource "opentelekomcloud_compute_bms_server_v2" "instance_1" {
  name = "instance_1"
  flavor_id = "physical.o2.medium"
  flavor_name = "physical.o2.medium"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }

  timeouts {
    create = "20m"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2BmsInstance_bootFromVolumeImage = fmt.Sprintf(`
resource "opentelekomcloud_compute_bms_server_v2" "instance_1" {
  name = "instance_1"
  flavor_id = "physical.h2.large"
  flavor_name = "physical.h2.large"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }

  block_device {
	uuid = "d50b4060-92cc-4d38-ae88-bd91bc3df00f"
	source_type = "image"
	volume_type = "SATA"
	volume_size = 100
	boot_index = 0
	destination_type = "volume"
	delete_on_termination = true
	device_name = "/dev/sda"
  }
  timeouts {
    create = "20m"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
