package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	throttlingpolicy "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/tr_policy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceApigwAssociateName = "opentelekomcloud_apigw_throttling_policy_associate_v2.associate"

func getAssociateFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud APIG v2 client: %s", err)
	}
	opt := throttlingpolicy.ListBoundOpts{
		GatewayID:  state.Primary.Attributes["gateway_id"],
		ThrottleID: state.Primary.Attributes["policy_id"],
	}
	return throttlingpolicy.ListAPIBoundPolicy(c, opt)
}

func TestAccThrottlingPolicyAssociate_basic(t *testing.T) {
	var apiDetails []throttlingpolicy.ApiThrottle
	name := fmt.Sprintf("apigw_acc_thpa%s", acctest.RandString(10))
	rc := common.InitResourceCheck(
		resourceApigwAssociateName,
		&apiDetails,
		getAssociateFunc,
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccThrottlingPolicyAssociateV2_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceApigwAssociateName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceApigwAssociateName, "policy_id"),
					resource.TestCheckResourceAttr(resourceApigwAssociateName, "publish_ids.#", "1"),
				),
			}, {
				Config: testAccThrottlingPolicyAssociateV2_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceApigwAssociateName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceApigwAssociateName, "policy_id"),
					resource.TestCheckResourceAttr(resourceApigwAssociateName, "publish_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceApigwAssociateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccThrottlingPolicyAssociateV2_basic(name string) string {
	relatedConfig := testAccApigwApi_basic(testAccApigwApi_base(name), name)
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_environment_v2" "env_two" {
  name        = "second_env_%[2]s"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"
}

resource "opentelekomcloud_apigw_api_publishment_v2" "pub_one" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_api_publishment_v2" "pub_two" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env_two.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_throttling_policy_v2" "policy" {
  instance_id      = opentelekomcloud_apigw_gateway_v2.gateway.id
  name             = "policy_%[2]s"
  type             = "API-based"
  period           = 15000
  period_unit      = "SECOND"
  max_api_requests = 100
}

resource "opentelekomcloud_apigw_throttling_policy_associate_v2" "associate" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  policy_id  = opentelekomcloud_apigw_throttling_policy_v2.policy.id

  publish_ids = [
    opentelekomcloud_apigw_api_publishment_v2.pub_one.publish_id
  ]
}
`, relatedConfig, name)
}

func testAccThrottlingPolicyAssociateV2_update(name string) string {
	relatedConfig := testAccApigwApi_basic(testAccApigwApi_base(name), name)
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_environment_v2" "env_two" {
  name        = "second_env_%[2]s"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"
}

resource "opentelekomcloud_apigw_api_publishment_v2" "pub_one" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_api_publishment_v2" "pub_two" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env_two.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_throttling_policy_v2" "policy" {
  instance_id      = opentelekomcloud_apigw_gateway_v2.gateway.id
  name             = "policy_%[2]s"
  type             = "API-based"
  period           = 15000
  period_unit      = "SECOND"
  max_api_requests = 100
}

resource "opentelekomcloud_apigw_throttling_policy_associate_v2" "associate" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  policy_id  = opentelekomcloud_apigw_throttling_policy_v2.policy.id

  publish_ids = [
    opentelekomcloud_apigw_api_publishment_v2.pub_two.publish_id
  ]
}
`, relatedConfig, name)
}
