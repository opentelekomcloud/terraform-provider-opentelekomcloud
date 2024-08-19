package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataSourceName = "data.opentelekomcloud_evs_volumes_v2.test"

func TestAccEvsVolumesDataSource_basic(t *testing.T) {
	dc := common.InitDataSourceCheck(dataSourceName)
	rName := fmt.Sprintf("evs-%s", acctest.RandString(5))
	rDescription := fmt.Sprintf("evs-description-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolumesDataSource_basic(rName, rDescription),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.0.bootable", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.0.description", rDescription),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.0.service_type", "EVS"),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.0.size", "12"),
				),
			},
			{
				Config: testAccEvsVolumesDataSource_updated(rName+"_updated", rDescription+"_updated"),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.0.name", rName+"_updated"),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.0.description", rDescription+"_updated"),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.0.tags.muh", "value-create"),
				),
			},
		},
	})
}

func testAccEvsVolumesDataSource_basic(rName, description string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "%[1]s"
  description       = "%[2]s"
  availability_zone = "%[3]s"
  volume_type       = "SATA"
  size              = 12
}

data "opentelekomcloud_evs_volumes_v2" "test" {
  name = opentelekomcloud_evs_volume_v3.volume_1.name
}

`, rName, description, env.OS_AVAILABILITY_ZONE)
}

func testAccEvsVolumesDataSource_updated(rName, description string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "%[1]s"
  description       = "%[2]s"
  availability_zone = "%[3]s"
  volume_type       = "SATA"
  size              = 12
  tags = {
    muh = "value-create"
    kuh = "value-test"
  }
}

data "opentelekomcloud_evs_volumes_v2" "test" {
  volume_id = opentelekomcloud_evs_volume_v3.volume_1.id
}

`, rName, description, env.OS_AVAILABILITY_ZONE)
}
