package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

const dataMemberIDsName = "data.opentelekomcloud_lb_member_ids_v2.this"

func TestAccLbMemberIDsV2DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.LoadBalancer, Count: 1},
				{Q: quotas.LbListener, Count: 1},
				{Q: quotas.LbPool, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccLbMemberIDsV2DataSourceInit,
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false, unfinished elb?
			},
			{
				Config: testAccLbMemberIDsV2DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccLbMemberIDsV3DataSourceID(dataMemberIDsName),
					resource.TestCheckResourceAttr(dataMemberIDsName, "ids.#", "1"),
				),
			},
		},
	})
}

func testAccLbMemberIDsV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find ELBv2 Member IDs data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

var testAccLbMemberIDsV2DataSourceInit = fmt.Sprintf(`
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
}

resource "opentelekomcloud_lb_member_v2" "member_2" {
  address       = "192.168.0.11"
  protocol_port = 8080
  pool_id       = opentelekomcloud_lb_pool_v2.pool_1.id
  subnet_id     = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
  weight        = 10
}
`, common.DataSourceSubnet)

var testAccLbMemberIDsV2DataSourceBasic = fmt.Sprintf(`
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
}

resource "opentelekomcloud_lb_member_v2" "member_1" {
  address       = "192.168.0.11"
  protocol_port = 8080
  pool_id       = opentelekomcloud_lb_pool_v2.pool_1.id
  subnet_id     = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
  weight        = 10
}

data "opentelekomcloud_lb_member_ids_v2" "this" {
  pool_id = opentelekomcloud_lb_pool_v2.pool_1.id
}
`, common.DataSourceSubnet)
