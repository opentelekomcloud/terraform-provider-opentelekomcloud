package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDmsRocketMQInstanceDiagnosesV2DataSource_basic(t *testing.T) {
	dataName := "data.opentelekomcloud_dms_rocketmq_instance_diagnoses_v2.test"
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))
	dc := common.InitDataSourceCheck(dataName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQInstanceDiagnosesV2DataSourceBasic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataName, "reports.#", "1"),
					resource.TestCheckResourceAttrPair(dataName, "reports.0.report_id",
						resourceRocketMQInstanceDiagnosisV2Name, "report_id"),
					resource.TestCheckResourceAttr(dataName, "reports.0.group_name", name),
					resource.TestCheckResourceAttrPair(dataName, "reports.0.status",
						resourceRocketMQInstanceDiagnosisV2Name, "status"),
					resource.TestCheckResourceAttrSet(dataName, "reports.0.created_at"),
					resource.TestCheckResourceAttrSet(dataName, "reports.0.consumer_nums"),
					resource.TestCheckResourceAttrSet(dataName, "reports.0.abnormal_item_sum"),
					resource.TestCheckResourceAttrSet(dataName, "reports.0.faulted_node_sum"),
				),
			},
		},
	})
}

func testAccDmsRocketMQInstanceDiagnosesV2DataSourceBasic(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_dms_rocketmq_instance_diagnoses_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id

  depends_on = [opentelekomcloud_dms_rocketmq_instance_diagnosis_v2.test]
}
`, testAccDmsRocketMQInstanceDiagnosisV2Basic(name))
}
