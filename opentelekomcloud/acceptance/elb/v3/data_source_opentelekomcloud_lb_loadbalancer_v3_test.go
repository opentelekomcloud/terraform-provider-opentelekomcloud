package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataLBName = "data.opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1"

func TestLoadBalancerV3_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testLoadBalancerV3Init,
			},
			{
				Config: testLoadBalancerV3ByID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataLBName, "id"),
					resource.TestCheckResourceAttrSet(dataLBName, "router_id"),
					resource.TestCheckResourceAttr(dataLBName, "availability_zones.#", "1"),
					resource.TestCheckResourceAttr(dataLBName, "name", "loadbalancer_1"),
				),
			},
			{
				Config: testLoadBalancerV3ByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataLBName, "id"),
					resource.TestCheckResourceAttrSet(dataLBName, "router_id"),
					resource.TestCheckResourceAttr(dataLBName, "availability_zones.#", "1"),
					resource.TestCheckResourceAttr(dataLBName, "name", "loadbalancer_1"),
				),
			},
		},
	})
}

var testLoadBalancerV3Init = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testLoadBalancerV3ByID = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

data "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testLoadBalancerV3ByName = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

data "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.name
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
