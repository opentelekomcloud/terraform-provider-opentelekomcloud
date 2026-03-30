package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestSwrRepositoryDSV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testSwrRepositoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testSwrRepositoryDSV2Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceRepoName1, "repositories.0.category", "linux"),
				),
			},
		},
	})
}

const dataSourceRepoName1 = "data.opentelekomcloud_swr_repository_v2.repo_1"

var testSwrRepositoryDSV2Basic = fmt.Sprintf(`
resource opentelekomcloud_swr_organization_v2 org_1 {
  name = "%[1]s"
}

resource opentelekomcloud_swr_repository_v2 repo {
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  name         = "%[1]s"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}

data opentelekomcloud_swr_repository_v2 repo_1 {
  depends_on = [opentelekomcloud_swr_repository_v2.repo]
  name         = "%[1]s"
}
`, name)
