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

func getBucketInventoryResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewObjectStorageClient(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OBS client: %s", err)
	}

	input := obs.GetBucketInventoryInput{
		BucketName:        state.Primary.Attributes["bucket"],
		InventoryConfigId: state.Primary.ID,
	}

	requestResp, err := client.GetBucketInventory(input)
	if err != nil {
		return nil, fmt.Errorf("error getting bucket inventory: %s", err)
	}

	return requestResp.BucketInventoryConfiguration, nil
}

func TestAccObsBucketInventory_basic(t *testing.T) {
	var (
		obsInventory  = obs.BucketInventoryConfiguration{}
		rInt          = acctest.RandIntRange(4, 200)
		inventoryName = "opentelekomcloud_obs_bucket_inventory.inventory"
		inventoryId   = acctest.RandString(5)
		rc            = common.InitResourceCheck(
			inventoryName,
			&obsInventory,
			getBucketInventoryResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccObsInventoryBasic(rInt, inventoryId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(inventoryName, "is_enabled", "false"),
					resource.TestCheckResourceAttr(inventoryName, "frequency", "Daily"),
					resource.TestCheckResourceAttr(inventoryName, "included_object_versions", "All"),
					resource.TestCheckResourceAttr(inventoryName, "destination.0.format", "CSV"),
				),
			},
			{
				Config: testAccObsInventoryUpdate(rInt, inventoryId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(inventoryName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(inventoryName, "frequency", "Weekly"),
					resource.TestCheckResourceAttr(inventoryName, "included_object_versions", "Current"),
					resource.TestCheckResourceAttr(inventoryName, "destination.0.format", "CSV"),
					resource.TestCheckResourceAttr(inventoryName, "destination.0.prefix", "test-"),
					resource.TestCheckResourceAttr(inventoryName, "filter_prefix", "test-filter-prefix"),
				),
			},
			{
				ResourceName:      inventoryName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOBSInventoryImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"configuration_id",
				},
			},
		},
	})
}

func testAccOBSInventoryImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var bucket string
		var inventory string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_obs_bucket" {
				bucket = rs.Primary.ID
			} else if rs.Type == "opentelekomcloud_obs_bucket_inventory" {
				inventory = rs.Primary.ID
			}
		}
		if bucket == "" || inventory == "" {
			return "", fmt.Errorf("resource not found: %s/%s", bucket, inventory)
		}
		return fmt.Sprintf("%s/%s", bucket, inventory), nil
	}
}

func testAccObsInventoryBasic(rInt int, configName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_obs_bucket_inventory" "inventory" {
  bucket           = opentelekomcloud_obs_bucket.bucket.bucket
  configuration_id = "%[2]s"
  is_enabled       = false
  frequency        = "Daily"
  destination {
    bucket = opentelekomcloud_obs_bucket.bucket.bucket
    format = "CSV"
  }
  included_object_versions = "All"
}
`, testAccObsBucketBasic(rInt), configName)
}

func testAccObsInventoryUpdate(rInt int, configName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_obs_bucket_inventory" "inventory" {
  bucket           = opentelekomcloud_obs_bucket.bucket.bucket
  configuration_id = "%[2]s"
  is_enabled       = true
  frequency        = "Weekly"
  destination {
    bucket = opentelekomcloud_obs_bucket.bucket.bucket
    format = "CSV"
    prefix = "test-"
  }
  filter_prefix            = "test-filter-prefix"
  included_object_versions = "Current"
}
`, testAccObsBucketBasic(rInt), configName)
}
