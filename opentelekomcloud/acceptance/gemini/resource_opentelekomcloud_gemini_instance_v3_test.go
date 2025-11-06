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
