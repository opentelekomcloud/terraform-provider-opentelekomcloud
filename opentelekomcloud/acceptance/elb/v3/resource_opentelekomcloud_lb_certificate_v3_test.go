package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/certificates"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
)

const (
	resourceCertificateName  = "opentelekomcloud_lb_certificate_v3.certificate_1"
	resourceCertificateName2 = "opentelekomcloud_lb_certificate_v3.certificate_ca"

	certificate = `<<EOT
-----BEGIN CERTIFICATE-----
MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNV
BAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQsw
CQYDVQQLEwJ4eDELMAkGA1UEAxMCeHgxGjAYBgkqhkiG9w0BCQEWC3h4eEAxNjMu
Y29tMB4XDTE3MTExMzAyMjYxM1oXDTIwMTExMjAyMjYxM1owajELMAkGA1UEBhMC
eHgxCzAJBgNVBAgTAnh4MQswCQYDVQQHEwJ4eDELMAkGA1UEChMCeHgxCzAJBgNV
BAsTAnh4MQswCQYDVQQDEwJ4eDEaMBgGCSqGSIb3DQEJARYLeHh4QDE2My5jb20w
gZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMU832iM+d3FILgTWmpZBUoYcIWV
cAAYE7FsZ9LNerOyjJpyi256oypdBvGs9JAUBN5WaFk81UQx29wAyNixX+bKa0DB
WpUDqr84V1f9vdQc75v9WoujcnlKszzpV6qePPC7igJJpu4QOI362BrWzJCYQbg4
Uzo1KYBhLFxl0TovAgMBAAGjgc8wgcwwHQYDVR0OBBYEFMbTvDyvE2KsRy9zPq/J
WOjovG+WMIGcBgNVHSMEgZQwgZGAFMbTvDyvE2KsRy9zPq/JWOjovG+WoW6kbDBq
MQswCQYDVQQGEwJ4eDELMAkGA1UECBMCeHgxCzAJBgNVBAcTAnh4MQswCQYDVQQK
EwJ4eDELMAkGA1UECxMCeHgxCzAJBgNVBAMTAnh4MRowGAYJKoZIhvcNAQkBFgt4
eHhAMTYzLmNvbYIJALV96mEtVF4EMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEF
BQADgYEAASkC/1iwiALa2RU3YCxqZFEEsZZvQxikrDkDbFeoa6Tk49Fnb1f7FCW6
PTtY3HPWl5ygsMsSy0Fi3xp3jmuIwzJhcQ3tcK5gC99HWp6Kw37RL8WoB8GWFU0Q
4tHLOjBIxkZROPRhH+zMIrqUexv6fsb3NWKhnlfh1Mj5wQE4Ldo=
-----END CERTIFICATE-----
EOT
`

	privateKey = `<<EOT
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotu
eqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqL
o3J5SrM86Veqnjzwu4oCSabuEDiN+tga1syQmEG4OFM6NSmAYSxcZdE6LwIDAQAB
AoGBAJvLzJCyIsCJcKHWL6onbSUtDtyFwPViD1QrVAtQYabF14g8CGUZG/9fgheu
TXPtTDcvu7cZdUArvgYW3I9F9IBb2lmF3a44xfiAKdDhzr4DK/vQhvHPuuTeZA41
r2zp8Cu+Bp40pSxmoAOK3B0/peZAka01Ju7c7ZChDWrxleHZAkEA/6dcaWHotfGS
eW5YLbSms3f0m0GH38nRl7oxyCW6yMIDkFHURVMBKW1OhrcuGo8u0nTMi5IH9gRg
5bH8XcujlQJBAMWBQgzCHyoSeryD3TFieXIFzgDBw6Ve5hyMjUtjvgdVKoxRPvpO
kclc39QHP6Dm2wrXXHEej+9RILxBZCVQNbMCQQC42i+Ut0nHvPuXN/UkXzomDHde
h1ySsOAO4H+8Y6OSI87l3HUrByCQ7stX1z3L0HofjHqV9Koy9emGTFLZEzSdAkB7
Ei6cUKKmztkYe3rr+RcATEmwAw3tEJOHmrW5ErApVZKr2TzLMQZ7WZpIPzQRCYnY
2ZZLDuZWFFG3vW+wKKktAkAaQ5GNzbwkRLpXF1FZFuNF7erxypzstbUmU/31b7tS
i5LmxTGKL/xRYtZEHjya4Ikkkgt40q1MrUsgIYbFYMf2
-----END RSA PRIVATE KEY-----
EOT
`
)

func TestAccLBV3Certificate_basic(t *testing.T) {
	var cert certificates.Certificate

	t.Parallel()
	quotas.BookOne(t, quotas.LbCertificate)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBv3CertificateConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2CertificateExists(resourceCertificateName, &cert),
				),
			},
			{
				Config: testAccLBv3CertificateConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceCertificateName, "name", "certificate_1_updated"),
				),
			},
			{
				Config: testAccLBv3ClientCertificateConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2CertificateExists(resourceCertificateName2, &cert),
				),
			},
			{
				Config: testAccLBv3ClientCertificateConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceCertificateName2, "name", "certificate_ca_updated"),
				),
			},
		},
	})
}

func TestAccLBv3Certificate_importBasic(t *testing.T) {
	t.Parallel()
	quotas.BookOne(t, quotas.LbCertificate)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBv3CertificateConfigBasic,
			},
			{
				ResourceName:      resourceCertificateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLBV3CertificateDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elbv3.ErrCreateClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_certificate_v3" {
			continue
		}

		_, err := certificates.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("certificate still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2CertificateExists(n string, cert *certificates.Certificate) resource.TestCheckFunc {
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

		found, err := certificates.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("certificate not found")
		}

		*cert = *found

		return nil
	}
}

var testAccLBv3CertificateConfigBasic = fmt.Sprintf(`
resource "opentelekomcloud_lb_certificate_v3" "certificate_1" {
  name        = "certificate_1"
  description = "terraform test certificate"
  domain      = "www.elb.com"
  private_key = %s
  certificate = %s
}
`, privateKey, certificate)

var testAccLBv3CertificateConfigUpdate = fmt.Sprintf(`
resource "opentelekomcloud_lb_certificate_v3" "certificate_1" {
  name        = "certificate_1_updated"
  description = "terraform test certificate"
  domain      = "www.elb.com"
  private_key = %s
  certificate = %s
}
`, privateKey, certificate)

var testAccLBv3ClientCertificateConfigBasic = fmt.Sprintf(`
resource "opentelekomcloud_lb_certificate_v3" "certificate_ca" {
  name        = "certificate_client"
  description = "terraform ca test certificate"
  type        = "client"
  certificate = %s
}
`, certificate)

var testAccLBv3ClientCertificateConfigUpdate = fmt.Sprintf(`
resource "opentelekomcloud_lb_certificate_v3" "certificate_ca" {
  name        = "certificate_ca_updated"
  description = "terraform ca test certificate"
  type        = "client"
  certificate = %s
}
`, certificate)
