package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acls "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/acl"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getAclPolicyFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return acls.Get(client, state.Primary.Attributes["gateway_id"], state.Primary.ID)
}

func TestAccAclPolicy_basic(t *testing.T) {
	var (
		policy acls.ApiAcl

		rName1 = "opentelekomcloud_apigw_acl_policy_v2.ip_rule"
		rName2 = "opentelekomcloud_apigw_acl_policy_v2.domain_rule"
		rName3 = "opentelekomcloud_apigw_acl_policy_v2.domain_id_rule"
		name   = fmt.Sprintf("tf_test_acl_%s", acctest.RandString(3))

		basicDomainNames  = strings.Join(common.GenerateRandomDomain(2, 4), ",")
		updateDomainNames = strings.Join(common.GenerateRandomDomain(2, 4), ",")
		basicDomainIds    = strings.Join(common.GenerateRandomDomain(2, 32), ",")
		updateDomainIds   = strings.Join(common.GenerateRandomDomain(2, 32), ",")

		rc1 = common.InitResourceCheck(rName1, &policy, getAclPolicyFunc)
		rc2 = common.InitResourceCheck(rName2, &policy, getAclPolicyFunc)
		rc3 = common.InitResourceCheck(rName3, &policy, getAclPolicyFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApigAclPolicy_basic(name, basicDomainNames, basicDomainIds),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_rule_ip"),
					resource.TestCheckResourceAttr(rName1, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName1, "entity_type", "IP"),
					resource.TestCheckResourceAttr(rName1, "value", "10.201.33.4,10.30.2.15"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_rule_domain"),
					resource.TestCheckResourceAttr(rName2, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName2, "entity_type", "DOMAIN"),
					resource.TestCheckResourceAttr(rName2, "value", basicDomainNames),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName3, "name", name+"_rule_domain_id"),
					resource.TestCheckResourceAttr(rName3, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName3, "entity_type", "DOMAIN_ID"),
					resource.TestCheckResourceAttr(rName3, "value", basicDomainIds),
				),
			},
			{
				Config: testAccApigAclPolicy_update(name, updateDomainNames, updateDomainIds),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_rule_ip_update"),
					resource.TestCheckResourceAttr(rName1, "type", "DENY"),
					resource.TestCheckResourceAttr(rName1, "entity_type", "IP"),
					resource.TestCheckResourceAttr(rName1, "value", "10.201.33.8,10.30.2.23"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_rule_domain_update"),
					resource.TestCheckResourceAttr(rName2, "type", "DENY"),
					resource.TestCheckResourceAttr(rName2, "entity_type", "DOMAIN"),
					resource.TestCheckResourceAttr(rName2, "value", updateDomainNames),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName3, "name", name+"_rule_domain_id_update"),
					resource.TestCheckResourceAttr(rName3, "type", "DENY"),
					resource.TestCheckResourceAttr(rName3, "entity_type", "DOMAIN_ID"),
					resource.TestCheckResourceAttr(rName3, "value", updateDomainIds),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyImportStateFunc(rName2),
			},
			{
				ResourceName:      rName3,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyImportStateFunc(rName3),
			},
		},
	})
}

func testAccAclPolicyImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["gateway_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<insgateway_idtance_id>/<id>', but '%s/%s'",
				rs.Primary.Attributes["gateway_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.ID), nil
	}
}

func testAccApigAclPolicy_base(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_apigw_gateway_v2" "gateway"{
  name                    		  = "%s"
  spec_id                 		  = "BASIC"
  vpc_id                  		  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               		  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id       		  = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zones      		  = ["eu-de-01", "eu-de-02"]
  description             		  = "test gateway"
  ingress_bandwidth_size          = 5
  ingress_bandwidth_charging_mode = "bandwidth"
  maintain_begin                  = "22:00:00"
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, name)
}

func testAccApigAclPolicy_basic(name, domainNames, domainIds string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_acl_policy_v2" "ip_rule" {
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%[2]s_rule_ip"
  type        = "PERMIT"
  entity_type = "IP"
  value       = "10.201.33.4,10.30.2.15"
}

resource "opentelekomcloud_apigw_acl_policy_v2" "domain_rule" {
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%[2]s_rule_domain"
  type        = "PERMIT"
  entity_type = "DOMAIN"
  value       = "%[3]s"
}

resource "opentelekomcloud_apigw_acl_policy_v2" "domain_id_rule" {
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%[2]s_rule_domain_id"
  type        = "PERMIT"
  entity_type = "DOMAIN_ID"
  value       = "%[4]s"
}
`, testAccApigAclPolicy_base(name), name, domainNames, domainIds)
}

func testAccApigAclPolicy_update(name, domainNames, domainIds string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_acl_policy_v2" "ip_rule" {
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%[2]s_rule_ip_update"
  type        = "DENY"
  entity_type = "IP"
  value       = "10.201.33.8,10.30.2.23"
}

resource "opentelekomcloud_apigw_acl_policy_v2" "domain_rule" {
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%[2]s_rule_domain_update"
  type        = "DENY"
  entity_type = "DOMAIN"
  value       = "%[3]s"
}

resource "opentelekomcloud_apigw_acl_policy_v2" "domain_id_rule" {
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%[2]s_rule_domain_id_update"
  type        = "DENY"
  entity_type = "DOMAIN_ID"
  value       = "%[4]s"
}
`, testAccApigAclPolicy_base(name), name, domainNames, domainIds)
}
