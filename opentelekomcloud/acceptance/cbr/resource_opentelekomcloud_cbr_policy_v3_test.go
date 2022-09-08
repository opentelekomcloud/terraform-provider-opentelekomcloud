package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourcePolicyName = "opentelekomcloud_cbr_policy_v3.policy"

func TestAccCBRPolicyV3_basic(t *testing.T) {
	var cbrPolicy policies.Policy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookOne(t, quotas.CBRPolicy)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRPolicyV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRPolicyV3Exists(resourcePolicyName, &cbrPolicy),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "test-policy"),
					resource.TestCheckResourceAttr(resourcePolicyName, "operation_type", "backup"),
					resource.TestCheckResourceAttr(resourcePolicyName, "enabled", "true"),
				),
			},
			{
				Config: testAccCBRPolicyV3Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRPolicyV3Exists(resourcePolicyName, &cbrPolicy),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "name2"),
					resource.TestCheckResourceAttr(resourcePolicyName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccCBRPolicyV3_minConfig(t *testing.T) {
	var cbrPolicy policies.Policy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookOne(t, quotas.CBRPolicy)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRPolicyV3MinConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRPolicyV3Exists(resourcePolicyName, &cbrPolicy),
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

func testAccCheckCBRPolicyV3Exists(n string, group *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CbrV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
		}

		found, err := policies.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("CBRv3 policy not found")
		}
		group = found

		return nil
	}
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
