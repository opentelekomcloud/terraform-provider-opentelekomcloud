package drs

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/drs/v3/public"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getDrsJobResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DrsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DRS client, err: %s", err)
	}
	detailResp, err := public.BatchListTaskDetails(client, public.BatchQueryTaskOpts{Jobs: []string{state.Primary.ID}})
	if err != nil {
		return nil, err
	}
	status := detailResp.Results[0].Status
	if status == "DELETED" {
		return nil, golangsdk.ErrDefault404{}
	}
	return detailResp, nil
}

func TestAccResourceDrsJob_basic(t *testing.T) {
	var obj public.BatchListTaskDetailsResponse
	resourceName := "opentelekomcloud_drs_task_v3.test"
	name := common.RandomAccResourceName()
	dbName := common.RandomAccResourceName()
	pwd := "TestDrs@123"
	startTime := strconv.FormatInt(time.Now().Add(time.Hour).UnixMilli(), 10)

	rc := common.InitResourceCheck(
		resourceName,
		&obj,
		getDrsJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJob_migrate_mysql(name, dbName, pwd, "", startTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "WAITING_FOR_START"),
				),
			},
			{
				Config: testAccDrsJob_migrate_mysql(name, dbName, pwd, "start", ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "migration"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.instance_id",
						"opentelekomcloud_rds_instance_v3.test2", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", name),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"source_db.0.password", "destination_db.0.password",
					"expired_days", "migrate_definer", "force_destroy", "status", "auto_renew", "updated_at", "policy_config",
					"source_db.0.ip", "destination_db.0.ip", "engine_type", "tags", "status", "net_type", "start_time", "action"},
			},
		},
	})
}

func testAccDrsJob_mysql(index int, name, pwd, ip string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_rds_instance_v3" "test%d" {
  depends_on = [
    opentelekomcloud_networking_secgroup_rule_v2.ingress,
    opentelekomcloud_networking_secgroup_rule_v2.egress,
  ]
  name                = "%s%d"
  flavor              = "rds.mysql.x1.large.2.ha"
  security_group_id   = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  private_ip          = "%s"
  ha_replication_mode = "semisync"

  availability_zone = ["eu-de-02", "eu-de-03"]

  db {
    password = "%s"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, index, name, index, ip, pwd)
}

const testAccSecgroupRule string = `
resource "opentelekomcloud_networking_secgroup_rule_v2" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 3306
  port_range_min    = 3306
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "egress" {
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
}
`

func testAccDrsJob_migrate_mysql(name, dbName, pwd, action, startTime string) string {
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")

	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s

%[5]s

resource "opentelekomcloud_drs_task_v3" "test" {
  name           = "%[6]s"
  type           = "migration"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "eip"
  migration_type = "FULL_INCR_TRANS"
  description    = "%[6]s"
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = opentelekomcloud_rds_instance_v3.test1.private_ips[0]
    port        = 3306
    user        = "root"
    password    = "%[7]s"
  }


  destination_db {
    ip          = opentelekomcloud_rds_instance_v3.test2.private_ips[0]
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "%[7]s"
    instance_id = opentelekomcloud_rds_instance_v3.test2.id
    subnet_id   = opentelekomcloud_rds_instance_v3.test2.subnet_id
  }

  tags = {
    key = "%[6]s"
  }

  action     = "%[8]s"
  start_time = "%[9]s"

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, testAccSecgroupRule, sourceDb, destDb, name, pwd, action, startTime)
}

func TestAccResourceDrsJob_down_migration(t *testing.T) {
	var obj public.BatchListTaskDetailsResponse
	resourceName := "opentelekomcloud_drs_task_v3.test"
	name := common.RandomAccResourceName()
	dbName := common.RandomAccResourceName()
	pwd := "_Hc142223"

	rc := common.InitResourceCheck(
		resourceName,
		&obj,
		getDrsJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJob_migrate_mysql_down(name, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "migration"),
					resource.TestCheckResourceAttr(resourceName, "direction", "down"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "vpc"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readonly", "false"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "source_db.0.instance_id",
						"opentelekomcloud_rds_instance_v3.test1", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.ads-service", "new_version"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "poc"),
				),
			},
		},
	})
}

