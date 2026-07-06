package cc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCcCentralNetworkPoliciesV3DataSource_basic(t *testing.T) {
	name := fmt.Sprintf("cc_acc_pol_ds_%s", acctest.RandString(5))
	asn := acctest.RandIntRange(64512, 65534)
	dataSourceName := "data.opentelekomcloud_cc_central_network_policies_v3.test"
	dc := common.InitDataSourceCheck(dataSourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			testAccPreCheckCcProjectID(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCcCentralNetworkPoliciesV3DataSource_basic(name, asn),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "policies.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.id"),
					resource.TestCheckResourceAttr(dataSourceName, "policies.0.state", "AVAILABLE"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.version"),
					resource.TestCheckResourceAttr(dataSourceName, "policies.0.document.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.document.0.er_instances.#"),
				),
			},
		},
	})
}

func testAccCcCentralNetworkPoliciesV3DataSource_basic(name string, asn int) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_cc_central_network_policies_v3" "test" {
  central_network_id = opentelekomcloud_cc_central_network_v3.test.id
  is_applied         = "false"

  depends_on = [opentelekomcloud_cc_central_network_policy_v3.test]
}
`, testAccCcCentralNetworkPolicyV3_basic(name, asn))
}
