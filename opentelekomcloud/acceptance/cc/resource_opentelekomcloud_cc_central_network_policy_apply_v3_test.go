package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/policy"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getCentralNetworkPolicyApplyResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CcV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CC v3 client: %s", err)
	}
	policyId := state.Primary.Attributes["policy_id"]
	resp, err := policy.List(client, policy.ListOpts{
		CentralNetworkId: state.Primary.ID,
		ID:               []string{policyId},
	})
	if err != nil {
		return nil, err
	}
	for i := range resp.CentralNetworkPolicies {
		if resp.CentralNetworkPolicies[i].ID == policyId && resp.CentralNetworkPolicies[i].IsApplied {
			return resp.CentralNetworkPolicies[i], nil
		}
	}
	return nil, fmt.Errorf("applied central network policy (%s) not found", policyId)
}

func TestAccCcCentralNetworkPolicyApplyV3_basic(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("cc_acc_apply_%s", acctest.RandString(5))
	asn := acctest.RandIntRange(64512, 65534)
	rName := "opentelekomcloud_cc_central_network_policy_apply_v3.test"

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getCentralNetworkPolicyApplyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			testAccPreCheckCcProjectID(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcCentralNetworkPolicyApplyV3_basic(name, asn),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "central_network_id",
						"opentelekomcloud_cc_central_network_v3.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"opentelekomcloud_cc_central_network_policy_v3.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "region"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCcCentralNetworkPolicyApplyImportStateIdFunc(rName),
			},
		},
	})
}

func testAccCcCentralNetworkPolicyApplyImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["policy_id"] == "" {
			return "", fmt.Errorf("missing ID or policy_id for resource (%s)", name)
		}
		return rs.Primary.ID + "/" + rs.Primary.Attributes["policy_id"], nil
	}
}

func testAccCcCentralNetworkPolicyApplyV3_basic(name string, asn int) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cc_central_network_policy_v3" "test" {
  central_network_id = opentelekomcloud_cc_central_network_v3.test.id

  er_instances {
    project_id           = "%[2]s"
    region_id            = "%[3]s"
    enterprise_router_id = opentelekomcloud_er_instance_v3.test.id
  }
}

resource "opentelekomcloud_cc_central_network_policy_apply_v3" "test" {
  central_network_id = opentelekomcloud_cc_central_network_v3.test.id
  policy_id          = opentelekomcloud_cc_central_network_policy_v3.test.id
}
`, testAccCcCentralNetworkPolicyBase(name, asn), env.OS_PROJECT_ID, env.OS_REGION_NAME)
}
