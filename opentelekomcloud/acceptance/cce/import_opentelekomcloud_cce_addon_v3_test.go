package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCEAddonV3ImportBasic(t *testing.T) {
	resourceName := "opentelekomcloud_cce_addon_v3.autoscaler"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3Basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCEAddonV3ImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"values",
				},
			},
		},
	})
}

func testAccCCEAddonV3ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterID string
		var addonID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_cce_cluster_v3" {
				clusterID = rs.Primary.ID
			} else if rs.Type == "opentelekomcloud_cce_addon_v3" {
				addonID = rs.Primary.ID
			}
		}
		if clusterID == "" || addonID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", clusterID, addonID)
		}
		return fmt.Sprintf("%s/%s", clusterID, addonID), nil
	}
}
