package acceptance

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/monitors"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
)

const resourceMonitorName = "opentelekomcloud_lb_monitor_v3.monitor"

func TestResourceMonitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
	}
	th.AssertNoErr(t, quotas.AcquireMultipleQuotas(qts, 5*time.Second))
	defer quotas.ReleaseMultipleQuotas(qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      checkMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceMonitorBasic,
				Check: resource.ComposeTestCheckFunc(
					checkMonitorExists(resourceMonitorName, &monitor),
					resource.TestCheckResourceAttr(resourceMonitorName, "url_path", "/"),
					resource.TestCheckResourceAttr(resourceMonitorName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceMonitorName, "expected_codes", "200"),
				),
			},
			{
				Config: testResourceMonitorUpdated,
				Check: resource.ComposeTestCheckFunc(
					checkMonitorExists(resourceMonitorName, &monitor),
					resource.TestCheckResourceAttr(resourceMonitorName, "type", "TCP"),
					resource.TestCheckResourceAttr(resourceMonitorName, "max_retries", "2"),
					resource.TestCheckResourceAttr(resourceMonitorName, "max_retries_down", "1"),
				),
			},
		},
	})
}

func TestResourceMonitor_import(t *testing.T) {
	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
	}
	th.AssertNoErr(t, quotas.AcquireMultipleQuotas(qts, 5*time.Second))
	defer quotas.ReleaseMultipleQuotas(qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      checkMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceMonitorBasic,
			},
			{
				ResourceName:      resourceMonitorName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkMonitorDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elbv3.ErrCreateClient, err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_monitor_v3" {
			continue
		}

		_, err := monitors.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("monitor still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func checkMonitorExists(n string, monitor *monitors.Monitor) resource.TestCheckFunc {
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

		found, err := monitors.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("monitor not found")
		}

		monitor = found

		return nil
	}
}

var testResourceMonitorBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  protocol        = "HTTP"
  lb_algorithm    = "ROUND_ROBIN"
}

resource "opentelekomcloud_lb_monitor_v3" "monitor" {
  pool_id      = opentelekomcloud_lb_pool_v3.pool.id
  type         = "HTTP"
  delay        = 3
  timeout      = 30
  monitor_port = 8080

  max_retries      = 5
  max_retries_down = 1
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testResourceMonitorUpdated = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  protocol        = "HTTP"
  lb_algorithm    = "ROUND_ROBIN"
}

resource "opentelekomcloud_lb_monitor_v3" "monitor" {
  pool_id      = opentelekomcloud_lb_pool_v3.pool.id
  type         = "TCP"
  delay        = 1
  timeout      = 5
  monitor_port = 8080

  max_retries      = 2
  max_retries_down = 1
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
