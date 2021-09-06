package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenStackIdentityV3ProjectDataSource_basic(t *testing.T) {
	projectName := fmt.Sprintf("tf_test_%s", acctest.RandString(5))
	projectDescription := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityProjectV3DataSource_project(projectName, projectDescription),
			},
			{
				Config: testAccOpenStackIdentityProjectV3DataSource_basic(projectName, projectDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectDataSourceID("data.opentelekomcloud_identity_project_v3.project_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_project_v3.project_1", "name", projectName),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_identity_project_v3.project_1", "description", projectDescription),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_identity_project_v3.project_1", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_identity_project_v3.project_1", "is_domain", "false"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ProjectDataSourceID(n string) resource.TestCheckFunc {
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

func testAccOpenStackIdentityProjectV3DataSource_project(name, description string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_project_v3" "project_1" {
  name        = "%s"
  description = "%s"
}
`, name, description)
}

func testAccOpenStackIdentityProjectV3DataSource_basic(name, description string) string {
	return fmt.Sprintf(`
	%s

data "opentelekomcloud_identity_project_v3" "project_1" {
  name = opentelekomcloud_identity_project_v3.project_1.name
}
`, testAccOpenStackIdentityProjectV3DataSource_project(name, description))
}
