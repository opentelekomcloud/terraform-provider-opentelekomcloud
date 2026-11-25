package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/resolverrules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/dns"
)

const dnsResolverRuleResourceName = "opentelekomcloud_dns_resolver_rule_v2.test"

func TestAccDNSV2ResolverRule_basic(t *testing.T) {
	var rule resolverrules.ResolverRule
	name := common.RandomAccResourceName()
	domainName := fmt.Sprintf("%s.example.com", acctest.RandString(5))
	rc := common.InitResourceCheck(dnsResolverRuleResourceName, &rule, getDNSResolverRuleFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2ResolverRuleConfig(name, domainName, `["1.1.1.1"]`, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dnsResolverRuleResourceName, "name", name),
					resource.TestCheckResourceAttr(dnsResolverRuleResourceName, "domain_name", domainName+"."),
					resource.TestCheckResourceAttr(dnsResolverRuleResourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dnsResolverRuleResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(dnsResolverRuleResourceName, "vpc_ids.#", "1"),
					resource.TestCheckResourceAttrPair(dnsResolverRuleResourceName, "endpoint_id",
						"opentelekomcloud_dns_endpoint_v2.test", "id"),
					resource.TestCheckTypeSetElemAttrPair(dnsResolverRuleResourceName, "vpc_ids.*",
						"data.opentelekomcloud_vpc_subnet_v1.shared_subnet", "vpc_id"),
				),
			},
			{
				ResourceName:      dnsResolverRuleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDNSV2ResolverRuleConfig(name+"-upd", domainName, `["1.1.1.1", "8.8.8.8"]`, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dnsResolverRuleResourceName, "name", name+"-upd"),
					resource.TestCheckResourceAttr(dnsResolverRuleResourceName, "ip_addresses.#", "2"),
					resource.TestCheckTypeSetElemAttr(dnsResolverRuleResourceName, "ip_addresses.*", "8.8.8.8"),
					resource.TestCheckResourceAttr(dnsResolverRuleResourceName, "vpc_ids.#", "0"),
				),
			},
			{
				Config: testAccDNSV2ResolverRuleConfig(name+"-upd", domainName, `["8.8.8.8"]`, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dnsResolverRuleResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(dnsResolverRuleResourceName, "vpc_ids.#", "1"),
				),
			},
		},
	})
}

func getDNSResolverRuleFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DnsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS v2 client: %w", err)
	}
	return dns.GetDNSResolverRuleV2(client, state.Primary.ID)
}

func testAccDNSV2ResolverRuleConfig(name, domainName, ips string, withVPC bool) string {
	vpcIDs := "[]"
	if withVPC {
		vpcIDs = "[data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id]"
	}
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dns_endpoint_v2" "test" {
  name      = "%[2]s"
  direction = "outbound"

  ip_addresses {
    subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  ip_addresses {
    subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_dns_resolver_rule_v2" "test" {
  name         = "%[2]s"
  domain_name  = "%[3]s"
  endpoint_id  = opentelekomcloud_dns_endpoint_v2.test.id
  ip_addresses = %[4]s
  vpc_ids      = %[5]s
}
`, common.DataSourceSubnet, name, domainName, ips, vpcIDs)
}
