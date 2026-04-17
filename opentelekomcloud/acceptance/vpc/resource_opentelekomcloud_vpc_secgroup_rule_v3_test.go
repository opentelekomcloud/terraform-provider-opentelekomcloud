package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/security/rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const vpcSecGroupRuleResourceName = "opentelekomcloud_vpc_secgroup_rule_v3.rule_1"

func getVpcSecGroupRuleV3(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NetworkingV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v3 Client: %s", err)
	}
	return rules.Get(client, state.Primary.ID)
}

func TestAccVpcSecGroupRuleV3_basic(t *testing.T) {
	var rule rules.SecurityGroupRule
	rc := common.InitResourceCheck(
		vpcSecGroupRuleResourceName,
		&rule,
		getVpcSecGroupRuleV3,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSecGroupRuleV3Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(vpcSecGroupRuleResourceName, "description", "created-rule"),
					resource.TestCheckResourceAttr(vpcSecGroupRuleResourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(vpcSecGroupRuleResourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(vpcSecGroupRuleResourceName, "action", "allow"),
					resource.TestCheckResourceAttr(vpcSecGroupRuleResourceName, "priority", "1"),
					resource.TestCheckResourceAttr(vpcSecGroupRuleResourceName, "multi_port", "8080"),
					resource.TestCheckResourceAttr(vpcSecGroupRuleResourceName, "remote_ip_prefix", "10.10.0.0/16"),
				),
			},
			{
				ResourceName:      vpcSecGroupRuleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccVpcSecGroupRuleV3Basic = `
resource "opentelekomcloud_vpc_secgroup_v3" "group_1" {
  name = "test-acc-sec-group-v3"
}

resource "opentelekomcloud_vpc_secgroup_rule_v3" "rule_1" {
  security_group_id = opentelekomcloud_vpc_secgroup_v3.group_1.id
  description       = "created-rule"
  direction         = "ingress"
  protocol          = "tcp"
  action            = "allow"
  priority          = 1
  multi_port        = "8080"
  remote_ip_prefix  = "10.10.0.0/16"
}
`
