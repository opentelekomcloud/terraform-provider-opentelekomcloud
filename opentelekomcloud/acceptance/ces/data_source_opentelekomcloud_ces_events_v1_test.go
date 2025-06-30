package acceptance

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const eventsdataSourceName = "data.opentelekomcloud_ces_events_v1.events_1"

func TestAccCESEventsV1DataSource_basic(t *testing.T) {
	if os.Getenv("RUN_CES_EVENTS") == "" {
		t.Skip("CES tests are not suitable for CI runner unless specifically flagged")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCESEventsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCESEventsV1DataSourceID(eventsdataSourceName),
					resource.TestCheckResourceAttr(eventsdataSourceName, "meta_data.0.total", "1"),
					resource.TestCheckResourceAttr(eventsdataSourceName, "events.0.latest_event_source", "SYS.ECS"),
				),
			},
		},
	})
}

func testAccCheckCESEventsV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find CES Events data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CES events data source ID not set")
		}

		return nil
	}
}

var testCESEventsBasic = fmt.Sprintf(`
resource "opentelekomcloud_ces_event_report_v1" "event_report_1" {
  event_name   = "Test_Acc_tf_event"
  event_source = "SYS.ECS"
  time         = %d
  detail {
    content     = "This is a test event run by TF tests"
    event_state = "normal"
    event_level = "Info"
  }
}

data "opentelekomcloud_ces_events_v1" "events_1" {
  depends_on = [opentelekomcloud_ces_event_report_v1.event_report_1]
  event_type = "EVENT.CUSTOM"
  from       = %d
  to         = %d
  limit      = 1
}
`, time.Now().Unix()*1000, time.Now().Unix()*1000-100000, time.Now().Unix()*1000+10000)
