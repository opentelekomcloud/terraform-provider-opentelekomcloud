package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccCBRPolicyV3_basic(t *testing.T) {
	var cbrPolicy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccFlavorPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRPolicyV3_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRPolicyV3Exists("opentelekomcloud_cbr_policy_v3.policy", &cbrPolicy),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_policy_v3.policy", "name", "test-policy"),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_policy_v3.policy", "operation_type", "backup"),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_policy_v3.policy", "enabled", "true"),
				),
			},
			{
				Config: testCBRPolicyV3_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRPolicyV3Exists("opentelekomcloud_cbr_policy_v3.policy", &cbrPolicy),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_policy_v3.policy", "name", "name2"),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_policy_v3.policy", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccCBRPolicyV3_minConfig(t *testing.T) {
	var cbrPolicy policies.Policy
	policyRes := "opentelekomcloud_cbr_policy_v3.policy"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccFlavorPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRPolicyV3_minConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRPolicyV3Exists(policyRes, &cbrPolicy),
					resource.TestCheckResourceAttr(policyRes, "name", "some-policy-min"),
					resource.TestCheckResourceAttr(policyRes, "operation_type", "backup"),
					resource.TestCheckResourceAttr(policyRes, "enabled", "true"),
				),
			},
			{
				Config: testCBRPolicyV3_minOperationDef,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(policyRes, "operation_definition.0.timezone", "UTC+03:00"),
				),
			},
		},
	})
}

func testAccCheckCBRPolicyV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*cfg.Config)
	asClient, err := config.CbrV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cbr_policy_v3" {
			continue
		}

		_, err := policies.Get(asClient, rs.Primary.ID).Extract()
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

		config := testAccProvider.Meta().(*cfg.Config)
		client, err := config.CbrV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
		}

		found, err := policies.Get(client, rs.Primary.ID).Extract()
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
	testCBRPolicyV3_basic = `
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
	testCBRPolicyV3_update = `
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
	testCBRPolicyV3_minConfig = `
resource opentelekomcloud_cbr_policy_v3 policy {
  name            = "some-policy-min"
  operation_type  = "backup"
  trigger_pattern = ["FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00"]
}
`

	testCBRPolicyV3_minOperationDef = `
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
