package acceptance

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	domains "github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/hosts"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const wafdDomainResourceName = "opentelekomcloud_waf_dedicated_domain_v1.domain_1"

func TestAccWafDedicatedDomainV1_basic(t *testing.T) {
	var domain domains.Host
	var hostName = fmt.Sprintf("wafd%s", acctest.RandString(5))
	log.Printf("[DEBUG] The opentelekomcloud Waf dedicated instance test running in '%s' region.", env.OS_REGION_NAME)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicateDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDomainV1_basic(hostName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(
						wafdDomainResourceName, &domain),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "domain", fmt.Sprintf("www.%s.com", hostName)),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.client_protocol", "HTTP"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.address", "192.168.0.10"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "timeout_config.0.connect_timeout", "20"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_basicUpdate(hostName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(
						wafdDomainResourceName, &domain),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "domain", fmt.Sprintf("www.%s.com", hostName)),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.client_protocol", "HTTP"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.address", "192.168.0.10"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.1.client_protocol", "HTTP"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.1.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.1.port", "80"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.1.address", "192.168.0.11"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.1.type", "ipv4"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_cert(hostName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(
						wafdDomainResourceName, &domain),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "domain", fmt.Sprintf("www.%s.com", hostName)),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "tls", "TLS v1.1"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "cipher", "cipher_1"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.address", "192.168.0.20"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.type", "ipv4"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_certUpdate(hostName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(
						wafdDomainResourceName, &domain),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "domain", fmt.Sprintf("www.%s.com", hostName)),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "cipher", "cipher_2"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "pci_3ds", "true"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "pci_dss", "true"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.0.address", "192.168.0.20"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.1.address", "192.168.0.21"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.1.port", "8443"),
				),
			},
			{
				ResourceName:            wafdDomainResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_policy"},
			},
		},
	})
}

func TestAccWafDedicatedDomainV1_timeoutConfig(t *testing.T) {
	var domain domains.Host
	var hostName = fmt.Sprintf("wafd%s", acctest.RandString(5))
	log.Printf("[DEBUG] The opentelekomcloud Waf dedicated instance test running in '%s' region.", env.OS_REGION_NAME)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicateDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDomainV1_basicUpdate(hostName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(
						wafdDomainResourceName, &domain),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "domain", fmt.Sprintf("www.%s.com", hostName)),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.#", "2"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_timeoutConfig(hostName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(
						wafdDomainResourceName, &domain),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "domain", fmt.Sprintf("www.%s.com", hostName)),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "timeout_config.0.connect_timeout", "200"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "timeout_config.0.read_timeout", "100"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "timeout_config.0.send_timeout", "50"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_timeoutConfigUpdate(hostName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(
						wafdDomainResourceName, &domain),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "domain", fmt.Sprintf("www.%s.com", hostName)),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "timeout_config.0.connect_timeout", "150"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "timeout_config.0.read_timeout", "200"),
					resource.TestCheckResourceAttr(wafdDomainResourceName, "timeout_config.0.send_timeout", "100"),
				),
			},
		},
	})
}

func testAccCheckWafDedicateDomainV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_domain_v1" {
			continue
		}
		_, err = domains.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated domain (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckWafDedicatedDomainV1Exists(n string, instance *domains.Host) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
		if err != nil {
			return err
		}

		var found *domains.Host
		found, err = domains.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		*instance = *found

		return nil
	}
}

func testAccWafDedicatedDomainV1_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name            = "domain_policy_1"
  protection_mode = "log"
  full_detection  = false
  level           = 2

  options {
    crawler    = true
    web_attack = true
  }
}

resource "opentelekomcloud_waf_dedicated_domain_v1" "domain_1" {
  domain      = "www.%s.com"
  keep_policy = true
  proxy       = true

  policy_id = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id

  server {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "192.168.0.10"
    port            = 8080
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }
  timeout_config {
    connect_timeout = 20
  }
}
`, common.DataSourceSubnet, name)
}

func testAccWafDedicatedDomainV1_basicUpdate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name            = "domain_policy_1"
  protection_mode = "log"
  full_detection  = false
  level           = 2

  options {
    crawler    = true
    web_attack = true
  }
}

resource "opentelekomcloud_waf_dedicated_domain_v1" "domain_1" {
  domain      = "www.%s.com"
  keep_policy = false
  proxy       = false

  server {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "192.168.0.10"
    port            = 8080
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }
  server {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "192.168.0.11"
    port            = 80
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }
}
`, common.DataSourceSubnet, name)
}

func testAccWafDedicatedDomainV1_cert(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_waf_dedicated_domain_v1" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = opentelekomcloud_waf_dedicated_certificate_v1.cert_1.id
  keep_policy    = false
  proxy          = false
  tls            = "TLS v1.1"
  cipher         = "cipher_1"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "192.168.0.20"
    port            = 8080
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }

  depends_on = [
    opentelekomcloud_waf_dedicated_certificate_v1.cert_1
  ]
}

%s
`, common.DataSourceSubnet, name, testAccWafDedicatedCertificateV1Basic)
}

func testAccWafDedicatedDomainV1_certUpdate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_waf_dedicated_domain_v1" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = opentelekomcloud_waf_dedicated_certificate_v1.cert_1.id
  keep_policy    = false
  proxy          = true
  tls            = "TLS v1.2"
  cipher         = "cipher_2"
  pci_3ds        = true
  pci_dss        = true

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "192.168.0.20"
    port            = 8443
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "192.168.0.21"
    port            = 8443
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }

  depends_on = [
    opentelekomcloud_waf_dedicated_certificate_v1.cert_1
  ]
}

%s
`, common.DataSourceSubnet, name, testAccWafDedicatedCertificateV1Basic)
}

func testAccWafDedicatedDomainV1_timeoutConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name            = "domain_policy_1"
  protection_mode = "log"
  full_detection  = false
  level           = 2

  options {
    crawler    = true
    web_attack = true
  }
}

resource "opentelekomcloud_waf_dedicated_domain_v1" "domain_1" {
  domain      = "www.%s.com"
  keep_policy = false
  proxy       = false

  server {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "192.168.0.10"
    port            = 8080
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }
  server {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "192.168.0.11"
    port            = 80
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }
  timeout_config {
    connect_timeout = 200
    read_timeout    = 100
    send_timeout    = 50
  }
}
`, common.DataSourceSubnet, name)
}

func testAccWafDedicatedDomainV1_timeoutConfigUpdate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name            = "domain_policy_1"
  protection_mode = "log"
  full_detection  = false
  level           = 2

  options {
    crawler    = true
    web_attack = true
  }
}

resource "opentelekomcloud_waf_dedicated_domain_v1" "domain_1" {
  domain      = "www.%s.com"
  keep_policy = false
  proxy       = false

  server {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "192.168.0.10"
    port            = 8080
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }
  server {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "192.168.0.11"
    port            = 80
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }
  timeout_config {
    connect_timeout = 150
    read_timeout    = 200
    send_timeout    = 100
  }
}
`, common.DataSourceSubnet, name)
}