func testAccDrsJob_migrate_mysql_down(name, dbName, pwd string) string {
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")

	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s

%[5]s

resource "opentelekomcloud_drs_task_v3" "test" {
  name          = "%[6]s"
  type          = "migration"
  engine_type   = "mysql"
  direction     = "down"
  net_type      = "vpc"
  force_destroy = true

  source_db {
    engine_type = "mysql"
    ip          = opentelekomcloud_rds_instance_v3.test1.private_ips[0]
    port        = 3306
    user        = "root"
    password    = "%[7]s"
    instance_id = opentelekomcloud_rds_instance_v3.test1.id
    subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  destination_db {
    engine_type = "mysql"
    ip          = opentelekomcloud_rds_instance_v3.test2.private_ips[0]
    port        = 3306
    user        = "root"
    password    = "%[7]s"
    subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
    vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }

  destination_db_readonly = false

  tags = {
    "ads-service" = "new_version"
    "env"         = "poc"
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, testAccSecgroupRule, sourceDb, destDb, name, pwd)
}

func TestAccResourceDrsJob_sync(t *testing.T) {
	var obj public.BatchListTaskDetailsResponse
	resourceName := "opentelekomcloud_drs_task_v3.test"
	name := common.RandomAccResourceName()
	dbName := common.RandomAccResourceName()
	pwd := "TestDrs@123"

	rc := common.InitResourceCheck(
		resourceName,
		&obj,
		getDrsJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJob_synchronize_mysql(name, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "sync"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "vpc"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "progress"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					waitForJobStatus(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"source_db.0.password", "destination_db.0.password",
					"expired_days", "migrate_definer", "force_destroy", "status", "auto_renew", "updated_at", "policy_config",
					"source_db.0.ip", "destination_db.0.ip", "engine_type", "tags", "status", "net_type", "start_time", "action"},
			},
		},
	})
}

func testAccDrsJob_synchronize_mysql(name, dbName, pwd string) string {
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")

	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s

%[5]s

resource "opentelekomcloud_drs_task_v3" "test" {
  name           = "%[6]s"
  type           = "sync"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "vpc"
  migration_type = "FULL_INCR_TRANS"
  description    = "%[6]s"
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = opentelekomcloud_rds_instance_v3.test1.private_ips[0]
    port        = 3306
    user        = "root"
    password    = "%[7]s"
    vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
    subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  destination_db {
    region      = "eu-de"
    ip          = opentelekomcloud_rds_instance_v3.test2.private_ips[0]
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "%[7]s"
    instance_id = opentelekomcloud_rds_instance_v3.test2.id
    subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  limit_speed {
    speed      = "15"
    start_time = "16:00"
    end_time   = "17:59"
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, testAccSecgroupRule, sourceDb, destDb, name, pwd)
}

func waitForJobStatus(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceName)
		}

		conf := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := conf.DrsV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating DRS client, err: %s", err)
		}

		// record the start time
		startTime := time.Now()
		for {
			respBody, err := public.BatchListTaskDetails(client, public.BatchQueryTaskOpts{Jobs: []string{rs.Primary.ID}})
			if err != nil {
				return fmt.Errorf("error querying job (%s): %s", rs.Primary.ID, err)
			}
			if len(respBody.Results) == 0 {
				return fmt.Errorf("error querying job (%s): results not found", rs.Primary.ID)
			}
			status := respBody.Results[0].Status
			if status == "INCRE_TRANSFER_STARTED" {
				return nil
			}

			if time.Since(startTime) > 30*time.Minute {
				return fmt.Errorf("error waiting for job status becoming INCRE_TRANSFER_STARTED, time out")
			}
			// lintignore:R018
			time.Sleep(30 * time.Second)
		}
	}
}
