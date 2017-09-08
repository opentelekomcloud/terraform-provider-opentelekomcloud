package hwcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHWCloudImagesV2ImageDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:             testAccHWCloudImagesV2ImageDataSource_cirros,
				ExpectNonEmptyPlan: true,
			},
			resource.TestStep{
				Config:             testAccHWCloudImagesV2ImageDataSource_basic,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.hwcloud_images_image_v2.image_1"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_images_image_v2.image_1", "name", "CirrOS-tf"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_images_image_v2.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_images_image_v2.image_1", "disk_format", "qcow2"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_images_image_v2.image_1", "min_disk_gb", "0"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_images_image_v2.image_1", "min_ram_mb", "0"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_images_image_v2.image_1", "protected", "false"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_images_image_v2.image_1", "visibility", "private"),
				),
			},
		},
	})
}

func TestAccHWCloudImagesV2ImageDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:             testAccHWCloudImagesV2ImageDataSource_cirros,
				ExpectNonEmptyPlan: true,
			},
			resource.TestStep{
				Config:             testAccHWCloudImagesV2ImageDataSource_queryTag,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.hwcloud_images_image_v2.image_1"),
				),
			},
			resource.TestStep{
				Config: testAccHWCloudImagesV2ImageDataSource_querySizeMin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.hwcloud_images_image_v2.image_1"),
				),
			},
			resource.TestStep{
				Config:             testAccHWCloudImagesV2ImageDataSource_querySizeMax,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.hwcloud_images_image_v2.image_1"),
				),
			},
			resource.TestStep{
				Config:             testAccHWCloudImagesV2ImageDataSource_cirros,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckImagesV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Image data source ID not set")
		}

		return nil
	}
}

// Standard CirrOS image
const testAccHWCloudImagesV2ImageDataSource_cirros = `
resource "hwcloud_images_image_v2" "image_1" {
	name = "CirrOS-tf"
	container_format = "bare"
	disk_format = "qcow2"
	image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
	tags = ["cirros-tf"]
}
`

var testAccHWCloudImagesV2ImageDataSource_basic = fmt.Sprintf(`
%s

data "hwcloud_images_image_v2" "image_1" {
	most_recent = true
	name = "${hwcloud_images_image_v2.image_1.name}"
}
`, testAccHWCloudImagesV2ImageDataSource_cirros)

var testAccHWCloudImagesV2ImageDataSource_queryTag = fmt.Sprintf(`
%s

data "hwcloud_images_image_v2" "image_1" {
	most_recent = true
	visibility = "private"
	tag = "cirros-tf"
}
`, testAccHWCloudImagesV2ImageDataSource_cirros)

var testAccHWCloudImagesV2ImageDataSource_querySizeMin = fmt.Sprintf(`
%s

data "hwcloud_images_image_v2" "image_1" {
	most_recent = true
	visibility = "private"
	size_min = "13000000"
}
`, testAccHWCloudImagesV2ImageDataSource_cirros)

var testAccHWCloudImagesV2ImageDataSource_querySizeMax = fmt.Sprintf(`
%s

data "hwcloud_images_image_v2" "image_1" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
}
`, testAccHWCloudImagesV2ImageDataSource_cirros)
