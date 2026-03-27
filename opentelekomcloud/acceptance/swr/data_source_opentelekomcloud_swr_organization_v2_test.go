package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestSwrOrganizationV2DS_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testSwrOrganizationV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testSwrOrganizationV2DSBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceOrgName, "organizations.0.auth", "7"),
					resource.TestCheckResourceAttrSet(dataSourceOrg2Name, "organizations.0.organization_id"),
				),
			},
		},
	})
}

const dataSourceOrgName = "data.opentelekomcloud_swr_organization_v2.org_1"
const dataSourceOrg2Name = "data.opentelekomcloud_swr_organization_v2.org_2"

var testSwrOrganizationV2DSBasic = fmt.Sprintf(
	`
resource opentelekomcloud_swr_organization_v2 org {
  name = "%[1]s"
}

data opentelekomcloud_swr_organization_v2 org_1 {
  depends_on = [opentelekomcloud_swr_organization_v2.org]
  name = "%[1]s"
}

data opentelekomcloud_swr_organization_v2 org_2 {}
`, name)
