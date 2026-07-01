package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gcb "github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/global_connection_bandwidth"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getGlobalConnectionBandwidthResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CcV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CC v3 client: %s", err)
	}
	return gcb.Get(client, state.Primary.ID)
}

func TestAccCcGlobalConnectionBandwidthV3_basic(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("cc_acc_gcb_%s", acctest.RandString(5))
	updateName := name + "_updated"
	rName := "opentelekomcloud_cc_global_connection_bandwidth_v3.test"

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getGlobalConnectionBandwidthResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcGlobalConnectionBandwidthV3_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform acceptance test"),
					resource.TestCheckResourceAttr(rName, "type", "Region"),
					resource.TestCheckResourceAttr(rName, "bordercross", "false"),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bwd"),
					resource.TestCheckResourceAttr(rName, "size", "5"),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttrSet(rName, "admin_state"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "region"),
				),
			},
			{
				Config: testAccCcGlobalConnectionBandwidthV3_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "updated by terraform acceptance test"),
					resource.TestCheckResourceAttr(rName, "size", "10"),
					resource.TestCheckResourceAttr(rName, "sla_level", "Au"),
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

func testAccCcGlobalConnectionBandwidthV3_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cc_global_connection_bandwidth_v3" "test" {
  name        = "%s"
  description = "created by terraform acceptance test"
  type        = "Region"
  bordercross = false
  charge_mode = "bwd"
  size        = 5
}
`, name)
}

func testAccCcGlobalConnectionBandwidthV3_update(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cc_global_connection_bandwidth_v3" "test" {
  name        = "%s"
  description = "updated by terraform acceptance test"
  type        = "Region"
  bordercross = false
  charge_mode = "bwd"
  size        = 10
  sla_level   = "Au"
}
`, name)
}
