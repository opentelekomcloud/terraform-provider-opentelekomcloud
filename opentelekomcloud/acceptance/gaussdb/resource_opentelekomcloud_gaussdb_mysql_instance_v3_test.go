package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const instanceV3ResourceName = "opentelekomcloud_gaussdb_mysql_instance_v3.instance"

func TestAccMysqlGaussdbInstanceV3Basic(t *testing.T) {
	name := "tf_gaussdb_instance" + acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussdbMySqlInstanceV3Basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", name),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "gaussdb.mysql.xlarge.x86.8"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "availability_zone_mode", "multi"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "port", "3306"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "nodes.0.ram", "32"),
				),
			},
			{
				Config: testAccGaussdbMySqlInstanceV3Update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", name),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "gaussdb.mysql.2xlarge.x86.8"),
				),
			},
		},
	})
}

func TestAccGaussDBMySqlV3_importBasic(t *testing.T) {
	name := "tf_gaussdb_import" + acctest.RandString(3)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussdbMySqlInstanceV3Basic(name),
			},

			{
				ResourceName:      instanceV3ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func testAccGaussdbMySqlInstanceV3Basic(postfix string) string {
	return fmt.Sprintf(`
%s
%s
resource "opentelekomcloud_gaussdb_mysql_instance_v3" "instance" {
  name                     = "%s"
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  security_group_id        = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  flavor                   = "gaussdb.mysql.xlarge.x86.8"
  password                 = "Test123!@#"
  availability_zone_mode   = "multi"
  master_availability_zone = "eu-de-01"
  read_replicas            = 1
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, postfix)
}

func testAccGaussdbMySqlInstanceV3Update(postfix string) string {
	return fmt.Sprintf(`
%s
%s
resource "opentelekomcloud_gaussdb_mysql_instance_v3" "instance" {
  name                     = "%s"
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  security_group_id        = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  flavor                   = "gaussdb.mysql.2xlarge.x86.8"
  password                 = "Test123!@#"
  availability_zone_mode   = "multi"
  master_availability_zone = "eu-de-01"
  read_replicas            = 1
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, postfix)
}
