package acceptance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const resourcePolicyName = "opentelekomcloud_waf_policy_v1.policy_1"

func ruleImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		policyID := s.RootModule().Resources[resourcePolicyName].Primary.ID
		ccRuleID := s.RootModule().Resources[resourceName].Primary.ID
		return fmt.Sprintf("%s/%s", policyID, ccRuleID), nil
	}
}

func stepWAFRuleImport(resourceName string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: ruleImportStateIDFunc(resourceName),
	}
}

func dedicatedRuleImportStateIDFunc(resourceName, resourcePolicyName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		policyID := s.RootModule().Resources[resourcePolicyName].Primary.ID
		ccRuleID := s.RootModule().Resources[resourceName].Primary.ID
		return fmt.Sprintf("%s/%s", policyID, ccRuleID), nil
	}
}
