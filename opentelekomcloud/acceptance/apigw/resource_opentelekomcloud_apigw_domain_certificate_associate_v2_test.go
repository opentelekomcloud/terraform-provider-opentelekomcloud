package acceptance

import (
	"fmt"
	"os"
	"testing"

	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/domain"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	accenv "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/apigw"
)

const resourceNameDomainCertAssoc = "opentelekomcloud_apigw_domain_certificate_associate_v2.assoc"

func getDomainCertificateAssociateFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.APIGWV2Client(accenv.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	gatewayId, groupId, domainId, certificateId, err := apigw.ParseDomainCertificateAssociateId(state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error parsing association ID: %s", err)
	}
	return domain.GetCertificate(client, domain.CertificateOpts{
		GatewayID:     gatewayId,
		GroupID:       groupId,
		DomainID:      domainId,
		CertificateID: certificateId,
	})
}

func TestAccAPIGWv2DomainCertificateAssociate_basic(t *testing.T) {
	gatewayID := os.Getenv("OS_GATEWAY_ID")
	if gatewayID == "" {
		t.Skip("`OS_GATEWAY_ID` needs to be defined")
	}
	domainName := fmt.Sprintf("terraform-acc-%s.example.com", acctest.RandString(5))
	certName := fmt.Sprintf("apigw_domain_cert_%s", acctest.RandString(5))
	certContent, privateKey, err := openstack.GenerateTestCertKeyPair(domainName)
	th.AssertNoErr(t, err)

	rc := common.InitResourceCheck(
		resourceNameDomainCertAssoc,
		nil,
		getDomainCertificateAssociateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2DomainCertificateAssociateBasic(gatewayID, certName, certContent, privateKey, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNameDomainCertAssoc, "bound", "true"),
				),
			},
			{
				ResourceName:      resourceNameDomainCertAssoc,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAPIGWv2DomainCertificateAssociateBasic(gatewayId, certName, certContent, privateKey, domainName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_apigw_group_v2" "group" {
  instance_id = "%[1]s"
  name        = "test_acc_group"
}

resource "opentelekomcloud_apigw_certificate_v2" "cert" {
  name        = "%[2]s"
  content     = <<-EOT
%[3]s
EOT
  private_key = <<-EOT
%[4]s
EOT
}

resource "opentelekomcloud_apigw_domain_v2" "domain" {
  gateway_id             = "%[1]s"
  group_id               = opentelekomcloud_apigw_group_v2.group.id
  name                   = "%[5]s"
}

resource "opentelekomcloud_apigw_domain_certificate_associate_v2" "assoc" {
  gateway_id     = "%[1]s"
  group_id       = opentelekomcloud_apigw_group_v2.group.id
  domain_id      = opentelekomcloud_apigw_domain_v2.domain.id
  certificate_id = opentelekomcloud_apigw_certificate_v2.cert.id
}
`, gatewayId, certName, certContent, privateKey, domainName)
}
