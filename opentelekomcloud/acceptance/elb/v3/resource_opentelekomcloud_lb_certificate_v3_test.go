package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/certificates"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
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
MIIB4TCCAYugAwIBAgIUPXCpWJCiy5mI79NIfenl5KNWPzkwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMTExMDIxMDM3MjBaFw0yMTEy
MDIxMDM3MjBaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwXDANBgkqhkiG9w0BAQEF
AANLADBIAkEAu+qgVpV6mqbaGW1Qn6eDPzhwentQPPiXwG1665M9+gjW4pUQ0Rud
Bc0fkUU/O+Q0UMT8ZV/I2hSenCVyJoyPEwIDAQABo1MwUTAdBgNVHQ4EFgQUtItI
IAXZDIEfuvCX7AY3s//wlI8wHwYDVR0jBBgwFoAUtItIIAXZDIEfuvCX7AY3s//w
lI8wDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAANBAEkgP/JlpVKc4j+Z
KRcMa7RAXYJqCbRxtpqRU7OOAhDmBnldtS5CTMoh1r7TOGMfM1Npa+kGV5QnjRzI
FzFSymo=
-----END CERTIFICATE-----
EOT
`

	privateKey = `<<EOT
-----BEGIN RSA PRIVATE KEY-----
MIIBUwIBADANBgkqhkiG9w0BAQEFAASCAT0wggE5AgEAAkEAu+qgVpV6mqbaGW1Q
n6eDPzhwentQPPiXwG1665M9+gjW4pUQ0RudBc0fkUU/O+Q0UMT8ZV/I2hSenCVy
JoyPEwIDAQABAkAbyksEAv8qt9oxQHVX5xIF23bm5i2rlqf6kTZIeHIF89/NNJ2E
sejiqFIWqPc5a00Scn+ymdCvjC25JVyup9cBAiEA4a+7WhPmgS54yNHjwkG2pflz
cfH1V7qPqlBKIGLwZbMCIQDVKCsZ6eoNdQoLVmK0zii8XDCgL8HWMrm/bytbYM9B
IQIgVdcAXKebEeF6IW/rwDQ8Y2644UsVdTPJdw8o0p6vLw8CIDqm191EiPt09fOS
rIxVoc3ajCK3oV2ADa5IN6ToKX8hAiBPuNCCIYcZz0tAzWX7I1OYMI3UhJjtrESg
mYFrsJ4gHw==
-----END RSA PRIVATE KEY-----
EOT
`
)

func TestAccLBV3Certificate_basic(t *testing.T) {
	var cert certificates.Certificate

	t.Parallel()
	th.AssertNoErr(t, quotas.LbCertificate.Acquire())
	defer quotas.LbCertificate.Release()

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
	th.AssertNoErr(t, quotas.LbCertificate.Acquire())
	defer quotas.LbCertificate.Release()

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
