package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/backups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceDdsBackupName = "opentelekomcloud_dds_backup_v3.backup"

func getBackupFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DdsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud DDSv3 client client: %s", err)
	}
	backupId := state.Primary.ID
	listOpts := backups.ListBackupsOpts{
		BackupId: backupId,
	}
	r, err := backups.List(client, listOpts)
	if len(r.Backups) == 0 {
		return nil, fmt.Errorf("error retrieving DDS backup by backup ID: %s", backupId)
	}
	return r.Backups[0], err
}

func TestAccDdsBackupV3_basic(t *testing.T) {
	var backupResp interface{}

	name := fmt.Sprintf("dds_acc_backup_%s", acctest.RandString(5))
	rc := common.InitResourceCheck(
		resourceDdsBackupName,
		&backupResp,
		getBackupFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsBackupV3_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceDdsBackupName, "name", name),
					resource.TestCheckResourceAttr(resourceDdsBackupName, "description", "this is a test dds instance"),
					resource.TestCheckResourceAttr(resourceDdsBackupName, "type", "Manual"),
					resource.TestCheckResourceAttr(resourceDdsBackupName, "status", "COMPLETED"),
					resource.TestCheckResourceAttr(resourceDdsBackupName, "datastore.0.type", "DDS-Community"),
				),
			},
			{
				ResourceName:      resourceDdsBackupName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDdsBackupV3ImportStateFunc(resourceDdsBackupName),
			},
		},
	})
}

func testAccDdsBackupV3ImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testDdsBackupV3_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dds_backup_v3" "backup" {
  instance_id = opentelekomcloud_dds_instance_v3.instance.id
  name        = "%s"
  description = "this is a test dds instance"
}
`, TestAccDDSInstanceV3ConfigSingle, name)
}
