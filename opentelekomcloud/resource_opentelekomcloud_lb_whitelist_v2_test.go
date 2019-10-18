package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/whitelists"
)

func TestAccLBV2Whitelist_basic(t *testing.T) {
	var whitelist whitelists.Whitelist

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2WhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLBV2WhitelistConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2WhitelistExists("opentelekomcloud_lb_whitelist_v2.whitelist_1", &whitelist),
				),
			},
			{
				Config: TestAccLBV2WhitelistConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_lb_whitelist_v2.whitelist_1", "enable_whitelist", "true"),
				),
			},
		},
	})
}

func testAccCheckLBV2WhitelistDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_whitelist_v2" {
			continue
		}

		_, err := whitelists.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Whitelist still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2WhitelistExists(n string, whitelist *whitelists.Whitelist) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := whitelists.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Whitelist not found")
		}

		*whitelist = *found

		return nil
	}
}

var TestAccLBV2WhitelistConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "opentelekomcloud_lb_whitelist_v2" "whitelist_1" {
  enable_whitelist = true
  whitelist = "192.168.11.1,192.168.0.1/24"
  listener_id = "${opentelekomcloud_lb_listener_v2.listener_1.id}"
}
`, OS_SUBNET_ID)

var TestAccLBV2WhitelistConfig_update = fmt.Sprintf(`
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "opentelekomcloud_lb_whitelist_v2" "whitelist_1" {
  enable_whitelist = true
  whitelist = "192.168.11.1,192.168.0.1/24,192.168.201.18/8"
  listener_id = "${opentelekomcloud_lb_listener_v2.listener_1.id}"
}
`, OS_SUBNET_ID)
