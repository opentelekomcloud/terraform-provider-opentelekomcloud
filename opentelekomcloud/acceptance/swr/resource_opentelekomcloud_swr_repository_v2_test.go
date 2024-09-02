package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/repositories"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/swr"
)

func TestSwrRepositoryV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testSwrRepositoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testSwrRepositoryV2Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRepoName1, "category", "linux"),
					resource.TestCheckResourceAttr(resourceRepoName2, "name", "one/two"),
					resource.TestCheckResourceAttr(resourceRepoName3, "name", "one/two/three"),
				),
			},
			{
				Config: testSwrRepositoryV2Updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRepoName1, "category", "linux"),
					resource.TestCheckResourceAttr(resourceRepoName1, "description", "Test repository updated"),
				),
			},
		},
	})
}

func testSwrRepositoryV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SwrV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(swr.ClientError, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_swr_repository_v2" {
			continue
		}

		_, err := repositories.Get(client, rs.Primary.Attributes["organization"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("SWR repository still exists")
		}
	}

	return nil
}

const (
	resourceRepoName1 = "opentelekomcloud_swr_repository_v2.repo_1"
	resourceRepoName2 = "opentelekomcloud_swr_repository_v2.repo_2"
	resourceRepoName3 = "opentelekomcloud_swr_repository_v2.repo_3"

	testSwRepositoryV2BasicTemplate = `
resource opentelekomcloud_swr_organization_v2 org_1 {
  name = "%[1]s"
}

resource opentelekomcloud_swr_repository_v2 repo_1 {
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  name         = "%[1]s"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}

resource opentelekomcloud_swr_repository_v2 repo_2 {
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  name         = "one/two"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}

resource opentelekomcloud_swr_repository_v2 repo_3 {
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  name         = "one/two/three"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}
`
)

var (
	testSwrRepositoryV2Basic = fmt.Sprintf(testSwRepositoryV2BasicTemplate, name)
)

const (
	testSwRepositoryV2UpdateTemplate = `
resource opentelekomcloud_swr_organization_v2 org_1 {
  name = "%[1]s"
}

resource opentelekomcloud_swr_repository_v2 repo_1 {
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  name         = "%[1]s"
  description  = "Test repository updated"
  category     = "linux"
  is_public    = false
}
`
)

var (
	testSwrRepositoryV2Updated = fmt.Sprintf(testSwRepositoryV2UpdateTemplate, name)
)
