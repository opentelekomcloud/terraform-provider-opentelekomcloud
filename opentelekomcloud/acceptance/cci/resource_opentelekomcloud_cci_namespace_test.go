package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/namespace"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getV2NamespaceResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CciV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CCI v2 client: %s", err)
	}
	return namespace.Get(client, state.Primary.ID)
}

func TestAccV2Namespace_basic(t *testing.T) {
	var ns namespace.Namespace
	rName := fmt.Sprintf("cci-ns-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_cci_namespace_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&ns,
		getV2NamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Namespace_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "api_version", "cci/v2"),
					resource.TestCheckResourceAttr(resourceName, "kind", "Namespace"),
					resource.TestCheckResourceAttr(resourceName, "finalizers.0", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccV2Namespace_basic(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cci_namespace_v2" "test" {
  name = "%s"
}
`, rName)
}
