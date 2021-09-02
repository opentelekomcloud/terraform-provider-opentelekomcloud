package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/queues"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceQueueName = "opentelekomcloud_dms_queue_v1.queue_1"

func TestAccDmsQueuesV1_basic(t *testing.T) {
	var queue queues.Queue
	var queueName = fmt.Sprintf("dms_queue_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV1QueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1QueueBasic(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1QueueExists(resourceQueueName, queue),
					resource.TestCheckResourceAttr(resourceQueueName, "name", queueName),
					resource.TestCheckResourceAttr(resourceQueueName, "queue_mode", "NORMAL"),
				),
			},
		},
	})
}

func TestAccDmsQueuesV1_FIFOMode(t *testing.T) {
	var queue queues.Queue
	var queueName = fmt.Sprintf("dms_queue_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV1QueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1QueueFIFOMode(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1QueueExists(resourceQueueName, queue),
					resource.TestCheckResourceAttr(resourceQueueName, "name", queueName),
					resource.TestCheckResourceAttr(resourceQueueName, "description", "test create dms queue"),
					resource.TestCheckResourceAttr(resourceQueueName, "queue_mode", "FIFO"),
					resource.TestCheckResourceAttr(resourceQueueName, "redrive_policy", "enable"),
					resource.TestCheckResourceAttr(resourceQueueName, "max_consume_count", "80"),
				),
			},
		},
	})
}

func testAccCheckDmsV1QueueDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DmsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DMSv1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dms_queue_v1" {
			continue
		}

		_, err := queues.Get(client, rs.Primary.ID, false).Extract()
		if err == nil {
			return fmt.Errorf("dms queue still exists")
		}
	}
	return nil
}

func testAccCheckDmsV1QueueExists(n string, queue queues.Queue) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DmsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DMSv1 client: %w", err)
		}

		v, err := queues.Get(client, rs.Primary.ID, false).Extract()
		if err != nil {
			return fmt.Errorf("error getting OpenTelekomCloud DMS queue: %w", err)
		}
		if v.ID != rs.Primary.ID {
			return fmt.Errorf("dms queue not found")
		}
		queue = *v
		return nil
	}
}

func testAccDmsV1QueueBasic(queueName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dms_queue_v1" "queue_1" {
  name = "%s"
}
`, queueName)
}

func testAccDmsV1QueueFIFOMode(queueName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dms_queue_v1" "queue_1" {
  name              = "%s"
  description       = "test create dms queue"
  queue_mode        = "FIFO"
  redrive_policy    = "enable"
  max_consume_count = 80
}
`, queueName)
}
