package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/central_network"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getCentralNetworkResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CcV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CC v3 client: %s", err)
	}
	return central_network.Get(client, state.Primary.ID)
}

func TestAccCcCentralNetworkV3_basic(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("cc_acc_cn_%s", acctest.RandString(5))
	updateName := name + "_updated"
	rName := "opentelekomcloud_cc_central_network_v3.test"

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getCentralNetworkResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcCentralNetworkV3_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform acceptance test"),
					resource.TestCheckResourceAttr(rName, "state", "AVAILABLE"),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "region"),
				),
			},
			{
				Config: testAccCcCentralNetworkV3_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "updated by terraform acceptance test"),
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

func testAccCcCentralNetworkV3_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cc_central_network_v3" "test" {
  name        = "%s"
  description = "created by terraform acceptance test"
}
`, name)
}

func testAccCcCentralNetworkV3_update(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cc_central_network_v3" "test" {
  name        = "%s"
  description = "updated by terraform acceptance test"
}
`, name)
}
