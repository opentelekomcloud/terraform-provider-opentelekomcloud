package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataBackupName = "data.opentelekomcloud_cbr_backup_v3.cbr"

func TestAccCBRBackupV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRBackupV3DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRBackupV3DataSourceID(dataBackupName),
				),
			},
		},
	})
}

func testAccCheckCBRBackupV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find backup data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("backup data source ID not set ")
		}

		return nil
	}
}

func testAccCBRBackupV3DataSourceBasic() string {
	return fmt.Sprintf(`
data "opentelekomcloud_cbr_backup_v3" "cbr" {
	checkpoint_id = "29b011d9-e0fa-4869-b727-e6184575081e"
}
`)
}
