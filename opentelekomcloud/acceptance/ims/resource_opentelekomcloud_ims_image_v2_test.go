package acceptance

import (
	"fmt"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/cloudimages"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/tags"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/ims"
)

const resourceImageName = "opentelekomcloud_ims_image_v2.image_1"

func TestAccImsImageV2_basic(t *testing.T) {
	var image cloudimages.Image

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImsImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImageV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageV2Exists(resourceImageName, &image),
					testAccCheckImsImageV2Tags(resourceImageName, "foo", "bar"),
					testAccCheckImsImageV2Tags(resourceImageName, "key", "value"),
					resource.TestCheckResourceAttr(resourceImageName, "name", "TFTest_image"),
				),
			},
			{
				Config: testAccImsImageV2Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageV2Exists(resourceImageName, &image),
					testAccCheckImsImageV2Tags(resourceImageName, "foo", "bar"),
					testAccCheckImsImageV2Tags(resourceImageName, "key", "value1"),
					testAccCheckImsImageV2Tags(resourceImageName, "key2", "value2"),
					resource.TestCheckResourceAttr(resourceImageName, "name", "TFTest_image_update"),
				),
			},
		},
	})
}

func testAccCheckImsImageV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	imageClient, err := config.ImageV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_ims_image_v2" {
			continue
		}

		_, err := ims.GetCloudImage(imageClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("image still exists")
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

func testAccCheckImsImageV2Tags(n string, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("IMS Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		imageClient, err := config.ImageV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud image client: %s", err)
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
			return fmt.Errorf("bad value for %s: %s", k, tag.Value)
		}
		return fmt.Errorf("tag not found: %s", k)
	}
}

var testAccImsImageV2Basic = fmt.Sprintf(`
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
}

resource "opentelekomcloud_ims_image_v2" "image_1" {
  name        = "TFTest_image"
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
  description = "created by TerraformAccTest"
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccImsImageV2Update = fmt.Sprintf(`
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
}

resource "opentelekomcloud_ims_image_v2" "image_1" {
  name        = "TFTest_image_update"
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
  description = "created by TerraformAccTest"
  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
