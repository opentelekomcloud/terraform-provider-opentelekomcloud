package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gpfs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	accenv "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceGPFSAccessRuleName = "opentelekomcloud_gpfs_access_rule_v3.rule"

func getGPFSAccessRuleFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewGPFSClient(accenv.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GPFS client: %s", err)
	}
	output, err := client.GetFSAccessRules(state.Primary.Attributes["fs_name"])
	if err != nil {
		return nil, err
	}
	for i := range output.Statement {
		if output.Statement[i].Condition.SourceVPC == state.Primary.Attributes["vpc_id"] {
			return &output.Statement[i], nil
		}
	}
	return nil, fmt.Errorf("GPFS access rule for VPC %s not found", state.Primary.Attributes["vpc_id"])
}

func TestAccGPFSAccessRuleV3_basic(t *testing.T) {
	fsName := fmt.Sprintf("terraform-acc-%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceGPFSAccessRuleName,
		&gpfs.AccessRule{},
		getGPFSAccessRuleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGPFSAccessRuleV3Basic(fsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceGPFSAccessRuleName, "fs_name", fsName),
					resource.TestCheckResourceAttr(resourceGPFSAccessRuleName, "action", "FullControl"),
				),
			},
			{
				ResourceName:      resourceGPFSAccessRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGPFSAccessRuleV3Basic(fsName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_gpfs_file_system_v3" "test" {
  name = "%s"
}

resource "opentelekomcloud_gpfs_access_rule_v3" "rule" {
  fs_name = opentelekomcloud_gpfs_file_system_v3.test.name
  vpc_id  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  action  = "FullControl"
}
`, common.DataSourceSubnet, fsName)
}
