package dcaas

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	dcep "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcaas/v2/dc-endpoint-group"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestDCEndpointGroupV2Resource_basic(t *testing.T) {
	DCegName := fmt.Sprintf("dceg-%s", acctest.RandString(5))
	DCegNameUpdated := fmt.Sprintf("dceg-updated-%s", acctest.RandString(5))

	const dceg = "opentelekomcloud_dc_endpoint_group_v2.dc_endpoint_group"

	tenantID := env.OS_PROJECT_ID

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDCegV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCegV2Resource_basic(DCegName, tenantID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dceg, "description", "first"),
					resource.TestCheckResourceAttrSet(dceg, "id"),
				),
			},
			{
				Config: testAccDCegV2ResourceUpdate_basic(DCegNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dceg, "description", "second"),
					resource.TestCheckResourceAttr(dceg, "name", DCegNameUpdated),
					resource.TestCheckResourceAttrSet(dceg, "id"),
				),
			},
		},
	})
}

func testAccDCegV2Resource_basic(dcegName string, tenantID string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dc_endpoint_group_v2" "dc_endpoint_group" {
  name        = "%s"
  type        = "cidr"
  endpoints   = ["10.2.0.0/24", "10.3.0.0/24"]
  description = "first"
  tenant_id   = "%s"
}
`, dcegName, tenantID)
}

func testAccDCegV2ResourceUpdate_basic(dcegNameUpdated string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dc_endpoint_group_v2" "dc_endpoint_group" {
  name        = "%s"
  type        = "cidr"
  description = "second"
  tenant_id   = "%s"
  endpoints   = ["10.2.0.0/24", "10.3.0.0/24"]
}
`, dcegNameUpdated, env.OS_PROJECT_ID)
}

func testAccCheckDCegV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DCaaSV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DCaasV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dc_endpoint_group_v2" {
			continue
		}

		dceg, _ := dcep.Get(client, rs.Primary.ID)
		if dceg != nil {
			return fmt.Errorf("DC endpoint group still exists")
		}
		var errDefault404 golangsdk.ErrDefault404
		if !errors.As(err, &errDefault404) {
			return err
		}
	}
	return nil
}
