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

func TestAccIdentityV3Project_basic(t *testing.T) {
	var project projects.Project
	var projectName = fmt.Sprintf("%s_%s", env.OS_REGION_NAME, acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Project_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("opentelekomcloud_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_project_v3.project_1", "name", &project.Name),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_project_v3.project_1", "description", &project.Description),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_identity_project_v3.project_1", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_identity_project_v3.project_1", "is_domain", "false"),
				),
			},
			{
				Config: testAccIdentityV3Project_update(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("opentelekomcloud_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_project_v3.project_1", "name", &project.Name),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_project_v3.project_1", "description", &project.Description),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_identity_project_v3.project_1", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_identity_project_v3.project_1", "is_domain", "false"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ProjectDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_project_v3" {
			continue
		}

		_, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("project still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3ProjectExists(n string, project *projects.Project) resource.TestCheckFunc {
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

		found, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if (found.ID != rs.Primary.ID) || (found.Enabled == false) {
			return fmt.Errorf("project not found")
		}

		*project = *found

		return nil
	}
}

func testAccIdentityV3Project_basic(projectName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_project_v3" "project_1" {
  name        = "%s"
  description = "A project"
}
  `, projectName)
}

func testAccIdentityV3Project_update(projectName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_project_v3" "project_1" {
  name        = "%s"
  description = "Some project"
}
  `, projectName)
}
