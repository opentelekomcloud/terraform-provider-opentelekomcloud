package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourcePolicyName = "opentelekomcloud_cbr_policy_v3.policy"

func getPolicyResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CbrV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR v3 client: %s", err)
	}
	return policies.Get(c, state.Primary.ID)
}

func TestAccCBRPolicyV3_basic(t *testing.T) {
	var cbrPolicy policies.Policy
	rc := common.InitResourceCheck(resourcePolicyName, &cbrPolicy, getPolicyResourceFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookOne(t, quotas.CBRPolicy)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCBRPolicyV3Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "test-policy"),
					resource.TestCheckResourceAttr(resourcePolicyName, "operation_type", "backup"),
					resource.TestCheckResourceAttr(resourcePolicyName, "enabled", "true"),
				),
			},
			{
				Config: testAccCBRPolicyV3Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "name2"),
					resource.TestCheckResourceAttr(resourcePolicyName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccCBRPolicyV3_minConfig(t *testing.T) {
	var cbrPolicy policies.Policy
	rc := common.InitResourceCheck(resourcePolicyName, &cbrPolicy, getPolicyResourceFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookOne(t, quotas.CBRPolicy)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCBRPolicyV3MinConfig,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "some-policy-min"),
					resource.TestCheckResourceAttr(resourcePolicyName, "operation_type", "backup"),
					resource.TestCheckResourceAttr(resourcePolicyName, "enabled", "true"),
				),
			},
			{
				Config: testAccCBRPolicyV3MinOperationDef,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePolicyName, "operation_definition.0.timezone", "UTC+03:00"),
				),
			},
		},
	})
}

func TestAccPolicy_replication(t *testing.T) {
	var policy policies.Policy
	name := fmt.Sprintf("cbr_acc_policy_%s", acctest.RandString(5))

	rc := common.InitResourceCheck(resourcePolicyName, &policy, getPolicyResourceFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckReplication(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicy_replication(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", name),
					resource.TestCheckResourceAttr(resourcePolicyName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourcePolicyName, "operation_type", "replication"),
					resource.TestCheckResourceAttr(resourcePolicyName, "destination_region", env.OS_DEST_REGION),
					resource.TestCheckResourceAttr(resourcePolicyName, "destination_project_id", env.OS_DEST_PROJECT_ID),
				),
			},
			{
				ResourceName:      resourcePolicyName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCBRPolicyV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.CbrV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cbr_policy_v3" {
			continue
		}

		_, err := policies.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("CBRv3 policy still exists")
		}
	}

	return nil
}

const (
	testAccCBRPolicyV3Basic = `
resource opentelekomcloud_cbr_policy_v3 policy {
  name                 = "test-policy"
  operation_type       = "backup"
  trigger_pattern      = ["FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00"]
  operation_definition {
    day_backups   = 1
    week_backups  = 2
    year_backups  = 3
    month_backups = 4
    max_backups   = 10
    timezone      = "UTC+03:00"
  }
  enabled = "true"
}
`
	testAccCBRPolicyV3Update = `
resource opentelekomcloud_cbr_policy_v3 policy {
  name                 = "name2"
  operation_type       = "backup"
  trigger_pattern      = ["FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00"]
  operation_definition {
    day_backups   = 1
    week_backups  = 2
    year_backups  = 3
    month_backups = 4
    max_backups   = 10
    timezone      = "UTC+03:00"
  }
  enabled = "false"
}
`
	testAccCBRPolicyV3MinConfig = `
resource opentelekomcloud_cbr_policy_v3 policy {
  name            = "some-policy-min"
  operation_type  = "backup"
  trigger_pattern = ["FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00"]
}
`

	testAccCBRPolicyV3MinOperationDef = `
resource opentelekomcloud_cbr_policy_v3 policy {
  name            = "some-policy-min"
  operation_type  = "backup"
  trigger_pattern = ["FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00"]

  operation_definition {
    timezone                = "UTC+03:00"
  }
}
`
)

func testAccPolicy_replication(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cbr_policy_v3" "policy" {
  name                   = "%[1]s"
  operation_type         = "replication"
  destination_region     = "%[2]s"
  destination_project_id = "%[3]s"

  trigger_pattern = ["FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00"]

  operation_definition {
    day_backups   = 1
    week_backups  = 2
    year_backups  = 3
    month_backups = 4
    max_backups   = 10
    timezone      = "UTC+03:00"
  }
}
`, name, env.OS_DEST_REGION, env.OS_DEST_PROJECT_ID)
}
