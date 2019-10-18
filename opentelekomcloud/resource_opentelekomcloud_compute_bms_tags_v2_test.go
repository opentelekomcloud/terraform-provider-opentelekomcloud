package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/bms/v2/tags"
)

func TestAccOTCBMSTagsV2_basic(t *testing.T) {
	var tags tags.Tags

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccBmsFlavorPreCheck(t) },
		Providers:    testAccProviders,
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
		Providers:    testAccProviders,
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
	config := testAccProvider.Meta().(*Config)
	bmsClient, err := config.bmsClient(OS_REGION_NAME)
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

		config := testAccProvider.Meta().(*Config)
		bmsClient, err := config.bmsClient(OS_REGION_NAME)
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
  server_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
  tags = ["foo","bar"]
}`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

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
  server_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
  tags = ["foo","bar"]
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
