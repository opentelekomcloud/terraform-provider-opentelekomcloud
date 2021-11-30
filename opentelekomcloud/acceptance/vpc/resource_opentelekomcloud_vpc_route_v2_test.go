package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/routes"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceVPCRouteName = "opentelekomcloud_vpc_route_v2.route_1"

func TestAccVpcRouteV2_basic(t *testing.T) {
	var route routes.Route

	t.Parallel()
	quotas.BookMany(t, multipleRouters(2))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRouteV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteV2Exists(resourceVPCRouteName, &route),
					resource.TestCheckResourceAttr(resourceVPCRouteName, "destination", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCRouteName, "type", "peering"),
				),
			},
		},
	})
}
func TestAccVpcRouteV2_import(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, multipleRouters(2))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRouteV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteV2Import,
			},
			{
				ResourceName:      resourceVPCRouteName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVpcRouteV2_timeout(t *testing.T) {
	var route routes.Route

	t.Parallel()
	quotas.BookMany(t, multipleRouters(2))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRouteV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteV2Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteV2Exists(resourceVPCRouteName, &route),
				),
			},
		},
	})
}

func testAccCheckRouteV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_route_v2" {
			continue
		}

		_, err := routes.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("route still exists")
		}
	}

	return nil
}

func testAccCheckRouteV2Exists(n string, route *routes.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
		}

		found, err := routes.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.RouteID != rs.Primary.ID {
			return fmt.Errorf("route not found")
		}

		*route = *found

		return nil
	}
}

const testAccRouteV2Basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test1"
  cidr = "192.168.0.0/16"
}
resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name        = "opentelekomcloud_peering"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}
resource "opentelekomcloud_vpc_route_v2" "route_1" {
  type        = "peering"
  nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering_1.id
  destination = "192.168.0.0/16"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id

}
`

const testAccRouteV2Import = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test_rt_imp"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test_rt_imp1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name        = "opentelekomcloud_peering_imp"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_v2" "route_1" {
  type        = "peering"
  nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering_1.id
  destination = "192.168.0.0/16"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
}
`

const testAccRouteV2Timeout = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test_rt_t"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test_rt_t1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name        = "opentelekomcloud_peering"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_v2" "route_1" {
  type        = "peering"
  nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering_1.id
  destination = "192.168.0.0/16"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
