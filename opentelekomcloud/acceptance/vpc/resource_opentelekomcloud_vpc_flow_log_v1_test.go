package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/flowlogs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	fl_name        = "vpc_flow_log1"
	fl_update_name = "vpc_flow_log_update"
)

func TestAccVpcFlowLogV1_basic(t *testing.T) {
	var flowlog flowlogs.FlowLog

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcFlowLogV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcFlowLogV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcFlowLogV1Exists("opentelekomcloud_vpc_flow_log_v1.flow_logl", &flowlog),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_flow_log_v1.flow_logl", "name", fl_name),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_flow_log_v1.flow_logl", "resource_type", "vpc"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_flow_log_v1.flow_logl", "traffic_type", "all"),
				),
			},
			{
				Config: testAccVpcFlowLogV1_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_flow_log_v1.flow_logl", "name", fl_update_name),
				),
			},
		},
	})
}

func testAccCheckVpcFlowLogV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_flow_log_v1" {
			continue
		}

		_, err := flowlogs.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VPC flow log still exists")
		}
	}

	return nil
}

func testAccCheckVpcFlowLogV1Exists(n string, flowlog *flowlogs.FlowLog) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		vpcClient, err := config.NetworkingV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud Vpc client: %s", err)
		}

		found, err := flowlogs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("VPC flow log not found")
		}

		*flowlog = *found

		return nil
	}
}

var testAccVpcFlowLogV1_basic = fmt.Sprintf(`
resource "opentelekomcloud_logtank_group_v2" "log_group1" {
  group_name  = "vpc_group"
}

resource "opentelekomcloud_logtank_topic_v2" "log_topic1" {
  group_id = opentelekomcloud_logtank_group_v2.log_group1.id
  topic_name = "vpc_topic"
}

resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "172.16.0.0/16"
}

resource "opentelekomcloud_vpc_flow_log_v1" "flow_logl" {
  name = "%s"
  description   = "this is a flow log from testacc"
  resource_type = "vpc"
  resource_id   = opentelekomcloud_vpc_v1.vpc_1.id
  traffic_type  = "all"
  log_group_id  = opentelekomcloud_logtank_group_v2.log_group1.id
  log_topic_id  = opentelekomcloud_logtank_topic_v2.log_topic1.id
}
`, fl_name)

var testAccVpcFlowLogV1_update = fmt.Sprintf(`
resource "opentelekomcloud_logtank_group_v2" "log_group1" {
  group_name  = "vpc_group"
}

resource "opentelekomcloud_logtank_topic_v2" "log_topic1" {
  group_id = opentelekomcloud_logtank_group_v2.log_group1.id
  topic_name = "vpc_topic"
}

resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "172.16.0.0/16"
}

resource "opentelekomcloud_vpc_flow_log_v1" "flow_logl" {
  name = "%s"
  description   = "this is a flow log from testacc update"
  resource_type = "vpc"
  resource_id   = opentelekomcloud_vpc_v1.vpc_1.id
  traffic_type  = "all"
  log_group_id  = opentelekomcloud_logtank_group_v2.log_group1.id
  log_topic_id  = opentelekomcloud_logtank_topic_v2.log_topic1.id
}
`, fl_update_name)
