package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// PASS, but normally skip
func TestAccDNSV2RecordSet_importBasic(t *testing.T) {
	zoneName := randomZoneName()
	resourceName := "opentelekomcloud_dns_recordset_v2.recordset_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDNSV2RecordSet_basic(zoneName),
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
