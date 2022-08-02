package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/pools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	elb "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v2"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

var (
	resourceMemberName1 = "opentelekomcloud_lb_member_v2.member_1"
	resourceMemberName2 = "opentelekomcloud_lb_member_v2.member_2"
)

func TestAccLBV2Member_basic(t *testing.T) {
	var member1 pools.Member
	var member2 pools.Member

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.LoadBalancer, Count: 1},
				{Q: quotas.LbListener, Count: 1},
				{Q: quotas.LbPool, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV2MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccLBV2MemberConfigBasic,
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false, unfinished elb?
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MemberExists(resourceMemberName1, &member1),
					testAccCheckLBV2MemberExists(resourceMemberName2, &member2),
					resource.TestCheckResourceAttr(resourceMemberName1, "weight", "0"),
					resource.TestCheckResourceAttr(resourceMemberName2, "weight", "10"),
				),
			},
			{
				Config:             TestAccLBV2MemberConfigUpdate,
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false, unfinished elb?
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMemberName1, "weight", "10"),
					resource.TestCheckResourceAttr(resourceMemberName2, "weight", "0"),
				),
			},
		},
	})
}

func TestAccLBV2Member_import(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.LoadBalancer, Count: 1},
				{Q: quotas.LbListener, Count: 1},
				{Q: quotas.LbPool, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV2MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccLBV2MemberConfigBasic,
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false, unfinished elb?
			},
			{
				ResourceName:      resourceMemberName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: memberImportStateIDFunc(resourceMemberName1),
			},
		},
	})
}

func memberImportStateIDFunc(memberResourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		policyID := s.RootModule().Resources[resourcePoolName].Primary.ID
		ccRuleID := s.RootModule().Resources[memberResourceName].Primary.ID
		return fmt.Sprintf("%s/%s", policyID, ccRuleID), nil
	}
}

func testAccCheckLBV2MemberDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elb.ErrCreationV2Client, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_member_v2" {
			continue
		}

		poolID := rs.Primary.Attributes["pool_id"]
		_, err := pools.GetMember(client, poolID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("member still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2MemberExists(n string, member *pools.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ElbV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(elb.ErrCreationV2Client, err)
		}

		poolID := rs.Primary.Attributes["pool_id"]
		found, err := pools.GetMember(client, poolID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("member not found")
		}

		*member = *found

		return nil
	}
}

var testAccLBV2MemberConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_1.id
}

resource "opentelekomcloud_lb_member_v2" "member_1" {
  address       = "10.0.0.10"
  protocol_port = 8080
  pool_id       = opentelekomcloud_lb_pool_v2.pool_1.id
  subnet_id     = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
  weight        = 0
}

resource "opentelekomcloud_lb_member_v2" "member_2" {
  address       = "10.0.0.11"
  protocol_port = 8080
  pool_id       = opentelekomcloud_lb_pool_v2.pool_1.id
  subnet_id     = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
  weight        = 10
}
`, common.DataSourceSubnet)

var TestAccLBV2MemberConfigUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_1.id
}

resource "opentelekomcloud_lb_member_v2" "member_1" {
  address        = "10.0.0.10"
  protocol_port  = 8080
  weight         = 10
  admin_state_up = true
  pool_id        = opentelekomcloud_lb_pool_v2.pool_1.id
  subnet_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
}

resource "opentelekomcloud_lb_member_v2" "member_2" {
  address        = "10.0.0.11"
  protocol_port  = 8080
  weight         = 0
  admin_state_up = true
  pool_id        = opentelekomcloud_lb_pool_v2.pool_1.id
  subnet_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
}
`, common.DataSourceSubnet)
