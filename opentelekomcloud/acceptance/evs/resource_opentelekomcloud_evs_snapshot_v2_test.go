package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v2/snapshots"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getVolumeSnapshotResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.BlockStorageV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud BlockStorage V2 client: %s", err)
	}
	return snapshots.Get(c, state.Primary.ID)
}
func TestAccEvsSnapshot_basic(t *testing.T) {
	var (
		snapshot     snapshots.Snapshot
		rName        = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
		resourceName = "opentelekomcloud_evs_snapshot_v2.test"
	)

	rc := common.InitResourceCheck(
		resourceName,
		&snapshot,
		getVolumeSnapshotResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsSnapshotV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Daily backup"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "size"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccEvsSnapshotV2_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Daily backup update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata", "force",
				},
			},
		},
	})
}

func testAccEvsSnapshotV2_base(rName string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_compute_availability_zones_v2" "test" {}

resource "opentelekomcloud_evs_volume_v3" "test" {
  name              = "%s"
  description       = "Created by acc test"
  availability_zone = data.opentelekomcloud_compute_availability_zones_v2.test.names[0]
  volume_type       = "SAS"
  size              = 12
}
`, rName)
}

func testAccEvsSnapshotV2_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_evs_snapshot_v2" "test" {
  volume_id   = opentelekomcloud_evs_volume_v3.test.id
  name        = "%[2]s"
  description = "Daily backup"
  metadata = {
    foo = "bar"
    key = "value"
  }
}
`, testAccEvsSnapshotV2_base(rName), rName)
}

func testAccEvsSnapshotV2_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_evs_snapshot_v2" "test" {
  volume_id   = opentelekomcloud_evs_volume_v3.test.id
  name        = "%[2]s-update"
  description = "Daily backup update"
  metadata = {
    foo = "bar"
    key = "value"
  }
}
`, testAccEvsSnapshotV2_base(rName), rName)
}
