package acceptance

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/recordsets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/dns"
)

func randomZoneName() string {
	// TODO: why does back-end convert name to lowercase?
	return fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))
}

func TestAccDNSV2RecordSet_basic(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists(
						"opentelekomcloud_dns_recordset_v2.recordset_1", &recordset),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "name", zoneName),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "description", "a record set"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "type", "A"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "ttl", "3000"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "tags.key", "value"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_update(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "ttl", "6000"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "tags.key", "value_updated"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "description", "an updated record set"),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_undotted(t *testing.T) {
	zoneName := randomZoneName()
	zoneName = strings.TrimSuffix(zoneName, ".")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(zoneName),
			},
		},
	})
}

// TestAccDNSV2RecordSet_childFirst covers #847
func TestAccDNSV2RecordSet_childFirst(t *testing.T) {
	zoneName := randomZoneName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_childFirst1(zoneName),
			},
			{
				Config: testAccDNSV2RecordSet_childFirst2(zoneName),
			},
		},
	})
}

func TestAccDNSV2RecordSet_readTTL(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_readTTL(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("opentelekomcloud_dns_recordset_v2.recordset_1", &recordset),
					resource.TestMatchResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_timeout(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_timeout(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("opentelekomcloud_dns_recordset_v2.recordset_1", &recordset),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_shared(t *testing.T) {
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "shared", "false"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_reuse(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "shared", "false"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_2", "shared", "true"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_basic(zoneName),
			},
		},
	})
}

func testAccCheckDNSV2RecordSetDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	dnsClient, err := config.DnsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dns_recordset_v2" {
			continue
		}

		zoneID, recordsetID, err := dns.ParseDNSV2RecordSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err == nil {
			return fmt.Errorf("record set still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2RecordSetExists(n string, recordset *recordsets.RecordSet) resource.TestCheckFunc {
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

		zoneID, recordsetID, err := dns.ParseDNSV2RecordSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err != nil {
			return err
		}

		if found.ID != recordsetID {
			return fmt.Errorf("record set not found")
		}

		*recordset = *found

		return nil
	}
}

func testAccDNSV2RecordSet_basic(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%[1]s"
  email       = "email2@example.com"
  description = "a zone"
  ttl         = 6000
}

resource "opentelekomcloud_dns_recordset_v2" "recordset_1" {
  zone_id     = opentelekomcloud_dns_zone_v2.zone_1.id
  name        = "%[1]s"
  type        = "A"
  description = "a record set"
  ttl         = 3000
  records     = ["10.1.0.0"]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, zoneName)
}

func testAccDNSV2RecordSet_update(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%s"
  email       = "email2@example.com"
  description = "an updated zone"
  ttl         = 6000
}

resource "opentelekomcloud_dns_recordset_v2" "recordset_1" {
  zone_id     = opentelekomcloud_dns_zone_v2.zone_1.id
  name        = "%s"
  type        = "A"
  description = "an updated record set"
  ttl         = 6000
  records     = ["10.1.0.1"]

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, zoneName, zoneName)
}

func testAccDNSV2RecordSet_readTTL(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%[1]s"
  email       = "email2@example.com"
  description = "an updated zone"
  ttl         = 6000
}

resource "opentelekomcloud_dns_recordset_v2" "recordset_1" {
  zone_id = opentelekomcloud_dns_zone_v2.zone_1.id
  name    = "%[1]s"
  type    = "A"
  records = ["10.1.0.2"]
}
`, zoneName)
}

func testAccDNSV2RecordSet_timeout(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%[1]s"
  email       = "email2@example.com"
  description = "an updated zone"
  ttl         = 6000
}

resource "opentelekomcloud_dns_recordset_v2" "recordset_1" {
  zone_id = opentelekomcloud_dns_zone_v2.zone_1.id
  name    = "%[1]s"
  type    = "A"
  ttl     = 3000
  records = ["10.1.0.3", "10.1.0.2"]

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, zoneName)
}

func testAccDNSV2RecordSet_reuse(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%[1]s"
  email       = "email2@example.com"
  description = "a zone"
  ttl         = 6000
}

resource "opentelekomcloud_dns_recordset_v2" "recordset_1" {
  zone_id     = opentelekomcloud_dns_zone_v2.zone_1.id
  name        = "%[1]s"
  type        = "A"
  description = "a record set"
  ttl         = 3000
  records     = ["10.1.0.0"]

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "opentelekomcloud_dns_recordset_v2" "recordset_2" {
  zone_id     = opentelekomcloud_dns_zone_v2.zone_1.id
  name        = "%[1]s"
  type        = "A"
  description = "a record set"
  ttl         = 3000
  records     = ["10.1.0.0"]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, zoneName)
}

func testAccDNSV2RecordSet_childFirst1(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%[1]s"
  email       = "email2@example.com"
  description = "a zone"
  ttl         = 6000
}

resource "opentelekomcloud_dns_recordset_v2" "recordset" {
  zone_id     = opentelekomcloud_dns_zone_v2.zone_1.id
  name        = "test.test.%[1]s"
  type        = "A"
  description = "a record set"
  ttl         = 3000
  records     = ["10.1.0.0"]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, zoneName)
}

func testAccDNSV2RecordSet_childFirst2(zoneName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dns_zone_v2" "zone_1" {
  name        = "%[1]s"
  email       = "email2@example.com"
  description = "a zone"
  ttl         = 6000
}

resource "opentelekomcloud_dns_recordset_v2" "recordset" {
  zone_id     = opentelekomcloud_dns_zone_v2.zone_1.id
  name        = "test.test.%[1]s"
  type        = "A"
  description = "a record set"
  ttl         = 3000
  records     = ["10.1.0.0"]
}
resource "opentelekomcloud_dns_recordset_v2" "recordset_sup" {
  zone_id     = opentelekomcloud_dns_zone_v2.zone_1.id
  name        = "test.%[1]s"
  type        = "A"
  description = "a parent record set"
  ttl         = 3000
  records     = ["10.1.0.0"]
}
`, zoneName)
}
