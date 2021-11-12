package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

const dataCertificateName = "data.opentelekomcloud_lb_certificate_v3.certificate_1"

func TestAccLBCertificateV3_basic(t *testing.T) {
	t.Parallel()
	th.AssertNoErr(t, quotas.LbCertificate.Acquire())
	defer quotas.LbCertificate.Release()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLBv3CertificateConfigBasic,
			},
			{
				Config: testAccLBv3CertificateByID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataCertificateName, "id"),
					resource.TestCheckResourceAttr(dataCertificateName, "domain", "www.elb.com"),
					resource.TestCheckResourceAttr(dataCertificateName, "description", "terraform test certificate"),
					resource.TestCheckResourceAttr(dataCertificateName, "name", "certificate_1"),
				),
			},
		},
	})
}

var testAccLBv3CertificateByID = fmt.Sprintf(`
%s

data "opentelekomcloud_lb_certificate_v3" "certificate_1" {
  id = opentelekomcloud_lb_certificate_v3.certificate_1.id
}
`, testAccLBv3CertificateConfigBasic)
