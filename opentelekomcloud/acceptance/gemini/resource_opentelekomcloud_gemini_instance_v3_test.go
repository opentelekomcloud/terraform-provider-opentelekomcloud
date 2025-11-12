package gemini

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/gemini"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/instance"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccGeminiDBInstance_basic(t *testing.T) {
	var inst instance.ListResult

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_gemini_instance_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckGeminiDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstanceConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceExists(resourceName, &inst),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttr(resourceName, "node_num", "4"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "100"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "geminidb.cassandra.xlarge.8"),
				),
			},
			{
				Config: testAccGeminiDBInstanceConfigUpdate(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceExists(resourceName, &inst),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "node_num", "5"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "120"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "geminidb.cassandra.2xlarge.8"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "21"),
				),
			},
		},
	})
}

func TestAccGeminiDBInstance_withTemplate(t *testing.T) {
	var inst instance.ListResult

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_gemini_instance_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckGeminiDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstanceConfigWithTemplate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceExists(resourceName, &inst),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttr(resourceName, "node_num", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "configuration_id"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBInstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.GeminiDBV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_gemini_instance_v3" {
			continue
		}

		found, err := gemini.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}
		if found.Id != "" {
			return fmt.Errorf("instance <%s> still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckGeminiDBInstanceExists(n string, inst *instance.ListResult) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.GeminiDBV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating GeminiDB client: %s", err)
		}

		found, err := gemini.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("instance <%s> not found", rs.Primary.ID)
			}
			return err
		}
		if found.Id == "" {
			return fmt.Errorf("instance <%s> not found", rs.Primary.ID)
		}
		*inst = *found

		return nil
	}
}

func testAccGeminiDBInstanceConfigBasic(rName string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_gemini_instance_v3" "test" {
  name        = "%s"
  password    = "Test@12345678"
  flavor      = "geminidb.cassandra.xlarge.8"
  volume_size = 100
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  node_num    = 4

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone = "eu-de-01,eu-de-02,eu-de-03"

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, rName)
}

func testAccGeminiDBInstanceConfigUpdate(rName string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_gemini_instance_v3" "test" {
  name        = "%s"
  password    = "Test@12345678"
  flavor      = "geminidb.cassandra.2xlarge.8"
  volume_size = 120
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  node_num    = 5

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone = "eu-de-01,eu-de-02,eu-de-03"

  backup_strategy {
    start_time = "04:00-05:00"
    keep_days  = 21
  }

  tags = {
    foo     = "bar"
    key     = "value"
    updated = "tag"
  }
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, rName)
}

func testAccGeminiDBInstanceConfigWithTemplate(rName string) string {
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
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, rName, rName)
}
