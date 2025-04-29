package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/servicegroup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const serviceGroupMemberResourceName = "opentelekomcloud_cfw_service_group_member_v1.member_1"

func getServiceGroupMemberFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CfwV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CFW client: %s", err)
	}
	groupMembers, err := servicegroup.ListGroupMembers(client, state.Primary.Attributes["set_id"])
	if err != nil {
		return nil, fmt.Errorf("unable to fetch OpenTelekomCloud CFW service group member: %s", err)
	}
	for _, member := range groupMembers {
		if member.ItemID == state.Primary.ID {
			return member, nil
		}
	}
	return nil, fmt.Errorf("unable to find OpenTelekomCloud CFW service group member or member does not exist: %s", err)
}

func TestAccCFWServiceGroupMemberV1_basic(t *testing.T) {
	var member servicegroup.GroupMemberRecord
	rc := common.InitResourceCheck(
		serviceGroupMemberResourceName,
		&member,
		getServiceGroupMemberFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCFWServiceGroupMemberV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(serviceGroupMemberResourceName, "description", "Test611"),
					resource.TestCheckResourceAttrSet(serviceGroupMemberResourceName, "id"),
				),
			},
			{
				ResourceName:      serviceGroupMemberResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCFWServiceGroupMemberV1ImportStateIdFunc(),
			},
		},
	})
}

func testAccCFWServiceGroupMemberV1ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var setId, id string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_cfw_service_group_member_v1" {
				setId = rs.Primary.Attributes["set_id"]
				id = rs.Primary.ID
			}
		}
		if setId == "" || id == "" {
			return "", fmt.Errorf("resource not found: %s/%s", setId, id)
		}
		return fmt.Sprintf(" %s/%s", setId, id), nil
	}
}

var testAccCFWServiceGroupMemberV1Basic = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_service_group_v1" "group_1" {
  object_id = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  name      = "test-acc-tf-service-group"
}

resource "opentelekomcloud_cfw_service_group_member_v1" "member_1" {
  set_id      = opentelekomcloud_cfw_service_group_v1.group_1.id
  protocol    = 6
  source_port = "1"
  dest_port   = "1"
  description = "Test611"
}
`
