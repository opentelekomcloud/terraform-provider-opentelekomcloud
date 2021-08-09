package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vbs/v2/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccVBSBackupPolicyV2_basic(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("opentelekomcloud_vbs_backup_policy_v2.vbs", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "name", "policy_001"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "status", "ON"),
				),
			},
			{
				Config: testAccVBSBackupPolicyV2Update,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("opentelekomcloud_vbs_backup_policy_v2.vbs", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "name", "policy_001_update"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "status", "ON"),
				),
			},
		},
	})
}

func TestAccVBSBackupPolicyV2_rentention_day(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupPolicyV2RententionDay,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("opentelekomcloud_vbs_backup_policy_v2.vbs", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "name", "policy_002"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "status", "ON"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "rentention_day", "30"),
				),
			},
		},
	})
}

func testAccVBSBackupPolicyV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	vbsClient, err := config.VbsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vbs_backup_policy_v2" {
			continue
		}

		_, err := policies.List(vbsClient, policies.ListOpts{ID: rs.Primary.ID})
		if err != nil {
			return fmt.Errorf("backup Policy still exists")
		}
	}

	return nil
}

func testAccVBSBackupPolicyV2Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		vbsClient, err := config.VbsV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud vbs client: %s", err)
		}

		policyList, err := policies.List(vbsClient, policies.ListOpts{ID: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := policyList[0]
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("backup policy not found")
		}

		*policy = found

		return nil
	}
}

const testAccVBSBackupPolicyV2Basic = `
resource "opentelekomcloud_vbs_backup_policy_v2" "vbs" {
  name                = "policy_001"
  start_time          = "12:00"
  status              = "ON"
  retain_first_backup = "N"
  rentention_num      = 2
  frequency           = 1
  tags {
    key   = "k2"
    value = "v2"
  }
}
`

const testAccVBSBackupPolicyV2Update = `
resource "opentelekomcloud_vbs_backup_policy_v2" "vbs" {
  name                = "policy_001_update"
  start_time          = "12:00"
  status              = "ON"
  retain_first_backup = "N"
  rentention_num      = 2
  frequency           = 1
  tags {
    key   = "k2"
    value = "v2"
  }
}
`

const testAccVBSBackupPolicyV2RententionDay = `
resource "opentelekomcloud_vbs_backup_policy_v2" "vbs" {
  name                = "policy_002"
  start_time          = "00:00,12:00"
  retain_first_backup = "N"
  rentention_day      = 30
  week_frequency      = ["MON", "WED", "SAT"]
}
`
