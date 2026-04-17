package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/backup"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getTaurusDbMysqlBackupResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.TaurusDBV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDb V3 client: %s", err)
	}

	getResp, err := backup.List(client, backup.ListOpts{
		BackupId: state.Primary.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("error retrieving TaurusDb MySQL backup: %s", err)
	}

	if len(getResp) == 0 {
		return nil, fmt.Errorf("TaurusDb backup doesn't exist")
	}

	return getResp[0], nil
}

func TestAccTaurusDbMysqlBackup_basic(t *testing.T) {
	var obj interface{}

	name := common.RandomAccResourceName()
	rName := "opentelekomcloud_taurusdb_mysql_backup_v3.test"

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getTaurusDbMysqlBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBackup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"opentelekomcloud_taurusdb_mysql_instance_v3.instance", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttrSet(rName, "begin_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "take_up_time"),
					resource.TestCheckResourceAttrSet(rName, "size"),
					resource.TestCheckResourceAttrSet(rName, "datastore.0.type"),
					resource.TestCheckResourceAttrSet(rName, "datastore.0.version"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"end_time",
				},
			},
		},
	})
}

func testBackup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_taurusdb_mysql_backup_v3" "test" {
  instance_id = opentelekomcloud_taurusdb_mysql_instance_v3.instance.id
  name        = "%s"
  description = "test description"
}
`, testAccTaurusDBMySqlInstanceV3Basic(name), name)
}
