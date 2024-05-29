package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v2/instance"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const dcsV2InstanceName = "opentelekomcloud_dcs_instance_v2.instance_1"

func TestAccDcsInstancesV2_basic(t *testing.T) {
	var dcsInstance instance.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV2InstanceBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV2InstanceExists(dcsV2InstanceName, dcsInstance),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "tags.environment", "basic"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "tags.managed_by", "terraform"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "flavor", "redis.ha.xu1.tiny.r2.128"),
				),
			},
			{
				Config: testAccDcsV2InstanceUpdated(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dcsV2InstanceName, "backup_policy.0.begin_at", "01:00-02:00"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "backup_policy.0.save_days", "2"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "backup_policy.0.backup_at.#", "3"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "tags.environment", "update"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "tags.managed_by", "terraform"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "tags.user", "admin"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV2_privateIPs(t *testing.T) {
	var dcsInstance instance.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV2InstancePrivateIPs(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV2InstanceExists(dcsV2InstanceName, dcsInstance),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttrSet(dcsV2InstanceName, "private_ip"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV2_basicSingleInstance(t *testing.T) {
	var dcsInstance instance.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV2InstanceSingle(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV2InstanceExists(dcsV2InstanceName, dcsInstance),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine", "Redis"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV2_basicEngineV3Instance(t *testing.T) {
	var dcsInstance instance.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV2InstanceEngineV3(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV2InstanceExists(dcsV2InstanceName, dcsInstance),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "flavor", "dcs.master_standby"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV2_Whitelist(t *testing.T) {
	var dcsInstance instance.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV2InstanceWhitelist(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV2InstanceExists(dcsV2InstanceName, dcsInstance),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "enable_whitelist", "false"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "whitelist.0.group_name", "test-group-name"),
				),
			},
			{
				Config: testAccDcsV2InstanceWhitelistUpdate(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV2InstanceExists(dcsV2InstanceName, dcsInstance),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine", "Redis"),
				),
			},
			{
				Config: testAccDcsV2InstanceWhitelistSecondUpdate(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV2InstanceExists(dcsV2InstanceName, dcsInstance),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "enable_whitelist", "true"),
				),
			},
			{
				Config: testAccDcsV2InstanceWhitelistThirdUpdate(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV2InstanceExists(dcsV2InstanceName, dcsInstance),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "enable_whitelist", "false"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV2_SSL(t *testing.T) {
	var dcsInstance instance.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV2InstanceSSL(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV2InstanceExists(dcsV2InstanceName, dcsInstance),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine_version", "6.0"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "capacity", "0.125"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttrSet(dcsV2InstanceName, "private_ip"),
					resource.TestCheckResourceAttrSet(dcsV2InstanceName, "port"),
				),
			},
			{
				Config: testAccDcsV2InstanceSSLUpdate(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV2InstanceExists(dcsV2InstanceName, dcsInstance),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "engine_version", "6.0"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "capacity", "0.125"),
					resource.TestCheckResourceAttr(dcsV2InstanceName, "ssl_enable", "false"),
					resource.TestCheckResourceAttrSet(dcsV2InstanceName, "private_ip"),
					resource.TestCheckResourceAttrSet(dcsV2InstanceName, "port"),
				),
			},
		},
	})
}

func TestAccDCSInstanceV2_importBasic(t *testing.T) {
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV2InstanceBasic(instanceName),
			},

			{
				ResourceName:      dcsV2InstanceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"backup_policy", "parameters", "password",
					"bandwidth_info.0.current_time",
				},
			},
		},
	})
}

func testAccCheckDcsV2InstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DcsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating DCSv2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dcs_instance_v2" {
			continue
		}

		_, err := instance.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("DCS instance still exists")
		}
	}
	return nil
}

func testAccCheckDcsV2InstanceExists(n string, dcsInstance instance.DcsInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		dcsClient, err := config.DcsV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating DCSv2 client: %w", err)
		}

		v, err := instance.Get(dcsClient, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting instance (%s): %w", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("DCS instance not found")
		}
		dcsInstance = *v
		return nil
	}
}

var testBase = fmt.Sprintf(`
%s

%s

data "opentelekomcloud_compute_availability_zones_v2" "zones" {}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet)

func testAccDcsV2InstanceBasic(instanceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"

  backup_policy {
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  parameters {
    id    = "1"
    name  = "timeout"
    value = "100"
  }

  tags = {
    environment = "basic"
    managed_by  = "terraform"
  }
}
`, testBase, instanceName)
}
func testAccDcsV2InstanceUpdated(instanceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  backup_policy {
    backup_type = "manual"
    begin_at    = "01:00-02:00"
    period_type = "weekly"
    backup_at   = [1, 2, 4]
    save_days   = 2
  }

  parameters {
    id    = "1"
    name  = "timeout"
    value = "200"
  }

  tags = {
    environment = "update"
    managed_by  = "terraform"
    user        = "admin"
  }
}
`, testBase, instanceName)
}

func testAccDcsV2InstanceSingle(instanceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "4.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "redis.single.xu1.tiny.128"
}
`, testBase, instanceName)
}

func testAccDcsV2InstanceEngineV3(instanceName string) string {
	return fmt.Sprintf(`


%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "3.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id  = opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "dcs.master_standby"
}
`, testBase, instanceName)
}

func testAccDcsV2InstanceWhitelist(instanceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  backup_policy {
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  parameters {
    id    = "1"
    name  = "timeout"
    value = "100"
  }

  enable_whitelist = false
  whitelist {
    group_name = "test-group-name"
    ip_list    = ["10.10.10.1", "10.10.10.2"]
  }
}
`, testBase, instanceName)
}

func testAccDcsV2InstanceWhitelistUpdate(instanceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  backup_policy {
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }
}
`, testBase, instanceName)
}

func testAccDcsV2InstanceWhitelistSecondUpdate(instanceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  backup_policy {
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  parameters {
    id    = "1"
    name  = "timeout"
    value = "100"
  }

  enable_whitelist = true
  whitelist {
    group_name = "test-group-name"
    ip_list    = ["10.10.10.1", "10.10.10.2"]
  }
  whitelist {
    group_name = "test-group-name-2"
    ip_list    = ["10.10.10.11", "10.10.10.3", "10.10.10.4"]
  }
}
`, testBase, instanceName)
}

func testAccDcsV2InstanceWhitelistThirdUpdate(instanceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  backup_policy {
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  parameters {
    id    = "1"
    name  = "timeout"
    value = "100"
  }

  enable_whitelist = false
  whitelist {
    group_name = "test-group-name-2"
    ip_list    = ["10.10.10.11", "10.10.10.3", "10.10.10.4"]
  }
}
`, testBase, instanceName)
}

func testAccDcsV2InstancePrivateIPs(instanceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  private_ip         = cidrhost(data.opentelekomcloud_vpc_subnet_v1.shared_subnet.cidr, 6)
  flavor             = "redis.ha.xu1.tiny.r2.128"

  backup_policy {
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  parameters {
    id    = "1"
    name  = "timeout"
    value = "100"
  }
}
`, testBase, instanceName)
}

func testAccDcsV2InstanceSSL(instanceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "6.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  ssl_enable         = true

}
`, testBase, instanceName)
}

func testAccDcsV2InstanceSSLUpdate(instanceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
  name               = "%s"
  engine_version     = "6.0"
  password           = "Hungarian_rapsody"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  ssl_enable         = false

}
`, testBase, instanceName)
}
