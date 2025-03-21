package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	fl "github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/flow-logs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getFlowResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.ErV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Enterprise Router client: %s", err)
	}

	return fl.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccFlowLog_basic(t *testing.T) {
	var (
		flowLog interface{}
		rName   = "opentelekomcloud_er_flow_log_v3.test"
		rc      = common.InitResourceCheck(rName, &flowLog, getFlowResourceFunc)

		name       = fmt.Sprintf("er-acc-api%s", acctest.RandString(3))
		updateName = fmt.Sprintf("er-acc-api-update%s", acctest.RandString(3))
		baseConfig = testaccFlowLog_base(name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFlowLog_basic(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "log_store_type", "LTS"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "opentelekomcloud_logtank_group_v2.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "opentelekomcloud_logtank_topic_v2.test", "id"),
					resource.TestCheckResourceAttr(rName, "resource_type", "attachment"),
					resource.TestCheckResourceAttrPair(rName, "resource_id", "opentelekomcloud_er_vpc_attachment_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by script"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(rName, "state"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testFlowLog_basic_update(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFlowLogImportStateFunc(rName),
			},
		},
	})
}

func testAccFlowLogImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, flowLogId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER flow log is not found in the tfstate", rsName)
		}

		instanceId = rs.Primary.Attributes["instance_id"]
		flowLogId = rs.Primary.ID
		if instanceId == "" || flowLogId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, flowLogId)
		}
		return fmt.Sprintf("%s/%s", instanceId, flowLogId), nil
	}
}

func testaccFlowLog_base(name string) string {
	bgpAsNum := acctest.RandIntRange(64512, 65534)
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_instance_v3" "test" {
  availability_zones = ["eu-de-01", "eu-de-02"]
  name               = "%[2]s"
  asn                = %[3]d
}

resource "opentelekomcloud_logtank_group_v2" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 7
}

resource "opentelekomcloud_logtank_topic_v2" "test" {
  group_id   = opentelekomcloud_logtank_group_v2.test.id
  topic_name = "%[2]s"
}

resource "opentelekomcloud_er_vpc_attachment_v3" "test" {
  instance_id            = opentelekomcloud_er_instance_v3.test.id
  vpc_id                 = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  name                   = "%[2]s"
  auto_create_vpc_routes = true

  tags = {
    foo = "bar"
  }
}
`, common.DataSourceSubnet, name, bgpAsNum)
}

func testFlowLog_basic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_flow_log_v3" "test" {
  instance_id    = opentelekomcloud_er_instance_v3.test.id
  log_store_type = "LTS"
  log_group_id   = opentelekomcloud_logtank_group_v2.test.id
  log_stream_id  = opentelekomcloud_logtank_topic_v2.test.id
  resource_type  = "attachment"
  resource_id    = opentelekomcloud_er_vpc_attachment_v3.test.id
  name           = "%[2]s"
  description    = "Created by script"
  enabled        = false
}
`, baseConfig, name)
}

func testFlowLog_basic_update(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_flow_log_v3" "test" {
  instance_id    = opentelekomcloud_er_instance_v3.test.id
  log_store_type = "LTS"
  log_group_id   = opentelekomcloud_logtank_group_v2.test.id
  log_stream_id  = opentelekomcloud_logtank_topic_v2.test.id
  resource_type  = "attachment"
  resource_id    = opentelekomcloud_er_vpc_attachment_v3.test.id
  name           = "%[2]s"
  enabled        = true
}
`, baseConfig, name)
}
