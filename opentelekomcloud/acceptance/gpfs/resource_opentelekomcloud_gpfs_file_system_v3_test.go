package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gpfs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	accenv "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceGPFSFileSystemName = "opentelekomcloud_gpfs_file_system_v3.fs"

func getGPFSFileSystemFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewGPFSClient(accenv.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GPFS client: %s", err)
	}
	output, err := client.ListFS(&gpfs.ListFSInput{BucketType: "SFS"})
	if err != nil {
		return nil, err
	}
	for i := range output.Buckets {
		if output.Buckets[i].Name == state.Primary.ID {
			return &output.Buckets[i], nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func TestAccGPFSFileSystemV3_basic(t *testing.T) {
	name := fmt.Sprintf("terraform-acc-%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceGPFSFileSystemName,
		&gpfs.Bucket{},
		getGPFSFileSystemFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGPFSFileSystemV3Basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceGPFSFileSystemName, "name", name),
					resource.TestCheckResourceAttr(resourceGPFSFileSystemName, "redundancy", "3az"),
					resource.TestCheckResourceAttrSet(resourceGPFSFileSystemName, "creation_date"),
				),
			},
			{
				ResourceName:      resourceGPFSFileSystemName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGPFSFileSystemV3Basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_gpfs_file_system_v3" "fs" {
  name       = "%s"
  redundancy = "3az"
}
`, name)
}
