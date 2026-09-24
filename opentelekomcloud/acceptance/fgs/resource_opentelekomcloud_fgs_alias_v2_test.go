package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/alias"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getFgsAliasFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FuncGraphV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud FunctionGraph client: %s", err)
	}
	return alias.GetAlias(c, state.Primary.Attributes["function_urn"], state.Primary.Attributes["name"])
}

func TestAccFgsAlias_basic(t *testing.T) {
	var config alias.FuncAliasesResp
	name := fmt.Sprintf("fgs-alias-%s", acctest.RandString(5))
	rName := "opentelekomcloud_fgs_alias_v2.test"

	rc := common.InitResourceCheck(rName, &config, getFgsAliasFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFgsAliasBasic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "live"),
					resource.TestCheckResourceAttr(rName, "description", "test alias"),
					resource.TestCheckResourceAttrSet(rName, "alias_urn"),
					resource.TestCheckResourceAttrSet(rName, "last_modified"),
				),
			},
			{
				Config: testFgsAliasUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "live"),
					resource.TestCheckResourceAttr(rName, "description", "updated alias"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testFgsAliasBasic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "e42a37a22f4988ba7a681e3042e5c7d13c04e6c1"
}

resource "opentelekomcloud_fgs_publish_version_v2" "test" {
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  version      = "v1"
}

resource "opentelekomcloud_fgs_alias_v2" "test" {
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  name         = "live"
  version      = opentelekomcloud_fgs_publish_version_v2.test.version
  description  = "test alias"
}
`, name)
}

func testFgsAliasUpdate(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "e42a37a22f4988ba7a681e3042e5c7d13c04e6c1"
}

resource "opentelekomcloud_fgs_publish_version_v2" "test" {
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  version      = "v1"
}

resource "opentelekomcloud_fgs_alias_v2" "test" {
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  name         = "live"
  version      = opentelekomcloud_fgs_publish_version_v2.test.version
  description  = "updated alias"
}
`, name)
}
