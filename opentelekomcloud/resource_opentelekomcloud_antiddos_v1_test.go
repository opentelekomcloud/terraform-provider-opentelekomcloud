package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/antiddos/v1/antiddos"
)

func TestAccAntiDdosV1_basic(t *testing.T) {
	var antiddos antiddos.GetResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAntiDdosV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiDdosV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAntiDdosV1Exists("opentelekomcloud_antiddos_v1.antiddos_1", &antiddos),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_antiddos_v1.antiddos_1", "enable_l7", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_antiddos_v1.antiddos_1", "traffic_pos_id", "1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_antiddos_v1.antiddos_1", "http_request_pos_id", "2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_antiddos_v1.antiddos_1", "cleaning_access_pos_id", "1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_antiddos_v1.antiddos_1", "app_type_id", "0"),
				),
			},
			{
				Config: testAccAntiDdosV1_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_antiddos_v1.antiddos_1", "traffic_pos_id", "2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_antiddos_v1.antiddos_1", "http_request_pos_id", "1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_antiddos_v1.antiddos_1", "cleaning_access_pos_id", "2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_antiddos_v1.antiddos_1", "app_type_id", "1"),
				),
			},
		},
	})
}

func TestAccAntiDdosV1_timeout(t *testing.T) {
	var antiddos antiddos.GetResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAntiDdosV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiDdosV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAntiDdosV1Exists("opentelekomcloud_antiddos_v1.antiddos_1", &antiddos),
				),
			},
		},
	})
}

func testAccCheckAntiDdosV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	antiddosClient, err := config.antiddosV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating antiddos client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_antiddos_v1" {
			continue
		}

		_, err := antiddos.Get(antiddosClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("antiddos still exists")
		}
	}

	return nil
}

func testAccCheckAntiDdosV1Exists(n string, ddos *antiddos.GetResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		antiddosClient, err := config.antiddosV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating antiddos client: %s", err)
		}

		found, err := antiddos.Get(antiddosClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		*ddos = *found

		return nil
	}
}

const testAccAntiDdosV1_basic = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name = "test"
    size = 8
    share_type = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_antiddos_v1" "antiddos_1" {
  floating_ip_id = "${opentelekomcloud_vpc_eip_v1.eip_1.id}"
  enable_l7 = true
  traffic_pos_id = 1
  http_request_pos_id = 2
  cleaning_access_pos_id = 1
  app_type_id = 0
}
`
const testAccAntiDdosV1_update = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name = "test"
    size = 8
    share_type = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_antiddos_v1" "antiddos_1" {
  floating_ip_id = "${opentelekomcloud_vpc_eip_v1.eip_1.id}"
  enable_l7 = true
  traffic_pos_id = 2
  http_request_pos_id = 1
  cleaning_access_pos_id = 2
  app_type_id = 1
}
`

const testAccAntiDdosV1_timeout = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name = "test"
    size = 8
    share_type = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_antiddos_v1" "antiddos_1" {
  floating_ip_id = "${opentelekomcloud_vpc_eip_v1.eip_1.id}"
  enable_l7 = true
  traffic_pos_id = 1
  http_request_pos_id = 2
  cleaning_access_pos_id = 1
  app_type_id = 0

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
