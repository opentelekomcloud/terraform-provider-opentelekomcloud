package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccOpenStackIdentityV3ProjectDataSource_basic(t *testing.T) {
	fallbackProjectName := fmt.Sprintf("tf_test_%s", acctest.RandString(5))
	fallbackProjectDescription := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckIdentityProjectV3DataSource(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityProjectV3DataSource_basic(fallbackProjectName, fallbackProjectDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectDataSourceID("data.opentelekomcloud_identity_project_v3.project_1"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_identity_project_v3.project_1", "name"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_identity_project_v3.project_1", "domain_id"),
				),
			},
		},
	})
}

func TestAccOpenStackIdentityV3ProjectDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityProjectV3DataSource_project_empty(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_identity_project_v3.project_1", "id"),
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

func testAccOpenStackIdentityProjectV3DataSource_project_empty() string {
	return `
data "opentelekomcloud_identity_project_v3" "project_1" {
}
`
}

func testAccOpenStackIdentityProjectV3DataSource_basic(name, description string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_identity_projects_v3" "all" {}

locals {
  has_existing_project  = length(data.opentelekomcloud_identity_projects_v3.all.projects) > 0
  selected_project_name = local.has_existing_project ? data.opentelekomcloud_identity_projects_v3.all.projects[0].name : opentelekomcloud_identity_project_v3.project_1[0].name
}

resource "opentelekomcloud_identity_project_v3" "project_1" {
  count       = local.has_existing_project ? 0 : 1
  name        = "%s"
  description = "%s"
}

data "opentelekomcloud_identity_project_v3" "project_1" {
  name = local.selected_project_name
}
`, name, description)
}

func testAccOpenStackIdentityProjectV3DataSource_byID(name, description string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_identity_projects_v3" "all" {}

locals {
  has_existing_project = length(data.opentelekomcloud_identity_projects_v3.all.projects) > 0
  selected_project_id  = local.has_existing_project ? data.opentelekomcloud_identity_projects_v3.all.projects[0].project_id : opentelekomcloud_identity_project_v3.project_1[0].id
}

resource "opentelekomcloud_identity_project_v3" "project_1" {
  count       = local.has_existing_project ? 0 : 1
  name        = "%s"
  description = "%s"
}

data "opentelekomcloud_identity_project_v3" "project_1" {
  id = local.selected_project_id
}
`, name, description)
}

func TestAccOpenStackIdentityV3ProjectDataSource_byID(t *testing.T) {
	fallbackProjectName := fmt.Sprintf("tf_test_%s", acctest.RandString(5))
	fallbackProjectDescription := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckIdentityProjectV3DataSource(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityProjectV3DataSource_byID(fallbackProjectName, fallbackProjectDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectDataSourceID("data.opentelekomcloud_identity_project_v3.project_1"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_identity_project_v3.project_1", "name"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_identity_project_v3.project_1", "id"),
				),
			},
		},
	})
}

func testAccPreCheckIdentityProjectV3DataSource(t *testing.T) {
	common.TestAccPreCheck(t)

	if testAccHasAnyIdentityProject(t) {
		return
	}

	common.TestAccPreCheckAdminOnly(t)
}

func testAccHasAnyIdentityProject(t *testing.T) bool {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		t.Fatalf("error creating OpenStack identity client: %s", err)
	}

	allPages, err := projects.List(identityClient, projects.ListOpts{}).AllPages()
	if err != nil {
		t.Fatalf("unable to query projects: %s", err)
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		t.Fatalf("unable to retrieve projects: %s", err)
	}

	return len(allProjects) > 0
}
