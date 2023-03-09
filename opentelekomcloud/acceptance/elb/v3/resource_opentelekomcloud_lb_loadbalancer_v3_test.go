package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/loadbalancers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
)

const resourceLBName = "opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1"
const resourceBWName = "opentelekomcloud_vpc_bandwidth_v2.bw"

func TestAccLBV3LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer

	qts := lbQuotas()
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3LoadBalancerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3LoadBalancerExists(resourceLBName, &lb),
					resource.TestCheckResourceAttr(resourceLBName, "deletion_protection", "true"),
				),
			},
			{
				Config: testAccLBV3LoadBalancerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceLBName, "name", "loadbalancer_1_updated"),
					resource.TestCheckResourceAttr(resourceLBName, "deletion_protection", "false"),
				),
			},
		},
	})
}

func TestAccLBV3LoadBalancer_bandwidth(t *testing.T) {
	var lb loadbalancers.LoadBalancer

	qts := lbQuotas()
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3LoadBalancerConfigNewBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3LoadBalancerExists(resourceLBName, &lb),
					resource.TestCheckResourceAttr(resourceLBName, "public_ip.0.bandwidth_share_type", "PER"),
					resource.TestCheckResourceAttr(resourceLBName, "public_ip.0.bandwidth_size", "10"),
					resource.TestCheckResourceAttr(resourceBWName, "size", "20"),
				),
			},
		},
	})
}

func TestAccLBV3LoadBalancer_import(t *testing.T) {
	qts := lbQuotas()
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3LoadBalancerConfigBasic,
			},
			{
				ResourceName:            resourceLBName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"public_ip.0._managed"},
			},
		},
	})
}

func testAccCheckLBV3LoadBalancerDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elbv3.ErrCreateClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_loadbalancer_v3" {
			continue
		}

		_, err := loadbalancers.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("loadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV3LoadBalancerExists(n string, lb *loadbalancers.LoadBalancer) resource.TestCheckFunc {
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

		found, err := loadbalancers.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("loadbalancer not found")
		}

		*lb = *found

		return nil
	}
}

var testAccLBV3LoadBalancerConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]

  public_ip {
    ip_type              = "5_gray"
    bandwidth_name       = "lb_band"
    bandwidth_size       = 10
    bandwidth_share_type = "PER"
  }

  deletion_protection = true

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3LoadBalancerConfigUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1_updated"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]

  public_ip {
    ip_type              = "5_gray"
    bandwidth_name       = "lb_band"
    bandwidth_size       = 10
    bandwidth_share_type = "PER"
  }

  deletion_protection = false

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3LoadBalancerConfigNewBandwidth = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]

  public_ip {
    ip_type              = "5_gray"
    bandwidth_name       = "lb_band"
    bandwidth_size       = 10
    bandwidth_share_type = "PER"
  }

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}

resource "opentelekomcloud_vpc_bandwidth_v2" "bw" {
  name = "lb_band"
  size = 20
}

resource "opentelekomcloud_vpc_bandwidth_associate_v2" "associate" {
  bandwidth    = opentelekomcloud_vpc_bandwidth_v2.bw.id
  floating_ips = [opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.public_ip.0.id]
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
