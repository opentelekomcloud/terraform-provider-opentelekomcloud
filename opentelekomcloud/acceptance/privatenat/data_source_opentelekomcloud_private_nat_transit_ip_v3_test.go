package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const transitIpDataSourceName = "data.opentelekomcloud_private_nat_transit_ip_v3.transit_ip_1"

func TestAccPrivateTransitIpV3DS_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateTransitIpV3DSBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateNatTransitIpV3DataSourceID(transitIpDataSourceName),
					resource.TestCheckResourceAttrSet(transitIpDataSourceName, "transit_ips.#"),
				),
			},
		},
	})
}

func testAccCheckPrivateNatTransitIpV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Private NAT Transit IP data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Private NAT Transit IP data source ID not set")
		}

		return nil
	}
}

var testAccPrivateTransitIpV3DSBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_test" {
  virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}

data "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_1" {
  depends_on = ["opentelekomcloud_private_nat_transit_ip_v3.transit_ip_test"]
}
`, common.DataSourceSubnet)

func TestAccPrivateTransitIpV3DS_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateTransitIpV3DS_id,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(transitIpDataSourceName, "transit_ips.0.tags.kuh", "muh"),
				),
			},
		},
	})
}

var testAccPrivateTransitIpV3DS_id = fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_test" {
  virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  tags = {
    kuh = "muh"
  }
}

data "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_1" {
  depends_on = ["opentelekomcloud_private_nat_transit_ip_v3.transit_ip_test"]
  id         = opentelekomcloud_private_nat_transit_ip_v3.transit_ip_test.id
}
`, common.DataSourceSubnet)

func TestAccPrivateTransitIpV3DS_ipaddr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateTransitIpV3DS_ipaddr,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(transitIpDataSourceName, "transit_ips.0.tags.kuh", "nuh"),
				),
			},
		},
	})
}

var testAccPrivateTransitIpV3DS_ipaddr = fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_test" {
  virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  tags = {
    kuh = "nuh"
  }
}

data "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_1" {
  depends_on = ["opentelekomcloud_private_nat_transit_ip_v3.transit_ip_test"]
  ip_address = opentelekomcloud_private_nat_transit_ip_v3.transit_ip_test.ip_address
}
`, common.DataSourceSubnet)

func TestAccPrivateTransitIpV3DS_subnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateTransitIpV3DS_subnet,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(transitIpDataSourceName, "transit_ips.0.tags.kuh", "puh"),
				),
			},
		},
	})
}

var testAccPrivateTransitIpV3DS_subnet = fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_test" {
  virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  tags = {
    kuh = "puh"
  }
}

data "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_1" {
  depends_on   = ["opentelekomcloud_private_nat_transit_ip_v3.transit_ip_test"]
  virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}
`, common.DataSourceSubnet)
