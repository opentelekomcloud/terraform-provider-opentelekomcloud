package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccDataSourceObsBucket_basic(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test")
	dataSourceName := "data.opentelekomcloud_obs_bucket.bucket"

	var bucket *obs.BaseModel

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { common.TestAccPreCheck(t) },
		ProviderFactories:         common.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceObsBucketInit(randomName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists("opentelekomcloud_obs_bucket.bucket"),
				),
			},
			{
				Config: testAccDataSourceObsBucketBasic(randomName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketDataSourceExists(dataSourceName, bucket),
					resource.TestCheckResourceAttr(dataSourceName, "bucket", randomName),
				),
			},
		},
	})
}

func testAccCheckObsBucketDataSourceExists(n string, obj *obs.BaseModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find OBS bucket data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("OBS bucket data source ID not set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NewObjectStorageClient(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OBS client: %w", err)
		}
		out, err := client.HeadBucket(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed getting OBS Bucket (%s): %s", rs.Primary.ID, err)
		}

		obj = out

		return nil
	}
}

func testAccDataSourceObsBucketInit(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  storage_class = "STANDARD"
  acl           = "public-read"
}
`, name)
}

func testAccDataSourceObsBucketBasic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%[1]s"
  storage_class = "STANDARD"
  acl           = "public-read"
}

data "opentelekomcloud_obs_bucket" "bucket" {
  bucket = "%[1]s"
}
`, name)
}
