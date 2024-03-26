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

const resourcePolicyName = "opentelekomcloud_apigw_throttling_policy_v2.policy"

func TestAccAPIGWv2Policy_basic(t *testing.T) {
	var policyConfig throttlingpolicy.ThrottlingResp
	name := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2PolicyBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2PolicyExists(resourcePolicyName, &policyConfig),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "policy_"+name),
					resource.TestCheckResourceAttr(resourcePolicyName, "description", "created by tf"),
					resource.TestCheckResourceAttr(resourcePolicyName, "type", "API-based"),
					resource.TestCheckResourceAttr(resourcePolicyName, "period", "15000"),
					resource.TestCheckResourceAttr(resourcePolicyName, "period_unit", "SECOND"),
					resource.TestCheckResourceAttr(resourcePolicyName, "max_api_requests", "100"),
					resource.TestCheckResourceAttr(resourcePolicyName, "max_user_requests", "60"),
					resource.TestCheckResourceAttr(resourcePolicyName, "max_app_requests", "60"),
					resource.TestCheckResourceAttr(resourcePolicyName, "max_ip_requests", "60"),
				),
			},
			{
				Config: testAccAPIGWv2PolicyUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2PolicyExists(resourcePolicyName, &policyConfig),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "policy_"+name+"_updated"),
					resource.TestCheckResourceAttr(resourcePolicyName, "description", "Updated by tf"),
					resource.TestCheckResourceAttr(resourcePolicyName, "type", "API-shared"),
					resource.TestCheckResourceAttr(resourcePolicyName, "period", "10"),
					resource.TestCheckResourceAttr(resourcePolicyName, "period_unit", "MINUTE"),
					resource.TestCheckResourceAttr(resourcePolicyName, "max_api_requests", "70"),
					resource.TestCheckResourceAttr(resourcePolicyName, "max_user_requests", "45"),
					resource.TestCheckResourceAttr(resourcePolicyName, "max_app_requests", "45"),
					resource.TestCheckResourceAttr(resourcePolicyName, "max_ip_requests", "45"),
				),
			},
		},
	})
}

// work in progress, needs APP resource to work
func TestAccAPIGWv2Policy_special(t *testing.T) {
	var policyConfig throttlingpolicy.ThrottlingResp
	name := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2PolicyBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2PolicyExists(resourcePolicyName, &policyConfig),
				),
			},
			{
				Config: testAccAPIGWv2PolicySpecial(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2PolicyExists(resourcePolicyName, &policyConfig),
					resource.TestCheckResourceAttr(resourcePolicyName, "app_throttles.#", "1"),
				),
			},
			{
				Config: testAccAPIGWv2PolicySpecialUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2PolicyExists(resourcePolicyName, &policyConfig),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "policy_"+name+"_updated"),
					resource.TestCheckResourceAttr(resourcePolicyName, "description", "Updated by script"),
					resource.TestCheckResourceAttr(resourcePolicyName, "type", "API-shared"),
				),
			},
		},
	})
}

func testAccCheckAPIGWv2PolicyExists(n string, configuration *throttlingpolicy.ThrottlingResp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		rsgw, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rsgw.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.APIGWV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud APIGWv2 client: %w", err)
		}

		found, err := throttlingpolicy.Get(client, rsgw.Primary.ID, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("APIGW group not found")
		}
		configuration = found

		return nil
	}
}

func TestAccAPIGWPolicyV2ImportBasic(t *testing.T) {
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAPIGWv2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2GroupBasic(name),
			},
			{
				ResourceName:      resourceGroupName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAPIGWv2PolicyImportStateIdFunc(),
			},
		},
	})
}

func testAccAPIGWv2PolicyImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var gatewayID string
		var groupID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_apigw_gateway_v2" {
				gatewayID = rs.Primary.ID
			} else if rs.Type == "opentelekomcloud_apigw_group_v2" && rs.Primary.ID != "" {
				groupID = rs.Primary.ID
			}
		}
		if gatewayID == "" || groupID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", gatewayID, groupID)
		}
		return fmt.Sprintf("%s/%s", gatewayID, groupID), nil
	}
}

func testAccAPIGWv2PolicyBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_throttling_policy_v2" "policy" {
  instance_id       = opentelekomcloud_apigw_gateway_v2.gateway.id
  name              = "%s"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "created by tf"
}
`, testAccAPIGWv2GatewayBasic("gateway-"+name), "policy_"+name)
}

func testAccAPIGWv2PolicyUpdated(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_throttling_policy_v2" "policy" {
  instance_id       = opentelekomcloud_apigw_gateway_v2.gateway.id
  name              = "%s"
  type              = "API-shared"
  period            = 10
  period_unit       = "MINUTE"
  max_api_requests  = 70
  max_user_requests = 45
  max_app_requests  = 45
  max_ip_requests   = 45
  description       = "Updated by tf"
}
`, testAccAPIGWv2GatewayBasic("gateway-"+name), "policy_"+name+"_updated")
}

func testAccAPIGWv2PolicySpecial(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_throttling_policy_v2" "policy" {
  instance_id      = opentelekomcloud_apigw_gateway_v2.gateway.id
  name             = "%s"
  type             = "API-based"
  period           = 15000
  period_unit      = "SECOND"
  max_api_requests = 100
}
`, testAccAPIGWv2GatewayBasic("gateway-"+name), "policy_"+name)
}

func testAccAPIGWv2PolicySpecialUpdated(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_throttling_policy_v2" "policy" {
  instance_id      = opentelekomcloud_apigw_gateway_v2.gateway.id
  name             = "%s"
  type             = "API-based"
  period           = 15000
  period_unit      = "SECOND"
  max_api_requests = 100
}
`, testAccAPIGWv2GatewayBasic("gateway-"+name), "policy_"+name+"_updated")
}
