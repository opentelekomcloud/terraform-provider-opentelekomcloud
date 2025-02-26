package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccRdsPublicIpAssociateV3Basic(t *testing.T) {
	resName := "opentelekomcloud_rds_public_ip_associate_v3.public_ip"
	postfix := acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsPublicIpAssociateV3Basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
				),
			},
			{
				Config: testAccRdsPublicIpAssociateV3Update(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
				),
			},
		},
	})
}

func testAccRdsPublicIpAssociateV3Basic(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor = "rds.pg.n1.large.4"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 0
  }
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
  lower_case_table_names = "0"
}
resource "opentelekomcloud_rds_public_ip_associate_v3" "public_ip" {
  instance_id  = opentelekomcloud_rds_instance_v3.instance.id
  public_ip    = opentelekomcloud_compute_floatingip_v2.eip_1.address
  public_ip_id = opentelekomcloud_compute_floatingip_v2.eip_1.id
}

resource "opentelekomcloud_compute_floatingip_v2" "eip_1" {}

resource "opentelekomcloud_compute_floatingip_v2" "eip_2" {}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsPublicIpAssociateV3Update(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor = "rds.pg.n1.large.4"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 0
  }
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
  lower_case_table_names = "0"
}
resource "opentelekomcloud_rds_public_ip_associate_v3" "public_ip" {
  instance_id  = opentelekomcloud_rds_instance_v3.instance.id
  public_ip    = opentelekomcloud_compute_floatingip_v2.eip_2.address
  public_ip_id = opentelekomcloud_compute_floatingip_v2.eip_2.id
}

resource "opentelekomcloud_compute_floatingip_v2" "eip_1" {}

resource "opentelekomcloud_compute_floatingip_v2" "eip_2" {}


`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}
