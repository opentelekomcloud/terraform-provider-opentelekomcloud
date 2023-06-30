package acceptance

import (
	"fmt"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/image/v2/images"
	ims "github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/images"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccImagesImageV2_basic(t *testing.T) {
	var image ims.ImageInfo

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("opentelekomcloud_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_images_image_v2.image_1", "name", "Rancher TerraformAccTest"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_images_image_v2.image_1", "container_format", "bare"),
					/*resource.TestCheckResourceAttr(
					"opentelekomcloud_images_image_v2.image_1", "disk_format", "qcow2"), */
					resource.TestCheckResourceAttr(
						"opentelekomcloud_images_image_v2.image_1", "schema", "/v2/schemas/image"),
				),
			},
		},
	})
}

func TestAccImagesImageV2_name(t *testing.T) {
	var image ims.ImageInfo

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_name_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("opentelekomcloud_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_images_image_v2.image_1", "name", "Rancher TerraformAccTest"),
				),
			},
			{
				Config: testAccImagesImageV2_name_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("opentelekomcloud_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_images_image_v2.image_1", "name", "TerraformAccTest Rancher"),
				),
			},
		},
	})
}

func TestAccImagesImageV2_tags(t *testing.T) {
	var image ims.ImageInfo

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_tags_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("opentelekomcloud_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("opentelekomcloud_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("opentelekomcloud_images_image_v2.image_1", "bar"),
					testAccCheckImagesImageV2TagCount("opentelekomcloud_images_image_v2.image_1", 2),
				),
			},
			{
				Config: testAccImagesImageV2_tags_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("opentelekomcloud_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("opentelekomcloud_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("opentelekomcloud_images_image_v2.image_1", "bar"),
					testAccCheckImagesImageV2HasTag("opentelekomcloud_images_image_v2.image_1", "baz"),
					testAccCheckImagesImageV2TagCount("opentelekomcloud_images_image_v2.image_1", 3),
				),
			},
			{
				Config: testAccImagesImageV2_tags_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("opentelekomcloud_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("opentelekomcloud_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("opentelekomcloud_images_image_v2.image_1", "baz"),
					testAccCheckImagesImageV2TagCount("opentelekomcloud_images_image_v2.image_1", 2),
				),
			},
		},
	})
}

func TestAccImagesImageV2_visibility(t *testing.T) {
	var image ims.ImageInfo

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_visibility_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("opentelekomcloud_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_images_image_v2.image_1", "visibility", "private"),
				),
			},
			{
				Config: testAccImagesImageV2_visibility_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("opentelekomcloud_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_images_image_v2.image_1", "visibility", "public"),
				),
			},
		},
	})
}

func TestAccImagesImageV2_timeout(t *testing.T) {
	var image ims.ImageInfo

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("opentelekomcloud_images_image_v2.image_1", &image),
				),
			},
		},
	})
}

func testAccCheckImagesImageV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	imageClient, err := config.ImageV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_images_image_v2" {
			continue
		}

		_, err := images.Get(imageClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("image still exists")
		}
	}

	return nil
}

func testAccCheckImagesImageV2Exists(n string, image *ims.ImageInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		imageClient, err := config.ImageV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud Image: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("image not found")
		}

		*image = *found

		return nil
	}
}

func testAccCheckImagesImageV2HasTag(n, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		imageClient, err := config.ImageV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud Image: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("image not found")
		}

		for _, v := range found.Tags {
			if tag == v {
				return nil
			}
		}

		return fmt.Errorf("tag not found: %s", tag)
	}
}

func testAccCheckImagesImageV2TagCount(n string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		imageClient, err := config.ImageV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud Image: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("image not found")
		}

		if len(found.Tags) != expected {
			return fmt.Errorf("expecting %d tags, found %d", expected, len(found.Tags))
		}

		return nil
	}
}

var testAccImagesImageV2_basic = `
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "Rancher TerraformAccTest"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
}`

var testAccImagesImageV2_name_1 = `
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "Rancher TerraformAccTest"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
}`

var testAccImagesImageV2_name_2 = `
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "TerraformAccTest Rancher"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
}`

var testAccImagesImageV2_tags_1 = `
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "Rancher TerraformAccTest"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
  tags             = ["foo", "bar"]
}`

var testAccImagesImageV2_tags_2 = `
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "Rancher TerraformAccTest"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
  tags             = ["foo", "bar", "baz"]
}`

var testAccImagesImageV2_tags_3 = `
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "Rancher TerraformAccTest"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
  tags             = ["foo", "baz"]
}`

var testAccImagesImageV2_visibility_1 = `
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "Rancher TerraformAccTest"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
  visibility       = "private"
}`

var testAccImagesImageV2_visibility_2 = `
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "Rancher TerraformAccTest"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
  visibility       = "public"
}`

var testAccImagesImageV2_timeout = `
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "Rancher TerraformAccTest"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"

  timeouts {
    create = "10m"
  }
}`
