package hwcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHWCloudNetworkingNetworkV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHWCloudNetworkingNetworkV2DataSource_network,
			},
			resource.TestStep{
				Config: testAccHWCloudNetworkingNetworkV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.hwcloud_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_networking_network_v2.net", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccHWCloudNetworkingNetworkV2DataSource_subnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHWCloudNetworkingNetworkV2DataSource_network,
			},
			resource.TestStep{
				Config: testAccHWCloudNetworkingNetworkV2DataSource_subnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.hwcloud_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_networking_network_v2.net", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccHWCloudNetworkingNetworkV2DataSource_networkID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHWCloudNetworkingNetworkV2DataSource_network,
			},
			resource.TestStep{
				Config: testAccHWCloudNetworkingNetworkV2DataSource_networkID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.hwcloud_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_networking_network_v2.net", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkingNetworkV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find network data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Network data source ID not set")
		}

		return nil
	}
}

const testAccHWCloudNetworkingNetworkV2DataSource_network = `
resource "hwcloud_networking_network_v2" "net" {
        name = "tf_test_network"
        admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet" {
  name = "tf_test_subnet"
  cidr = "192.168.199.0/24"
  no_gateway = true
  network_id = "${hwcloud_networking_network_v2.net.id}"
}
`

var testAccHWCloudNetworkingNetworkV2DataSource_basic = fmt.Sprintf(`
%s

data "hwcloud_networking_network_v2" "net" {
	name = "${hwcloud_networking_network_v2.net.name}"
}
`, testAccHWCloudNetworkingNetworkV2DataSource_network)

var testAccHWCloudNetworkingNetworkV2DataSource_subnet = fmt.Sprintf(`
%s

data "hwcloud_networking_network_v2" "net" {
	matching_subnet_cidr = "${hwcloud_networking_subnet_v2.subnet.cidr}"
}
`, testAccHWCloudNetworkingNetworkV2DataSource_network)

var testAccHWCloudNetworkingNetworkV2DataSource_networkID = fmt.Sprintf(`
%s

data "hwcloud_networking_network_v2" "net" {
	network_id = "${hwcloud_networking_network_v2.net.id}"
}
`, testAccHWCloudNetworkingNetworkV2DataSource_network)
