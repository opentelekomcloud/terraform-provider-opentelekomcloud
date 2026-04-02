package swr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const policyDataSourceName = "data.opentelekomcloud_swr_policy_v2.policy_1"

func TestSwrPolicyV2DS_basic(t *testing.T) {
	dc := common.InitDataSourceCheck(policyDataSourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSwrPolicyV2DSBasic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(policyDataSourceName, "rules.0.params.days", "30"),
				),
			},
		},
	})
}

var testAccSwrPolicyV2DSBasic = `
resource "opentelekomcloud_swr_organization_v2" "org_1" {
  name = "test-acc-org-swr-pol"
}

resource "opentelekomcloud_swr_repository_v2" "repo_1" {
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  name         = "test-acc-repo-swr-pol"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}

resource "opentelekomcloud_swr_policy_v2" "policy" {
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  repository   = opentelekomcloud_swr_repository_v2.repo_1.name
  algorithm    = "or"
  rules {
    template = "date_rule"
    params = {
      days = "30"
    }
    tag_selector {
      kind    = "label"
      pattern = "v1"
    }
  }
}

data "opentelekomcloud_swr_policy_v2" "policy_1" {
  depends_on   = [opentelekomcloud_swr_policy_v2.policy]
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  repository   = opentelekomcloud_swr_repository_v2.repo_1.name
  policy_id    = opentelekomcloud_swr_policy_v2.policy.id
}
`
