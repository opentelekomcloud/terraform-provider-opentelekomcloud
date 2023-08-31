package dcaas

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	dcaas "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcaas/v2/direct-connect"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestDirectConnectV2Resource_basic(t *testing.T) {
	directConnectName := fmt.Sprintf("dc-%s", acctest.RandString(5))
	directConnectNameUpdated := fmt.Sprintf("dc-updated-%s", acctest.RandString(5))

	const dc = "opentelekomcloud_direct_connect_v2.direct_connect"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDirectConnectV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDirectConnectV2Resource_basic(directConnectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dc, "bandwidth", "50"),
					resource.TestCheckResourceAttrSet(dc, "id"),
				),
			},
			{
				Config: testAccDirectConnectV2ResourceUpdate_basic(directConnectNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dc, "bandwidth", "100"),
					resource.TestCheckResourceAttrSet(dc, "id"),
				),
			},
		},
	})
}

func testAccDirectConnectV2Resource_basic(directConnectName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_direct_connect_v2" "direct_connect" {
  name          = "%s"
  port_type     = "1G"
  location      = "Biere"
  bandwidth     = 50
  provider_name = "OTC"
}
`, directConnectName)
}

func testAccDirectConnectV2ResourceUpdate_basic(directConnectName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_direct_connect_v2" "direct_connect" {
  name          = "%s"
  port_type     = "1G"
  location      = "Biere"
  bandwidth     = 100
  provider_name = "OTC"
}
`, directConnectName)
}

func testAccCheckDirectConnectV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	dcaasClient, err := config.DCaaSV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating DCaaS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_direct_connect_v2" {
			continue
		}

		_, err := dcaas.Get(dcaasClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("DirectConnect still exists")
		}
		var errDefault404 golangsdk.ErrDefault404
		if !errors.As(err, &errDefault404) {
			return err
		}
	}
	return nil
}
