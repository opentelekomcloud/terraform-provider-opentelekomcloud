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

const addressGroupMemberResourceName = "opentelekomcloud_cfw_address_group_member_v1.member_1"

func getAddressGroupMemberFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CfwV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CFW client: %s", err)
	}
	return addressgroup.GetGroupMember(client, state.Primary.Attributes["set_id"], state.Primary.Attributes["address"])
}

func TestAccCFWAddressGroupMemberV1_basic(t *testing.T) {
	var aclRule addressgroup.AddressGroupData
	rc := common.InitResourceCheck(
		addressGroupMemberResourceName,
		&aclRule,
		getAddressGroupMemberFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCFWAddressGroupMemberV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(addressGroupMemberResourceName, "description", "Test1111"),
					resource.TestCheckResourceAttrSet(addressGroupMemberResourceName, "id"),
				),
			},
			{
				ResourceName:      addressGroupMemberResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCFWAddressGroupMemberV1ImportStateIdFunc(),
			},
		},
	})
}

func testAccCFWAddressGroupMemberV1ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var setId string
		var address string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_cfw_address_group_member_v1" {
				setId = rs.Primary.Attributes["set_id"]
				address = rs.Primary.Attributes["address"]
			}
		}
		if setId == "" || address == "" {
			return "", fmt.Errorf("resource not found: %s/%s", setId, address)
		}
		return fmt.Sprintf("%s/%s", setId, address), nil
	}
}

var testAccCFWAddressGroupMemberV1Basic = `
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

resource "opentelekomcloud_cfw_address_group_member_v1" "member_1" {
  set_id      = opentelekomcloud_cfw_address_group_v1.group_1.id
  address     = "1.1.1.1"
  description = "Test1111"
}
`
