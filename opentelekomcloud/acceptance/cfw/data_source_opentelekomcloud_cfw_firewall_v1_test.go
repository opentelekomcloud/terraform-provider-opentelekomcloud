package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const firewalldataSourceName = "data.opentelekomcloud_cfw_firewall_v1.firewall_1"

func TestAccCFWFirewallV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCFWFirewallV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCFWFirewallV1DataSourceID(firewalldataSourceName),
					resource.TestCheckResourceAttrSet(firewalldataSourceName, "protect_objects.0.object_id"),
					resource.TestCheckResourceAttr(firewalldataSourceName, "service_type", "0"),
				),
			},
		},
	})
}

func testAccCheckCFWFirewallV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find CFW firewall data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CFW firewall data source ID not set")
		}

		return nil
	}
}

var testAccDataSourceCFWFirewallV1Basic = `
data "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  id = "dexc310f-7a44-437e-b535-b1bd26b491d5"
}`
