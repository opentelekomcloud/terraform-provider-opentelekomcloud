package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/networks"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccNetworkingV2RouterRoute_basic(t *testing.T) {
	var router routers.Router
	var network [2]networks.Network
	var subnet [2]subnets.Subnet

	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.Router, Count: 1},
		{Q: quotas.Network, Count: 2},
		{Q: quotas.Subnet, Count: 2},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterRouteCreate,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists("opentelekomcloud_networking_router_v2.router_1", &router),
					TestAccCheckNetworkingV2NetworkExists(
						"opentelekomcloud_networking_network_v2.network_1", &network[0]),
					TestAccCheckNetworkingV2SubnetExists(
						"opentelekomcloud_networking_subnet_v2.subnet_1", &subnet[0]),
					TestAccCheckNetworkingV2NetworkExists(
						"opentelekomcloud_networking_network_v2.network_1", &network[1]),
					TestAccCheckNetworkingV2SubnetExists(
						"opentelekomcloud_networking_subnet_v2.subnet_1", &subnet[1]),
					TestAccCheckNetworkingV2RouterInterfaceExists(
						"opentelekomcloud_networking_router_interface_v2.int_1"),
					TestAccCheckNetworkingV2RouterInterfaceExists(
						"opentelekomcloud_networking_router_interface_v2.int_2"),
					testAccCheckNetworkingV2RouterRouteExists(
						"opentelekomcloud_networking_router_route_v2.router_route_1"),
				),
			},
			{
				Config: testAccNetworkingV2RouterRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterRouteExists(
						"opentelekomcloud_networking_router_route_v2.router_route_1"),
					testAccCheckNetworkingV2RouterRouteExists(
						"opentelekomcloud_networking_router_route_v2.router_route_2"),
				),
			},
			{
				Config: testAccNetworkingV2RouterRouteDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterRouteEmpty("opentelekomcloud_networking_router_v2.router_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterRoute_ecs(t *testing.T) {
	resourceName := "opentelekomcloud_networking_router_route_v2.router_route_1"
	name := fmt.Sprintf("router_acc_route%s", acctest.RandString(10))
	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.Router, Count: 1},
		{Q: quotas.Network, Count: 1},
		{Q: quotas.Subnet, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterRouteEcs(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterRouteExists(resourceName),
				),
			},
			{
				Config: testAccNetworkingV2RouterRouteEcsUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterRouteExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterRouteEmpty(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		router, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if router.ID != rs.Primary.ID {
			return fmt.Errorf("router not found")
		}

		if len(router.Routes) != 0 {
			return fmt.Errorf("invalid number of route entries: %d", len(router.Routes))
		}

		return nil
	}
}

func testAccCheckNetworkingV2RouterRouteExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		router, err := routers.Get(networkingClient, rs.Primary.Attributes["router_id"]).Extract()
		if err != nil {
			return err
		}

		if router.ID != rs.Primary.Attributes["router_id"] {
			return fmt.Errorf("router for route not found")
		}

		var found bool
		for _, r := range router.Routes {
			if r.DestinationCIDR == rs.Primary.Attributes["destination_cidr"] && r.NextHop == rs.Primary.Attributes["next_hop"] {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("could not find route for destination CIDR: %s, next hop: %s", rs.Primary.Attributes["destination_cidr"], rs.Primary.Attributes["next_hop"])
		}

		return nil
	}
}

const testAccNetworkingV2RouterRouteCreate = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1_rr"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1_rr"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.1"
  }
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_networking_network_v2" "network_2" {
  name           = "network_2_rr"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_2" {
  cidr       = "192.168.200.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_2.id
}

resource "opentelekomcloud_networking_port_v2" "port_2" {
  name           = "port_2"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_2.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_2.id
    ip_address = "192.168.200.1"
  }
}

resource "opentelekomcloud_networking_router_interface_v2" "int_2" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_2.id
}

