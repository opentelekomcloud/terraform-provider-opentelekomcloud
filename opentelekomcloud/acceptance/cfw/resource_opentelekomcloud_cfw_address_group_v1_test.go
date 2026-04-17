package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/addressgroup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const addressGroupResourceName = "opentelekomcloud_cfw_address_group_v1.group_1"

func getAddressGroupFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CfwV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CFW client: %s", err)
	}
	return addressgroup.GetAddressGroup(client, state.Primary.ID)
}

func TestAccCFWAddressGroupV1_basic(t *testing.T) {
	var group addressgroup.AddressGroupData
	rc := common.InitResourceCheck(
		addressGroupResourceName,
		&group,
		getAddressGroupFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCFWAddressGroupV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(addressGroupResourceName, "name", "test-acc-tf-address-group"),
					resource.TestCheckResourceAttrSet(addressGroupResourceName, "id"),
				),
			},
			{
				Config: testAccCFWAddressGroupV1Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(addressGroupResourceName, "name", "test-acc-tf-address-group-updated"),
				),
			},
			{
				ResourceName:      addressGroupResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"object_id",
				},
			},
		},
	})
}

var testAccCFWAddressGroupV1Basic = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_address_group_v1" "group_1" {
  object_id    = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  name         = "test-acc-tf-address-group"
  address_type = 0
}
`

var testAccCFWAddressGroupV1Update = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_address_group_v1" "group_1" {
  object_id    = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  name         = "test-acc-tf-address-group-updated"
  address_type = 0
}
`
