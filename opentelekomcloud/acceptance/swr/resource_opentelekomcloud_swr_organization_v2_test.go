package swr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/organizations"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/swr"
)

func TestSwrOrganizationV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testSwrOrganizationV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testSwrOrganizationV2Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "auth", "7"),
				),
			},
		},
	})
}

func TestSwrOrganizationV2_validateName(t *testing.T) {
	names := []string{"1-start-with-number", "end-with-2-dot.", "biGger-3-one"}
	steps := make([]resource.TestStep, len(names))
	for i, name := range names {
		steps[i] = resource.TestStep{
			Config:      fmt.Sprintf(testSwrOrganizationV2BasicTemplate, name),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`invalid value for name.+`),
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testSwrOrganizationV2Destroy,
		Steps:             steps,
	})
}

func testSwrOrganizationV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SwrV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(swr.ClientError, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_swr_organization_v2" {
			continue
		}

		_, err := organizations.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("SWR organization still exists")
		}
	}

	return nil
}

const (
	resourceName = "opentelekomcloud_swr_organization_v2.org_1"

	testSwrOrganizationV2BasicTemplate = `
resource opentelekomcloud_swr_organization_v2 org_1 {
  name = "%s"
}
`
)

var (
	name = fmt.Sprintf("test-organization-%d", tools.RandomInt(0, 99))

	testSwrOrganizationV2Basic = fmt.Sprintf(testSwrOrganizationV2BasicTemplate, name)
)
