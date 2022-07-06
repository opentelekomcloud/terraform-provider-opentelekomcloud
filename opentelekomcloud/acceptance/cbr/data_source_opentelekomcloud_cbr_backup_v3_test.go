package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataBackupName = "data.opentelekomcloud_cbr_backup_v3.cbr"

func TestAccCBRBackupV3DataSource_basic(t *testing.T) {
	if os.Getenv("CBR_DATA") == "" {
		t.Skip("this test is not a stable one")
	}
	vaultId := "insert_vaultId_here"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRBackupV3DataSourceBasic(vaultId),
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

func testAccCBRBackupV3DataSourceBasic(vaultId string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_cbr_backup_v3" "cbr" {
  checkpoint_id = "%s"
}
`, vaultId)
}
