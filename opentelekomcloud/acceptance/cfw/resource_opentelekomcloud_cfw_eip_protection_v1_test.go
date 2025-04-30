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

func getFirewallFuncForEIPtest(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
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

func TestAccCFWEipProtetionV1_basic(t *testing.T) {
	var firewall cfwmanagementv1.GetFirewallInstanceResponseRecord
	firewallResourceName := "opentelekomcloud_cfw_firewall_v1.firewall_1"
	rc := common.InitResourceCheck(
		firewallResourceName,
		&firewall,
		getFirewallFuncForEIPtest,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCFWEipProtectionV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

var testAccCFWEipProtectionV1Basic = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
    name = "test-acc-tf-cfw-eip"
  }
  bandwidth {
    name        = "test-acc-tf-cfw-band"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_cfw_eip_protection_v1" "protect_1" {
  firewall_id = opentelekomcloud_cfw_firewall_v1.firewall_1.id
  object_id   = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  status      = 0
  eip_id      = opentelekomcloud_vpc_eip_v1.eip_1.id
  public_ip   = opentelekomcloud_vpc_eip_v1.eip_1.publicip.0.ip_address
}
`
