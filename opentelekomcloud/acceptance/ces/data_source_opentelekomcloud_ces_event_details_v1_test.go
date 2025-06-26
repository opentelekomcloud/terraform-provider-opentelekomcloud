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

const eventDetailsdataSourceName = "data.opentelekomcloud_ces_event_details_v1.event_details_1"

func TestAccCESEventDetailsV1DataSource_basic(t *testing.T) {
	if os.Getenv("RUN_CES_EVENTS") == "" {
		t.Skip("CES tests are not suitable for CI runner unless specifically flagged")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCESEventDetailsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCESEventDetailsV1DataSourceID(eventDetailsdataSourceName),
					resource.TestCheckResourceAttr(eventDetailsdataSourceName, "event_info.0.event_source", "SYS.ECS"),
				),
			},
		},
	})
}

func testAccCheckCESEventDetailsV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find CES Event details data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CES event details data source ID not set")
		}

		return nil
	}
}

var testCESEventDetailsBasic = fmt.Sprintf(`
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

data "opentelekomcloud_ces_event_details_v1" "event_details_1" {
  depends_on = [opentelekomcloud_ces_event_report_v1.event_report_1]
  event_name = "Test_Acc_tf_event"
  event_type = "EVENT.CUSTOM"
  from       = %d
  to         = %d
}
`, time.Now().Unix()*1000, time.Now().Unix()*1000-100000, time.Now().Unix()*1000+10000)
