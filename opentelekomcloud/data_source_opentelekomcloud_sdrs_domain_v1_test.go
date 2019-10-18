package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSdrsDomainV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsDomainV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsDomainV1DataSourceID("data.opentelekomcloud_sdrs_domain_v1.domain_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_sdrs_domain_v1.domain_1", "name", "domain_001"),
				),
			},
		},
	})
}

func testAccCheckSdrsDomainV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find SDRS domain data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("SDRS domain data source ID not set")
		}

		return nil
	}
}

const testAccSdrsDomainV1DataSource_basic = `
data "opentelekomcloud_sdrs_domain_v1" "domain_1" {
	name = "domain_001"
}
`
