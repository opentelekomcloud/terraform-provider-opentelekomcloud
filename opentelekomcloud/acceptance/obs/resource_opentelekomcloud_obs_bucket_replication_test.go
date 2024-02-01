package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const resourceReplication = "opentelekomcloud_obs_bucket_replication.test"

func TestAccObsBucketReplication_basic(t *testing.T) {
	destBucket := os.Getenv("OS_DESTINATION_BUCKET")
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketReplicationBasic(rInt, destBucket),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReplication, "id", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceReplication, "agency", "test-obs-agency"),
					resource.TestCheckResourceAttr(resourceReplication, "destination_bucket", destBucket),
					resource.TestCheckResourceAttr(resourceReplication, "rule.0.prefix", "log"),
					resource.TestCheckResourceAttr(resourceReplication, "rule.1.prefix", "imgs/"),
					resource.TestCheckResourceAttr(resourceReplication, "rule.1.storage_class", "COLD"),
					resource.TestCheckResourceAttr(resourceReplication, "rule.1.enabled", "true"),
					resource.TestCheckResourceAttr(resourceReplication, "rule.1.history_enabled", "false"),
				),
			},
			{
				Config: testAccObsBucketReplicationUpdate(rInt, destBucket),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReplication, "id", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceReplication, "agency", "test-obs-agency"),
					resource.TestCheckResourceAttr(resourceReplication, "destination_bucket", destBucket),
					resource.TestCheckResourceAttr(resourceReplication, "rule.0.delete_data", "true"),
				),
			},
		},
	})
}

func TestAccOBSReplication_importBasic(t *testing.T) {
	destBucket := os.Getenv("OS_DESTINATION_BUCKET")
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketReplicationBasic(rInt, destBucket),
			},

			{
				ResourceName:      resourceReplication,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccObsBucketReplicationBasic(randInt int, destBucket string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-%d"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_identity_agency_v3" "agency_obs" {
  name                  = "test-obs-agency"
  delegated_domain_name = "op_svc_obs"
  project_role {
    project = "%s"
    roles = [
      "OBS Administrator",
    ]
  }
}

resource "opentelekomcloud_obs_bucket_replication" "test" {
  bucket             = opentelekomcloud_obs_bucket.bucket.bucket
  destination_bucket = "%s"
  agency             = opentelekomcloud_identity_agency_v3.agency_obs.name

  rule {
    prefix = "log"
  }

  rule {
    prefix          = "imgs/"
    storage_class   = "COLD"
    enabled         = true
    history_enabled = false
  }
}
`, randInt, env.OS_TENANT_NAME, destBucket)
}

func testAccObsBucketReplicationUpdate(randInt int, destBucket string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-%d"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_identity_agency_v3" "agency_obs" {
  name                  = "test-obs-agency"
  delegated_domain_name = "op_svc_obs"
  project_role {
    project = "%s"
    roles = [
      "OBS Administrator",
    ]
  }
}

resource "opentelekomcloud_obs_bucket_replication" "test" {
  bucket             = opentelekomcloud_obs_bucket.bucket.bucket
  destination_bucket = "%s"
  agency             = opentelekomcloud_identity_agency_v3.agency_obs.name

  rule {
    prefix      = "test"
    delete_data = true
  }
}
`, randInt, env.OS_TENANT_NAME, destBucket)
}
