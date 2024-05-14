package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	acls "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/acl"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceApigwAssociateAclName = "opentelekomcloud_apigw_acl_policy_associate_v2.associate"

func getAclPolicyAssociateFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	opt := acls.ListBoundOpts{
		GatewayID: state.Primary.Attributes["gateway_id"],
		ID:        state.Primary.Attributes["policy_id"],
	}
	resp, err := acls.ListAPIBoundPolicy(client, opt)
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, err
}

func TestAccAclPolicyAssociate_basic(t *testing.T) {
	var apiDetails []acls.ApiAcl
	name := fmt.Sprintf("apigw_acc_acl%s", acctest.RandString(10))
	rc := common.InitResourceCheck(
		resourceApigwAssociateAclName,
		&apiDetails,
		getAclPolicyAssociateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAclPolicyAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceApigwAssociateAclName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceApigwAssociateAclName, "policy_id"),
					resource.TestCheckResourceAttr(resourceApigwAssociateAclName, "publish_ids.#", "1"),
				),
			},
			{
				Config: testAccAclPolicyAssociate_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceApigwAssociateAclName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceApigwAssociateAclName, "policy_id"),
					resource.TestCheckResourceAttr(resourceApigwAssociateAclName, "publish_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceApigwAssociateAclName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyAssociateImportStateFunc(resourceApigwAssociateAclName),
			},
		},
	})
}

func testAccAclPolicyAssociateImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["gateway_id"] == "" || rs.Primary.Attributes["policy_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<policy_id>', but got '%s/%s'",
				rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["policy_id"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["policy_id"]), nil
	}
}

func testAccAclPolicyAssociate_basic(name string) string {
	relatedConfig := testAccApigwApi_basic(testAccApigwApi_base(name), name)
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_api_publishment_v2" "pub_one" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_environment_v2" "env_two" {
  name        = "second_env_%[2]s"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"
}

resource "opentelekomcloud_apigw_api_publishment_v2" "pub_two" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env_two.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_acl_policy_v2" "ip_rule" {
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%[2]s_rule_ip"
  type        = "PERMIT"
  entity_type = "IP"
  value       = "10.201.33.4,10.30.2.15"
}

resource "opentelekomcloud_apigw_acl_policy_associate_v2" "associate" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  policy_id  = opentelekomcloud_apigw_acl_policy_v2.ip_rule.id

  publish_ids = [
    opentelekomcloud_apigw_api_publishment_v2.pub_one.publish_id
  ]
}
`, relatedConfig, name)
}

func testAccAclPolicyAssociate_update(name string) string {
	relatedConfig := testAccApigwApi_basic(testAccApigwApi_base(name), name)
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_api_publishment_v2" "pub_one" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_environment_v2" "env_two" {
  name        = "second_env_%[2]s"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"
}

resource "opentelekomcloud_apigw_api_publishment_v2" "pub_two" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env_two.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_acl_policy_v2" "ip_rule" {
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%[2]s_rule_ip"
  type        = "PERMIT"
  entity_type = "IP"
  value       = "10.201.33.4,10.30.2.15"
}

resource "opentelekomcloud_apigw_acl_policy_associate_v2" "associate" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  policy_id  = opentelekomcloud_apigw_acl_policy_v2.ip_rule.id

  publish_ids = [
    opentelekomcloud_apigw_api_publishment_v2.pub_two.publish_id
  ]
}
`, relatedConfig, name)
}
