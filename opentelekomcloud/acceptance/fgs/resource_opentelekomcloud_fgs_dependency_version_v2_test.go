package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/dependency_version"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/fgs"
)

func getDependencyVersionFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.FuncGraphV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	dependId, dependVersion := fgs.ParseDependVersionResourceId(state.Primary.ID)

	return dependency_version.Get(client, dependId, dependVersion)
}

func TestAccDependencyVersion_basic(t *testing.T) {
	var (
		obj dependency_version.DepVersionResp

		rName        = fmt.Sprintf("fgs-dep-%s", acctest.RandString(5))
		resourceName = "opentelekomcloud_fgs_dependency_version_v2.test"
		rc           = common.InitResourceCheck(resourceName, &obj, getDependencyVersionFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDependencyVersion_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Python3.9"),
					resource.TestCheckResourceAttrSet(resourceName, "link"),
					resource.TestCheckResourceAttrSet(resourceName, "dependency_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"link",
					"file",
				},
			},
		},
	})
}

func TestAccDependencyVersion_Zip(t *testing.T) {
	var (
		obj dependency_version.DepVersionResp

		rName        = fmt.Sprintf("fgs-dep-%s", acctest.RandString(5))
		resourceName = "opentelekomcloud_fgs_dependency_version_v2.test"
		rc           = common.InitResourceCheck(resourceName, &obj, getDependencyVersionFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDependencyVersion_zip(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Python3.9"),
					resource.TestCheckResourceAttrSet(resourceName, "file"),
					resource.TestCheckResourceAttrSet(resourceName, "dependency_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"link",
					"file",
				},
			},
		},
	})
}

func testAccDependencyVersion_basic(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_dependency_version_v2" "test" {
  name        = "%s"
  description = "Created by terraform script"
  runtime     = "Python3.9"
  link        = "https://regr-func-graph.obs.eu-de.otc.t-systems.com/requirements.zip"
}
`, rName)
}

func testAccDependencyVersion_zip(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_dependency_version_v2" "test" {
  name        = "%[1]s"
  description = "Created by terraform script"
  runtime     = "Python3.9"
  file        = filebase64("%[2]s")
}
`, rName, env.OS_FGS_REQ_FILE)
}
