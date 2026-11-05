package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/aom/v1/monitoring"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	accenv "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/aom"
)

const resourceServiceDiscoveryRuleName = "opentelekomcloud_aom_service_discovery_rule_v1.rule"

func getServiceDiscoveryRuleFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.AomV1Client(accenv.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM v1 client: %s", err)
	}
	return aom.GetServiceDiscoveryRule(client, state.Primary.ID)
}

func TestAccAomServiceDiscoveryRuleV1_basic(t *testing.T) {
	ruleName := fmt.Sprintf("tf-svc-disc-%s", strings.ToLower(acctest.RandString(5)))

	rc := common.InitResourceCheck(
		resourceServiceDiscoveryRuleName,
		&monitoring.AppRule{},
		getServiceDiscoveryRuleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAomServiceDiscoveryRuleV1Basic(ruleName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceServiceDiscoveryRuleName, "name", ruleName),
					resource.TestCheckResourceAttr(resourceServiceDiscoveryRuleName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceServiceDiscoveryRuleName, "spec.0.priority", "9999"),
					resource.TestCheckResourceAttrSet(resourceServiceDiscoveryRuleName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceServiceDiscoveryRuleName, "create_time"),
				),
			},
			{
				Config: testAccAomServiceDiscoveryRuleV1Updated(ruleName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceServiceDiscoveryRuleName, "name", ruleName),
					resource.TestCheckResourceAttr(resourceServiceDiscoveryRuleName, "description", "updated by acceptance test"),
					resource.TestCheckResourceAttr(resourceServiceDiscoveryRuleName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceServiceDiscoveryRuleName, "spec.0.priority", "100"),
				),
			},
			{
				ResourceName:      resourceServiceDiscoveryRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAomServiceDiscoveryRuleV1Basic(ruleName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_aom_service_discovery_rule_v1" "rule" {
  name        = "%[1]s"
  enabled     = true
  description = "created by acceptance test"

  spec {
    detect_log      = false
    is_default_rule = false
    is_detect       = false
    priority        = 9999
    app_type        = "Java"
    log_file_fix    = ["log", "trace", "out"]

    discovery_rule {
      check_type    = "cmdLine"
      check_mode    = "contain"
      check_content = ["acceptance-test-process"]
    }

    name_rule {
      app_name_rule {
        name_type = "str"
        args      = ["-svc"]
      }

      application_name_rule {
        name_type = "str"
        args      = ["-app"]
      }
    }
  }
}
`, ruleName)
}

func testAccAomServiceDiscoveryRuleV1Updated(ruleName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_aom_service_discovery_rule_v1" "rule" {
  name        = "%[1]s"
  enabled     = false
  description = "updated by acceptance test"

  spec {
    detect_log      = false
    is_default_rule = false
    is_detect       = false
    priority        = 100
    app_type        = "Java"
    log_file_fix    = ["log", "trace", "out"]

    discovery_rule {
      check_type    = "cmdLine"
      check_mode    = "contain"
      check_content = ["acceptance-test-process"]
    }

    name_rule {
      app_name_rule {
        name_type = "str"
        args      = ["-svc"]
      }

      application_name_rule {
        name_type = "str"
        args      = ["-app"]
      }
    }
  }
}
`, ruleName)
}
