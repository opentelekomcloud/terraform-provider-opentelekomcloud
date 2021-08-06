package acceptance

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/zones"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccDNSV2Zone_basic(t *testing.T) {
	var zone zones.Zone
	// TODO: Why does it lowercase names in back-end?
	var zoneName = fmt.Sprintf("accepttest%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_basic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("opentelekomcloud_dns_zone_v2.zone_1", &zone),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_zone_v2.zone_1", "description", "a public zone"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_zone_v2.zone_1", "tags.foo", "bar"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_zone_v2.zone_1", "tags.key", "value"),
				),
			},
			{
				Config: testAccDNSV2Zone_update(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_dns_zone_v2.zone_1", "name", zoneName),
					resource.TestCheckResourceAttr("opentelekomcloud_dns_zone_v2.zone_1", "email", "email2@example.com"),
					resource.TestCheckResourceAttr("opentelekomcloud_dns_zone_v2.zone_1", "ttl", "6000"),
					// TODO: research why this is blank...
					// resource.TestCheckResourceAttr("opentelekomcloud_dns_zone_v2.zone_1", "type", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_zone_v2.zone_1", "description", "an updated zone"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_zone_v2.zone_1", "tags.key", "value_updated"),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_undotted(t *testing.T) {
	zoneName := randomZoneName()
	zoneName = strings.TrimSuffix(zoneName, ".")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_basic(zoneName),
			},
		},
	})
}

func TestAccDNSV2Zone_private(t *testing.T) {
	var zone zones.Zone
	// TODO: Why does it lowercase names in back-end?
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_private(zoneName),
				// ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("opentelekomcloud_dns_zone_v2.zone_1", &zone),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_zone_v2.zone_1", "description", "a private zone"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_zone_v2.zone_1", "type", "private"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_zone_v2.zone_1", "tags.foo", "bar"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_zone_v2.zone_1", "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_readTTL(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccDNSV2Zone_readTTL(zoneName),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("opentelekomcloud_dns_zone_v2.zone_1", &zone),
					// resource.TestCheckResourceAttr("opentelekomcloud_dns_zone_v2.zone_1", "type", "PRIMARY"),
					resource.TestMatchResourceAttr(
						"opentelekomcloud_dns_zone_v2.zone_1", "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_timeout(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccDNSV2Zone_timeout(zoneName),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("opentelekomcloud_dns_zone_v2.zone_1", &zone),
				),
			},
		},
	})
}

func testAccCheckDNSV2ZoneDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	dnsClient, err := config.DnsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dns_zone_v2" {
			continue
		}

		_, err := zones.Get(dnsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("zone still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2ZoneExists(n string, zone *zones.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		dnsClient, err := config.DnsV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
		}

		found, err := zones.Get(dnsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("zone not found")
		}

		*zone = *found

		return nil
	}
}

func testAccDNSV2Zone_basic(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%s"
  email       = "email1@example.com"
  description = "a public zone"
  ttl         = 3000
  type        = "public"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, zoneName)
}

func testAccDNSV2Zone_private(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%s"
  email       = "email1@example.com"
  description = "a private zone"
  ttl         = 3000
  type        = "private"

  router {
    router_id     = "%s"
    router_region = "%s"
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, zoneName, env.OS_VPC_ID, env.OS_REGION_NAME)
}

func testAccDNSV2Zone_update(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%s"
  email       = "email2@example.com"
  description = "an updated zone"
  ttl         = 6000
  type        = "public"

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, zoneName)
}

func testAccDNSV2Zone_readTTL(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name  = "%s"
  email = "email1@example.com"
}
`, zoneName)
}

func testAccDNSV2Zone_timeout(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name  = "%s"
  email = "email@example.com"
  ttl   = 3000

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, zoneName)
}
