package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/secgroups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccComputeV2Instance_basic(t *testing.T) {
	var instance servers.Server
	resourceName := "opentelekomcloud_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "name", "instance_1"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", env.OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-create"),
				),
			},
			{
				Config: testAccComputeV2InstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", "instance_2"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-update"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_multiSecgroup(t *testing.T) {
	var instance servers.Server
	var secGroup1, secGroup2 secgroups.SecurityGroup
	resourceName := "opentelekomcloud_compute_instance_v2.instance_1"
	secGroupName1 := "opentelekomcloud_compute_secgroup_v2.secgroup_1"
	secGroupName2 := "opentelekomcloud_compute_secgroup_v2.secgroup_2"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceMultiSecgroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists(secGroupName1, &secGroup1),
					testAccCheckComputeV2SecGroupExists(secGroupName2, &secGroup2),
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
				),
			},
			{
				Config: testAccComputeV2InstanceMultiSecgroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists(secGroupName1, &secGroup1),
					testAccCheckComputeV2SecGroupExists(secGroupName2, &secGroup2),
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromImage(t *testing.T) {
	var instance servers.Server
	resourceName := "opentelekomcloud_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBootFromVolumeImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolume(t *testing.T) {
	var instance servers.Server
	resourceName := "opentelekomcloud_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBootFromVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_changeFixedIP(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceFixedIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"opentelekomcloud_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeVolume(t *testing.T) {
	var instance servers.Server
	resourceName := "opentelekomcloud_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBootFromVolumeVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_stopBeforeDestroy(t *testing.T) {
	var instance servers.Server
	resourceName := "opentelekomcloud_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceStopBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_metadata(t *testing.T) {
	var instance servers.Server
	resourceName := "opentelekomcloud_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceMetadata,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceMetadata(&instance, "abc", "def"),
					resource.TestCheckResourceAttr(resourceName, "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "all_metadata.abc", "def"),
				),
			},
			{
				Config: testAccComputeV2InstanceMetadataUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceMetadata(&instance, "ghi", "jkl"),
					testAccCheckComputeV2InstanceNoMetadataKey(&instance, "abc"),
					resource.TestCheckResourceAttr(resourceName, "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "all_metadata.ghi", "jkl"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_timeout(t *testing.T) {
	var instance servers.Server
	resourceName := "opentelekomcloud_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_autoRecovery(t *testing.T) {
	var instance servers.Server
	resourceName := "opentelekomcloud_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_recovery", "true"),
				),
			},
			{
				Config: testAccComputeV2InstanceAutoRecovery,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_recovery", "false"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_crazyNICs(t *testing.T) {
	var instance servers.Server
	resourceName := "opentelekomcloud_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceCrazyNICs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_initialStateActive(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceActive,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccComputeV2InstanceShutoff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "power_state", "shutoff"),
					testAccCheckComputeV2InstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: testAccComputeV2InstanceActive,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_initialStateShutoff(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      TestAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceShutoff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "power_state", "shutoff"),
					testAccCheckComputeV2InstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: testAccComputeV2InstanceActive,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccComputeV2InstanceShutoff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "power_state", "shutoff"),
					testAccCheckComputeV2InstanceState(&instance, "shutoff"),
				),
			},
		},
	})
}

func testAccCheckComputeV2InstanceExists(n string, instance *servers.Server) resource.TestCheckFunc {
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

		found, err := servers.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckComputeV2InstanceMetadata(instance *servers.Server, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return fmt.Errorf("no metadata")
		}
		for key, value := range instance.Metadata {
			if k != key {
				continue
			}
			if v == value {
				return nil
			}

			return fmt.Errorf("bad value for %s: %s", k, value)
		}

		return fmt.Errorf("metadata not found: %s", k)
	}
}

func testAccCheckComputeV2InstanceNoMetadataKey(instance *servers.Server, k string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return nil
		}

		for key := range instance.Metadata {
			if k == key {
				return fmt.Errorf("metadata found: %s", k)
			}
		}

		return nil
	}
}

func testAccCheckComputeV2InstanceBootVolumeAttachment(instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var attachments []volumeattach.VolumeAttachment

		config := common.TestAccProvider.Meta().(*cfg.Config)
		computeClient, err := config.ComputeV2Client(env.OS_REGION_NAME)
		if err != nil {
			return err
		}

		err = volumeattach.List(computeClient, instance.ID).EachPage(
			func(page pagination.Page) (bool, error) {
				actual, err := volumeattach.ExtractVolumeAttachments(page)
				if err != nil {
					return false, fmt.Errorf("unable to lookup attachment: %s", err)
				}

				attachments = actual
				return true, nil
			})
		if err != nil {
			return fmt.Errorf("error listing attached volumes: %w", err)
		}

		if len(attachments) == 1 {
			return nil
		}

		return fmt.Errorf("no attached volume found")
	}
}

var testAccComputeV2InstanceBasic = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID)

