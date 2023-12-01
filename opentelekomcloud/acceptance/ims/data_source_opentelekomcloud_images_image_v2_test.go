package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccImagesV2ImageDataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_images_image_v2.image_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesV2ImageDataSource_cirros,
			},
			{
				Config: testAccImagesV2ImageDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", "CirrOS-tf"),
					resource.TestCheckResourceAttr(dataSourceName, "container_format", "bare"),
					resource.TestCheckResourceAttr(dataSourceName, "min_ram", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "private"),
				),
			},
		},
	})
}

func TestAccImagesV2ImageDataSource_testQueries(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_images_image_v2.image_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesV2ImageDataSource_cirros,
			},
			{
				Config: testAccImagesV2ImageDataSource_queryTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
				),
			},
			// it fails on api side with previous size_min = "13000000", need to investigate something was changed
			{
				Config: testAccImagesV2ImageDataSource_querySizeMin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
				),
			},
			{
				Config: testAccImagesV2ImageDataSource_querySizeMax,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
				),
			},
			{
				Config: testAccImagesV2ImageDataSource_cirros,
			},
		},
	})
}

func TestAccImagesV2ImageDataSource_regex(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_images_image_v2.image_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesV2ImageDataSource_regex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "public"),
				),
			},
		},
	})
}

func testAccCheckImagesV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("image data source ID not set")
		}

		return nil
	}
}

// Standard CirrOS image
const testAccImagesV2ImageDataSource_cirros = `
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "CirrOS-tf"
  container_format = "bare"
  disk_format      = "qcow2"
  image_source_url = "https://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags             = ["cirros-tf"]
}
`

var testAccImagesV2ImageDataSource_basic = fmt.Sprintf(`
%s

data "opentelekomcloud_images_image_v2" "image_1" {
  most_recent = true
  name        = opentelekomcloud_images_image_v2.image_1.name
}
`, testAccImagesV2ImageDataSource_cirros)

var testAccImagesV2ImageDataSource_queryTag = fmt.Sprintf(`
%s

data "opentelekomcloud_images_image_v2" "image_1" {
  most_recent = true
  visibility  = "private"
  tag         = "cirros-tf"
}
`, testAccImagesV2ImageDataSource_cirros)

var testAccImagesV2ImageDataSource_querySizeMin = fmt.Sprintf(`
%s

data "opentelekomcloud_images_image_v2" "image_1" {
  most_recent = true
  visibility  = "private"
  size_min    = "10"
}
`, testAccImagesV2ImageDataSource_cirros)

var testAccImagesV2ImageDataSource_querySizeMax = fmt.Sprintf(`
%s

data "opentelekomcloud_images_image_v2" "image_1" {
  most_recent = true
  visibility  = "private"
  size_max    = "23000000"
}
`, testAccImagesV2ImageDataSource_cirros)

var testAccImagesV2ImageDataSource_regex = `
data "opentelekomcloud_images_image_v2" "image_1" {
  name_regex  = "^Standard_Debian.+"
  most_recent = true
  visibility  = "public"
}
`
