package acceptance

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v3/volumes"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceVolumeV3Name = "opentelekomcloud_evs_volume_v3.volume_1"

func TestAccEvsStorageV3Volume_basic(t *testing.T) {
	var volume volumes.Volume
	t.Parallel()
	qts := []*quotas.ExpectedQuota{{Q: quotas.Volume, Count: 1}, {Q: quotas.VolumeSize, Count: 12}}
	th.AssertNoErr(t, quotas.AcquireMultipleQuotas(qts, 5*time.Second))
	defer quotas.ReleaseMultipleQuotas(qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEvsStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsStorageV3VolumeBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeExists(resourceVolumeV3Name, &volume),
					resource.TestCheckResourceAttr(resourceVolumeV3Name, "name", "volume_1"),
				),
			},
			{
				Config: testAccEvsStorageV3VolumeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeExists(resourceVolumeV3Name, &volume),
					resource.TestCheckResourceAttr(resourceVolumeV3Name, "name", "volume_1-updated"),
				),
			},
		},
	})
}

func TestAccEvsStorageV3Volume_tags(t *testing.T) {
	t.Parallel()
	qts := []*quotas.ExpectedQuota{{Q: quotas.Volume, Count: 1}, {Q: quotas.VolumeSize, Count: 12}}
	th.AssertNoErr(t, quotas.AcquireMultipleQuotas(qts, 5*time.Second))
	defer quotas.ReleaseMultipleQuotas(qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEvsStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsStorageV3VolumeTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVolumeV3Name, "tags.muh", "value-create"),
				),
			},
			{
				Config: testAccEvsStorageV3VolumeTagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVolumeV3Name, "tags.muh", "value-update"),
				),
			},
		},
	})
}

func TestAccEvsStorageV3Volume_image(t *testing.T) {
	var volume volumes.Volume
	t.Parallel()
	qts := []*quotas.ExpectedQuota{{Q: quotas.Volume, Count: 1}, {Q: quotas.VolumeSize, Count: 12}}
	th.AssertNoErr(t, quotas.AcquireMultipleQuotas(qts, 5*time.Second))
	defer quotas.ReleaseMultipleQuotas(qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEvsStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsStorageV3VolumeImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeExists(resourceVolumeV3Name, &volume),
					resource.TestCheckResourceAttr(resourceVolumeV3Name, "name", "volume_1"),
				),
			},
		},
	})
}

func TestAccEvsStorageV3Volume_timeout(t *testing.T) {
	var volume volumes.Volume
	t.Parallel()
	qts := []*quotas.ExpectedQuota{{Q: quotas.Volume, Count: 1}, {Q: quotas.VolumeSize, Count: 12}}
	th.AssertNoErr(t, quotas.AcquireMultipleQuotas(qts, 5*time.Second))
	defer quotas.ReleaseMultipleQuotas(qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEvsStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsStorageV3VolumeTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeExists(resourceVolumeV3Name, &volume),
				),
			},
		},
	})
}

func TestAccEvsStorageV3Volume_volumeType(t *testing.T) {
	t.Parallel()
	qts := []*quotas.ExpectedQuota{{Q: quotas.Volume, Count: 1}, {Q: quotas.VolumeSize, Count: 12}}
	th.AssertNoErr(t, quotas.AcquireMultipleQuotas(qts, 5*time.Second))
	defer quotas.ReleaseMultipleQuotas(qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEvsStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccEvsStorageV3VolumeVolumeType,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`volume type .+ doesn't exist`),
			},
		},
	})
}

func TestAccEvsStorageV3Volume_resize(t *testing.T) {
	var volume volumes.Volume
	var volumeUpScaled volumes.Volume
	t.Parallel()
	qts := []*quotas.ExpectedQuota{{Q: quotas.Volume, Count: 1}, {Q: quotas.VolumeSize, Count: 20}}
	th.AssertNoErr(t, quotas.AcquireMultipleQuotas(qts, 5*time.Second))
	defer quotas.ReleaseMultipleQuotas(qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEvsStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsStorageV3VolumeBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeExists(resourceVolumeV3Name, &volume),
					resource.TestCheckResourceAttr(resourceVolumeV3Name, "name", "volume_1"),
				),
			},
			{
				Config: testAccEvsStorageV3VolumeUpscale,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumePersists(resourceVolumeV3Name, &volumeUpScaled, &volume),
					resource.TestCheckResourceAttr(resourceVolumeV3Name, "size", "20"),
				),
			},
		},
	})
}

func testAccCheckEvsStorageV3VolumeDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.BlockStorageV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud BlockStorageV3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_evs_volume_v3" {
			continue
		}

		_, err := volumes.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("volume still exists")
		}
	}

	return nil
}

func testAccCheckEvsStorageV3VolumeExists(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.BlockStorageV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud BlockStorageV3 client: %s", err)
		}

		found, err := volumes.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("volume not found")
		}

		*volume = *found

		return nil
	}
}
func testAccCheckEvsStorageV3VolumePersists(n string, volume, oldVolume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.BlockStorageV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud BlockStorageV3 client: %s", err)
		}

		found, err := volumes.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("volume not found")
		}

		*volume = *found

		if found.ID != oldVolume.ID {
			return fmt.Errorf("volume was re-created")
		}

		return nil
	}
}

var (
	testAccEvsStorageV3VolumeBasic = fmt.Sprintf(`
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  description       = "first test volume"
  availability_zone = "%s"
  volume_type       = "SATA"
  size              = 12
}
`, env.OS_AVAILABILITY_ZONE)
	testAccEvsStorageV3VolumeUpdate = fmt.Sprintf(`
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1-updated"
  description       = "first test volume"
  availability_zone = "%s"
  volume_type       = "SATA"
  size              = 12
}
`, env.OS_AVAILABILITY_ZONE)
	testAccEvsStorageV3VolumeTags = fmt.Sprintf(`
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_tags"
  description       = "test volume with tags"
  availability_zone = "%s"
  volume_type       = "SATA"
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
  size = 12
}
`, env.OS_AVAILABILITY_ZONE)
	testAccEvsStorageV3VolumeTagsUpdate = fmt.Sprintf(`
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_tags-updated"
  description       = "test volume with tags"
  availability_zone = "%s"
  volume_type       = "SATA"
  tags = {
    muh = "value-update"
  }
  size = 12
}
`, env.OS_AVAILABILITY_ZONE)
	testAccEvsStorageV3VolumeImage = fmt.Sprintf(`
%s

resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  availability_zone = "%s"
  volume_type       = "SATA"
  size              = 12
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id
}
`, common.DataSourceImage, env.OS_AVAILABILITY_ZONE)
	testAccEvsStorageV3VolumeTimeout = fmt.Sprintf(`
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  description       = "first test volume"
  availability_zone = "%s"
  size              = 12
  volume_type       = "SATA"
  device_type       = "SCSI"
  timeouts {
    create = "10m"
    delete = "5m"
  }
}
`, env.OS_AVAILABILITY_ZONE)

	testAccEvsStorageV3VolumeVolumeType = fmt.Sprintf(`
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  description       = "first test volume"
  availability_zone = "%s"
  volume_type       = "asfddasf"
  size              = 12
}
`, env.OS_AVAILABILITY_ZONE)
	testAccEvsStorageV3VolumeUpscale = fmt.Sprintf(`
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  description       = "first test volume"
  availability_zone = "%s"
  volume_type       = "SATA"
  size              = 20
}
`, env.OS_AVAILABILITY_ZONE)
)
