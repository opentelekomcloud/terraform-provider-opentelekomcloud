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

func getCentralNetworkPolicyResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CcV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CC v3 client: %s", err)
	}
	resp, err := policy.List(client, policy.ListOpts{
		CentralNetworkId: state.Primary.Attributes["central_network_id"],
		ID:               []string{state.Primary.ID},
	})
	if err != nil {
		return nil, err
	}
	for i := range resp.CentralNetworkPolicies {
		if resp.CentralNetworkPolicies[i].ID == state.Primary.ID {
			return resp.CentralNetworkPolicies[i], nil
		}
	}
	return nil, fmt.Errorf("central network policy (%s) not found", state.Primary.ID)
}

func TestAccCcCentralNetworkPolicyV3_basic(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("cc_acc_pol_%s", acctest.RandString(5))
	asn := acctest.RandIntRange(64512, 65534)
	rName := "opentelekomcloud_cc_central_network_policy_v3.test"

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getCentralNetworkPolicyResourceFunc,
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
				Config: testAccCcCentralNetworkPolicyV3_basic(name, asn),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "central_network_id",
						"opentelekomcloud_cc_central_network_v3.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "er_instances.0.enterprise_router_id",
						"opentelekomcloud_er_instance_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "state", "AVAILABLE"),
					resource.TestCheckResourceAttr(rName, "is_applied", "false"),
					resource.TestCheckResourceAttrSet(rName, "document_template_version"),
					resource.TestCheckResourceAttrSet(rName, "version"),
					resource.TestCheckResourceAttrSet(rName, "region"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCcCentralNetworkPolicyImportStateIdFunc(rName),
			},
		},
	})
}

func testAccCcCentralNetworkPolicyImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		centralNetworkId := rs.Primary.Attributes["central_network_id"]
		if centralNetworkId == "" || rs.Primary.ID == "" {
			return "", fmt.Errorf("missing central_network_id or ID for resource (%s)", name)
		}
		return centralNetworkId + "/" + rs.Primary.ID, nil
	}
}

func testAccPreCheckCcProjectID(t *testing.T) {
	if env.OS_PROJECT_ID == "" {
		t.Skip("OS_PROJECT_ID must be set for CC central network policy acceptance tests")
	}
}

func testAccCcCentralNetworkPolicyBase(name string, asn int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_er_instance_v3" "test" {
  availability_zones = ["eu-de-01", "eu-de-02"]

  name = "%[1]s"
  asn  = %[2]d

  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true
}

resource "opentelekomcloud_cc_central_network_v3" "test" {
  name        = "%[1]s"
  description = "created by terraform acceptance test"
}
`, name, asn)
}

func testAccCcCentralNetworkPolicyV3_basic(name string, asn int) string {
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
`, testAccCcCentralNetworkPolicyBase(name, asn), env.OS_PROJECT_ID, env.OS_REGION_NAME)
}
