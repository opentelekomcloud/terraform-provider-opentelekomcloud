package acceptance

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cfwmanagementv1 "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/management"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const firewallResourceName = "opentelekomcloud_cfw_firewall_v1.firewall_1"

func getFirewallFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CfwV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CFW client: %s", err)
	}
	serviceType, err := strconv.Atoi(state.Primary.Attributes["service_type"])
	if err != nil {
		return nil, fmt.Errorf("error converting service type to integer: %s", err)
	}
	return cfwmanagementv1.Get(client, state.Primary.ID, serviceType)
}

func TestAccCFWFirewallV1_basic(t *testing.T) {
	var firewall cfwmanagementv1.GetFirewallInstanceResponseRecord
	rc := common.InitResourceCheck(
		firewallResourceName,
		&firewall,
		getFirewallFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCFWFirewallV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(firewallResourceName, "name", "test-acc-tf-firewall"),
					resource.TestCheckResourceAttr(firewallResourceName, "service_type", "0"),
				),
			},
			{
				ResourceName:      firewallResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCFWFirewallV1ImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"flavor.version",
					"charge_info",
				},
			},
		},
	})
}

func testAccCFWFirewallV1ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var serviceType string
		var id string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_cfw_firewall_v1" {
				id = rs.Primary.ID
				serviceType = rs.Primary.Attributes["service_type"]
			}
		}
		if id == "" || serviceType == "" {
			return "", fmt.Errorf("resource not found: %s/%s", id, serviceType)
		}
		return fmt.Sprintf("%s/%s", id, serviceType), nil
	}
}

var testAccCFWFirewallV1Basic = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}`