resource "opentelekomcloud_networking_router_route_v2" "router_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop         = "192.168.199.250"

  depends_on = ["opentelekomcloud_networking_router_interface_v2.int_1"]
  router_id  = opentelekomcloud_networking_router_v2.router_1.id
}
`

const testAccNetworkingV2RouterRouteUpdate = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1_rr"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1_rr"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.1"
  }
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_networking_network_v2" "network_2" {
  name           = "network_2_rr"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_2" {
  cidr       = "192.168.200.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_2.id
}

resource "opentelekomcloud_networking_port_v2" "port_2" {
  name           = "port_2"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_2.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_2.id
    ip_address = "192.168.200.1"
  }
}

resource "opentelekomcloud_networking_router_interface_v2" "int_2" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_2.id
}

resource "opentelekomcloud_networking_router_route_v2" "router_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop         = "192.168.199.250"

  depends_on = ["opentelekomcloud_networking_router_interface_v2.int_1"]
  router_id  = opentelekomcloud_networking_router_v2.router_1.id
}

resource "opentelekomcloud_networking_router_route_v2" "router_route_2" {
  destination_cidr = "10.0.2.0/24"
  next_hop         = "192.168.200.250"

  depends_on = ["opentelekomcloud_networking_router_interface_v2.int_2"]
  router_id  = opentelekomcloud_networking_router_v2.router_1.id
}
`

const testAccNetworkingV2RouterRouteDestroy = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1_rr"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1_rr"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.1"
  }
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_networking_network_v2" "network_2" {
  name           = "network_2_rr"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_2" {
  ip_version = 4
  cidr       = "192.168.200.0/24"
  network_id = opentelekomcloud_networking_network_v2.network_2.id
}

resource "opentelekomcloud_networking_port_v2" "port_2" {
  name           = "port_2"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_2.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_2.id
    ip_address = "192.168.200.1"
  }
}

resource "opentelekomcloud_networking_router_interface_v2" "int_2" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_2.id
}
`

func testAccNetworkingV2RouterRouteEcs(name string) string {
	return fmt.Sprintf(`
%[3]s

data "opentelekomcloud_images_image_v2" "other_image" {
  name        = "Standard_Debian_12_amd64_bios_latest"
  most_recent = true
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "%[1]s_router"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "%[1]s_network"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name       = "%[1]s_port"
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  }
}

resource "opentelekomcloud_networking_port_v2" "instance_port_1" {
  name       = "%[1]s_port"
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "%[1]s_instance"
  security_groups   = ["default"]
  availability_zone = "%[2]s"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id

  network {
    port = opentelekomcloud_networking_port_v2.instance_port_1.id
  }
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_networking_router_route_v2" "router_route_1" {
  destination_cidr = "192.168.254.254/32"
  next_hop         = opentelekomcloud_compute_instance_v2.instance_1.network[0].fixed_ip_v4

  depends_on = ["opentelekomcloud_networking_router_interface_v2.int_1"]
  router_id  = opentelekomcloud_networking_router_v2.router_1.id
}
`, name, env.OS_AVAILABILITY_ZONE, common.DataSourceImage)
}

func testAccNetworkingV2RouterRouteEcsUpdate(name string) string {
	return fmt.Sprintf(`
%[3]s


data "opentelekomcloud_images_image_v2" "other_image" {
  name        = "Standard_Debian_12_amd64_bios_latest"
  most_recent = true
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "%[1]s_router"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "%[1]s_network"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name       = "%[1]s_port"
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  }
}

resource "opentelekomcloud_networking_port_v2" "instance_port_1" {
  name       = "%[1]s_port"
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "%[1]s_instance"
  security_groups   = ["default"]
  availability_zone = "%[2]s"
  image_id          = data.opentelekomcloud_images_image_v2.other_image.id

  network {
    port = opentelekomcloud_networking_port_v2.instance_port_1.id
  }
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_networking_router_route_v2" "router_route_1" {
  destination_cidr = "192.168.254.254/32"
  next_hop         = opentelekomcloud_compute_instance_v2.instance_1.network[0].fixed_ip_v4

  depends_on = ["opentelekomcloud_networking_router_interface_v2.int_1"]
  router_id  = opentelekomcloud_networking_router_v2.router_1.id
}
`, name, env.OS_AVAILABILITY_ZONE, common.DataSourceImage)
}
