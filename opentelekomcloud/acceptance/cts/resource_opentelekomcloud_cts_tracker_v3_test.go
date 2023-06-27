package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cts/v3/tracker"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const trackerV3Resource = "opentelekomcloud_cts_tracker_v3.tracker_v3"

func TestAccCTSTrackerV3_basic(t *testing.T) {
	var ctsTracker tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV3Basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV3Exists(trackerV3Resource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerV3Resource, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(trackerV3Resource, "file_prefix_name", "yO8Q"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_lts_enabled", "false"),
					resource.TestCheckResourceAttr(trackerV3Resource, "status", "disabled"),
				),
			},
			{
				Config: testAccCTSTrackerV3Update(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV3Exists(trackerV3Resource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerV3Resource, "file_prefix_name", "yO8Q1"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_lts_enabled", "true"),
					resource.TestCheckResourceAttr(trackerV3Resource, "status", "enabled"),
				),
			},
		},
	})
}

func TestAccCTSTrackerV3_importBasic(t *testing.T) {
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV3ImportBasic(bucketName),
			},

			{
				ResourceName:      trackerV3Resource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCTSTrackerV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	ctsClient, err := config.CtsV3Client(env.OS_TENANT_NAME)
	if err != nil {
		return fmt.Errorf("error creating cts client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cts_tracker_v3" {
			continue
		}

		ctsTracker, err := tracker.List(ctsClient, "system")
		if err != nil {
			return fmt.Errorf("failed to retrieve CTS list: %s", err)
		}

		if len(ctsTracker) != 0 {
			return fmt.Errorf("failed to delete CTS tracker")
		}
	}

	return nil
}

func testAccCheckCTSTrackerV3Exists(n string, trackers *tracker.Tracker, projectName cfg.ProjectName) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CtsV3Client(projectName)
		if err != nil {
			return fmt.Errorf("error creating cts client: %s", err)
		}

		ctsTracker, err := tracker.List(client, "system")
		if err != nil {
			return err
		}

		if ctsTracker[0].TrackerName != rs.Primary.ID {
			return fmt.Errorf("CTS tracker not found")
		}

		trackers = &ctsTracker[0]

		return nil
	}
}

func testAccCTSTrackerV3Basic(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket = "%s"
  acl    = "public-read"
}

resource "opentelekomcloud_cts_tracker_v3" "tracker_v3" {
  bucket_name      = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name = "yO8Q"
  status           = "disabled"
}
`, bucketName)
}

func testAccCTSTrackerV3Update(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_cts_tracker_v3" "tracker_v3" {
  bucket_name      = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name = "yO8Q1"
  is_lts_enabled   = true
  status           = "enabled"
}
`, bucketName)
}

func testAccCTSTrackerV3ImportBasic(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_cts_tracker_v3" "tracker_v3" {
  bucket_name      = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name = "yO8Q1"
  is_lts_enabled   = false
  status           = "enabled"
}
`, bucketName)
}
