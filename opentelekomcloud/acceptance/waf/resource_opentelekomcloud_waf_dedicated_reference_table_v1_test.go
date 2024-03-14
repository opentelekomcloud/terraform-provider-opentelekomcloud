package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const wafdRefTableResourceName = "opentelekomcloud_waf_dedicated_reference_table_v1.table"

func TestAccWafDedicatedReferenceTableV1_basic(t *testing.T) {
	var refTable rules.ReferenceTable
	var name = fmt.Sprintf("wafd_rt.v1-%s", acctest.RandString(5))
	updateName := name + "_update"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedRefTableV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedRefTablesV1_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedRefTableV1Exists(wafdRefTableResourceName, &refTable),
					resource.TestCheckResourceAttr(wafdRefTableResourceName, "name", name),
					resource.TestCheckResourceAttr(wafdRefTableResourceName, "description", ""),
					resource.TestCheckResourceAttr(wafdRefTableResourceName, "type", "url"),
					resource.TestCheckResourceAttr(wafdRefTableResourceName, "conditions.#", "2"),
				),
			},
			{
				Config: testAccWafReferenceTableV1_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedRefTableV1Exists(wafdRefTableResourceName, &refTable),
					resource.TestCheckResourceAttr(wafdRefTableResourceName, "name", updateName),
					resource.TestCheckResourceAttr(wafdRefTableResourceName, "type", "url"),
					resource.TestCheckResourceAttr(wafdRefTableResourceName, "conditions.#", "2"),
					resource.TestCheckResourceAttr(wafdRefTableResourceName, "description", "new description"),
					resource.TestCheckResourceAttrSet(wafdRefTableResourceName, "created_at"),
				),
			},
			{
				ResourceName:      wafdRefTableResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckWafDedicatedRefTableV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_reference_table_v1" {
			continue
		}
		_, err = rules.GetReferenceTable(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated reference table (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckWafDedicatedRefTableV1Exists(n string, instance *rules.ReferenceTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
		if err != nil {
			return err
		}

		var found *rules.ReferenceTable
		found, err = rules.GetReferenceTable(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		*instance = *found

		return nil
	}
}

func testAccWafDedicatedRefTablesV1_basic(name string) string {
	return fmt.Sprintf(`


resource "opentelekomcloud_waf_dedicated_reference_table_v1" "table" {
  name = "%s"
  type = "url"

  conditions = [
    "/admin",
    "/manage"
  ]
}
`, name)
}

func testAccWafReferenceTableV1_update(name string) string {
	return fmt.Sprintf(`

resource "opentelekomcloud_waf_dedicated_reference_table_v1" "table" {
  name        = "%s"
  type        = "url"
  description = "new description"

  conditions = [
    "/bill",
    "/sql"
  ]
}
`, name)
}
