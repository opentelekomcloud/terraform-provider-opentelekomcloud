package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/listeners"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
)

const resourceListenerName = "opentelekomcloud_lb_listener_v3.listener_1"

func TestAccLBV3Listener_basic(t *testing.T) {
	var listener listeners.Listener

	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbCertificate, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.LbListener, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3ListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3ListenerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3ListenerExists(resourceListenerName, &listener),
					resource.TestCheckResourceAttr(resourceListenerName, "name", "listener_1"),
					resource.TestCheckResourceAttr(resourceListenerName, "description", "some interesting description"),
					resource.TestCheckResourceAttr(resourceListenerName, "tls_ciphers_policy", "tls-1-2-fs"),
					resource.TestCheckResourceAttr(resourceListenerName, "advanced_forwarding", "true"),
					resource.TestCheckResourceAttr(resourceListenerName, "sni_match_algo", "wildcard"),
					resource.TestCheckResourceAttr(resourceListenerName, "security_policy_id", ""),
				),
			},
			{
				Config: testAccLBV3ListenerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceListenerName, "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(resourceListenerName, "description", ""),
					resource.TestCheckResourceAttr(resourceListenerName, "tls_ciphers_policy", "tls-1-2-fs-with-1-3"),
					resource.TestCheckResourceAttr(resourceListenerName, "advanced_forwarding", "true"),
					resource.TestCheckResourceAttr(resourceListenerName, "sni_match_algo", "longest_suffix"),
					resource.TestCheckResourceAttr(resourceListenerName, "security_policy_id", ""),
				),
			},
		},
	})
}

func TestAccLBV3Listener_TCP(t *testing.T) {
	var listener listeners.Listener

	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbCertificate, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.LbListener, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3ListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3ListenerConfigTCP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3ListenerExists(resourceListenerName, &listener),
					resource.TestCheckResourceAttr(resourceListenerName, "name", "listener_1"),
					resource.TestCheckResourceAttr(resourceListenerName, "description", "some interesting description"),
				),
			},
			{
				Config: testAccLBV3ListenerConfigTCPUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3ListenerExists(resourceListenerName, &listener),
					resource.TestCheckResourceAttr(resourceListenerName, "name", "listener_1"),
					resource.TestCheckResourceAttr(resourceListenerName, "description", "other description"),
				),
			},
		},
	})
}

func TestAccLBV3Listener_HTTP_to_TCP(t *testing.T) {
	var listener listeners.Listener

	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbCertificate, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.LbListener, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3ListenerDestroy,
		Steps: []resource.TestStep{

			{
				Config: testAccLBV3ListenerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3ListenerExists(resourceListenerName, &listener),
					resource.TestCheckResourceAttr(resourceListenerName, "name", "listener_1"),
					resource.TestCheckResourceAttr(resourceListenerName, "description", "some interesting description"),
				),
			},
			{
				Config: testAccLBV3ListenerConfigTCP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3ListenerExists(resourceListenerName, &listener),
					resource.TestCheckResourceAttr(resourceListenerName, "name", "listener_1"),
					resource.TestCheckResourceAttr(resourceListenerName, "description", "some interesting description"),
				),
			},
		},
	})
}

func TestAccLBV3Listener_ipGroup(t *testing.T) {
	var listener listeners.Listener

	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbCertificate, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.LbListener, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3ListenerDestroy,
		Steps: []resource.TestStep{

			{
				Config: testAccLBV3ListenerConfigIpGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3ListenerExists(resourceListenerName, &listener),
					resource.TestCheckResourceAttr(resourceListenerName, "name", "listener_1"),
					resource.TestCheckResourceAttr(resourceListenerName, "description", "some interesting description"),
					resource.TestCheckResourceAttr(resourceListenerName, "ip_group.#", "1"),
					resource.TestCheckResourceAttr(resourceListenerName, "ip_group.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceListenerName, "ip_group.0.type", "white"),
				),
			},
			{
				Config: testAccLBV3ListenerConfigIpGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3ListenerExists(resourceListenerName, &listener),
					resource.TestCheckResourceAttr(resourceListenerName, "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(resourceListenerName, "ip_group.#", "1"),
					resource.TestCheckResourceAttr(resourceListenerName, "ip_group.0.enable", "false"),
					resource.TestCheckResourceAttr(resourceListenerName, "ip_group.0.type", "white"),
				),
			},
			{
				Config: testAccLBV3ListenerConfigIpGroupRemoveAllIpAddresses,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3ListenerExists(resourceListenerName, &listener),
					resource.TestCheckResourceAttr(resourceListenerName, "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(resourceListenerName, "ip_group.#", "1"),
					resource.TestCheckResourceAttr(resourceListenerName, "ip_group.0.enable", "false"),
					resource.TestCheckResourceAttr(resourceListenerName, "ip_group.0.type", "black"),
				),
			},
		},
	})
}

