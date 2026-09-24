package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/alias"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/function"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getPublishVersionFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FuncGraphV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud FunctionGraph client: %s", err)
	}
	listResp, err := alias.ListVersion(c, alias.ListVersionOpts{FuncUrn: state.Primary.Attributes["function_urn"]})
	if err != nil {
		return nil, fmt.Errorf("error fetching versions list: %s", err)
	}
	for _, v := range listResp.Functions {
		if v.Version == state.Primary.Attributes["version"] {
			return v, nil
		}
	}
	return nil, fmt.Errorf("error fetching fgs function version (%s): %s", state.Primary.Attributes["version"], err)
}

func TestAccFgsPublishVersion_basic(t *testing.T) {
	var config function.FuncGraph
	name := fmt.Sprintf("fgs-publish-version-%s", acctest.RandString(5))
	rName := "opentelekomcloud_fgs_publish_version_v2.test"

	rc := common.InitResourceCheck(rName, &config, getPublishVersionFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFgsPublishVersionBasic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "version", "v1"),
					resource.TestCheckResourceAttrSet(rName, "func_name"),
					resource.TestCheckResourceAttrSet(rName, "last_modified"),
				),
			},
		},
	})
}

func testFgsPublishVersionBasic(name string) string {
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
`, name)
}
