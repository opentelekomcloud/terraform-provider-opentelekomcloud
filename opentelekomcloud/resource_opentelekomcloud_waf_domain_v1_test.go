package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/waf/v1/domains"
)

func TestAccWafDomainV1_basic(t *testing.T) {
	var domain domains.Domain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1Exists("opentelekomcloud_waf_domain_v1.domain_1", &domain),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_domain_v1.domain_1", "hostname", "www.b.com"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_domain_v1.domain_1", "sip_header_name", "default"),
				),
			},
			{
				Config: testAccWafDomainV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1Exists("opentelekomcloud_waf_domain_v1.domain_1", &domain),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_domain_v1.domain_1", "hostname", "www.b.com"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_domain_v1.domain_1", "sip_header_name", ""),
				),
			},
		},
	})
}

func testAccCheckWafDomainV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.wafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_domain_v1" {
			continue
		}

		_, err := domains.Get(wafClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf domain still exists")
		}
	}

	return nil
}

func testAccCheckWafDomainV1Exists(n string, domain *domains.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		wafClient, err := config.wafV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
		}

		found, err := domains.Get(wafClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Waf domain not found")
		}

		*domain = *found

		return nil
	}
}

const testAccWafDomainV1_basic = `
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
}

resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
	name = "cert_1"
	content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
	key = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}

resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
	options {
		webattack = true
		crawler = true
	}
	full_detection = false
}

resource "opentelekomcloud_waf_domain_v1" "domain_1" {
	hostname = "www.b.com"
	server {
		front_protocol = "HTTPS"
		back_protocol = "HTTP"
		address = "${opentelekomcloud_networking_floatingip_v2.fip_1.address}"
		port = "8080"
	}
	certificate_id = "${opentelekomcloud_waf_certificate_v1.certificate_1.id}"
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	proxy = "true"
	sip_header_name = "default"
	sip_header_list = ["X-Forwarded-For"]
}
`

const testAccWafDomainV1_update = `
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
}

resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
	name = "cert_1"
	content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
	key = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}

resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
	options {
		webattack = true
		crawler = true
	}
	full_detection = false
}

resource "opentelekomcloud_waf_domain_v1" "domain_1" {
	hostname = "www.b.com"
	server {
		front_protocol = "HTTPS"
		back_protocol = "HTTP"
		address = "${opentelekomcloud_networking_floatingip_v2.fip_1.address}"
		port = "80"
	}
	certificate_id = "${opentelekomcloud_waf_certificate_v1.certificate_1.id}"
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	proxy = "false"
}
`
