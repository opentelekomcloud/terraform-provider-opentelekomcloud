package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceWafDedicatedRefTablesV1_basic(t *testing.T) {
	var name = fmt.Sprintf("wafd_rt_%s", acctest.RandString(5))
	dataSourceName := "data.opentelekomcloud_waf_dedicated_reference_tables_v1.table"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedRefTablesV1_ds(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReferenceTablesId(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.conditions.0"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.created_at"),
				),
			},
		},
	})
}

func testAccCheckReferenceTablesId(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("can't find OpenTelekomCloud WAF reference tables data source: %s", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("the OpenTelekomCloud WAF reference tables data source ID not set")
		}
		return nil
	}
}

func testAccWafDedicatedRefTablesV1_ds(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_waf_dedicated_reference_tables_v1" "table" {
  depends_on = [opentelekomcloud_waf_dedicated_reference_table_v1.table]
}
`, testAccWafDedicatedRefTablesV1_basic(name))
}