func TestAccLBV3Listener_import(t *testing.T) {
	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbCertificate, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.LbListener, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3ListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3ListenerConfigBasic,
			},
			{
				ResourceName:      resourceListenerName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLBV3ListenerDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elbv3.ErrCreateClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_listener_v3" {
			continue
		}

		_, err := listeners.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("listener still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV3ListenerExists(n string, listener *listeners.Listener) resource.TestCheckFunc {
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

		found, err := listeners.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("listener not found")
		}

		*listener = *found

		return nil
	}
}

var testAccLBV3ListenerConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_certificate_v3" "certificate_1" {
  name        = "certificate_1"
  type        = "server"
  private_key = %s
  certificate = %s
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name                      = "listener_1"
  description               = "some interesting description"
  loadbalancer_id           = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol                  = "HTTPS"
  protocol_port             = 443
  default_tls_container_ref = opentelekomcloud_lb_certificate_v3.certificate_1.id
  tls_ciphers_policy        = "tls-1-2-fs"

  advanced_forwarding = true
  sni_match_algo      = "wildcard"

  insert_headers {
    forwarded_host = true
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, privateKey, certificate)

var testAccLBV3ListenerConfigUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1_updated"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_certificate_v3" "certificate_1" {
  name        = "certificate_1"
  type        = "server"
  private_key = %s
  certificate = %s
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name                      = "listener_1_updated"
  loadbalancer_id           = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol                  = "HTTPS"
  protocol_port             = 443
  default_tls_container_ref = opentelekomcloud_lb_certificate_v3.certificate_1.id
  tls_ciphers_policy        = "tls-1-2-fs-with-1-3"

  sni_match_algo = "longest_suffix"
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, privateKey, certificate)

var testAccLBV3ListenerConfigTCP = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name            = "listener_1"
  description     = "some interesting description"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol        = "TCP"
  protocol_port   = 5000
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3ListenerConfigTCPUpdated = fmt.Sprintf(`
%s
resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name            = "listener_1"
  description     = "other description"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol        = "TCP"
  protocol_port   = 5360
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3ListenerConfigIpGroup = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_1"
  description = "some interesting description 1"

  ip_list {
    ip          = "192.168.10.10"
    description = "first"
  }
}

resource "opentelekomcloud_lb_ipgroup_v3" "group_2" {
  name        = "group_2"
  description = "some interesting description 2"

  ip_list {
    ip          = "192.168.10.11"
    description = "second"
  }
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name            = "listener_1"
  description     = "some interesting description"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol        = "HTTP"
  protocol_port   = 8080

  advanced_forwarding = true
  sni_match_algo      = "wildcard"

  insert_headers {
    forwarded_host = true
  }

  ip_group {
    id     = opentelekomcloud_lb_ipgroup_v3.group_1.id
    enable = true
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3ListenerConfigIpGroupUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1_updated"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_1"
  description = "some interesting description 1"

  ip_list {
    ip          = "192.168.10.10"
    description = "first"
  }
}

resource "opentelekomcloud_lb_ipgroup_v3" "group_2" {
  name        = "group_2"
  description = "some interesting description 2"

  ip_list {
    ip          = "192.168.10.11"
    description = "second"
  }
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name            = "listener_1_updated"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol        = "HTTP"
  protocol_port   = 8080

  sni_match_algo = "longest_suffix"

  ip_group {
    id     = opentelekomcloud_lb_ipgroup_v3.group_2.id
    enable = false
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3ListenerConfigIpGroupRemoveAllIpAddresses = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1_updated"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_1"
  description = "some interesting description 1"

  ip_list {
    ip          = "192.168.10.10"
    description = "first"
  }
}

resource "opentelekomcloud_lb_ipgroup_v3" "group_2" {
  name        = "group_2_empty"
  description = "some interesting description 2"
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name            = "listener_1_updated"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol        = "HTTP"
  protocol_port   = 8080

  sni_match_algo = "longest_suffix"

  ip_group {
    id     = opentelekomcloud_lb_ipgroup_v3.group_2.id
    enable = false
    type   = "black"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
