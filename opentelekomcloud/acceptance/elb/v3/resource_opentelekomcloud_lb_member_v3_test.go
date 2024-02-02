package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/members"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
)

const resourceMemberName = "opentelekomcloud_lb_member_v3.member"

func TestLBMemberV3_basic(t *testing.T) {
	var member members.Member
	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testLBMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testLBMemberV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testLBMemberV3Exists(resourceMemberName, &member),
					resource.TestCheckResourceAttr(resourceMemberName, "name", "member-1"),
					resource.TestCheckResourceAttr(resourceMemberName, "weight", "1"),
				),
			},
			{
				Config: testLBMemberV3Updated,
				Check: resource.ComposeTestCheckFunc(
					testLBMemberV3Exists(resourceMemberName, &member),
					resource.TestCheckResourceAttr(resourceMemberName, "name", ""),
					resource.TestCheckResourceAttr(resourceMemberName, "weight", "0"),
				),
			},
		},
	})
}

func TestLBMemberV3_ZeroWeight(t *testing.T) {
	var member members.Member
	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testLBMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testLBMemberV3ZeroWeight,
				Check: resource.ComposeTestCheckFunc(
					testLBMemberV3Exists(resourceMemberName, &member),
					resource.TestCheckResourceAttr(resourceMemberName, "name", "member-1"),
					resource.TestCheckResourceAttr(resourceMemberName, "weight", "0"),
				),
			},
			{
				Config: testLBMemberV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testLBMemberV3Exists(resourceMemberName, &member),
					resource.TestCheckResourceAttr(resourceMemberName, "weight", "0"),
				),
			},
		},
	})
}

func TestLBMemberV3_import(t *testing.T) {
	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testLBMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testLBMemberV3Basic,
			},
			{
				ResourceName:      resourceMemberName,
				ImportStateVerify: true,
				ImportState:       true,
			},
		},
	})
}

func testLBMemberV3Exists(n string, member *members.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ElbV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(elbv3.ErrCreateClient, err)
		}

		part := strings.Split(rs.Primary.ID, "/")
		found, err := members.Get(client, part[0], part[1]).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.Attributes["member_id"] {
			return fmt.Errorf("loadbalancer member not found")
		}

		*member = *found

		return nil
	}
}

func testLBMemberDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elbv3.ErrCreateClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_member_v3" {
			continue
		}

		part := strings.Split(rs.Primary.ID, "/")
		_, err := members.Get(client, part[0], part[1]).Extract()
		if err == nil {
			return fmt.Errorf("loadbalancer member still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

var testLBMemberV3Basic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  ip_target_enable = true

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "TCP"
}

resource "opentelekomcloud_lb_member_v3" "member" {
  name          = "member-1"
  pool_id       = opentelekomcloud_lb_pool_v3.pool.id
  address       = cidrhost(data.opentelekomcloud_vpc_subnet_v1.shared_subnet.cidr, 3)
  protocol_port = 8080
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testLBMemberV3Updated = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  ip_target_enable = true

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "TCP"
}

resource "opentelekomcloud_lb_member_v3" "member" {
  name          = ""
  pool_id       = opentelekomcloud_lb_pool_v3.pool.id
  address       = cidrhost(data.opentelekomcloud_vpc_subnet_v1.shared_subnet.cidr, 3)
  protocol_port = 8080
  weight        = 0
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testLBMemberV3ZeroWeight = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  ip_target_enable = true

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "TCP"
}

resource "opentelekomcloud_lb_member_v3" "member" {
  name          = "member-1"
  pool_id       = opentelekomcloud_lb_pool_v3.pool.id
  address       = cidrhost(data.opentelekomcloud_vpc_subnet_v1.shared_subnet.cidr, 3)
  protocol_port = 8080
  weight        = 0
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
