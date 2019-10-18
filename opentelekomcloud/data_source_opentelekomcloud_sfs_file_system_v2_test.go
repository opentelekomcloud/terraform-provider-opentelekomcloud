package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccOTCSFSFileSystemV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2DataSourceID("data.opentelekomcloud_sfs_file_system_v2.shares"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_sfs_file_system_v2.shares", "name", "sfs-c2c-1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_sfs_file_system_v2.shares", "status", "available"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_sfs_file_system_v2.shares", "size", "1"),
				),
			},
		},
	})
}

func testAccCheckSFSFileSystemV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find share file data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("share file data source ID not set ")
		}

		return nil
	}
}

var testAccSFSFileSystemV2DataSource_basic = `
resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
	share_proto = "NFS"
	size=1
	name="sfs-c2c-1"
  	availability_zone="eu-de-01"
	access_to="%s"
  	access_type="cert"
  	access_level="rw"
	description="sfs_c2c_test-file"
}
data "opentelekomcloud_sfs_file_system_v2" "shares" {
  id = "${opentelekomcloud_sfs_file_system_v2.sfs_1.id}"
}
`
