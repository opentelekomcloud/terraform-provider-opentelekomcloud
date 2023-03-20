package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataSourceListenerName = "data.opentelekomcloud_lb_listener_v3.listener"

func TestDataSourceListenerV3_basic(t *testing.T) {
	t.Parallel()
	qts := quotas.MultipleQuotas{
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.LbListener, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceListenerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceListenerName, "protocol_port", "80"),
					resource.TestCheckResourceAttrSet(dataSourceListenerName, "loadbalancer_id"),
					resource.TestCheckResourceAttr(dataSourceListenerName, "advanced_forwarding", "false"),
					resource.TestCheckResourceAttr(dataSourceListenerName, "sni_match_algo", ""),
					resource.TestCheckResourceAttr(dataSourceListenerName, "security_policy_id", ""),
				),
			},
		},
	})
}

func TestDataSourceListenerV3_byID(t *testing.T) {
	t.Parallel()
	qts := quotas.MultipleQuotas{
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.LbListener, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceListenerConfigID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceListenerName, "protocol_port", "443"),
				),
			},
		},
	})
}

var testDataSourceListenerConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name            = "%s"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol        = "HTTP"
  protocol_port   = 80
}

data "opentelekomcloud_lb_listener_v3" "listener" {
  loadbalancer_id = opentelekomcloud_lb_listener_v3.listener_1.loadbalancer_id
  name            = "%[3]s"
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, tools.RandomString("lst-", 4))

var testDataSourceListenerConfigID = fmt.Sprintf(`
%s

data "opentelekomcloud_lb_listener_v3" "listener" {
  id = opentelekomcloud_lb_listener_v3.listener_1.id
}
`, testAccLBV3ListenerConfigBasic)
