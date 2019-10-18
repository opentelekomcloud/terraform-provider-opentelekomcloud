package opentelekomcloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/identity/v3/agency"
)

func TestAccIdentityV3Agency_basic(t *testing.T) {
	var a agency.Agency

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIdentityV3AgencyPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3AgencyDestroy,
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
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.identityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Opentelekomcloud identity client: %s", err)
	}
	identityClient.Endpoint = strings.Replace(identityClient.Endpoint, "v3", "v3.0", 1)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_agency_v3" {
			continue
		}

		_, err := agency.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Agency still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3AgencyExists(n string, a *agency.Agency) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.identityV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Opentelekomcloud identity client: %s", err)
		}
		identityClient.Endpoint = strings.Replace(identityClient.Endpoint, "v3", "v3.0", 1)
		found, err := agency.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Agency not found")
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
    }`, OS_TENANT_NAME)
