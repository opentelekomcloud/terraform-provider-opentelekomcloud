package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/pools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourcePoolName = "opentelekomcloud_lb_pool_v2.pool_1"

func TestAccLBV2Pool_basic(t *testing.T) {
	var pool pools.Pool

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV2PoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2PoolConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2PoolExists(resourcePoolName, &pool),
					resource.TestCheckResourceAttr(resourcePoolName, "name", "pool_1"),
					resource.TestCheckResourceAttr(resourcePoolName, "lb_method", "ROUND_ROBIN"),
				),
			},
			{
				Config: testAccLBV2PoolConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePoolName, "name", "pool_1_updated"),
					resource.TestCheckResourceAttr(resourcePoolName, "lb_method", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttr(resourcePoolName, "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccLBV2Pool_persistenceNull(t *testing.T) {
	var pool pools.Pool
	resourceName := "opentelekomcloud_lb_pool_v2.pool_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV2PoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2PoolConfigPersistence,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2PoolExists(resourceName, &pool),
					resource.TestCheckResourceAttr(resourceName, "name", "pool_1"),
				),
			},
		},
	})
}

func testAccCheckLBV2PoolDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elb.ErrCreationV2Client, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_pool_v2" {
			continue
		}

		_, err := pools.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("pool still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2PoolExists(n string, pool *pools.Pool) resource.TestCheckFunc {
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

		found, err := pools.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("member not found")
		}

		*pool = *found

		return nil
	}
}

var testAccLBV2PoolConfigBasic = fmt.Sprintf(`
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

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, common.DataSourceSubnet)

var testAccLBV2PoolConfigUpdate = fmt.Sprintf(`
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

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name           = "pool_1_updated"
  protocol       = "HTTP"
  lb_method      = "LEAST_CONNECTIONS"
  admin_state_up = "true"
  listener_id    = opentelekomcloud_lb_listener_v2.listener_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, common.DataSourceSubnet)

var testAccLBV2PoolConfigPersistence = fmt.Sprintf(`
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

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_1.id

  persistence {
    type        = null
    cookie_name = null
  }
}
`, common.DataSourceSubnet)
