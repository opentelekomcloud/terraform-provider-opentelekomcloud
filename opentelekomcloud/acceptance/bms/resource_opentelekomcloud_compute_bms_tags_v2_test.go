package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/bms/v2/tags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccOTCBMSTagsV2_basic(t *testing.T) {
	var tags tags.Tags

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccBmsFlavorPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckOTCBMSTagsV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSTagsV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCBMSTagsV2Exists("opentelekomcloud_compute_bms_tags_v2.tags_1", &tags),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_bms_tags_v2.tags_1", "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccOTCBMSTagsV2_timeout(t *testing.T) {
	var tags tags.Tags

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccBmsFlavorPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckOTCBMSTagsV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSTagsV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCBMSTagsV2Exists("opentelekomcloud_compute_bms_tags_v2.tags_1", &tags),
				),
			},
		},
	})
}

func testAccCheckOTCBMSTagsV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	bmsClient, err := config.ComputeV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud bms client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_bms_tags_v2" {
			continue
		}

		_, err := tags.Get(bmsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("tags still exists")
		}
	}

	return nil
}

func testAccCheckOTCBMSTagsV2Exists(n string, tag *tags.Tags) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		bmsClient, err := config.ComputeV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud bms client: %s", err)
		}

		found, err := tags.Get(bmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		*tag = *found

		return nil
	}
}

var testAccBMSTagsV2_basic = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "BMSinstance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "physical.o2.medium"
  flavor_name = "physical.o2.medium"
  tags = ["foo","bar"]
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_compute_bms_tags_v2" "tags_1" {
  server_id = opentelekomcloud_compute_instance_v2.instance_1.id
  tags = ["foo","bar"]
}`, env.OS_IMAGE_ID, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID)

var testAccBMSTagsV2_timeout = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "BMSinstance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "physical.o2.medium"
  flavor_name = "physical.o2.medium"
  tags = ["foo","bar"]
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_compute_bms_tags_v2" "tags_1" {
  server_id = opentelekomcloud_compute_instance_v2.instance_1.id
  tags = ["foo","bar"]
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, env.OS_IMAGE_ID, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID)
