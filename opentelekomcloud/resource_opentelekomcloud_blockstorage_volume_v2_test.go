package opentelekomcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/blockstorage/v2/volumes"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v2/tags"
)

func TestAccBlockStorageV2Volume_basic(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("opentelekomcloud_blockstorage_volume_v2.volume_1", &volume),
					testAccCheckBlockStorageV2VolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_blockstorage_volume_v2.volume_1", "name", "volume_1"),
				),
			},
			{
				Config: testAccBlockStorageV2Volume_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("opentelekomcloud_blockstorage_volume_v2.volume_1", &volume),
					testAccCheckBlockStorageV2VolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_blockstorage_volume_v2.volume_1", "name", "volume_1-updated"),
				),
			},
		},
	})
}

func TestAccBlockStorageV2Volume_upscaleDownScale(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_basic,
				Check:  testAccCheckBlockStorageV2VolumeExists("opentelekomcloud_blockstorage_volume_v2.volume_1", &volume),
			},
			{
				Config: testAccBlockStorageV2Volume_bigger,
				Check:  testAccCheckBlockStorageV2VolumeSame("opentelekomcloud_blockstorage_volume_v2.volume_1", &volume),
			},
			{
				Config: testAccBlockStorageV2Volume_basic,
				Check:  testAccCheckBlockStorageV2VolumeNew("opentelekomcloud_blockstorage_volume_v2.volume_1", &volume),
			},
		},
	})
}
func TestAccBlockStorageV2Volume_upscaleDownScaleAssigned(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_assigned(10),
				Check:  testAccCheckBlockStorageV2VolumeExists("opentelekomcloud_blockstorage_volume_v2.volume_1", &volume),
			},
			{
				Config: testAccBlockStorageV2Volume_assigned(12),
				Check:  testAccCheckBlockStorageV2VolumeSame("opentelekomcloud_blockstorage_volume_v2.volume_1", &volume),
			},
			{
				Config: testAccBlockStorageV2Volume_assigned(10),
				Check:  testAccCheckBlockStorageV2VolumeNew("opentelekomcloud_blockstorage_volume_v2.volume_1", &volume),
			},
		},
	})
}

func TestAccBlockStorageV2Volume_policy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testPolicyPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_policy(os.Getenv("OS_KMS_KEY")),
			},
		},
	})
}

func testPolicyPreCheck(t *testing.T) {
	if os.Getenv("OS_KMS_KEY") == "" {
		t.Errorf("OS_KMS_KEY should be set for this test to existing KMS key alias")
	}
}

func TestAccBlockStorageV2Volume_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_tags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeTags("opentelekomcloud_blockstorage_volume_v2.volume_1", "foo", "bar"),
					testAccCheckBlockStorageV2VolumeTags("opentelekomcloud_blockstorage_volume_v2.volume_1", "key", "value"),
				),
			},
			{
				Config: testAccBlockStorageV2Volume_tags_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeTags("opentelekomcloud_blockstorage_volume_v2.volume_1", "foo2", "bar2"),
					testAccCheckBlockStorageV2VolumeTags("opentelekomcloud_blockstorage_volume_v2.volume_1", "key2", "value2"),
				),
			},
		},
	})
}

func TestAccBlockStorageV2Volume_image(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_image,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("opentelekomcloud_blockstorage_volume_v2.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_blockstorage_volume_v2.volume_1", "name", "volume_1"),
				),
			},
		},
	})
}

func TestAccBlockStorageV2Volume_timeout(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("opentelekomcloud_blockstorage_volume_v2.volume_1", &volume),
				),
			},
		},
	})
}

func testAccCheckBlockStorageV2VolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.blockStorageV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_blockstorage_volume_v2" {
			continue
		}

		_, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func testAccCheckBlockStorageV2VolumeExists(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.blockStorageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud block storage client: %s", err)
		}

		found, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		*volume = *found

		return nil
	}
}

func testAccCheckBlockStorageV2VolumeDoesNotExist(t *testing.T, n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.blockStorageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud block storage client: %s", err)
		}

		_, err = volumes.Get(blockStorageClient, volume.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		return fmt.Errorf("Volume still exists")
	}
}

func testAccCheckBlockStorageV2VolumeMetadata(
	volume *volumes.Volume, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if volume.Metadata == nil {
			return fmt.Errorf("No metadata")
		}

		for key, value := range volume.Metadata {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, value)
		}

		return fmt.Errorf("Metadata not found: %s", k)
	}
}

