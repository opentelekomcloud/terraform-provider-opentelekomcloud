package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/configurations"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const pgResourceName = "opentelekomcloud_rds_parametergroup_v3.pg_1"

func TestAccRdsConfigurationV3_basic(t *testing.T) {
	var config configurations.Configuration

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsConfigV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsConfigV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsConfigV3Exists(pgResourceName, &config),
					resource.TestCheckResourceAttr(pgResourceName, "name", "pg_create"),
					resource.TestCheckResourceAttr(pgResourceName, "description", "some description"),
				),
			},
			{
				Config: testAccRdsConfigV3Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsConfigV3Exists(pgResourceName, &config),
					resource.TestCheckResourceAttr(pgResourceName, "name", "pg_update"),
					resource.TestCheckResourceAttr(pgResourceName, "description", "updated description"),
				),
			},
		},
	})
}

func TestAccRdsConfigurationV3_invalidDbVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsConfigV3Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccRdsConfigV3InvalidDataStoreVersion,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't find version.+`),
			},
		},
	})
}

func testAccCheckRdsConfigV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	rdsClient, err := config.RdsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud RDSv3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rds_parametergroup_v3" {
			continue
		}

		_, err := configurations.Get(rdsClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("RDSv3 configuration still exists")
		}
	}

	return nil
}

func testAccCheckRdsConfigV3Exists(n string, configuration *configurations.Configuration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.RdsV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud RDSv3 client: %s", err)
		}

		found, err := configurations.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("RDSv3 configuration not found")
		}

		*configuration = *found

		return nil
	}
}

const (
	testAccRdsConfigV3Basic = `
resource "opentelekomcloud_rds_parametergroup_v3" "pg_1" {
  name        = "pg_create"
  description = "some description"

  values = {
    max_connections = "10"
    autocommit      = "OFF"
  }

  datastore {
    type    = "mysql"
    version = "5.6"
  }
}
`
	testAccRdsConfigV3Update = `
resource "opentelekomcloud_rds_parametergroup_v3" "pg_1" {
  name        = "pg_update"
  description = "updated description"

  values = {
    max_connections = "10"
    autocommit      = "OFF"
  }

  datastore {
    type    = "mysql"
    version = "5.6"
  }
}
`
	testAccRdsConfigV3InvalidDataStoreVersion = `
resource "opentelekomcloud_rds_parametergroup_v3" "pg_1" {
  name        = "pg_update"
  description = "updated description"

  values = {
    max_connections = "10"
    autocommit      = "OFF"
  }

  datastore {
    type    = "mysql"
    version = "1"
  }
}
`
)
