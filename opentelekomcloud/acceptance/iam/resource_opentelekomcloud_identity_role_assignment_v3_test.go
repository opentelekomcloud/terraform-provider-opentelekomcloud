package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"

	"github.com/opentelekomcloud/gophertelekomcloud/pagination"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/roles"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/iam"
)

func TestAccIdentityV3RoleAssignment_basic(t *testing.T) {
	if os.Getenv("OS_TEST_ROLENAME") == "" || os.Getenv("OS_TEST_DOMAINID") == "" {
		t.Skip("OS envs not set for tests.")
	}
	var role roles.Role
	var group groups.Group
	var project projects.Project
	groupName := acctest.RandomWithPrefix("group")
	roleName := os.Getenv("OS_TEST_ROLENAME")
	domainId := os.Getenv("OS_TEST_DOMAINID")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3RoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRoleAssignProjectV3_basic("eu-de", groupName, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleAssignmentExists("opentelekomcloud_identity_role_assignment_v3.role_assignment_1", &role, &group, &project),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_role_assignment_v3.role_assignment_1", "project_id", &project.ID),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_role_assignment_v3.role_assignment_1", "group_id", &group.ID),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_role_assignment_v3.role_assignment_1", "role_id", &role.ID),
				),
			},
			{
				Config: testAccIdentityRoleAssignDomainV3_basic(domainId, groupName, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleAssignmentExistsShort("opentelekomcloud_identity_role_assignment_v3.role_assignment_1"),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_role_assignment_v3.role_assignment_1", "group_id", &group.ID),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_role_assignment_v3.role_assignment_1", "role_id", &role.ID),
				),
			},
			{
				Config: testAccIdentityRoleAssignAllProjectsV3_basic(domainId, groupName, roleName),
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrPtr(
					"opentelekomcloud_identity_role_assignment_v3.role_assignment_1", "group_id", &group.ID),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_role_assignment_v3.role_assignment_1", "role_id", &role.ID),
				),
			},
		},
	})
}

func testAccCheckIdentityV3RoleAssignmentDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_role_assignment_v3" {
			continue
		}

		_, err := roles.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("role assignment still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3RoleAssignmentExists(n string, role *roles.Role, group *groups.Group, project *projects.Project) resource.TestCheckFunc {
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
			return fmt.Errorf("error creating OpenTelekomCloud identity client: %s", err)
		}

		domainID, projectID, groupID, roleID := iam.ExtractRoleAssignmentID(rs.Primary.ID)

		opts := roles.ListAssignmentsOpts{
			GroupID:        groupID,
			ScopeDomainID:  domainID,
			ScopeProjectID: projectID,
		}

		pager := roles.ListAssignments(identityClient, opts)
		var assignment roles.RoleAssignment

		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			assignmentList, err := roles.ExtractRoleAssignments(page)
			if err != nil {
				return false, err
			}

			for _, a := range assignmentList {
				if a.ID == roleID {
					assignment = a
					return false, nil
				}
			}

			return true, nil
		})
		if err != nil {
			return err
		}

		p, err := projects.Get(identityClient, projectID).Extract()
		if err != nil {
			return fmt.Errorf("project not found")
		}
		*project = *p
		g, err := groups.Get(identityClient, groupID).Extract()
		if err != nil {
			return fmt.Errorf("group not found")
		}
		*group = *g
		r, err := roles.Get(identityClient, assignment.ID).Extract()
		if err != nil {
			return fmt.Errorf("role not found")
		}
		*role = *r

		return nil
	}
}

func testAccCheckIdentityV3RoleAssignmentExistsShort(n string) resource.TestCheckFunc {
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
			return fmt.Errorf("error creating OpenTelekomCloud identity client: %s", err)
		}

		domainID, projectID, groupID, roleID := iam.ExtractRoleAssignmentID(rs.Primary.ID)

		opts := roles.ListAssignmentsOpts{
			GroupID:        groupID,
			ScopeDomainID:  domainID,
			ScopeProjectID: projectID,
		}

		pager := roles.ListAssignments(identityClient, opts)
		var assignment roles.RoleAssignment

		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			assignmentList, err := roles.ExtractRoleAssignments(page)
			if err != nil {
				return false, err
			}

			for _, a := range assignmentList {
				if a.ID == roleID {
					assignment = a
					return false, nil
				}
			}

			return true, nil
		})
		if err != nil {
			return err
		}

		_, err = roles.Get(identityClient, assignment.ID).Extract()
		if err != nil {
			return fmt.Errorf("role not found")
		}

		return nil
	}
}

func testAccIdentityRoleAssignProjectV3_basic(project_name string, group_name string, role_name string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_identity_project_v3" "project_1" {
  name = "%s"
}

resource "opentelekomcloud_identity_group_v3" "group_1" {
  name = "%s"
}

data "opentelekomcloud_identity_role_v3" "role_1" {
  name = "%s"
}

resource "opentelekomcloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id   = opentelekomcloud_identity_group_v3.group_1.id
  project_id = data.opentelekomcloud_identity_project_v3.project_1.id
  role_id    = data.opentelekomcloud_identity_role_v3.role_1.id
}
`, project_name, group_name, role_name)
}

func testAccIdentityRoleAssignDomainV3_basic(domain_id string, group_name string, role_name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_group_v3" "group_1" {
  name = "%s"
}

data "opentelekomcloud_identity_role_v3" "role_1" {
  name = "%s"
}

resource "opentelekomcloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id  = opentelekomcloud_identity_group_v3.group_1.id
  domain_id = "%s"
  role_id   = data.opentelekomcloud_identity_role_v3.role_1.id
}
`, group_name, role_name, domain_id)
}

func testAccIdentityRoleAssignAllProjectsV3_basic(domain_id string, group_name string, role_name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_group_v3" "group_1" {
  name = "%s"
}

data "opentelekomcloud_identity_role_v3" "role_1" {
  name = "%s"
}

resource "opentelekomcloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id     = opentelekomcloud_identity_group_v3.group_1.id
  domain_id    = "%s"
  role_id      = data.opentelekomcloud_identity_role_v3.role_1.id
  all_projects = true
}
`, group_name, role_name, domain_id)
}
