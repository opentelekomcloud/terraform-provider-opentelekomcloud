package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccIdentityV3Projects_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectsIdentityV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectsExist("data.opentelekomcloud_identity_projects_v3.all"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ProjectsExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenStack identity client: %s", err)
		}

		allPages, err := projects.List(identityClient, projects.ListOpts{}).AllPages()
		if err != nil {
			return fmt.Errorf("unable to query projects: %s", err)
		}

		allProjects, err := projects.ExtractProjects(allPages)
		if err != nil {
			return fmt.Errorf("unable to retrieve projects: %s", err)
		}

		if len(allProjects) == 0 {
			return fmt.Errorf("project not found")
		}

		return nil
	}
}

const testAccProjectsIdentityV3DataSource_basic = `
data "opentelekomcloud_identity_projects_v3" "all" {
}
`
