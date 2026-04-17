package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/transfers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccLtsV2Transfer_basic(t *testing.T) {
	var (
		transfer transfers.Transfer
		name     = fmt.Sprintf("lts_transfer%s", acctest.RandString(3))
		obsName  = fmt.Sprintf("lts-obs%s", acctest.RandString(3))
		rName    = "opentelekomcloud_lts_transfer_v2.transfer"
		rc       = common.InitResourceCheck(rName, &transfer, getLtsTransferResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLtsTransferV2_basic(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "opentelekomcloud_lts_group_v2.group", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_id",
						"opentelekomcloud_lts_stream_v2.stream", "id"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_type", "OBS"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_mode", "cycle"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period", "3"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period_unit", "hour"),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.obs_bucket_name",
						"opentelekomcloud_obs_bucket.output", "bucket"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_dir_prefix_name", "lts_transfer_obs_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_prefix_name", "obs_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone", "UTC"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone_id", "Etc/GMT"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
				),
			},
			{
				Config: testLtsTransferV2_update(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "opentelekomcloud_lts_group_v2.group", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_id",
						"opentelekomcloud_lts_stream_v2.stream", "id"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_type", "OBS"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_mode", "cycle"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "DISABLE"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period", "2"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period_unit", "min"),
					resource.TestCheckResourceAttrPair(rName, "log_transfer_info.0.log_transfer_detail.0.obs_bucket_name",
						"opentelekomcloud_obs_bucket.output", "bucket"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_dir_prefix_name", "lts_transfer_obs_2_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_prefix_name", "obs_2_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone", "UTC-02:00"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone_id", "Etc/GMT+2"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testLtsTransferV2_basic(name, obsName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}
resource "opentelekomcloud_lts_stream_v2" "stream" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_obs_bucket" "output" {
  bucket        = "%[2]s"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_lts_transfer_v2" "transfer" {
  log_group_id = opentelekomcloud_lts_group_v2.group.id

  log_streams {
    log_stream_id = opentelekomcloud_lts_stream_v2.stream.id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      obs_period          = 3
      obs_period_unit     = "hour"
      obs_bucket_name     = opentelekomcloud_obs_bucket.output.bucket
      obs_dir_prefix_name = "lts_transfer_obs_"
      obs_prefix_name     = "obs_"
      obs_time_zone       = "UTC"
      obs_time_zone_id    = "Etc/GMT"
    }
  }
}
`, name, obsName)
}

func testLtsTransferV2_update(name, obsName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}
resource "opentelekomcloud_lts_stream_v2" "stream" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_obs_bucket" "output" {
  bucket        = "%[2]s"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_lts_transfer_v2" "transfer" {
  log_group_id = opentelekomcloud_lts_group_v2.group.id

  log_streams {
    log_stream_id = opentelekomcloud_lts_stream_v2.stream.id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "DISABLE"

    log_transfer_detail {
      obs_period          = 2
      obs_period_unit     = "min"
      obs_bucket_name     = opentelekomcloud_obs_bucket.output.bucket
      obs_dir_prefix_name = "lts_transfer_obs_2_"
      obs_prefix_name     = "obs_2_"
      obs_time_zone       = "UTC-02:00"
      obs_time_zone_id    = "Etc/GMT+2"
    }
  }
}
`, name, obsName)
}
