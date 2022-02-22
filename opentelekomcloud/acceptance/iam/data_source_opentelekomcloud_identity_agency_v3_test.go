package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/agency"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataSourceAgencyName = "data.opentelekomcloud_identity_agency_v3.agency_1"

func TestDataSourceIdentityAgencyV3_basic(t *testing.T) {
	var a agency.Agency
	name := tools.RandomString("ag", 5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			TestAccIdentityV3AgencyPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3AgencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceIdentityAgencyV3Basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3AgencyExists(resourceAgencyName, &a),
					resource.TestCheckResourceAttr(resourceAgencyName, "name", name),

					testAccCheckIdentityV3AgencyExists(dataSourceAgencyName, &a),
					resource.TestCheckResourceAttr(dataSourceAgencyName, "name", name),
				),
			},
		},
	})
}

var delegatedDomainName = os.Getenv("OS_DELEGATED_DOMAIN_NAME")

func testDataSourceIdentityAgencyV3Basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_agency_v3" "agency_1" {
  name                  = "%[1]s"
  delegated_domain_name = "%s"
  project_role {
    project = "%s"
    roles = [
      "KMS Administrator",
    ]
  }
}

data "opentelekomcloud_identity_agency_v3" "agency_1" {
  name = "%[1]s"
}
`, name, delegatedDomainName, env.OS_TENANT_NAME)
}
