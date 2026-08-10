package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/security/group"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const vpcSecGroupResourceName = "opentelekomcloud_vpc_secgroup_v3.group_1"

func getVpcSecGroupV3(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NetworkingV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v3 Client: %s", err)
	}
	return group.Get(client, state.Primary.ID)
}

func TestAccVpcSecGroupV3_basic(t *testing.T) {
	var group group.SecurityGroup
	rc := common.InitResourceCheck(
		vpcSecGroupResourceName,
		&group,
		getVpcSecGroupV3,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSecGroupV3Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(vpcSecGroupResourceName, "name", "test-acc-sec-group-v3"),
					resource.TestCheckResourceAttr(vpcSecGroupResourceName, "description", "created"),
				),
			},
			{
				Config: testAccVpcSecGroupV3Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(vpcSecGroupResourceName, "name", "test-acc-sec-group-v3"),
					resource.TestCheckResourceAttr(vpcSecGroupResourceName, "description", "updated"),
				),
			},
			{
				ResourceName:      vpcSecGroupResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_default_rules",
				},
			},
		},
	})
}

var testAccVpcSecGroupV3Basic = `
resource "opentelekomcloud_vpc_secgroup_v3" "group_1" {
  name                 = "test-acc-sec-group-v3"
  description          = "created"
  delete_default_rules = true
}
`

var testAccVpcSecGroupV3Update = `
resource "opentelekomcloud_vpc_secgroup_v3" "group_1" {
  name                 = "test-acc-sec-group-v3"
  description          = "updated"
  delete_default_rules = true
}
`
