package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccRdsReadReplicaV3Basic(t *testing.T) {
	postfix := tools.RandomString("rr", 3)
	var rdsInstance instances.InstanceResponse

	resName := "opentelekomcloud_rds_read_replica_v3.replica"

	secondAZ := "eu-de-03"

	if env.OS_AVAILABILITY_ZONE == secondAZ {
		t.Skip("OS_AVAILABILITY_ZONE should be set to value !=", secondAZ)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsReadReplicaV3BasicNoIP(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resName, &rdsInstance),
					resource.TestCheckResourceAttr(resName, "availability_zone", secondAZ),
					resource.TestCheckResourceAttr(resName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resName, "public_ips.#", "0"),
				),
			},
			{
				Config: testAccRdsReadReplicaV3Basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resName, &rdsInstance),
					resource.TestCheckResourceAttr(resName, "availability_zone", secondAZ),
					resource.TestCheckResourceAttr(resName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resName, "public_ips.#", "1"),
				),
			},
			{
				Config: testAccRdsReadReplicaV3BasicNoIP(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resName, &rdsInstance),
					resource.TestCheckResourceAttr(resName, "availability_zone", secondAZ),
					resource.TestCheckResourceAttr(resName, "volume.0.size", "40"),
				),
			},
		},
	})
}

func testAccRdsReadReplicaV3Basic(postfix string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.medium"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  tag = {
    foo = "bar"
    key = "value"
  }
}

resource "opentelekomcloud_compute_floatingip_v2" "eip" {}

resource "opentelekomcloud_rds_read_replica_v3" "replica" {
  name          = "test-replica"
  replica_of_id = opentelekomcloud_rds_instance_v3.instance.id
  flavor_ref    = "${opentelekomcloud_rds_instance_v3.instance.flavor}.rr"

  availability_zone = "eu-de-03"

  public_ips = [opentelekomcloud_compute_floatingip_v2.eip.address]

  volume {
    type = "COMMON"
  }
}

`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsReadReplicaV3BasicNoIP(postfix string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.medium"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  tag = {
    foo = "bar"
    key = "value"
  }
}

resource "opentelekomcloud_compute_floatingip_v2" "eip" {}

resource "opentelekomcloud_rds_read_replica_v3" "replica" {
  name          = "test-replica"
  replica_of_id = opentelekomcloud_rds_instance_v3.instance.id
  flavor_ref    = "${opentelekomcloud_rds_instance_v3.instance.flavor}.rr"

  availability_zone = "eu-de-03"

  public_ips = []

  volume {
    type = "COMMON"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}
