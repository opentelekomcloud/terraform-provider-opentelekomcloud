package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/connection"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getConnectionResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.EvpnV5Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud EVPN v5 client: %s", err)
	}
	return connection.Get(client, state.Primary.ID)
}

func TestAccConnection_basic(t *testing.T) {
	var conn connection.Connection

	rName := "opentelekomcloud_enterprise_vpn_connection_v5.conn"
	ipAddress := "172.16.1.2"

	name := fmt.Sprintf("evpn_acc_conn_%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		rName,
		&conn,
		getConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEvpnConnection_basic(name, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "vpn_type", "STATIC"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.authentication_algorithm", "sha2-256"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.encryption_algorithm", "aes-128"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.lifetime_seconds", "86400"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.authentication_algorithm", "sha2-256"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.encryption_algorithm", "aes-128"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.lifetime_seconds", "3600"),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
				),
			},
			{
				Config: testEvpnConnection_update(name, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.authentication_algorithm", "sha2-512"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.encryption_algorithm", "aes-256"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.lifetime_seconds", "172800"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.authentication_algorithm", "sha2-512"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.encryption_algorithm", "aes-256"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.lifetime_seconds", "7200"),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"psk",
				},
			},
		},
	})
}

func TestAccConnection_policy(t *testing.T) {
	var conn connection.Connection

	rName := "opentelekomcloud_enterprise_vpn_connection_v5.pol"
	ipAddress := "172.16.1.3"

	name := fmt.Sprintf("evpn_acc_conn_%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		rName,
		&conn,
		getConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEvpnConnection_policy(name, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "vpn_type", "POLICY"),
					resource.TestCheckResourceAttr(rName, "policy_rules.0.source", "192.168.11.0/24"),
					resource.TestCheckResourceAttr(rName, "policy_rules.0.destination.0", "192.168.12.0/24"),
					resource.TestCheckResourceAttr(rName, "policy_rules.0.destination.1", "192.168.13.0/24"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.authentication_algorithm", "sha2-512"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.encryption_algorithm", "aes-256"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.lifetime_seconds", "172800"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.authentication_algorithm", "sha2-512"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.encryption_algorithm", "aes-256"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.lifetime_seconds", "7200"),
				),
			},
		},
	})
}

func TestAccConnection_haRole(t *testing.T) {
	var conn connection.Connection

	rName := "opentelekomcloud_enterprise_vpn_connection_v5.ha"
	ipAddress := "172.16.1.5"

	name := fmt.Sprintf("evpn_acc_conn_%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		rName,
		&conn,
		getConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEvpnConnection_haRole(name, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "vpn_type", "POLICY"),
					resource.TestCheckResourceAttr(rName, "policy_rules.0.source", "192.168.11.0/24"),
					resource.TestCheckResourceAttr(rName, "policy_rules.0.destination.0", "192.168.12.0/24"),
					resource.TestCheckResourceAttr(rName, "policy_rules.0.destination.1", "192.168.13.0/24"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.authentication_algorithm", "sha2-512"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.encryption_algorithm", "aes-256"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.lifetime_seconds", "172800"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.authentication_algorithm", "sha2-512"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.encryption_algorithm", "aes-256"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.lifetime_seconds", "7200"),
					resource.TestCheckResourceAttr(rName, "ha_role", "slave"),
				),
			},
		},
	})
}

func testEvpnConnection_basic(name, ip string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_enterprise_vpn_connection_v5" "conn" {
  name                = "%s"
  gateway_id          = opentelekomcloud_enterprise_vpn_gateway_v5.gw_1.id
  gateway_ip          = opentelekomcloud_vpc_eip_v1.eip_1.id
  customer_gateway_id = opentelekomcloud_enterprise_vpn_customer_gateway_v5.cgw_1.id
  peer_subnets        = ["192.168.55.0/24"]
  vpn_type            = "static"
  psk                 = "Test@123"
  enable_nqa          = true

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, testEvpnGateway_basic(name), testCustomerGateway_basic(name, ip), name)
}

func testEvpnConnection_update(name, ip string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_enterprise_vpn_connection_v5" "conn" {
  name                = "%s_update"
  gateway_id          = opentelekomcloud_enterprise_vpn_gateway_v5.gw_1.id
  gateway_ip          = opentelekomcloud_vpc_eip_v1.eip_1.id
  customer_gateway_id = opentelekomcloud_enterprise_vpn_customer_gateway_v5.cgw_1.id
  peer_subnets        = ["192.168.55.0/24"]
  vpn_type            = "static"
  psk                 = "Test@123"

  ikepolicy {
    authentication_algorithm = "sha2-512"
    encryption_algorithm     = "aes-256"
    lifetime_seconds         = 172800
  }

  ipsecpolicy {
    authentication_algorithm = "sha2-512"
    encryption_algorithm     = "aes-256"
    lifetime_seconds         = 7200
  }

  tags = {
    key = "val"
    foo = "bar-update"
  }
}
`, testEvpnGateway_basic(name), testCustomerGateway_basic(name, ip), name)
}

func testEvpnConnection_policy(name, ip string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_enterprise_vpn_connection_v5" "pol" {
  name                = "%s"
  gateway_id          = opentelekomcloud_enterprise_vpn_gateway_v5.gw_1.id
  gateway_ip          = opentelekomcloud_vpc_eip_v1.eip_1.id
  customer_gateway_id = opentelekomcloud_enterprise_vpn_customer_gateway_v5.cgw_1.id
  peer_subnets        = ["192.168.55.0/24"]
  vpn_type            = "policy"
  psk                 = "Test@123"

  policy_rules {
    source      = "192.168.11.0/24"
    destination = ["192.168.12.0/24", "192.168.13.0/24"]
  }

  ikepolicy {
    authentication_algorithm = "sha2-512"
    encryption_algorithm     = "aes-256"
    lifetime_seconds         = 172800
  }

  ipsecpolicy {
    authentication_algorithm = "sha2-512"
    encryption_algorithm     = "aes-256"
    lifetime_seconds         = 7200
  }
}
`, testEvpnGateway_basic(name), testCustomerGateway_basic(name, ip), name)
}

func testEvpnConnection_haRole(name, ip string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_enterprise_vpn_connection_v5" "ha" {
  name                = "%s"
  gateway_id          = opentelekomcloud_enterprise_vpn_gateway_v5.gw_1.id
  gateway_ip          = opentelekomcloud_enterprise_vpn_gateway_v5.gw_1.eip2.0.id
  customer_gateway_id = opentelekomcloud_enterprise_vpn_customer_gateway_v5.cgw_1.id
  peer_subnets        = ["192.168.55.0/24"]
  vpn_type            = "policy"
  psk                 = "Test@123"
  ha_role             = "slave"

  policy_rules {
    source      = "192.168.11.0/24"
    destination = ["192.168.12.0/24", "192.168.13.0/24"]
  }

  ikepolicy {
    authentication_algorithm = "sha2-512"
    encryption_algorithm     = "aes-256"
    lifetime_seconds         = 172800
  }

  ipsecpolicy {
    authentication_algorithm = "sha2-512"
    encryption_algorithm     = "aes-256"
    lifetime_seconds         = 7200
  }
}
`, testEvpnGateway_activeStandbyHAMode(name), testCustomerGateway_basic(name, ip), name)
}
