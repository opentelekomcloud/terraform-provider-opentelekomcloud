package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/diagnosis"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceRocketMQInstanceDiagnosisV2Name = "opentelekomcloud_dms_rocketmq_instance_diagnosis_v2.test"

func getDmsRocketMQInstanceDiagnosisFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DmsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud DMSv2 client: %s", err)
	}
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	v, err := diagnosis.Get(client, parts[1])
	if err != nil {
		return nil, err
	}
	if v.Status == "deleted" {
		return nil, golangsdk.ErrDefault404{}
	}
	return v, nil
}

func TestAccDmsRocketMQInstanceDiagnosisV2_basic(t *testing.T) {
	var report diagnosis.Report
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))

	rc := common.InitResourceCheck(resourceRocketMQInstanceDiagnosisV2Name, &report, getDmsRocketMQInstanceDiagnosisFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQInstanceDiagnosisV2Basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceRocketMQInstanceDiagnosisV2Name, "instance_id",
						resourceRocketMQInstanceV2Name, "id"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceDiagnosisV2Name, "group_name", name),
					resource.TestCheckResourceAttrSet(resourceRocketMQInstanceDiagnosisV2Name, "report_id"),
					resource.TestCheckResourceAttrSet(resourceRocketMQInstanceDiagnosisV2Name, "status"),
					resource.TestCheckResourceAttrSet(resourceRocketMQInstanceDiagnosisV2Name, "created_at"),
					resource.TestCheckResourceAttrSet(resourceRocketMQInstanceDiagnosisV2Name, "abnormal_item_sum"),
					resource.TestCheckResourceAttrSet(resourceRocketMQInstanceDiagnosisV2Name, "faulted_node_sum"),
					resource.TestCheckResourceAttrSet(resourceRocketMQInstanceDiagnosisV2Name, "consumer_nums"),
				),
			},
			{
				ResourceName:            resourceRocketMQInstanceDiagnosisV2Name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_id_list"},
			},
		},
	})
}

func testAccDmsRocketMQInstanceDiagnosisV2Basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_instance_diagnosis_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  group_name  = opentelekomcloud_dms_rocketmq_consumer_group_v2.test.name
}
`, testAccDmsRocketMQConsumerGroupV2Basic(name))
}
