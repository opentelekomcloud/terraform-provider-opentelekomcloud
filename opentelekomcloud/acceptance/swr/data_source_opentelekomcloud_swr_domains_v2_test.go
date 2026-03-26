package swr

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestSwrDomainsV2DataSourceBasic(t *testing.T) {
	domainToShare := os.Getenv("OS_DOMAIN_NAME_2")
	if domainToShare == "" {
		t.Skip("OS_DOMAIN_NAME_2 is empty")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testSwrDomainV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testSwrDomainsV2DSBasic(name, domainToShare),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceDomainsName, "organization", name),
					resource.TestCheckResourceAttrSet(dataSourceDomainsName, "access_domains.0.status"),
				),
			},
		},
	})
}

const (
	dataSourceDomainsName = "data.opentelekomcloud_swr_domains_v2.domains"
)

func testSwrDomainsV2DSBasic(name, domainToShare string) string {
	return fmt.Sprintf(`
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

resource opentelekomcloud_swr_domain_v2 domain {
  organization  = opentelekomcloud_swr_organization_v2.org_1.name
  repository    = opentelekomcloud_swr_repository_v2.repo_1.name
  access_domain = "%[2]s"
  permission    = "read"
  deadline      = "forever"
}

data opentelekomcloud_swr_domains_v2 domains {
  depends_on = [opentelekomcloud_swr_domain_v2.domain]
  organization  = opentelekomcloud_swr_organization_v2.org_1.name
  repository    = opentelekomcloud_swr_repository_v2.repo_1.name
}
`, name, domainToShare)
}