var testAccComputeV2InstanceUpdate = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_2"
  security_groups   = ["default"]
  availability_zone = "%s"

  network {
    uuid = "%s"
  }

  tags = {
    muh = "value-update"
  }
}
`, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID)

var testAccComputeV2InstanceMultiSecgroup = fmt.Sprintf(`
resource "opentelekomcloud_compute_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "a security group"
  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_compute_secgroup_v2" "secgroup_2" {
  name        = "secgroup_2"
  description = "another security group"
  rule {
    from_port   = 80
    to_port     = 80
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}
`, env.OS_NETWORK_ID)

var testAccComputeV2InstanceMultiSecgroupUpdate = fmt.Sprintf(`
resource "opentelekomcloud_compute_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "a security group"
  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_compute_secgroup_v2" "secgroup_2" {
  name        = "secgroup_2"
  description = "another security group"
  rule {
    from_port   = 80
    to_port     = 80
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = [
    "default",
    opentelekomcloud_compute_secgroup_v2.secgroup_1.name,
    opentelekomcloud_compute_secgroup_v2.secgroup_2.name,
  ]
  network {
    uuid = "%s"
  }
}
`, env.OS_NETWORK_ID)

var testAccComputeV2InstanceBootFromVolumeImage = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  security_groups   = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
  block_device {
    uuid                  = "%s"
    source_type           = "image"
    volume_size           = 50
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }
}
`, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID, env.OS_IMAGE_ID)

var testAccComputeV2InstanceBootFromVolumeVolume = fmt.Sprintf(`
resource "opentelekomcloud_blockstorage_volume_v2" "vol_1" {
  name     = "vol_1"
  size     = 50
  image_id = "%s"
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  block_device {
    uuid                  = opentelekomcloud_blockstorage_volume_v2.vol_1.id
    source_type           = "volume"
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }
}
`, env.OS_IMAGE_ID, env.OS_NETWORK_ID)

var testAccComputeV2InstanceBootFromVolume = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  block_device {
    uuid                  = "%s"
    source_type           = "image"
    volume_size           = 50
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }
}
`, env.OS_NETWORK_ID, env.OS_IMAGE_ID)

var testAccComputeV2InstanceFixedIP = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid        = "%s"
    fixed_ip_v4 = "192.168.0.24"
  }
}
`, env.OS_NETWORK_ID)

var testAccComputeV2InstanceStopBeforeDestroy = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  stop_before_destroy = true
}
`, env.OS_NETWORK_ID)

var testAccComputeV2InstanceMetadata = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  metadata = {
    foo = "bar"
    abc = "def"
  }
}
`, env.OS_NETWORK_ID)

var testAccComputeV2InstanceMetadataUpdate = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  metadata = {
    foo = "bar"
    ghi = "jkl"
  }
}
`, env.OS_NETWORK_ID)

var testAccComputeV2InstanceTimeout = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }

  timeouts {
    create = "10m"
  }
}
`, env.OS_NETWORK_ID)

var testAccComputeV2InstanceAutoRecovery = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  security_groups   = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
  auto_recovery = false
}
`, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID)

var testAccComputeV2InstanceCrazyNICs = fmt.Sprintf(`
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
}
resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name        = "subnet_1"
  network_id  = opentelekomcloud_networking_network_v2.network_1.id
  cidr        = "192.168.1.0/24"
  ip_version  = 4
  enable_dhcp = true
  no_gateway  = true
}
resource "opentelekomcloud_networking_network_v2" "network_2" {
  name = "network_2"
}
resource "opentelekomcloud_networking_subnet_v2" "subnet_2" {
  name        = "subnet_2"
  network_id  = opentelekomcloud_networking_network_v2.network_2.id
  cidr        = "192.168.2.0/24"
  ip_version  = 4
  enable_dhcp = true
  no_gateway  = true
}
resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  admin_state_up = "true"
  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.1.103"
  }
}
resource "opentelekomcloud_networking_port_v2" "port_2" {
  name           = "port_2"
  network_id     = opentelekomcloud_networking_network_v2.network_2.id
  admin_state_up = "true"
  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_2.id
    ip_address = "192.168.2.103"
  }
}
resource "opentelekomcloud_networking_port_v2" "port_3" {
  name           = "port_3"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  admin_state_up = "true"
  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.1.104"
  }
}
resource "opentelekomcloud_networking_port_v2" "port_4" {
  name           = "port_4"
  network_id     = opentelekomcloud_networking_network_v2.network_2.id
  admin_state_up = "true"
  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_2.id
    ip_address = "192.168.2.104"
  }
}
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  depends_on = [
    "opentelekomcloud_networking_subnet_v2.subnet_1",
    "opentelekomcloud_networking_subnet_v2.subnet_2",
    "opentelekomcloud_networking_port_v2.port_1",
    "opentelekomcloud_networking_port_v2.port_2",
  ]
  name            = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
  }
  network {
    uuid        = opentelekomcloud_networking_network_v2.network_1.id
    fixed_ip_v4 = "192.168.1.100"
  }
  network {
    uuid        = opentelekomcloud_networking_network_v2.network_2.id
    fixed_ip_v4 = "192.168.2.100"
  }
  network {
    uuid        = opentelekomcloud_networking_network_v2.network_1.id
    fixed_ip_v4 = "192.168.1.101"
  }
  network {
    uuid        = opentelekomcloud_networking_network_v2.network_2.id
    fixed_ip_v4 = "192.168.2.101"
  }
  network {
    port = opentelekomcloud_networking_port_v2.port_1.id
  }
  network {
    port = opentelekomcloud_networking_port_v2.port_2.id
  }
  network {
    port = opentelekomcloud_networking_port_v2.port_3.id
  }
  network {
    port = opentelekomcloud_networking_port_v2.port_4.id
  }
}
`, env.OS_NETWORK_ID)

var testAccComputeV2InstanceActive = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  power_state = "active"
  network {
    uuid = "%s"
  }
}
`, env.OS_NETWORK_ID)

var testAccComputeV2InstanceShutoff = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  power_state = "shutoff"
  network {
    uuid = "%s"
  }
}
`, env.OS_NETWORK_ID)

func testAccCheckComputeV2InstanceState(
	instance *servers.Server, state string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if strings.ToLower(instance.Status) != state {
			return fmt.Errorf("instance state is not match")
		}
		return nil
	}
}
