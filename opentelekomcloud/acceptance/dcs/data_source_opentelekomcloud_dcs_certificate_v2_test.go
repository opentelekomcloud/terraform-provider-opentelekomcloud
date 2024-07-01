package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataCert = "data.opentelekomcloud_dcs_certificate_v2.cert"

func TestAccDcsCertificateV2DataSource_basic(t *testing.T) {
	instanceName := fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsCertificateV2DataSourceBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsAZV1DataSourceID(dataCert),
					resource.TestCheckResourceAttrSet(dataCert, "file_name"),
					resource.TestCheckResourceAttrSet(dataCert, "link"),
					resource.TestCheckResourceAttrSet(dataCert, "bucket_name"),
					resource.TestCheckResourceAttrSet(dataCert, "certificate"),
				),
			},
		},
	})
}
func testAccDcsCertificateV2DataSourceBasic(instanceName string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_dcs_certificate_v2" "cert" {
  instance_id = opentelekomcloud_dcs_instance_v2.instance_1.id
}
`, testAccDcsV2InstanceBasic(instanceName))
}
