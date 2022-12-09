package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/antiddos/v1/antiddos"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

var unstableMessage = "Unstable service, tests often fails"

func TestAccAntiDdosV1_basic(t *testing.T) {
	t.Log(unstableMessage)
	supportedRegions := []string{"eu-de"}
	var antiddosElement antiddos.GetResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckServiceAvailability(t, testServiceV1, supportedRegions)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAntiDdosV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiDdosV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAntiDdosV1Exists("opentelekomcloud_antiddos_v1.antiddos_1", &antiddosElement),
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
	t.Log(unstableMessage)
	supportedRegions := []string{"eu-de"}
	var antiddosElement antiddos.GetResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckServiceAvailability(t, testServiceV1, supportedRegions)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAntiDdosV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiDdosV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAntiDdosV1Exists("opentelekomcloud_antiddos_v1.antiddos_1", &antiddosElement),
				),
			},
		},
	})
}

func TestAccAntiDdosV1_importBasic(t *testing.T) {
	t.Log(unstableMessage)
	supportedRegions := []string{"eu-de"}
	resourceName := "opentelekomcloud_antiddos_v1.antiddos_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckServiceAvailability(t, testServiceV1, supportedRegions)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAntiDdosV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiDdosV1_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAntiDdosV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	antiddosClient, err := config.AntiddosV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating antiddos client: %s", err)
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
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		antiddosClient, err := config.AntiddosV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating antiddos client: %s", err)
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
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_antiddos_v1" "antiddos_1" {
  floating_ip_id         = opentelekomcloud_vpc_eip_v1.eip_1.id
  enable_l7              = true
  traffic_pos_id         = 1
  http_request_pos_id    = 2
  cleaning_access_pos_id = 1
  app_type_id            = 0
}
`
const testAccAntiDdosV1_update = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_antiddos_v1" "antiddos_1" {
  floating_ip_id         = opentelekomcloud_vpc_eip_v1.eip_1.id
  enable_l7              = true
  traffic_pos_id         = 2
  http_request_pos_id    = 1
  cleaning_access_pos_id = 2
  app_type_id            = 1
}
`

const testAccAntiDdosV1_timeout = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_antiddos_v1" "antiddos_1" {
  floating_ip_id         = opentelekomcloud_vpc_eip_v1.eip_1.id
  enable_l7              = true
  traffic_pos_id         = 1
  http_request_pos_id    = 2
  cleaning_access_pos_id = 1
  app_type_id            = 0

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
