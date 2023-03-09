package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/pools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
)

const resourcePoolName = "opentelekomcloud_lb_pool_v3.pool"

func TestLBPoolV3_basic(t *testing.T) {
	var pool pools.Pool
	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testLBPoolV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testLBPoolV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testLBPoolV3Exists(resourcePoolName, &pool),
					resource.TestCheckResourceAttr(resourcePoolName, "name", ""),
					resource.TestCheckResourceAttr(resourcePoolName, "session_persistence.#", "1"),
					resource.TestCheckResourceAttr(resourcePoolName, "ip_version", "dualstack"),
					resource.TestCheckResourceAttr(resourcePoolName, "type", "instance"),
					resource.TestCheckResourceAttr(resourcePoolName, "member_deletion_protection", "true"),
				),
			},
			{
				Config: testLBPoolV3Updated,
				Check: resource.ComposeTestCheckFunc(
					testLBPoolV3Exists(resourcePoolName, &pool),
					resource.TestCheckResourceAttr(resourcePoolName, "name", "pool_1"),
					resource.TestCheckResourceAttr(resourcePoolName, "member_deletion_protection", "false"),
				),
			},
			{
				Config: testLBPoolV3HTTPSBasic,
				Check: resource.ComposeTestCheckFunc(
					testLBPoolV3Exists(resourcePoolName, &pool),
					resource.TestCheckResourceAttr(resourcePoolName, "name", ""),
					resource.TestCheckResourceAttr(resourcePoolName, "protocol", "HTTPS"),
				),
			},
		},
	})
}

func TestLBPoolV3_import(t *testing.T) {
	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testLBPoolV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testLBPoolV3Basic,
			},
			{
				ResourceName:      resourcePoolName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testLBPoolV3Exists(n string, pool *pools.Pool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ElbV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(elbv3.ErrCreateClient, err)
		}

		found, err := pools.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("loadbalancer pool not found")
		}

		*pool = *found

		return nil
	}
}

func testLBPoolV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elbv3.ErrCreateClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_pool_v3" {
			continue
		}

		_, err := pools.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("LB Pool still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

var (
	testLBPoolV3Basic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  lb_algorithm    = "QUIC_CID"
  protocol        = "QUIC"

  session_persistence {
    type = "SOURCE_IP"
  }
  type   = "instance"
  vpc_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  member_deletion_protection = true
}

`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
	testLBPoolV3Updated = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  name            = "pool_1"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  lb_algorithm    = "QUIC_CID"
  protocol        = "QUIC"

  session_persistence {
    type                = "SOURCE_IP"
    persistence_timeout = "30"
  }

  member_deletion_protection = false
}

`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

	testLBPoolV3HTTPSBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTPS"

  session_persistence {
    type = "HTTP_COOKIE"
  }

  member_deletion_protection = false
}

`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
)
