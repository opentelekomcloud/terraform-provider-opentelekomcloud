package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/ptrrecords"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourcePtrRecordName = "opentelekomcloud_dns_ptrrecord_v2.ptr_1"

func TestAccDNSV2PtrRecord_basic(t *testing.T) {
	var ptr ptrrecords.Ptr
	ptrName := fmt.Sprintf("acc-test-%s.com.", acctest.RandString(3))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2PtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2PtrRecordBasic(ptrName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2PtrRecordExists(resourcePtrRecordName, &ptr),
					resource.TestCheckResourceAttr(resourcePtrRecordName, "description", "a ptr record"),
					resource.TestCheckResourceAttr(resourcePtrRecordName, "tags.muh", "value-create"),
				),
			},
			{
				Config: testAccDNSV2PtrRecordUpdate(ptrName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2PtrRecordExists(resourcePtrRecordName, &ptr),
					resource.TestCheckResourceAttr(resourcePtrRecordName, "description", "ptr record updated"),
					resource.TestCheckResourceAttr(resourcePtrRecordName, "tags.muh", "value-update"),
				),
			},
		},
	})
}

func TestAccDNSV2PtrRecord_unDotted(t *testing.T) {
	zoneName := randomZoneName()
	zoneName = strings.TrimSuffix(zoneName, ".")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2PtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2PtrRecordBasic(zoneName),
			},
		},
	})
}

func testAccCheckDNSV2PtrRecordDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DnsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dns_ptrrecord_v2" {
			continue
		}

		_, err = ptrrecords.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ptr record still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2PtrRecordExists(n string, ptr *ptrrecords.Ptr) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DnsV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
		}

		found, err := ptrrecords.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ptr record not found")
		}

		*ptr = *found

		return nil
	}
}

func testAccDNSV2PtrRecordBasic(ptrName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_dns_ptrrecord_v2" "ptr_1" {
  name          = "%s"
  description   = "a ptr record"
  floatingip_id = opentelekomcloud_networking_floatingip_v2.fip_1.id
  ttl           = 6000
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, ptrName)
}

func testAccDNSV2PtrRecordUpdate(ptrName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_dns_ptrrecord_v2" "ptr_1" {
  name          = "%s"
  description   = "ptr record updated"
  floatingip_id = opentelekomcloud_networking_floatingip_v2.fip_1.id
  ttl           = 6000
  tags = {
    muh = "value-update"
  }
}
`, ptrName)
}
