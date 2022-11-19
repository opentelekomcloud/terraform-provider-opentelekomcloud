package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cts/v1/tracker"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const trackerResource = "opentelekomcloud_cts_tracker_v1.tracker_v1"

func TestAccCTSTrackerV1_basic(t *testing.T) {
	var ctsTracker tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1Basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists(trackerResource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerResource, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(trackerResource, "file_prefix_name", "yO8Q"),
					resource.TestCheckResourceAttr(trackerResource, "is_lts_enabled", "false"),
				),
			},
			{
				Config: testAccCTSTrackerV1Update(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists(trackerResource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerResource, "file_prefix_name", "yO8Q1"),
					resource.TestCheckResourceAttr(trackerResource, "is_lts_enabled", "true"),
				),
			},
		},
	})
}

func TestAccCTSTrackerV1_importBasic(t *testing.T) {
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1ImportBasic(bucketName),
			},

			{
				ResourceName:      trackerResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCTSTrackerV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	ctsClient, err := config.CtsV1Client(env.OS_TENANT_NAME)
	if err != nil {
		return fmt.Errorf("error creating cts client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cts_tracker_v1" {
			continue
		}

		ctsTracker, err := tracker.Get(ctsClient, "system")
		if err != nil {
			return fmt.Errorf("failed to retrieve CTS list: %s", err)
		}

		if ctsTracker.BucketName != "" {
			return fmt.Errorf("failed to delete CTS tracker")
		}
	}

	return nil
}

func testAccCheckCTSTrackerV1Exists(n string, trackers *tracker.Tracker, projectName cfg.ProjectName) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CtsV1Client(projectName)
		if err != nil {
			return fmt.Errorf("error creating cts client: %s", err)
		}

		ctsTracker, err := tracker.Get(client, "system")
		if err != nil {
			return err
		}

		if ctsTracker.TrackerName != rs.Primary.ID {
			return fmt.Errorf("CTS tracker not found")
		}

		trackers = ctsTracker

		return nil
	}
}

func testAccCTSTrackerV1Basic(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket = "%s"
  acl    = "public-read"
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name = "yO8Q"
}
`, bucketName)
}

func testAccCTSTrackerV1Update(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name = "yO8Q1"
  is_lts_enabled   = true
}
`, bucketName)
}

func testAccCTSTrackerV1ImportBasic(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name = "yO8Q1"
  is_lts_enabled   = false
}
`, bucketName)
}
