package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dis/v2/checkpoints"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceCheckpointName = "opentelekomcloud_dis_checkpoint_v2.checkpoint_1"

func TestAccDisCheckpointV2_basic(t *testing.T) {
	var streamName = fmt.Sprintf("dis_stream_%s", acctest.RandString(5))
	var appName = fmt.Sprintf("dis_app_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDisV2CheckpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDisV2CheckpointBasic(streamName, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDisV2CheckpointExists(resourceCheckpointName),
					resource.TestCheckResourceAttr(resourceCheckpointName, "metadata", "my_first_checkpoint"),
				),
			},
		},
	})
}

func testAccCheckDisV2CheckpointDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DisV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating DISv2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dis_checkpoint_v2" {
			continue
		}

		_, err := checkpoints.GetCheckpoint(client, checkpoints.GetCheckpointOpts{
			StreamName:     rs.Primary.ID,
			AppName:        rs.Primary.Attributes["app_name"],
			PartitionId:    rs.Primary.Attributes["partition_id"],
			CheckpointType: rs.Primary.Attributes["checkpoint_type"],
		})
		if err == nil {
			return fmt.Errorf("DIS checkpoint still exists")
		}
	}
	return nil
}

func testAccCheckDisV2CheckpointExists(n string) resource.TestCheckFunc {
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

		v, err := checkpoints.GetCheckpoint(client, checkpoints.GetCheckpointOpts{
			StreamName:     rs.Primary.ID,
			AppName:        rs.Primary.Attributes["app_name"],
			PartitionId:    rs.Primary.Attributes["partition_id"],
			CheckpointType: rs.Primary.Attributes["checkpoint_type"],
		})
		if err != nil {
			return fmt.Errorf("error getting checkpoint (%s): %w", rs.Primary.ID, err)
		}

		if v.SequenceNumber == "" {
			return fmt.Errorf("DIS checkpoint not found")
		}
		return nil
	}
}

func testAccDisV2CheckpointBasic(streamName string, appName string) string {
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

resource "opentelekomcloud_dis_checkpoint_v2" "checkpoint_1" {
  app_name        = opentelekomcloud_dis_app_v2.app_1.name
  stream_name     = opentelekomcloud_dis_stream_v2.stream_1.name
  partition_id    = "0"
  sequence_number = "0"
  metadata        = "my_first_checkpoint"
}
`, streamName, appName)
}
