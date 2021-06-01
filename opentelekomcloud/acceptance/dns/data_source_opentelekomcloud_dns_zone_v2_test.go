package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenStackDNSZoneV2DataSource_basic(t *testing.T) {
	zone := randomZoneName()
	randZoneTag := fmt.Sprintf("value-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheckRequiredEnvVars(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackDNSZoneV2DataSource_zone(zone, randZoneTag),
			},
			{
				Config: testAccOpenStackDNSZoneV2DataSource_basic(zone, randZoneTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSZoneV2DataSourceID("data.opentelekomcloud_dns_zone_v2.z1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dns_zone_v2.z1", "name", zone),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dns_zone_v2.z1", "ttl", "3000"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dns_zone_v2.z1", "zone_type", "public"),
				),
			},
			{
				Config: testAccOpenStackDNSZoneV2DataSource_zone(zone, randZoneTag),
			},
		},
	})
}

func TestAccOpenStackDNSZoneV2DataSource_byTag(t *testing.T) {
	zone := randomZoneName()
	randZoneTag := fmt.Sprintf("value-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheckRequiredEnvVars(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackDNSZoneV2DataSource_zone(zone, randZoneTag),
			},
			{
				Config: testAccOpenStackDNSZoneV2DataSource_byTag(zone, randZoneTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSZoneV2DataSourceID("data.opentelekomcloud_dns_zone_v2.z1"),
				),
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

const zoneTemplate = `
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%s"
  email       = "email1@example.com"
  description = "a public zone"
  ttl         = 3000
  type        = "public"

  tags = {
    key = "%s"
  }
}
`

func testAccOpenStackDNSZoneV2DataSource_zone(zoneName, zoneTagValue string) string {
	return fmt.Sprintf(zoneTemplate, zoneName, zoneTagValue)
}

func testAccOpenStackDNSZoneV2DataSource_basic(zoneName, zoneTagValue string) string {
	base := testAccOpenStackDNSZoneV2DataSource_zone(zoneName, zoneTagValue)
	return fmt.Sprintf(`
%s
data "opentelekomcloud_dns_zone_v2" "z1" {
  name = "%s"
}
`, base, zoneName)
}

func testAccOpenStackDNSZoneV2DataSource_byTag(zoneName, zoneTagValue string) string {
	base := testAccOpenStackDNSZoneV2DataSource_zone(zoneName, zoneTagValue)
	return fmt.Sprintf(`
%s
data "opentelekomcloud_dns_zone_v2" "z1" {
  tags = {
    key = "%s"
  }
}
`, base, zoneName)
}
