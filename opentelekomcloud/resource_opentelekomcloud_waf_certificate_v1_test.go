package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/waf/v1/certificates"
)

func TestAccWafCertificateV1_basic(t *testing.T) {
	var certificate certificates.Certificate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafCertificateV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCertificateV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafCertificateV1Exists("opentelekomcloud_waf_certificate_v1.certificate_1", &certificate),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_certificate_v1.certificate_1", "name", "cert_1"),
				),
			},
			{
				Config: testAccWafCertificateV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafCertificateV1Exists("opentelekomcloud_waf_certificate_v1.certificate_1", &certificate),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_waf_certificate_v1.certificate_1", "name", "cert_update"),
				),
			},
		},
	})
}

func testAccCheckWafCertificateV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.wafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_certificate_v1" {
			continue
		}

		_, err := certificates.Get(wafClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf certificate still exists")
		}
	}

	return nil
}

func testAccCheckWafCertificateV1Exists(n string, certificate *certificates.Certificate) resource.TestCheckFunc {
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

		found, err := certificates.Get(wafClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Waf certificate not found")
		}

		*certificate = *found

		return nil
	}
}

const testAccWafCertificateV1_basic = `
resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
	name = "cert_1"
	content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
	key = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}
`

const testAccWafCertificateV1_update = `
resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
	name = "cert_update"
	content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
	key = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}
`
