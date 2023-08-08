package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/whitelists"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	elb "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v2"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceWhitelistName = "opentelekomcloud_lb_whitelist_v2.whitelist_1"

func TestAccLBV2Whitelist_basic(t *testing.T) {
	var whitelist whitelists.Whitelist

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.LoadBalancer, Count: 1},
				{Q: quotas.LbListener, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV2WhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2WhitelistConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2WhitelistExists(resourceWhitelistName, &whitelist),
				),
			},
			{
				Config: testAccLBV2WhitelistConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceWhitelistName, "enable_whitelist", "true"),
				),
			},
		},
	})
}

func TestAccLBV2Whitelist_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV2WhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2WhitelistConfigBasic,
			},
			{
				ResourceName:      resourceWhitelistName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLBV2WhitelistDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elb.ErrCreationV2Client, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_whitelist_v2" {
			continue
		}

		_, err := whitelists.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("whitelist still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2WhitelistExists(n string, whitelist *whitelists.Whitelist) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ElbV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(elb.ErrCreationV2Client, err)
		}

		found, err := whitelists.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("whitelist not found")
		}

		*whitelist = *found

		return nil
	}
}

var testAccLBV2WhitelistConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_whitelist_v2" "whitelist_1" {
  enable_whitelist = true
  whitelist        = "192.168.11.1,192.168.0.1/24"
  listener_id      = opentelekomcloud_lb_listener_v2.listener_1.id
}
`, common.DataSourceSubnet)

var testAccLBV2WhitelistConfigUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_whitelist_v2" "whitelist_1" {
  enable_whitelist = true
  whitelist        = "192.168.11.1,192.168.0.1/24,192.168.201.18/8"
  listener_id      = opentelekomcloud_lb_listener_v2.listener_1.id
}
`, common.DataSourceSubnet)
