package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/ims/v2/cloudimages"
	"github.com/huaweicloud/golangsdk/openstack/ims/v2/tags"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccImsImageV2_basic(t *testing.T) {
	var image cloudimages.Image

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImsImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImageV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageV2Exists("opentelekomcloud_ims_image_v2.image_1", &image),
					testAccCheckImsImageV2Tags("opentelekomcloud_ims_image_v2.image_1", "foo", "bar"),
					testAccCheckImsImageV2Tags("opentelekomcloud_ims_image_v2.image_1", "key", "value"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_ims_image_v2.image_1", "name", "TFTest_image"),
				),
			},
			{
				Config: testAccImsImageV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageV2Exists("opentelekomcloud_ims_image_v2.image_1", &image),
					testAccCheckImsImageV2Tags("opentelekomcloud_ims_image_v2.image_1", "foo", "bar"),
					testAccCheckImsImageV2Tags("opentelekomcloud_ims_image_v2.image_1", "key", "value1"),
					testAccCheckImsImageV2Tags("opentelekomcloud_ims_image_v2.image_1", "key2", "value2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_ims_image_v2.image_1", "name", "TFTest_image_update"),
				),
			},
		},
	})
}

func testAccCheckImsImageV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	imageClient, err := config.imageV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_ims_image_v2" {
			continue
		}

		_, err := getCloudimage(imageClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Image still exists")
		}
	}

	return nil
}

func testAccCheckImsImageV2Exists(n string, image *cloudimages.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("IMS Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		imageClient, err := config.imageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud Image: %s", err)
		}

		found, err := getCloudimage(imageClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		*image = *found
		return nil
	}
}

func testAccCheckImsImageV2Tags(n string, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("IMS Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		imageClient, err := config.imageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud image client: %s", err)
		}

		found, err := tags.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Tags == nil {
			return fmt.Errorf("IMS Tags not found")
		}

		for _, tag := range found.Tags {
			if k != tag.Key {
				continue
			}

			if v == tag.Value {
				return nil
			}
			return fmt.Errorf("Bad value for %s: %s", k, tag.Value)
		}
		return fmt.Errorf("Tag not found: %s", k)
	}
}

var testAccImsImageV2_basic = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}

resource "opentelekomcloud_ims_image_v2" "image_1" {
  name   = "TFTest_image"
  instance_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
  description = "created by TerraformAccTest"
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccImsImageV2_update = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}

resource "opentelekomcloud_ims_image_v2" "image_1" {
  name   = "TFTest_image_update"
  instance_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
  description = "created by TerraformAccTest"
  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
