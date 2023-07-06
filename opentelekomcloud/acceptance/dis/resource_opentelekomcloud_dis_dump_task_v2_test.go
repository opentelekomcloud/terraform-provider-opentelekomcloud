package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dis/v2/dump"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceDumpName = "opentelekomcloud_dis_dump_task_v2.task_1"

func TestAccDisDumpV2_basic(t *testing.T) {
	var cls dump.GetTransferTaskResponse
	var streamName = fmt.Sprintf("dis_stream_%s", acctest.RandString(5))
	var appName = fmt.Sprintf("dis_app_%s", acctest.RandString(5))
	var taskName = fmt.Sprintf("dis_task_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDisV2DumpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDisV2DumpBasic(streamName, appName, taskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDisV2DumpExists(resourceDumpName, &cls),
					resource.TestCheckResourceAttr(resourceDumpName, "name", taskName),
					resource.TestCheckResourceAttr(resourceDumpName, "destination", "OBS"),
					resource.TestCheckResourceAttr(resourceDumpName, "status", "PAUSED"),
				),
			},
			{
				Config: testAccDisV2SDumpBasicUpdated(streamName, appName, taskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDisV2DumpExists(resourceDumpName, &cls),
					resource.TestCheckResourceAttr(resourceDumpName, "name", taskName),
					resource.TestCheckResourceAttr(resourceDumpName, "destination", "OBS"),
					resource.TestCheckResourceAttr(resourceDumpName, "status", "RUNNING"),
				),
			},
		},
	})
}

func testAccCheckDisV2DumpDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DisV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating DISv2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dis_dump_task_v2" {
			continue
		}

		_, err := dump.GetTransferTask(client, dump.GetTransferTaskOpts{
			StreamName: rs.Primary.Attributes["stream_name"],
			TaskName:   rs.Primary.ID,
		})
		if err == nil {
			return fmt.Errorf("DIS dump task still exists")
		}
	}
	return nil
}

func testAccCheckDisV2DumpExists(n string, cls *dump.GetTransferTaskResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DisV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating DISv2 client: %w", err)
		}

		v, err := dump.GetTransferTask(client, dump.GetTransferTaskOpts{
			StreamName: rs.Primary.Attributes["stream_name"],
			TaskName:   rs.Primary.ID,
		})
		if err != nil {
			return fmt.Errorf("error getting dump task (%s): %w", rs.Primary.ID, err)
		}

		if v.TaskName != rs.Primary.ID {
			return fmt.Errorf("DIS dump task not found")
		}
		*cls = *v
		return nil
	}
}

func testAccDisV2DumpBasic(streamName string, appName string, taskName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dis_stream_v2" "stream_1" {
  name                           = "%s"
  partition_count                = 3
  stream_type                    = "COMMON"
  retention_period               = 24
  auto_scale_min_partition_count = 1
  auto_scale_max_partition_count = 4
  compression_format             = "zip"

  data_type = "BLOB"

  tags = {
    foo = "bar"
  }
}

resource "opentelekomcloud_dis_app_v2" "app_1" {
  name = "%s"
}

resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "tf-test-dis-bucket"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_dis_dump_task_v2" "task_1" {
  stream_name = opentelekomcloud_dis_stream_v2.stream_1.name
  destination = "OBS"
  action      = "stop"

  obs_destination_descriptor {
    task_name             = "%s"
    agency_name           = "dis_admin_agency"
    deliver_time_interval = 30
    consumer_strategy     = "LATEST"
    file_prefix           = "_pf"
    partition_format      = "yyyy/MM/dd/HH/mm"
    obs_bucket_path       = opentelekomcloud_obs_bucket.bucket.bucket
    destination_file_type = "text"
    record_delimiter      = "|"
  }
}
`, streamName, appName, taskName)
}

func testAccDisV2SDumpBasicUpdated(streamName string, appName string, taskName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dis_stream_v2" "stream_1" {
  name                           = "%s"
  partition_count                = 3
  stream_type                    = "COMMON"
  retention_period               = 24
  auto_scale_min_partition_count = 1
  auto_scale_max_partition_count = 4
  compression_format             = "zip"

  data_type = "BLOB"

  tags = {
    foo = "bar"
  }
}

resource "opentelekomcloud_dis_app_v2" "app_1" {
  name = "%s"
}

resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "tf-test-dis-bucket"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_dis_dump_task_v2" "task_1" {
  stream_name = opentelekomcloud_dis_stream_v2.stream_1.name
  destination = "OBS"
  action      = "start"

  obs_destination_descriptor {
    task_name             = "%s"
    agency_name           = "dis_admin_agency"
    deliver_time_interval = 30
    consumer_strategy     = "LATEST"
    file_prefix           = "_pf"
    partition_format      = "yyyy/MM/dd/HH/mm"
    obs_bucket_path       = opentelekomcloud_obs_bucket.bucket.bucket
    destination_file_type = "text"
    record_delimiter      = "|"
  }
}
`, streamName, appName, taskName)
}
