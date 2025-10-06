package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/sqlfilter"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const sqlControlRuleResourceName = "opentelekomcloud_taurusdb_mysql_sql_control_rule_v3.test"

func getTaurusDbMysqlSqlControlRuleResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.TaurusDBV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid id format")
	}

	instanceID := parts[0]
	nodeID := parts[1]
	sqlType := parts[2]
	pattern := parts[3]

	opts := sqlfilter.GetSqlFilterRulesOpts{
		NodeId: nodeID,
	}

	resp, err := sqlfilter.GetSqlFilterRules(client, instanceID, opts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SQL control rule: %s", err)
	}

	for _, rule := range resp.SqlFilterRules {
		if rule.SqlType == sqlType {
			for _, p := range rule.Patterns {
				if p.Pattern == pattern {
					return p, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("TaurusDB MySQL SQL control rule not found")
}

func TestAccTaurusDBMySQLSqlControlRule_basic(t *testing.T) {
	var obj interface{}

	name := "tf_taurusdb_" + acctest.RandString(3)

	rc := common.InitResourceCheck(
		sqlControlRuleResourceName,
		&obj,
		getTaurusDbMysqlSqlControlRuleResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBMySQLSqlControlRuleBasic(name, 20),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(sqlControlRuleResourceName, "instance_id",
						"opentelekomcloud_taurusdb_mysql_instance_v3.test", "id"),
					resource.TestCheckResourceAttrPair(sqlControlRuleResourceName, "node_id",
						"opentelekomcloud_taurusdb_mysql_instance_v3.test", "nodes.0.id"),
					resource.TestCheckResourceAttr(sqlControlRuleResourceName, "sql_type", "SELECT"),
					resource.TestCheckResourceAttr(sqlControlRuleResourceName, "pattern", "select~from~t1"),
					resource.TestCheckResourceAttr(sqlControlRuleResourceName, "max_concurrency", "20"),
				),
			},
			{
				Config: testAccTaurusDBMySQLSqlControlRuleBasic(name, 30),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(sqlControlRuleResourceName, "max_concurrency", "30"),
				),
			},
			{
				ResourceName:      sqlControlRuleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTaurusDBMySQLSqlControlRuleBasic(name string, maxConcurrency int) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_taurusdb_mysql_sql_control_rule_v3" "test" {
  instance_id     = opentelekomcloud_taurusdb_mysql_instance_v3.test.id
  node_id         = opentelekomcloud_taurusdb_mysql_instance_v3.test.nodes[0].id
  sql_type        = "SELECT"
  pattern         = "select~from~t1"
  max_concurrency = %d
}
`, testAccTaurusDBMySqlInstanceForSqlControl(name), maxConcurrency)
}

func testAccTaurusDBMySqlInstanceForSqlControl(name string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_taurusdb_mysql_instance_v3" "test" {
  name                     = "%s"
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  security_group_id        = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  flavor                   = "gaussdb.mysql.xlarge.x86.8"
  password                 = "Test123!@#"
  availability_zone_mode   = "multi"
  master_availability_zone = "eu-de-01"
  read_replicas            = 1

  datastore {
    engine  = "gaussdb-mysql"
    version = "8.0"
  }

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 7
  }
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, name)
}
