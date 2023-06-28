package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dis/v2/apps"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceAppName = "opentelekomcloud_dis_app_v2.app_1"

func TestAccDisAppV2_basic(t *testing.T) {
	var cls apps.GetAppResponse
	var appName = fmt.Sprintf("dis_app_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDisV2AppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDisV2AppBasic(appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDisV2AppExists(resourceAppName, &cls),
					resource.TestCheckResourceAttr(resourceAppName, "app_name", appName),
				),
			},
			{
				ResourceName:      resourceAppName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDisV2AppDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DisV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating DISv2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dis_app_v2" {
			continue
		}

		_, err := apps.GetApp(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("DIS app still exists")
		}
	}
	return nil
}

func testAccCheckDisV2AppExists(n string, cls *apps.GetAppResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DisV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating DISv2 client: %w", err)
		}

		v, err := apps.GetApp(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting app (%s): %w", rs.Primary.ID, err)
		}

		if v.AppName != rs.Primary.ID {
			return fmt.Errorf("DIS stream not found")
		}
		*cls = *v
		return nil
	}
}

func testAccDisV2AppBasic(appName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dis_app_v2" "app_1" {
  app_name = "%s"
}
`, appName)
}
