package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDNSV2ResolverRulesDataSource_basic(t *testing.T) {
	byEndpoint := "data.opentelekomcloud_dns_resolver_rules_v2.by_endpoint"
	byName := "data.opentelekomcloud_dns_resolver_rules_v2.by_name"
	name := common.RandomAccResourceName()
	domainName := fmt.Sprintf("%s.example.com", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2ResolverRulesDataSourceConfig(name, domainName),
				Check: resource.ComposeTestCheckFunc(
					common.InitDataSourceCheck(byEndpoint).CheckResourceExists(),
					resource.TestCheckResourceAttr(byEndpoint, "resolver_rules.#", "1"),
					resource.TestCheckResourceAttrPair(byEndpoint, "resolver_rules.0.id",
						"opentelekomcloud_dns_resolver_rule_v2.test", "id"),
					resource.TestCheckResourceAttr(byEndpoint, "resolver_rules.0.name", name),
					resource.TestCheckResourceAttr(byEndpoint, "resolver_rules.0.domain_name", domainName+"."),
					resource.TestCheckResourceAttr(byEndpoint, "resolver_rules.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(byEndpoint, "resolver_rules.0.ip_address_count", "1"),
					resource.TestCheckResourceAttr(byEndpoint, "resolver_rules.0.vpc_ids.#", "1"),
					resource.TestCheckResourceAttrPair(byEndpoint, "resolver_rules.0.vpc_ids.0",
						"data.opentelekomcloud_vpc_subnet_v1.shared_subnet", "vpc_id"),

					common.InitDataSourceCheck(byName).CheckResourceExists(),
					resource.TestCheckResourceAttr(byName, "resolver_rules.#", "1"),
					resource.TestCheckResourceAttrPair(byName, "resolver_rules.0.id",
						"opentelekomcloud_dns_resolver_rule_v2.test", "id"),
				),
			},
		},
	})
}

func testAccDNSV2ResolverRulesDataSourceConfig(name, domainName string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_dns_resolver_rules_v2" "by_endpoint" {
  depends_on = [opentelekomcloud_dns_resolver_rule_v2.test]

  endpoint_id = opentelekomcloud_dns_endpoint_v2.test.id
}

data "opentelekomcloud_dns_resolver_rules_v2" "by_name" {
  depends_on = [opentelekomcloud_dns_resolver_rule_v2.test]

  name        = "%s"
  domain_name = "%s"
}
`, testAccDNSV2ResolverRuleConfig(name, domainName, `["1.1.1.1"]`, true), name, domainName)
}
