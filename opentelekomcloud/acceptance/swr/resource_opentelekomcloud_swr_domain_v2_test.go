package swr

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/domains"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/swr"
)

func TestSwrDomainV2Basic(t *testing.T) {
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
				Config: testSwrDomainV2Basic(name, domainToShare),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceDomainName, "permission", "read"),
					resource.TestCheckResourceAttr(resourceDomainName, "access_domain", domainToShare),
				),
			},
		},
	})
}

func TestSwrDomainV2Slashes(t *testing.T) {
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
				Config: testSwrDomainV2Slashes(name, "grafana/grafana", domainToShare),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceDomainName, "permission", "read"),
					resource.TestCheckResourceAttr(resourceDomainName, "access_domain", domainToShare),
					resource.TestCheckResourceAttr(resourceDomainName, "repository", "grafana/grafana"),
				),
			},
		},
	})
}

func testSwrDomainV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SwrV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(swr.ClientError, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_swr_domain_v2" {
			continue
		}

		_, err := domains.Get(client, domains.GetOpts{
			Namespace:    rs.Primary.Attributes["organization"],
			Repository:   rs.Primary.Attributes["repository"],
			AccessDomain: rs.Primary.ID,
		})
		if err == nil {
			return fmt.Errorf("SWR domain still exists")
		}
	}

	return nil
}

const (
	resourceDomainName = "opentelekomcloud_swr_domain_v2.domain_1"
)

func testSwrDomainV2Basic(name, domainToShare string) string {
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

resource opentelekomcloud_swr_domain_v2 domain_1 {
  organization  = opentelekomcloud_swr_organization_v2.org_1.name
  repository    = opentelekomcloud_swr_repository_v2.repo_1.name
  access_domain = "%[2]s"
  permission    = "read"
  deadline      = "forever"
}
`, name, domainToShare)
}

func testSwrDomainV2Slashes(orgName string, repoName string, domainToShare string) string {
	return fmt.Sprintf(`
resource opentelekomcloud_swr_organization_v2 org_1 {
  name = "%[1]s"
}

resource opentelekomcloud_swr_repository_v2 repo_1 {
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  name         = "%[2]s"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}

resource opentelekomcloud_swr_domain_v2 domain_1 {
  organization  = opentelekomcloud_swr_organization_v2.org_1.name
  repository    = opentelekomcloud_swr_repository_v2.repo_1.name
  access_domain = "%[3]s"
  permission    = "read"
  deadline      = "forever"
}
`, orgName, repoName, domainToShare)
}
