package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/bms/v2/tags"
)

func TestAccOTCBMSTagsV2_basic(t *testing.T) {
	var tags tags.Tags

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckRequiredEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCBMSTagsV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
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

// PASS
func TestAccOTCBMSTagsV2_timeout(t *testing.T) {
	var tags tags.Tags

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckRequiredEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCBMSTagsV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
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
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_otc12391"
  cidr = "192.168.0.0/16"
}
resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "sub_otc12391"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  availability_zone = "%s"
}
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "BMSinstance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "physical.o2.medium"
  tags = ["foo","bar"]
  metadata {
    foo = "bar"
  }
  network {
    uuid = "${opentelekomcloud_vpc_subnet_v1.subnet_1.id}"
  }
}
resource "opentelekomcloud_compute_bms_tags_v2" "tags_1" {
  server_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
  tags = ["foo","bar"]
}`, OS_AVAILABILITY_ZONE, OS_IMAGE_ID, OS_AVAILABILITY_ZONE)

var testAccBMSTagsV2_timeout = fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_otc12391"
  cidr = "192.168.0.0/16"
}
resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "sub_otc12391"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  availability_zone = "%s"
}
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "BMSinstance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "physical.o2.medium"
  tags = ["foo","bar"]
  metadata {
    foo = "bar"
  }
  network {
    uuid = "${opentelekomcloud_vpc_subnet_v1.subnet_1.id}"
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
`, OS_AVAILABILITY_ZONE, OS_IMAGE_ID, OS_AVAILABILITY_ZONE)
