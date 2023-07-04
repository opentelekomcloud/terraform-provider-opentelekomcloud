package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dis/v2/streams"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceName = "opentelekomcloud_dis_stream_v2.stream_1"

func TestAccDisStreamV2_basic(t *testing.T) {
	var cls streams.DescribeStreamResponse
	var streamName = fmt.Sprintf("dis_stream_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDisV2StreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDisV2StreamBasic(streamName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDisV2StreamExists(resourceName, &cls),
					resource.TestCheckResourceAttr(resourceName, "name", streamName),
					resource.TestCheckResourceAttr(resourceName, "partitions.#", "3"),
				),
			},
			{
				Config: testAccDisV2StreamBasicUpdated(streamName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDisV2StreamExists(resourceName, &cls),
					resource.TestCheckResourceAttr(resourceName, "name", streamName),
					resource.TestCheckResourceAttr(resourceName, "partitions.#", "4"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDisV2StreamDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DisV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating DISv2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dis_stream_v2" {
			continue
		}

		_, err := streams.GetStream(client, streams.GetStreamOpts{StreamName: rs.Primary.ID})
		if err == nil {
			return fmt.Errorf("DIS stream still exists")
		}
	}
	return nil
}

func testAccCheckDisV2StreamExists(n string, cls *streams.DescribeStreamResponse) resource.TestCheckFunc {
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

		v, err := streams.GetStream(client, streams.GetStreamOpts{StreamName: rs.Primary.ID})
		if err != nil {
			return fmt.Errorf("error getting stream (%s): %w", rs.Primary.ID, err)
		}

		if v.StreamName != rs.Primary.ID {
			return fmt.Errorf("DIS stream not found")
		}
		*cls = *v
		return nil
	}
}

func testAccDisV2StreamBasic(streamName string) string {
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
`, streamName)
}

func testAccDisV2StreamBasicUpdated(streamName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dis_stream_v2" "stream_1" {
  name                           = "%s"
  partition_count                = 4
  stream_type                    = "COMMON"
  retention_period               = 24
  auto_scale_min_partition_count = 1
  auto_scale_max_partition_count = 4
  compression_format             = "zip"

  data_type = "BLOB"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, streamName)
}
