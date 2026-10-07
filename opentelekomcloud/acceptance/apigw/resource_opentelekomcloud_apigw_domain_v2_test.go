package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/group"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	accenv "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/apigw"
)

const resourceNameDomain = "opentelekomcloud_apigw_domain_v2.domain"

func getDomainFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.APIGWV2Client(accenv.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return apigw.GetDomain(client, state.Primary.Attributes["gateway_id"], state.Primary.Attributes["group_id"], state.Primary.ID)
}

func TestAccAPIGWv2Domain_basic(t *testing.T) {
	gatewayID := os.Getenv("OS_GATEWAY_ID")
	if gatewayID == "" {
		t.Skip("`OS_GATEWAY_ID` needs to be defined")
	}
	groupName := fmt.Sprintf("group-%s", acctest.RandString(5))
	domainName := fmt.Sprintf("terraform-acc-%s.example.com", acctest.RandString(3))

	rc := common.InitResourceCheck(
		resourceNameDomain,
		&group.UrlDomains{},
		getDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2DomainBasic(gatewayID, groupName, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNameDomain, "name", domainName),
					resource.TestCheckResourceAttr(resourceNameDomain, "min_ssl_version", "TLSv1.2"),
					resource.TestCheckResourceAttr(resourceNameDomain, "http_redirect_to_https", "false"),
				),
			},
			{
				Config: testAccAPIGWv2DomainUpdated(gatewayID, groupName, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNameDomain, "name", domainName),
					resource.TestCheckResourceAttr(resourceNameDomain, "min_ssl_version", "TLSv1.1"),
					resource.TestCheckResourceAttr(resourceNameDomain, "http_redirect_to_https", "true"),
				),
			},
			{
				ResourceName:      resourceNameDomain,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name",
					"http_redirect_to_https",
				},
				ImportStateIdFunc: testAccAPIGWv2DomainImportStateIdFunc(),
			},
		},
	})
}

func testAccAPIGWv2DomainImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		gatewayID := os.Getenv("OS_GATEWAY_ID")
		var groupID, domainID string
		for _, rs := range s.RootModule().Resources {
			switch rs.Type {
			case "opentelekomcloud_apigw_group_v2":
				groupID = rs.Primary.ID
			case "opentelekomcloud_apigw_domain_v2":
				domainID = rs.Primary.ID
			}
		}
		if gatewayID == "" || groupID == "" || domainID == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", gatewayID, groupID, domainID)
		}
		return fmt.Sprintf("%s/%s/%s", gatewayID, groupID, domainID), nil
	}
}

func testAccAPIGWv2DomainBasic(gatewayId, groupName, domainName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_apigw_group_v2" "group" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "opentelekomcloud_apigw_domain_v2" "domain" {
  gateway_id             = "%[1]s"
  group_id               = opentelekomcloud_apigw_group_v2.group.id
  name                   = "%[3]s"
  min_ssl_version        = "TLSv1.2"
  http_redirect_to_https = false
}
`, gatewayId, groupName, domainName)
}

func testAccAPIGWv2DomainUpdated(gatewayId, groupName, domainName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_apigw_group_v2" "group" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "opentelekomcloud_apigw_domain_v2" "domain" {
  gateway_id             = "%[1]s"
  group_id               = opentelekomcloud_apigw_group_v2.group.id
  name                   = "%[3]s"
  min_ssl_version        = "TLSv1.1"
  http_redirect_to_https = true
}
`, gatewayId, groupName, domainName)
}
