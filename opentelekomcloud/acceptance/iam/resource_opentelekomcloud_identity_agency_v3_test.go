package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/agency"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccIdentityV3Agency_basic(t *testing.T) {
	var a agency.Agency

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			TestAccIdentityV3AgencyPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3AgencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Agency_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3AgencyExists("opentelekomcloud_identity_agency_v3.agency_1", &a),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_identity_agency_v3.agency_1", "name", "test"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_identity_agency_v3.agency_1", "delegated_domain_name", "op_svc_evs"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3AgencyDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomcloud identity client: %s", err)
	}
	identityClient.Endpoint = strings.Replace(identityClient.Endpoint, "v3", "v3.0", 1)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_agency_v3" {
			continue
		}

		_, err := agency.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("agency still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3AgencyExists(n string, a *agency.Agency) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomcloud identity client: %s", err)
		}
		identityClient.Endpoint = strings.Replace(identityClient.Endpoint, "v3", "v3.0", 1)
		found, err := agency.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("agency not found")
		}

		*a = *found

		return nil
	}
}

var testAccIdentityV3Agency_basic = fmt.Sprintf(`
    resource "opentelekomcloud_identity_agency_v3" "agency_1" {
      name = "test"
      delegated_domain_name = "op_svc_evs"
      project_role = [
      {
      		project = "%s"
      		roles = [
        		"KMS Administrator",
      		]
    	}
		]
    }`, env.OS_TENANT_NAME)
