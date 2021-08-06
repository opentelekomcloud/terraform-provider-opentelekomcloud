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

func TestAccDmsQueuesV1_basic(t *testing.T) {
	var queue queues.Queue
	var queueName = fmt.Sprintf("dms_queue_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV1QueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Queue_basic(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1QueueExists("opentelekomcloud_dms_queue_v1.queue_1", queue),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dms_queue_v1.queue_1", "name", queueName),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dms_queue_v1.queue_1", "queue_mode", "NORMAL"),
				),
			},
		},
	})
}

func TestAccDmsQueuesV1_FIFOmode(t *testing.T) {
	var queue queues.Queue
	var queueName = fmt.Sprintf("dms_queue_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV1QueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Queue_FIFOmode(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1QueueExists("opentelekomcloud_dms_queue_v1.queue_1", queue),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dms_queue_v1.queue_1", "name", queueName),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dms_queue_v1.queue_1", "description", "test create dms queue"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dms_queue_v1.queue_1", "queue_mode", "FIFO"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dms_queue_v1.queue_1", "redrive_policy", "enable"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dms_queue_v1.queue_1", "max_consume_count", "80"),
				),
			},
		},
	})
}

func testAccCheckDmsV1QueueDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	dmsClient, err := config.DmsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud queue client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dms_queue_v1" {
			continue
		}

		_, err := queues.Get(dmsClient, rs.Primary.ID, false).Extract()
		if err == nil {
			return fmt.Errorf("the Dms queue still exists.")
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
		dmsClient, err := config.DmsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud queue client: %s", err)
		}

		v, err := queues.Get(dmsClient, rs.Primary.ID, false).Extract()
		if err != nil {
			return fmt.Errorf("error getting OpenTelekomCloud queue: %s, err: %s", rs.Primary.ID, err)
		}
		if v.ID != rs.Primary.ID {
			return fmt.Errorf("the Dms queue not found.")
		}
		queue = *v
		return nil
	}
}

func testAccDmsV1Queue_basic(queueName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dms_queue_v1" "queue_1" {
  name = "%s"
}
	`, queueName)
}

func testAccDmsV1Queue_FIFOmode(queueName string) string {
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
