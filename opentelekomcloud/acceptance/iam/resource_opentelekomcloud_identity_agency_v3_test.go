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

const resourceAgencyName = "opentelekomcloud_identity_agency_v3.agency_1"

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
				Config: testAccIdentityV3AgencyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3AgencyExists(resourceAgencyName, &a),
					resource.TestCheckResourceAttr(resourceAgencyName, "name", "test"),
					resource.TestCheckResourceAttr(resourceAgencyName, "delegated_domain_name", "op_svc_evs"),
				),
			},
			{
				Config: testAccIdentityV3AgencyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3AgencyExists(resourceAgencyName, &a),
					resource.TestCheckResourceAttr(resourceAgencyName, "name", "test"),
					resource.TestCheckResourceAttr(resourceAgencyName, "delegated_domain_name", "op_svc_evs"),
				),
			},
		},
	})
}

func TestAccIdentityV3Agency_importBasic(t *testing.T) {
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
				Config: testAccIdentityV3AgencyBasic,
			},

			{
				ResourceName:      resourceAgencyName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIdentityV3AgencyDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomcloud identity client: %s", err)
	}
	client.Endpoint = strings.Replace(client.Endpoint, "v3", "v3.0", 1)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_agency_v3" {
			continue
		}

		_, err := agency.Get(client, rs.Primary.ID).Extract()
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
		client, err := config.IdentityV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomcloud identity client: %s", err)
		}
		client.Endpoint = strings.Replace(client.Endpoint, "v3", "v3.0", 1)
		found, err := agency.Get(client, rs.Primary.ID).Extract()
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

var testAccIdentityV3AgencyBasic = fmt.Sprintf(`
resource "opentelekomcloud_identity_agency_v3" "agency_1" {
  name                  = "test"
  delegated_domain_name = "op_svc_evs"
  project_role {
    project = "%[1]s"
    roles = [
      "KMS Administrator",
      "CCE ReadOnlyAccess",
    ]
  }
  project_role {
    all_projects = true
    roles = [
      "CES Administrator",
      "ER ReadOnlyAccess",
    ]
  }
}`, env.OS_TENANT_NAME)

var testAccIdentityV3AgencyUpdate = fmt.Sprintf(`
resource "opentelekomcloud_identity_agency_v3" "agency_1" {
  name                  = "test"
  delegated_domain_name = "op_svc_evs"
  project_role {
    project = "%[1]s"
    roles = [
      "CCE ReadOnlyAccess",
      "WAF FullAccess",
    ]
  }
  project_role {
    all_projects = true
    roles = [
      "CES Administrator",
      "DRS FullWithOutDeletePermission",
    ]
  }

}`, env.OS_TENANT_NAME)
