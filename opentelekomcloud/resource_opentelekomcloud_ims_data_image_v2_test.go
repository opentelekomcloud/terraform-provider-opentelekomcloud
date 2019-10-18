package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccImsDataImageV2_basic(t *testing.T) {
	var image cloudimages.Image

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImsDataImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsDataImageV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsDataImageV2Exists("opentelekomcloud_ims_data_image_v2.image_1", &image),
					testAccCheckImsImageV2Tags("opentelekomcloud_ims_data_image_v2.image_1", "foo", "bar"),
					testAccCheckImsImageV2Tags("opentelekomcloud_ims_data_image_v2.image_1", "key", "value"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_ims_data_image_v2.image_1", "name", "TFTest_data_image"),
				),
			},
			{
				Config: testAccImsDataImageV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsDataImageV2Exists("opentelekomcloud_ims_data_image_v2.image_1", &image),
					testAccCheckImsImageV2Tags("opentelekomcloud_ims_data_image_v2.image_1", "foo", "bar"),
					testAccCheckImsImageV2Tags("opentelekomcloud_ims_data_image_v2.image_1", "key", "value1"),
					testAccCheckImsImageV2Tags("opentelekomcloud_ims_data_image_v2.image_1", "key2", "value2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_ims_data_image_v2.image_1", "name", "TFTest_data_image_update"),
				),
			},
		},
	})
}

func testAccCheckImsDataImageV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	imageClient, err := config.imageV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_ims_data_image_v2" {
			continue
		}

		_, err := getCloudimage(imageClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Image still exists")
		}
	}

	return nil
}

func testAccCheckImsDataImageV2Exists(n string, image *cloudimages.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
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

var testAccImsDataImageV2_basic = fmt.Sprintf(`
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
  block_device {
    boot_index = 0
    delete_on_termination = true
    destination_type = "volume"
	volume_size = 40
    source_type = "image"
    uuid = "%s"
  }
  block_device {
    boot_index = 1
    delete_on_termination = true
    destination_type = "volume"
    source_type = "blank"
    volume_size = 1
  }
}

resource "opentelekomcloud_ims_data_image_v2" "image_1" {
  name   = "TFTest_data_image"
  description = "created by TerraformAccTest"
  volume_id = "${opentelekomcloud_compute_instance_v2.instance_1.volume_attached.1.id}"
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_IMAGE_ID)

var testAccImsDataImageV2_update = fmt.Sprintf(`
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
  block_device {
    boot_index = 0
    delete_on_termination = true
    destination_type = "volume"
	volume_size = 40
    source_type = "image"
    uuid = "%s"
  }
  block_device {
    boot_index = 1
    delete_on_termination = true
    destination_type = "volume"
    source_type = "blank"
    volume_size = 1
  }
}

resource "opentelekomcloud_ims_data_image_v2" "image_1" {
  name   = "TFTest_data_image_update"
  description = "created by TerraformAccTest"
  volume_id = "${opentelekomcloud_compute_instance_v2.instance_1.volume_attached.1.id}"
  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_IMAGE_ID)
