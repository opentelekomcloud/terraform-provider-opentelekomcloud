package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/policy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const policyResourceName = "opentelekomcloud_swr_policy_v2.policy_1"

func getPolicyFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.SwrV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud SWR client: %s", err)
	}
	return policy.Get(client, state.Primary.Attributes["organization"], state.Primary.Attributes["repository"], state.Primary.ID)
}

func TestSwrPolicyV2_basic(t *testing.T) {
	var retentionPolicy policy.ImageRetentionPolicy
	rc := common.InitResourceCheck(
		policyResourceName,
		&retentionPolicy,
		getPolicyFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(policyResourceName, "rules.0.params.days", "30"),
				),
			},
			{
				Config: testAccSwrPolicyV2Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(policyResourceName, "rules.0.params.days", "45"),
				),
			},
			{
				ResourceName:      policyResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSwrPolicyV2ImportStateIdFunc(),
			},
		},
	})
}

func testAccSwrPolicyV2ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var org, repo, id string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_swr_policy_v2" {
				org = rs.Primary.Attributes["organization"]
				repo = rs.Primary.Attributes["repository"]
				id = rs.Primary.ID
			}
		}
		if org == "" || repo == "" || id == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", org, repo, id)
		}
		return fmt.Sprintf("%s/%s/%s", org, repo, id), nil
	}
}

var testAccSwrPolicyV2Basic = `
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

resource "opentelekomcloud_swr_policy_v2" "policy_1" {
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
`

var testAccSwrPolicyV2Update = `
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

resource "opentelekomcloud_swr_policy_v2" "policy_1" {
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  repository   = opentelekomcloud_swr_repository_v2.repo_1.name
  algorithm    = "or"
  rules {
    template = "date_rule"
    params = {
      days = "45"
    }
    tag_selector {
      kind    = "label"
      pattern = "v1"
    }
  }
}
`
