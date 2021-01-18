package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestRDSVersionsV3_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRDSVersionsV3_basic,
				Check:  testAccCheckRdsFlavorV3DataSourceID("data.opentelekomcloud_rds_versions_v3.versions"),
			},
		},
	})
}

var testRDSVersionsV3_basic = `
data "opentelekomcloud_rds_versions_v3" "versions" {
  database_name = "mysql"
}
`
