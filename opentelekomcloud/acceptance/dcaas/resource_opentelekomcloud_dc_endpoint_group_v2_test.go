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

const dceg = "opentelekomcloud_dc_endpoint_group_v2.dc_endpoint_group"

func TestDCEndpointGroupV2Resource_basic(t *testing.T) {
	var group dcep.DCEndpointGroup
	DCegName := fmt.Sprintf("dceg-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDCegV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCegV2Resource_basic(DCegName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCegV2Exists(dceg, &group),
					resource.TestCheckResourceAttr(dceg, "description", "first"),
					resource.TestCheckResourceAttrSet(dceg, "id"),
				),
			},
			{
				ResourceName:      dceg,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDCegV2Exists(n string, group *dcep.DCEndpointGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DCaaSV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DCaasV2 client: %s", err)
		}

		found, err := dcep.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("DC endpoint group not found")
		}
		group = found

		return nil
	}
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

func testAccDCegV2Resource_basic(dcegName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dc_endpoint_group_v2" "dc_endpoint_group" {
  name        = "%s"
  type        = "cidr"
  endpoints   = ["10.2.0.0/24", "10.3.0.0/24"]
  description = "first"
  project_id  = "%s"
}
`, dcegName, env.OS_PROJECT_ID)
}
