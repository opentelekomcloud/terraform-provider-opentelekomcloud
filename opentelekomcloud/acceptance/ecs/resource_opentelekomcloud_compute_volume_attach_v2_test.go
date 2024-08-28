package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/volumeattach"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/ecs"
)

const resourceVolumeAttach = "opentelekomcloud_compute_volume_attach_v2.va_1"

func TestAccComputeV2VolumeAttach_basic(t *testing.T) {
	var va volumeattach.VolumeAttachment
	qts := serverQuotas(4+1, env.OsFlavorID)
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2VolumeAttachBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists(resourceVolumeAttach, &va),
				),
			},
			{
				ResourceName:      resourceVolumeAttach,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeV2VolumeAttach_device(t *testing.T) {
	var va volumeattach.VolumeAttachment
	qts := serverQuotas(4+1, env.OsFlavorID)
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2VolumeAttachDevice,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists(resourceVolumeAttach, &va),
					resource.TestCheckResourceAttr(resourceVolumeAttach, "device", "/dev/vdc"),
				),
			},
		},
	})
}

func TestAccComputeV2VolumeAttach_MultipleDevices(t *testing.T) {
	var va volumeattach.VolumeAttachment
	name1 := "opentelekomcloud_compute_volume_attach_v2.va_1"
	name2 := "opentelekomcloud_compute_volume_attach_v2.va_2"
	name3 := "opentelekomcloud_compute_volume_attach_v2.va_3"
	qts := serverQuotas(4+1, env.OsFlavorID)
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2VolumeAttachMultipleDevices,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists(name1, &va),
					resource.TestCheckResourceAttr(name1, "device", "/dev/vdb"),
					testAccCheckComputeV2VolumeAttachExists(name2, &va),
					resource.TestCheckResourceAttr(name2, "device", "/dev/vdc"),
					testAccCheckComputeV2VolumeAttachExists(name3, &va),
					resource.TestCheckResourceAttr(name3, "device", "/dev/vdd"),
				),
			},
		},
	})
}

func TestAccComputeV2VolumeAttach_timeout(t *testing.T) {
	var va volumeattach.VolumeAttachment
	qts := serverQuotas(4+1, env.OsFlavorID)
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2VolumeAttachTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists(resourceVolumeAttach, &va),
				),
			},
		},
	})
}

func testAccCheckComputeV2VolumeAttachDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ComputeV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_volume_attach_v2" {
			continue
		}

		instanceId, volumeId, err := ecs.ParseComputeVolumeAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = volumeattach.Get(client, instanceId, volumeId).Extract()
		if err == nil {
			return fmt.Errorf("volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2VolumeAttachExists(n string, va *volumeattach.VolumeAttachment) resource.TestCheckFunc {
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

		instanceID, volumeId, err := ecs.ParseComputeVolumeAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := volumeattach.Get(client, instanceID, volumeId).Extract()
		if err != nil {
			return err
		}

		if found.ServerID != instanceID || found.VolumeID != volumeId {
			return fmt.Errorf("volumeAttach not found")
		}

		*va = *found

		return nil
	}
}

var testAccComputeV2VolumeAttachBasic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  image_name      = data.opentelekomcloud_images_image_v2.latest_image.name
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_compute_volume_attach_v2" "va_1" {
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
  volume_id   = opentelekomcloud_blockstorage_volume_v2.volume_1.id
}
`, common.DataSourceSubnet, common.DataSourceImage)

var testAccComputeV2VolumeAttachDevice = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  image_name      = data.opentelekomcloud_images_image_v2.latest_image.name
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_compute_volume_attach_v2" "va_1" {
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
  volume_id   = opentelekomcloud_blockstorage_volume_v2.volume_1.id
  device      = "/dev/vdc"
}
`, common.DataSourceSubnet, common.DataSourceImage)

var testAccComputeV2VolumeAttachTimeout = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  image_name      = data.opentelekomcloud_images_image_v2.latest_image.name
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_compute_volume_attach_v2" "va_1" {
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
  volume_id   = opentelekomcloud_blockstorage_volume_v2.volume_1.id

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, common.DataSourceSubnet, common.DataSourceImage)

var testAccComputeV2VolumeAttachMultipleDevices = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_ecs_instance_v1" "ecs" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  password                    = "Password@123"
  availability_zone           = "%s"
  auto_recovery               = true
  delete_disks_on_termination = true

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}


resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  availability_zone = opentelekomcloud_ecs_instance_v1.ecs.availability_zone
  name              = "first"
  volume_type       = "SSD"
  size              = "10"
}

resource "opentelekomcloud_evs_volume_v3" "volume_2" {
  availability_zone = opentelekomcloud_ecs_instance_v1.ecs.availability_zone
  name              = "second"
  volume_type       = "SSD"
  size              = "15"
}

resource "opentelekomcloud_evs_volume_v3" "volume_3" {
  availability_zone = opentelekomcloud_ecs_instance_v1.ecs.availability_zone
  name              = "third"
  volume_type       = "SAS"
  size              = "25"
}

resource "opentelekomcloud_compute_volume_attach_v2" "va_1" {
  instance_id = opentelekomcloud_ecs_instance_v1.ecs.id
  volume_id   = opentelekomcloud_evs_volume_v3.volume_1.id
  device      = "/dev/vdb"
}

resource "opentelekomcloud_compute_volume_attach_v2" "va_2" {
  instance_id = opentelekomcloud_ecs_instance_v1.ecs.id
  volume_id   = opentelekomcloud_evs_volume_v3.volume_2.id
  device      = "/dev/vdc"
}

resource "opentelekomcloud_compute_volume_attach_v2" "va_3" {
  instance_id = opentelekomcloud_ecs_instance_v1.ecs.id
  volume_id   = opentelekomcloud_evs_volume_v3.volume_3.id
  device      = "/dev/vdd"
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
