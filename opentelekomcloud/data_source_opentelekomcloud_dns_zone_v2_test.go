package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))

func TestAccOpenStackDNSZoneV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckRequiredEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackDNSZoneV2DataSource_zone,
			},
			{
				Config: testAccOpenStackDNSZoneV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSZoneV2DataSourceID("data.opentelekomcloud_dns_zone_v2.z1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dns_zone_v2.z1", "name", zoneName),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dns_zone_v2.z1", "ttl", "3000"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dns_zone_v2.z1", "zone_type", "public"),
				),
			},
			{
				Config: testAccOpenStackDNSZoneV2DataSource_byTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSZoneV2DataSourceID("data.opentelekomcloud_dns_zone_v2.z1"),
				),
			},
			{
				Config: testAccOpenStackDNSZoneV2DataSource_zone,
			},
		},
	})
}

func testAccCheckDNSZoneV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find DNS Zone data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DNS Zone data source ID not set")
		}

		return nil
	}
}

var testAccOpenStackDNSZoneV2DataSource_zone = testAccDNSV2Zone_basic(zoneName)

var testAccOpenStackDNSZoneV2DataSource_basic = fmt.Sprintf(`
%s
data "opentelekomcloud_dns_zone_v2" "z1" {
	name = "%s"
}
`, testAccDNSV2Zone_basic(zoneName), zoneName)

var testAccOpenStackDNSZoneV2DataSource_byTag = fmt.Sprintf(`
%s
data "opentelekomcloud_dns_zone_v2" "z1" {
	tags = {
		key = "value"
	}
}
`, testAccDNSV2Zone_basic(zoneName))
