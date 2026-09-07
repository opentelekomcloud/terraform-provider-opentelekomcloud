package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcQuotasV1DataSource_basic(t *testing.T) {
	const (
		allQuotas      = "data.opentelekomcloud_vpc_quotas_v1.all"
		publicIPQuotas = "data.opentelekomcloud_vpc_quotas_v1.public_ip"
	)

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcQuotasV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(allQuotas, "quotas.#", regexp.MustCompile(`^[1-9][0-9]*$`)),
					resource.TestCheckResourceAttrSet(allQuotas, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(allQuotas, "quotas.0.used"),
					resource.TestCheckResourceAttrSet(allQuotas, "quotas.0.quota"),
					resource.TestCheckResourceAttrSet(allQuotas, "quotas.0.min"),
					resource.TestCheckResourceAttr(publicIPQuotas, "type", "publicIp"),
					resource.TestCheckResourceAttr(publicIPQuotas, "quotas.#", "1"),
					resource.TestCheckResourceAttr(publicIPQuotas, "quotas.0.type", "publicIp"),
				),
			},
		},
	})
}

const testAccVpcQuotasV1DataSourceBasic = `
data "opentelekomcloud_vpc_quotas_v1" "all" {}

data "opentelekomcloud_vpc_quotas_v1" "public_ip" {
  type = "publicIp"
}
`
