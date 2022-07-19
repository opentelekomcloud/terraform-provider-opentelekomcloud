package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenStackIdentityV3ProjectsDataSource_basic(t *testing.T) {
	userRegion := "eu-de"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIdentityProjectsV3DataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectsDataSourceID("data.opentelekomcloud_identity_projects_v3.projects_data"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_identity_projects_v3.projects_data", "region", userRegion),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ProjectsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find project data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("project data source ID not set")
		}

		return nil
	}
}

func testIdentityProjectsV3DataSource_basic() string {
	return `data "opentelekomcloud_identity_projects_v3" "projects_data" {}`
}
