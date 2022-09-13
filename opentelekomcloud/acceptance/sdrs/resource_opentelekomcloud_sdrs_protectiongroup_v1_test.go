package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/protectiongroups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const pgResourceName = "opentelekomcloud_sdrs_protectiongroup_v1.group_1"

func TestAccSdrsProtectionGroupV1_basic(t *testing.T) {
	var group protectiongroups.Group

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSdrsProtectionGroupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsProtectionGroupV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsProtectionGroupV1Exists(pgResourceName, &group),
					resource.TestCheckResourceAttr(pgResourceName, "name", "group_1"),
				),
			},
			{
				Config: testAccSdrsProtectionGroupV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsProtectionGroupV1Exists(pgResourceName, &group),
					resource.TestCheckResourceAttr(pgResourceName, "name", "group_updated"),
				),
			},
		},
	})
}

func testAccCheckSdrsProtectionGroupV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SdrsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud SDRS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_sdrs_protectiongroup_v1" {
			continue
		}

		_, err := protectiongroups.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SDRS protectiongroup still exists")
		}
	}

	return nil
}

func testAccCheckSdrsProtectionGroupV1Exists(n string, group *protectiongroups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.SdrsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud SDRS client: %s", err)
		}

		found, err := protectiongroups.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("SDRS protectiongroup not found")
		}

		*group = *found

		return nil
	}
}

var testAccSdrsProtectionGroupV1Basic = fmt.Sprintf(`
%s

data "opentelekomcloud_sdrs_domain_v1" "domain_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_1"
  description              = "test description"
  source_availability_zone = "eu-de-02"
  target_availability_zone = "eu-de-01"
  domain_id                = data.opentelekomcloud_sdrs_domain_v1.domain_1.id
  source_vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  dr_type                  = "migration"
}
`, common.DataSourceSubnet)

var testAccSdrsProtectionGroupV1Update = fmt.Sprintf(`
%s

data "opentelekomcloud_sdrs_domain_v1" "domain_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_updated"
  description              = "test description"
  source_availability_zone = "eu-de-02"
  target_availability_zone = "eu-de-01"
  domain_id                = data.opentelekomcloud_sdrs_domain_v1.domain_1.id
  source_vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  dr_type                  = "migration"
}
`, common.DataSourceSubnet)
