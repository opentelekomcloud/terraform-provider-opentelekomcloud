package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/endpoints"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const dnsEndpointResourceName = "opentelekomcloud_dns_endpoint_v2.test"

func TestAccDNSV2Endpoint_basic(t *testing.T) {
	var ep endpoints.Endpoint
	name := common.RandomAccResourceName()
	rc := common.InitResourceCheck(dnsEndpointResourceName, &ep, getDNSEndpointFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2EndpointConfig(name, 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dnsEndpointResourceName, "name", name),
					resource.TestCheckResourceAttr(dnsEndpointResourceName, "direction", "inbound"),
					resource.TestCheckResourceAttr(dnsEndpointResourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dnsEndpointResourceName, "ip_addresses.#", "2"),
					resource.TestCheckResourceAttrSet(dnsEndpointResourceName, "ip_addresses.0.ip"),
					resource.TestCheckResourceAttrSet(dnsEndpointResourceName, "ip_addresses.0.ip_address_id"),
					resource.TestCheckResourceAttrPair(dnsEndpointResourceName, "vpc_id",
						"data.opentelekomcloud_vpc_subnet_v1.shared_subnet", "vpc_id"),
				),
			},
			{
				ResourceName:      dnsEndpointResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDNSV2EndpointConfig(name+"-upd", 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dnsEndpointResourceName, "name", name+"-upd"),
					resource.TestCheckResourceAttr(dnsEndpointResourceName, "ip_addresses.#", "3"),
					resource.TestCheckResourceAttrSet(dnsEndpointResourceName, "ip_addresses.2.ip_address_id"),
				),
			},
			{
				Config: testAccDNSV2EndpointConfig(name+"-upd", 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dnsEndpointResourceName, "ip_addresses.#", "2"),
				),
			},
		},
	})
}

func getDNSEndpointFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DnsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS v2 client: %w", err)
	}
	return endpoints.Get(client, state.Primary.ID)
}

func testAccDNSV2EndpointConfig(name string, ipCount int) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dns_endpoint_v2" "test" {
  name      = "%s"
  direction = "inbound"

  dynamic "ip_addresses" {
    for_each = range(%d)
    content {
      subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
    }
  }
}
`, common.DataSourceSubnet, name, ipCount)
}
