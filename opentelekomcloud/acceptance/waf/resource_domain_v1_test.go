package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/domains"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceDomainName = "opentelekomcloud_waf_domain_v1.domain_1"

func TestAccWafDomainV1Basic(t *testing.T) {
	var domain domains.Domain

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1Exists(resourceDomainName, &domain),
					resource.TestCheckResourceAttr(resourceDomainName, "hostname", "www.b.com"),
					resource.TestCheckResourceAttr(resourceDomainName, "sip_header_name", "default"),
					resource.TestCheckResourceAttr(resourceDomainName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceDomainName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceDomainName, "cipher", "cipher_1"),
				),
			},
			{
				Config: testAccWafDomainV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1Exists(resourceDomainName, &domain),
					resource.TestCheckResourceAttr(resourceDomainName, "hostname", "www.b.com"),
					resource.TestCheckResourceAttr(resourceDomainName, "sip_header_name", ""),
					resource.TestCheckResourceAttr(resourceDomainName, "cipher", "cipher_default"),
				),
			},
			{
				Config: testAccWafDomainV1UpdateCertificate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1CertificateChanged(resourceDomainName, &domain),
				),
			},
		},
	})
}

func TestAccWafDomain_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainV1Basic,
			},
			{
				ResourceName:      resourceDomainName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckWafDomainV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	wafClient, err := config.WafV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud WAF client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_domain_v1" {
			continue
		}

		_, err := domains.Get(wafClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("WAF domain still exists")
		}
	}

	return nil
}

func testAccCheckWafDomainV1Exists(n string, domain *domains.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.WafV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud WAF client: %w", err)
		}

		found, err := domains.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("WAF domain not found")
		}

		*domain = *found

		return nil
	}
}

func testAccCheckWafDomainV1CertificateChanged(n string, domain *domains.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.WafV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud WAF client: %w", err)
		}

		found, err := domains.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("WAF domain not found")
		}

		if found.CertificateId == domain.CertificateId {
			return fmt.Errorf("certificate has not changed")
		}

		*domain = *found

		return nil
	}
}

const testAccWafDomainV1Basic = `
variable "content" {
	default = "<h1>Hello world</h1>"
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
  name    = "cert_1"
  content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
  key     = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}

resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_1"
  options {
    webattack = true
    crawler   = true
  }
  full_detection = false
}

resource "opentelekomcloud_waf_domain_v1" "domain_1" {
  hostname = "www.b.com"
  cipher   = "cipher_1"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = opentelekomcloud_networking_floatingip_v2.fip_1.address
    port            = "8080"
  }
  certificate_id  = opentelekomcloud_waf_certificate_v1.certificate_1.id
  policy_id       = opentelekomcloud_waf_policy_v1.policy_1.id
  proxy           = true
  sip_header_name = "default"
  sip_header_list = ["X-Forwarded-For"]

  block_page {
  	template = "custom"
	status_code = "200"
    content_type = "application/json"
    content = var.content
  }
}
`

const testAccWafDomainV1Update = `
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
  name    = "cert_1"
  content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
  key     = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}

resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_1"
  options {
    webattack = true
    crawler   = true
  }
  full_detection = false
}

resource "opentelekomcloud_waf_policy_v1" "policy_2" {
  name           = "policy_2"
  full_detection = true
}

resource "opentelekomcloud_waf_domain_v1" "domain_1" {
  hostname = "www.b.com"
  cipher   = "cipher_default"
  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = opentelekomcloud_networking_floatingip_v2.fip_1.address
    port            = 80
  }
  certificate_id = opentelekomcloud_waf_certificate_v1.certificate_1.id
  policy_id      = opentelekomcloud_waf_policy_v1.policy_2.id
  proxy          = false
  block_page {
  	template = "default"
  }
}
`

const testAccWafDomainV1UpdateCertificate = `
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
  name    = "cert_1"
  content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
  key     = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}

resource "opentelekomcloud_waf_certificate_v1" "certificate_2" {
  name    = "cert_2"
  content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
  key     = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}

resource "opentelekomcloud_waf_policy_v1" "policy_2" {
  name           = "policy_2"
  full_detection = true
}

resource "opentelekomcloud_waf_domain_v1" "domain_1" {
  hostname = "www.b.com"
  cipher   = "cipher_default"
  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = opentelekomcloud_networking_floatingip_v2.fip_1.address
    port            = 80
  }
  certificate_id = opentelekomcloud_waf_certificate_v1.certificate_2.id
  policy_id      = opentelekomcloud_waf_policy_v1.policy_2.id
  proxy          = false
}
`
