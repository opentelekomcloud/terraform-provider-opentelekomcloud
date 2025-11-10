package gemini

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccGeminiDBInstancesDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstancesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstancesDataSourceID("data.opentelekomcloud_gemini_instances_v3.test"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_gemini_instances_v3.test", "instances.#", "1"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBInstancesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GeminiDB instance data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GeminiDB instances data source ID not set ")
		}

		return nil
	}
}

func testAccGeminiDBInstancesDataSource_basic(rName string) string {
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

data "opentelekomcloud_gemini_instances_v3" "test" {
  name = opentelekomcloud_gemini_instance_v3.test.name
  depends_on = [
    opentelekomcloud_gemini_instance_v3.test,
  ]
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, rName)
}
