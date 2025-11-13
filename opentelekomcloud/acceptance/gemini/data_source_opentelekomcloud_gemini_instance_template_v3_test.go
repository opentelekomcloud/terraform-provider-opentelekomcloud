package gemini

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceGeminiDBInstanceTemplate_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSource := "data.opentelekomcloud_gemini_instance_template_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckGeminiDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGeminiDBInstanceTemplate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceTemplateDataSourceID(dataSource),
					resource.TestCheckResourceAttrSet(dataSource, "id"),
					resource.TestCheckResourceAttrSet(dataSource, "datastore_version_name"),
					resource.TestCheckResourceAttrSet(dataSource, "datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "mode"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.#"),
					testAccCheckGeminiDBInstanceTemplateParameterValue(dataSource, "write_request_timeout_in_ms", "5000"),
					testAccCheckGeminiDBInstanceTemplateParameterValue(dataSource, "slow_query_log_timeout_in_ms", "10000"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBInstanceTemplateDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GeminiDB instance template data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GeminiDB instance template data source ID not set")
		}

		return nil
	}
}

func testAccCheckGeminiDBInstanceTemplateParameterValue(n, paramName, expectedValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GeminiDB instance template data source: %s", n)
		}

		found := false
		for key := range rs.Primary.Attributes {
			if key == fmt.Sprintf("configuration_parameters.%d.name", 0) {
				continue
			}

			// Check each parameter
			for i := 0; ; i++ {
				nameKey := fmt.Sprintf("configuration_parameters.%d.name", i)
				valueKey := fmt.Sprintf("configuration_parameters.%d.value", i)

				name, nameExists := rs.Primary.Attributes[nameKey]
				value, valueExists := rs.Primary.Attributes[valueKey]

				if !nameExists {
					break
				}

				if name == paramName && valueExists {
					if value != expectedValue {
						return fmt.Errorf("parameter %s has value %s, expected %s", paramName, value, expectedValue)
					}
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			return fmt.Errorf("parameter %s not found in configuration_parameters", paramName)
		}

		return nil
	}
}

func testDataSourceGeminiDBInstanceTemplate_basic(rName string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_gemini_template_v3" "test" {
  name           = "%s_template"
  description    = "test configuration"
  instance_type  = "cassandra"
  engine_version = "3.11"

  parameters {
    name  = "write_request_timeout_in_ms"
    value = "5000"
  }

  parameters {
    name  = "slow_query_log_timeout_in_ms"
    value = "10000"
  }
}

resource "opentelekomcloud_gemini_instance_v3" "test" {
  name             = "%s"
  password         = "Test@12345678"
  flavor           = "geminidb.cassandra.xlarge.8"
  volume_size      = 100
  vpc_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id        = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  node_num         = 3
  configuration_id = opentelekomcloud_gemini_template_v3.test.id

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone = "eu-de-01,eu-de-02,eu-de-03"

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 7
  }

  tags = {
    foo = "bar"
  }
}

data "opentelekomcloud_gemini_instance_template_v3" "test" {
  instance_id = opentelekomcloud_gemini_instance_v3.test.id
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, rName, rName)
}
