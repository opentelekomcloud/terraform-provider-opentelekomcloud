package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/gateway"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	resourceEvpnGatewayName     = "opentelekomcloud_enterprise_vpn_gateway_v5.gw_1"
	resourceEvpnGatewayEip1Name = "opentelekomcloud_vpc_eip_v1.eip_1"
	resourceEvpnGatewayEip2Name = "opentelekomcloud_vpc_eip_v1.eip_2"
)

func getGatewayResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.EvpnV5Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud EVPN v5 client: %s", err)
	}
	return gateway.Get(client, state.Primary.ID)
}

func TestAccGateway_basic(t *testing.T) {
	var gw gateway.Gateway
	name := fmt.Sprintf("evpn_acc_gw_%s", acctest.RandString(5))
	updateName := fmt.Sprintf("evpn_acc_gw_up_%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceEvpnGatewayName,
		&gw,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEvpnGateway_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "name", name),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "ha_mode", "active-active"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceEvpnGatewayName, "eip1.0.id", resourceEvpnGatewayEip1Name, "id"),
					resource.TestCheckResourceAttrPair(resourceEvpnGatewayName, "eip2.0.id", resourceEvpnGatewayEip2Name, "id"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "tags.foo", "bar"),
				),
			},
			{
				Config: testEvpnGateway_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "name", updateName),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "local_subnets.1", "192.168.2.0/24"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "tags.foo", "bar-update"),
				),
			},
			{
				ResourceName:      resourceEvpnGatewayName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGateway_activeStandbyHAMode(t *testing.T) {
	var gw gateway.Gateway
	name := fmt.Sprintf("evpn_acc_gw_%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceEvpnGatewayName,
		&gw,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEvpnGateway_activeStandbyHAMode(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "name", name),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "ha_mode", "active-standby"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "eip1.0.bandwidth_name", "evpn-gw-bw-1"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "eip2.0.bandwidth_name", "evpn-gw-bw-2"),
				),
			},
			{
				ResourceName:      resourceEvpnGatewayName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGateway_withER(t *testing.T) {
	var gw gateway.Gateway
	name := fmt.Sprintf("evpn_acc_gw_%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceEvpnGatewayName,
		&gw,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEvpnGateway_withER(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "name", name),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "network_type", "private"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "attachment_type", "er"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "access_private_ip_1", "172.16.0.99"),
					resource.TestCheckResourceAttr(resourceEvpnGatewayName, "access_private_ip_2", "172.16.0.100"),
				),
			},
			{
				ResourceName:      resourceEvpnGatewayName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testEvpnGateway_base(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet" {
  name       = "%[1]s"
  vpc_id     = opentelekomcloud_vpc_v1.vpc.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s-1"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_vpc_eip_v1" "eip_2" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s-2"
    size        = 9
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, name)
}

func testEvpnGateway_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_enterprise_vpn_gateway_v5" "gw_1" {
  name           = "%s"
  vpc_id         = opentelekomcloud_vpc_v1.vpc.id
  local_subnets  = [opentelekomcloud_vpc_subnet_v1.subnet.cidr]
  connect_subnet = opentelekomcloud_vpc_subnet_v1.subnet.id

  availability_zones = [
    "eu-de-01",
    "eu-de-02"
  ]

  eip1 {
    id = opentelekomcloud_vpc_eip_v1.eip_1.id
  }

  eip2 {
    id = opentelekomcloud_vpc_eip_v1.eip_2.id
  }

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, testEvpnGateway_base(name), name)
}

func testEvpnGateway_update(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_enterprise_vpn_gateway_v5" "gw_1" {
  name   = "%s"
  vpc_id = opentelekomcloud_vpc_v1.vpc.id
  local_subnets = [
    opentelekomcloud_vpc_subnet_v1.subnet.cidr,
    "192.168.2.0/24"
  ]
  connect_subnet = opentelekomcloud_vpc_subnet_v1.subnet.id

  availability_zones = [
    "eu-de-01",
    "eu-de-02"
  ]

  eip1 {
    id = opentelekomcloud_vpc_eip_v1.eip_1.id
  }

  eip2 {
    id = opentelekomcloud_vpc_eip_v1.eip_2.id
  }

  tags = {
    key = "val"
    foo = "bar-update"
  }
}
`, testEvpnGateway_base(name), name)
}

func testEvpnGateway_activeStandbyHAMode(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet" {
  name       = "%[1]s"
  vpc_id     = opentelekomcloud_vpc_v1.vpc.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "opentelekomcloud_enterprise_vpn_gateway_v5" "gw_1" {
  name           = "%[1]s"
  ha_mode        = "active-standby"
  vpc_id         = opentelekomcloud_vpc_v1.vpc.id
  local_subnets  = [opentelekomcloud_vpc_subnet_v1.subnet.cidr]
  connect_subnet = opentelekomcloud_vpc_subnet_v1.subnet.id

  availability_zones = [
    "eu-de-01",
    "eu-de-02"
  ]

  eip1 {
    bandwidth_name = "evpn-gw-bw-1"
    type           = "5_bgp"
    bandwidth_size = 5
    charge_mode    = "traffic"
  }

  eip2 {
    bandwidth_name = "evpn-gw-bw-2"
    type           = "5_bgp"
    bandwidth_size = 5
    charge_mode    = "traffic"
  }
}
`, name)
}

func testEvpnGateway_withER(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc_er" {
  name = "%[1]s"
  cidr = "172.16.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_er" {
  name       = "%[1]s"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_er.id
  cidr       = "172.16.0.0/24"
  gateway_ip = "172.16.0.1"
}

resource "opentelekomcloud_er_instance_v3" "er_1" {
  availability_zones = ["eu-de-01", "eu-de-02"]

  name        = "%[1]s"
  asn         = "65000"
  description = "evpn test"
}

resource "opentelekomcloud_enterprise_vpn_gateway_v5" "gw_1" {
  name            = "%[1]s"
  network_type    = "private"
  attachment_type = "er"
  er_id           = opentelekomcloud_er_instance_v3.er_1.id

  availability_zones = [
    "eu-de-01",
    "eu-de-02"
  ]

  access_vpc_id    = opentelekomcloud_vpc_v1.vpc_er.id
  access_subnet_id = opentelekomcloud_vpc_subnet_v1.subnet_er.id

  access_private_ip_1 = "172.16.0.99"
  access_private_ip_2 = "172.16.0.100"
}
`, name)
}