func testAccCheckBlockStorageV2VolumeTags(n string, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.blockStorageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud block storage client: %s", err)
		}

		found, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		client, err := config.blockStorageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud block storage client: %s", err)
		}
		taglist, err := tags.Get(client, "volumes", found.ID).Extract()
		for key, value := range taglist.Tags {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, value)
		}

		return fmt.Errorf("Tag not found: %s", k)
	}
}

func testAccCheckBlockStorageV2VolumeSame(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID != volume.ID {
			return fmt.Errorf("volume ID changed")
		}
		return nil
	}
}

func testAccCheckBlockStorageV2VolumeNew(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := testAccCheckBlockStorageV2VolumeSame(n, volume)(s)
		if err == nil {
			return fmt.Errorf("volume ID not changed")
		}
		return nil
	}
}

const testAccBlockStorageV2Volume_basic = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
}
`

const testAccBlockStorageV2Volume_bigger = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 2
}
`

func testAccBlockStorageV2Volume_assigned(size int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = %d
  image_id = "%s"
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
	uuid = "%s"
  }
  block_device {
	uuid = opentelekomcloud_blockstorage_volume_v2.volume_1.id
	source_type = "volume"
	boot_index = 0
	destination_type = "volume"
	delete_on_termination = true
  }
}
`, size, OS_IMAGE_ID, OS_NETWORK_ID)
}

const testAccBlockStorageV2Volume_update = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1-updated"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
}
`

const testAccBlockStorageV2Volume_tags = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  tags = {
    foo = "bar"
	key = "value"
  }
  size = 1
}
`

const testAccBlockStorageV2Volume_tags_update = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1-updated"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  tags = {
    foo2 = "bar2"
	key2 = "value2"
  }
  size = 1
}
`

var testAccBlockStorageV2Volume_image = fmt.Sprintf(`
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 12
  image_id = "%s"
}
`, OS_IMAGE_ID)

const testAccBlockStorageV2Volume_timeout = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  size = 1
  device_type = "SCSI"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

func testAccBlockStorageV2Volume_policy(kmsKeyAlias string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_kms_key_v1" key {
  key_alias = "%s"
}

data opentelekomcloud_compute_availability_zones_v2 available {}

locals {
  count = 1
}

resource "opentelekomcloud_blockstorage_volume_v2" "volume" {
  count             = local.count
  availability_zone = data.opentelekomcloud_compute_availability_zones_v2.available.names[count.index]
  name              = "test-vol0${count.index + 1}-datadisk"
  size              = 40
  tags              =   {
    generator   = "terraform"
  }
  metadata = {
    __system__encrypted = "1"
    __system__cmkid     = data.opentelekomcloud_kms_key_v1.key.id
    attached_mode       = "rw"
    readonly            = "False"
  }
}

resource "opentelekomcloud_vbs_backup_policy_v2" "vbs_policy1" {
  name   = "policy_001"
  status = "ON"

  start_time          = "12:00"
  retain_first_backup = "N"
  rentention_num      = 7
  frequency           = 1

  resources = opentelekomcloud_blockstorage_volume_v2.volume[*].id

}
`, kmsKeyAlias)
}

func testAccBlockStorageV2Volume_policyUpdate(kmsKeyAlias string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_kms_key_v1" key {
  key_alias = "%s"
}

data opentelekomcloud_compute_availability_zones_v2 available {}

locals {
  count = 1
}

resource "opentelekomcloud_blockstorage_volume_v2" "volume" {
  count             = local.count
  availability_zone = data.opentelekomcloud_compute_availability_zones_v2.available.names[count.index]
  name              = "test-vol0${count.index + 1}-datadisk"
  size              = 60
  tags              =   {
    generator   = "terraform"
  }
  metadata = {
    __system__encrypted = "1"
    __system__cmkid     = data.opentelekomcloud_kms_key_v1.key.id
    attached_mode       = "rw"
    readonly            = "False"
  }
}

resource "opentelekomcloud_vbs_backup_policy_v2" "vbs_policy1" {
  name   = "policy_001"
  status = "ON"

  start_time          = "12:00"
  retain_first_backup = "N"
  rentention_num      = 7
  frequency           = 1

  resources = opentelekomcloud_blockstorage_volume_v2.volume[*].id

}
`, kmsKeyAlias)
}
