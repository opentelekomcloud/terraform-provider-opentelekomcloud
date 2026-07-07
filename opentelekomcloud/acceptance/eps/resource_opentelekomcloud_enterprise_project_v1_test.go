package eps

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/eps/v1/projects"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getResourceEnterpriseProject(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.EpsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("unable to create EPS client: %s", err)
	}

	return projects.Get(client, state.Primary.ID)
}

func TestAccEnterpriseProject_basic(t *testing.T) {
	var project projects.EnterpriseProject
	rName := common.RandomAccResourceName()
	updateName := rName + "update"
	resourceName := "opentelekomcloud_enterprise_project_v1.test"

	rc := common.InitResourceCheck(
		resourceName,
		&project,
		getResourceEnterpriseProject,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEnterpriseProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnterpriseProject_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
				),
			},
			{
				Config: testAccEnterpriseProject_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_flag"},
			},
		},
	})
}

func testAccCheckEnterpriseProjectDestroy(s *terraform.State) error {
	conf := common.TestAccProvider.Meta().(*cfg.Config)
	epsClient, err := conf.EpsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("unable to create EPS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_enterprise_project_v1" {
			continue
		}

		project, err := projects.Get(epsClient, rs.Primary.ID)
		if err == nil {
			if project.Status != 2 {
				return fmt.Errorf("project still active")
			}
		}
	}

	return nil
}

func testAccEnterpriseProject_basic(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_enterprise_project_v1" "test" {
  name        = "%s"
  description = "terraform test"
}`, rName)
}

func testAccEnterpriseProject_update(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_enterprise_project_v1" "test" {
  name        = "%s"
  description = "terraform test update"
}`, rName)
}
