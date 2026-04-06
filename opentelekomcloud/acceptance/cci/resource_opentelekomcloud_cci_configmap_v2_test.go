package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/configmap"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getV2ConfigMapResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CciV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CCI v2 client: %s", err)
	}
	return configmap.Get(client, state.Primary.Attributes["namespace"], state.Primary.Attributes["name"])
}

func TestAccV2ConfigMap_basic(t *testing.T) {
	var cm configmap.ConfigMap
	rName := fmt.Sprintf("cci-configmap-%s", acctest.RandString(5))
	nsName := fmt.Sprintf("cci-ns-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_cci_configmap_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&cm,
		getV2ConfigMapResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2ConfigMap_basic(nsName, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "api_version", "cci/v2"),
					resource.TestCheckResourceAttr(resourceName, "kind", "ConfigMap"),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
					resource.TestCheckResourceAttrSet(resourceName, "data.%"),
					resource.TestCheckResourceAttr(resourceName, "data.key1", "value1"),
				),
			},
			{
				Config: testAccV2ConfigMap_update(nsName, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "data.key1", "updated_value1"),
					resource.TestCheckResourceAttr(resourceName, "data.key2", "value2"),
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

func testAccV2ConfigMap_basic(nsName, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cci_configmap_v2" "test" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"

  data = {
    "key1" = "value1"
  }
}
`, testAccV2Namespace_basic(nsName), rName)
}

func testAccV2ConfigMap_update(nsName, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cci_configmap_v2" "test" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"

  data = {
    "key1" = "updated_value1"
    "key2" = "value2"
  }
}
`, testAccV2Namespace_basic(nsName), rName)
}
