package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDNSV2EndpointsDataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_dns_endpoints_v2.test"
	name := common.RandomAccResourceName()
	domainName := fmt.Sprintf("%s.example.com", acctest.RandString(5))
	dc := common.InitDataSourceCheck(dataSourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2EndpointsDataSourceConfig(name, domainName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "endpoints.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "endpoints.0.id",
						"opentelekomcloud_dns_endpoint_v2.test", "id"),
					resource.TestCheckResourceAttr(dataSourceName, "endpoints.0.name", name),
					resource.TestCheckResourceAttr(dataSourceName, "endpoints.0.direction", "outbound"),
					resource.TestCheckResourceAttr(dataSourceName, "endpoints.0.status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(dataSourceName, "endpoints.0.vpc_id",
						"data.opentelekomcloud_vpc_subnet_v1.shared_subnet", "vpc_id"),
					resource.TestCheckResourceAttr(dataSourceName, "endpoints.0.ip_address_count", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "endpoints.0.resolver_rule_count", "1"),
				),
			},
		},
	})
}

func testAccDNSV2EndpointsDataSourceConfig(name, domainName string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_dns_endpoints_v2" "test" {
  depends_on = [opentelekomcloud_dns_resolver_rule_v2.test]

  direction = "outbound"
  name      = opentelekomcloud_dns_endpoint_v2.test.name
}
`, testAccDNSV2ResolverRuleConfig(name, domainName, `["1.1.1.1"]`, true))
}
