package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					resource.TestCheckResourceAttr(
						"opentelekomcloud_ims_image_v2.image_1", "name", "TFTest_image"),
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
  image_tags = {
    foo = "bar"
    key = "value"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
