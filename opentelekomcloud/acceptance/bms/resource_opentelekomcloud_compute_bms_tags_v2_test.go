package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/bms/v2/tags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceTagsName = "opentelekomcloud_compute_bms_tags_v2.tags_1"

func TestAccBMSTagsV2_basic(t *testing.T) {
	var tagList []string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBMSTagsV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSTagsV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSTagsV2Exists(resourceTagsName, tagList),
					resource.TestCheckResourceAttr(resourceTagsName, "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccBMSTagsV2_timeout(t *testing.T) {
	var tagList []string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBMSTagsV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSTagsV2Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSTagsV2Exists(resourceTagsName, tagList),
				),
			},
		},
	})
}

func testAccCheckBMSTagsV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ComputeV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud CombuteV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_bms_tags_v2" {
			continue
		}

		_, err := tags.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("tags still exists")
		}
	}

	return nil
}

func testAccCheckBMSTagsV2Exists(n string, tag []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ComputeV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud CombuteV2 client: %s", err)
		}

		found, err := tags.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		tag = found
		return nil
	}
}

var testAccBMSTagsV2Basic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "BMSinstance_1"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "physical.o2.medium"
  flavor_name       = "physical.o2.medium"
  tags              = {
    foo = "bar"
    john = "doe"
  }
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
resource "opentelekomcloud_compute_bms_tags_v2" "tags_1" {
  server_id = opentelekomcloud_compute_instance_v2.instance_1.id
  tags      = ["foo", "bar"]
}`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccBMSTagsV2Timeout = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "BMSinstance_1"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "physical.o2.medium"
  flavor_name       = "physical.o2.medium"
  tags              = ["foo", "bar"]
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
resource "opentelekomcloud_compute_bms_tags_v2" "tags_1" {
  server_id = opentelekomcloud_compute_instance_v2.instance_1.id
  tags      = ["foo", "bar"]
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
