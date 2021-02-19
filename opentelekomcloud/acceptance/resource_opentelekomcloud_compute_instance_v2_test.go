package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/secgroups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ecs/v1/tags"

	"github.com/opentelekomcloud/gophertelekomcloud/pagination"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccComputeV2Instance_basic(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "availability_zone", OS_AVAILABILITY_ZONE),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_tags(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_withTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceTagsV1(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceTagsV1(&instance, "key", "value"),
				),
			},
			{
				Config: testAccComputeV2Instance_updateTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceTagsV1(&instance, "foo2", "bar2"),
					testAccCheckComputeV2InstanceTagsV1(&instance, "key", "value2"),
				),
			},
			{
				Config: testAccComputeV2Instance_withoutTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceNoTagV1(&instance),
				),
			},
			{
				Config: testAccComputeV2Instance_withTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceTagsV1(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceTagsV1(&instance, "key", "value"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_tag(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceTagsV1(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceTagsV1(&instance, "key", "value"),
				),
			},
			{
				Config: testAccComputeV2Instance_updateTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceTagsV1(&instance, "foo", "bar2"),
					testAccCheckComputeV2InstanceTagsV1(&instance, "key1", "value"),
				),
			},
			{
				Config: testAccComputeV2Instance_withoutTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceNoTagV1(&instance),
				),
			},
			{
				Config: testAccComputeV2Instance_tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceTagsV1(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceTagsV1(&instance, "key", "value"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_multiSecgroup(t *testing.T) {
	var instance servers.Server
	var firstSecGroup, secondSecGroup secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_multiSecgroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists(
						"opentelekomcloud_compute_secgroup_v2.secgroup_1", &firstSecGroup),
					testAccCheckComputeV2SecGroupExists(
						"opentelekomcloud_compute_secgroup_v2.secgroup_2", &secondSecGroup),
					testAccCheckComputeV2InstanceExists(
						"opentelekomcloud_compute_instance_v2.instance_1", &instance),
				),
			},
			{
				Config: testAccComputeV2Instance_multiSecgroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists(
						"opentelekomcloud_compute_secgroup_v2.secgroup_1", &firstSecGroup),
					testAccCheckComputeV2SecGroupExists(
						"opentelekomcloud_compute_secgroup_v2.secgroup_2", &secondSecGroup),
					testAccCheckComputeV2InstanceExists(
						"opentelekomcloud_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeImage(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeVolume(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeForceNew(t *testing.T) {
	var instance servers.Server
	var newInstance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"opentelekomcloud_compute_instance_v2.instance_1", &instance),
				),
			},
			{
				Config: testAccComputeV2Instance_bootFromVolumeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"opentelekomcloud_compute_instance_v2.instance_1", &newInstance),
					testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(&instance, &newInstance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_changeFixedIP(t *testing.T) {
	var instance servers.Server
	var newInstance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_fixedIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"opentelekomcloud_compute_instance_v2.instance_1", &instance),
				),
			},
			{
				Config: testAccComputeV2Instance_fixedIPUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"opentelekomcloud_compute_instance_v2.instance_1", &newInstance),
					testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(&instance, &newInstance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_stopBeforeDestroy(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_stopBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_metadata(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_metadata,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceMetadata(&instance, "abc", "def"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "all_metadata.abc", "def"),
				),
			},
			{
				Config: testAccComputeV2Instance_metadataUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceMetadata(&instance, "ghi", "jkl"),
					testAccCheckComputeV2InstanceNoMetadataKey(&instance, "abc"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "all_metadata.ghi", "jkl"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_timeout(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_autoRecovery(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "auto_recovery", "true"),
				),
			},
			{
				Config: testAccComputeV2Instance_autoRecovery,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_instance_v2.instance_1", "auto_recovery", "false"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_crazyNICs(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_crazyNICs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func testAccCheckComputeV2InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*cfg.Config)
	computeClient, err := config.ComputeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_instance_v2" {
			continue
		}

		server, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "SOFT_DELETED" {
				return fmt.Errorf("instance still exists")
			}
		}
	}

	return nil
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

		config := testAccProvider.Meta().(*cfg.Config)
		computeClient, err := config.ComputeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
		}

		found, err := servers.Get(computeClient, rs.Primary.ID).Extract()
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

