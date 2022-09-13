package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const domainDataName = "data.opentelekomcloud_sdrs_domain_v1.domain_1"

func TestAccSdrsDomainV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsDomainV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsDomainV1DataSourceID(domainDataName),
					resource.TestCheckResourceAttr(domainDataName, "name", "domain_001"),
				),
			},
		},
	})
}

func TestAccSdrsDomainV1DataSource_withoutName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsDomainV1DataSource_withoutName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsDomainV1DataSourceID(domainDataName),
					resource.TestCheckResourceAttr(domainDataName, "name", "domain_001"),
				),
			},
		},
	})
}

func testAccCheckSdrsDomainV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find SDRS domain data source: %s", n)
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

const testAccSdrsDomainV1DataSource_withoutName = `
data "opentelekomcloud_sdrs_domain_v1" "domain_1" {
}
`
