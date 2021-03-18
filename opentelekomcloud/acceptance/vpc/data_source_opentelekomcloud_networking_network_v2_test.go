package acceptance

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccNetworkingNetworkV2DataSource_basic(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	network := fmt.Sprintf("acc_test_network-%06x", rand.Int31n(1000000))
	cidr := fmt.Sprintf("192.168.%d.0/24", rand.Intn(200))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { common.TestAccPreCheck(t) },
		Providers: common.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkV2DataSource_basic(network, cidr),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingNetworkV2DataSourceID("data.opentelekomcloud_networking_network_v2.net_by_name"),
					TestAccCheckNetworkingNetworkV2DataSourceID("data.opentelekomcloud_networking_network_v2.net_by_id"),
					TestAccCheckNetworkingNetworkV2DataSourceID("data.opentelekomcloud_networking_network_v2.net_by_cidr"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_network_v2.net_by_name", "name", network),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_network_v2.net_by_id", "name", network),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_network_v2.net_by_cidr", "name", network),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_network_v2.net_by_name", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_network_v2.net_by_id", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_network_v2.net_by_cidr", "matching_subnet_cidr", cidr),
				),
			},
		},
	})
}

func testAccNetworkingNetworkV2DataSource_basic(name, cidr string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_network_v2" "net" {
  name = "%s"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet" {
  name = "opentelekomcloud_test_subnet"
  cidr = "%s"
  no_gateway = true
  network_id = opentelekomcloud_networking_network_v2.net.id
}

data "opentelekomcloud_networking_network_v2" "net_by_name" {
	name = opentelekomcloud_networking_network_v2.net.name
}

data "opentelekomcloud_networking_network_v2" "net_by_id" {
	network_id = opentelekomcloud_networking_network_v2.net.id
}

data "opentelekomcloud_networking_network_v2" "net_by_cidr" {
	matching_subnet_cidr = opentelekomcloud_networking_subnet_v2.subnet.cidr
}
`, name, cidr)
}