func testAccCheckComputeV2InstanceTagsV1(instance *servers.Server, k, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*cfg.Config)
		client, err := config.ComputeV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud ComputeV1 client: %s", err)
		}

		tagsList, err := tags.Get(client, instance.ID).Extract()
		if err != nil {
			return err
		}
		for _, val := range tagsList.Tags {
			if k != val.Key {
				continue
			}
			if v == val.Value {
				return nil
			}

			return fmt.Errorf("bad value for %s: %s", k, val.Value)
		}

		return fmt.Errorf("tag not found: %s", k)
	}
}

func testAccCheckComputeV2InstanceNoTagV1(instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*cfg.Config)
		client, err := config.ComputeV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud ComputeV1 client: %s", err)
		}

		tagList, err := tags.Get(client, instance.ID).Extract()
		if err != nil {
			return err
		}

		if tagList.Tags == nil {
			return nil
		}
		if len(tagList.Tags) == 0 {
			return nil
		}

		return fmt.Errorf("expected no tags, but found %v", tagList.Tags)
	}
}

func testAccCheckComputeV2InstanceBootVolumeAttachment(instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var attachments []volumeattach.VolumeAttachment

		config := testAccProvider.Meta().(*cfg.Config)
		computeClient, err := config.ComputeV2Client(OS_REGION_NAME)
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

		if len(attachments) == 1 {
			return nil
		}

		return fmt.Errorf("no attached volume found")
	}
}

func testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(instance1, instance2 *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance1.ID == instance2.ID {
			return fmt.Errorf("instance was not recreated")
		}

		return nil
	}
}

var testAccComputeV2Instance_basic = fmt.Sprintf(`
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
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2Instance_withoutTags = fmt.Sprintf(`
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
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2Instance_withTags = fmt.Sprintf(`
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
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2Instance_updateTags = fmt.Sprintf(`
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
  tags = {
    foo2 = "bar2"
    key  = "value2"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2Instance_tag = fmt.Sprintf(`
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
  tag = {
    foo = "bar"
    key = "value"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2Instance_updateTag = fmt.Sprintf(`
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
  tag = {
    foo  = "bar2"
    key1 = "value"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2Instance_multiSecgroup = fmt.Sprintf(`
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
`, OS_NETWORK_ID)

var testAccComputeV2Instance_multiSecgroupUpdate = fmt.Sprintf(`
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
`, OS_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolumeImage = fmt.Sprintf(`
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
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_IMAGE_ID)

var testAccComputeV2Instance_bootFromVolumeVolume = fmt.Sprintf(`
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
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolume = fmt.Sprintf(`
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
`, OS_NETWORK_ID, OS_IMAGE_ID)

var testAccComputeV2Instance_bootFromVolumeUpdate = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  block_device {
    uuid                  = "%s"
    source_type           = "image"
    volume_size           = 51
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }
}
`, OS_NETWORK_ID, OS_IMAGE_ID)

var testAccComputeV2Instance_fixedIP = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid        = "%s"
    fixed_ip_v4 = "192.168.0.24"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_fixedIPUpdate = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid        = "%s"
    fixed_ip_v4 = "192.168.0.25"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_stopBeforeDestroy = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  stop_before_destroy = true
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_metadata = fmt.Sprintf(`
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
`, OS_NETWORK_ID)

var testAccComputeV2Instance_metadataUpdate = fmt.Sprintf(`
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
`, OS_NETWORK_ID)

var testAccComputeV2Instance_timeout = fmt.Sprintf(`
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
`, OS_NETWORK_ID)

var testAccComputeV2Instance_autoRecovery = fmt.Sprintf(`
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
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2Instance_crazyNICs = fmt.Sprintf(`
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
`, OS_NETWORK_ID)
