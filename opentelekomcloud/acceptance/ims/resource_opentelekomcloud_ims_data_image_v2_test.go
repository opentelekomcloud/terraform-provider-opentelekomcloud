package acceptance

import (
	"fmt"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/images"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/ims"
)

const resourceDataImageName = "opentelekomcloud_ims_data_image_v2.image_1"

func TestAccImsDataImageV2_basic(t *testing.T) {
	var image images.ImageInfo

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImsDataImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsDataImageV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsDataImageV2Exists(resourceDataImageName, &image),
					testAccCheckImsImageV2Tags(resourceDataImageName, "foo", "bar"),
					testAccCheckImsImageV2Tags(resourceDataImageName, "key", "value"),
					resource.TestCheckResourceAttr(resourceDataImageName, "name", "TFTest_data_image"),
				),
			},
			{
				Config: testAccImsDataImageV2Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsDataImageV2Exists(resourceDataImageName, &image),
					testAccCheckImsImageV2Tags(resourceDataImageName, "foo", "bar"),
					testAccCheckImsImageV2Tags(resourceDataImageName, "key", "value1"),
					testAccCheckImsImageV2Tags(resourceDataImageName, "key2", "value2"),
					resource.TestCheckResourceAttr(resourceDataImageName, "name", "TFTest_data_image_update"),
				),
			},
		},
	})
}

func testAccCheckImsDataImageV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	imageClient, err := config.ImageV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_ims_data_image_v2" {
			continue
		}

		_, err := ims.GetCloudImage(imageClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("image still exists")
		}
	}

	return nil
}

func testAccCheckImsDataImageV2Exists(n string, image *images.ImageInfo) resource.TestCheckFunc {
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

		found, err := ims.GetCloudImage(imageClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		*image = *found
		return nil
	}
}

var testAccImsDataImageV2Basic = fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  security_groups   = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "volume"
    volume_size           = 40
    source_type           = "image"
    uuid                  = data.opentelekomcloud_images_image_v2.latest_image.id
  }
  block_device {
    boot_index            = 1
    delete_on_termination = true
    destination_type      = "volume"
    source_type           = "blank"
    volume_size           = 1
  }
}

resource "opentelekomcloud_ims_data_image_v2" "image_1" {
  name        = "TFTest_data_image"
  description = "created by TerraformAccTest"
  volume_id   = opentelekomcloud_compute_instance_v2.instance_1.volume_attached.1.id
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.DataSourceSubnet, common.DataSourceImage, env.OS_AVAILABILITY_ZONE)

var testAccImsDataImageV2Update = fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  security_groups   = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "volume"
    volume_size           = 40
    source_type           = "image"
    uuid                  = data.opentelekomcloud_images_image_v2.latest_image.id
  }
  block_device {
    boot_index            = 1
    delete_on_termination = true
    destination_type      = "volume"
    source_type           = "blank"
    volume_size           = 1
  }
}

resource "opentelekomcloud_ims_data_image_v2" "image_1" {
  name        = "TFTest_data_image_update"
  description = "created by TerraformAccTest"
  volume_id   = opentelekomcloud_compute_instance_v2.instance_1.volume_attached.1.id
  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, common.DataSourceSubnet, common.DataSourceImage, env.OS_AVAILABILITY_ZONE)
