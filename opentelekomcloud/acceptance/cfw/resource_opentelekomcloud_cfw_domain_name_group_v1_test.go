package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	group "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/dns"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const domainNameResourceName = "opentelekomcloud_cfw_domain_name_group_v1.group_1"

func getDomainNameFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CfwV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CFW client: %s", err)
	}
	return group.GetDomainNameGroup(client, state.Primary.Attributes["name"], state.Primary.Attributes["firewall_id"], state.Primary.Attributes["object_id"])
}

func TestAccCFWDomainNameGroupV1_basic(t *testing.T) {
	var group group.DomainSetVO
	rc := common.InitResourceCheck(
		domainNameResourceName,
		&group,
		getDomainNameFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCFWDomainNameV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(domainNameResourceName, "name", "test-acc-tf-domain-group"),
					resource.TestCheckResourceAttrSet(domainNameResourceName, "id"),
				),
			},
			{
				Config: testAccCFWDomainNameV1Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(domainNameResourceName, "name", "test-acc-tf-domain-group-updated"),
					resource.TestCheckResourceAttr(domainNameResourceName, "domain_names.0.domain_name", "www.testaccupdated.com"),
				),
			},
			{
				ResourceName:      domainNameResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCFWDomainNameGroupV1ImportStateIdFunc(),
			},
		},
	})
}

func testAccCFWDomainNameGroupV1ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var firewallId, objectId, name string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_cfw_domain_name_group_v1" {
				firewallId = rs.Primary.Attributes["firewall_id"]
				objectId = rs.Primary.Attributes["object_id"]
				name = rs.Primary.Attributes["name"]
			}
		}
		if firewallId == "" || objectId == "" || name == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", firewallId, objectId, name)
		}
		return fmt.Sprintf("%s/%s/%s", firewallId, objectId, name), nil
	}
}

var testAccCFWDomainNameV1Basic = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_domain_name_group_v1" "group_1" {
  firewall_id = opentelekomcloud_cfw_firewall_v1.firewall_1.id
  object_id   = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  name        = "test-acc-tf-domain-group"
  domain_names {
    domain_name = "www.testacctf.com"
  }
}
`

var testAccCFWDomainNameV1Update = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_domain_name_group_v1" "group_1" {
  firewall_id = opentelekomcloud_cfw_firewall_v1.firewall_1.id
  object_id   = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  name        = "test-acc-tf-domain-group-updated"
  domain_names {
    domain_name = "www.testaccupdated.com"
  }
}
`
