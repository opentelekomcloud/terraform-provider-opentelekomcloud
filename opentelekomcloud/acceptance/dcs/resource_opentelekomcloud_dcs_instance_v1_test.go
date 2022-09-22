package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v1/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceInstanceName = "opentelekomcloud_dcs_instance_v1.instance_1"

func TestAccDcsInstancesV1_basic(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1InstanceBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceInstanceName, instance),
					resource.TestCheckResourceAttr(resourceInstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceInstanceName, "engine", "Redis"),
				),
			},
			{
				Config: testAccDcsV1InstanceUpdated(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceInstanceName, "backup_policy.0.begin_at", "01:00-02:00"),
					resource.TestCheckResourceAttr(resourceInstanceName, "backup_policy.0.save_days", "2"),
					resource.TestCheckResourceAttr(resourceInstanceName, "backup_policy.0.backup_at.#", "3"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV1_basicSingleInstance(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1InstanceSingle(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceInstanceName, instance),
					resource.TestCheckResourceAttr(resourceInstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceInstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceInstanceName, "resource_spec_code", "redis.single.xu1.tiny.128"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV1_basicEngineV3Instace(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1InstanceEngineV3(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceInstanceName, instance),
					resource.TestCheckResourceAttr(resourceInstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceInstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceInstanceName, "resource_spec_code", "dcs.master_standby"),
				),
			},
		},
	})
}

func TestAccASV1Configuration_invalidSecurityGroup(t *testing.T) {
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDcsV1InstanceValidation(instanceName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`'security_group_id' should be set.+`),
			},
		},
	})
}

func TestAccASV1Configuration_WhitelistValidation(t *testing.T) {
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDcsV1InstanceWhitelistValidation(instanceName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`DCS Redis 3.0 instance does not support.+`),
			},
		},
	})
}

func TestAccDcsInstancesV1_Whitelist(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1InstanceWhitelist(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceInstanceName, instance),
					resource.TestCheckResourceAttr(resourceInstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceInstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceInstanceName, "enable_whitelist", "false"),
					resource.TestCheckResourceAttr(resourceInstanceName, "whitelist.0.group_name", "test-group-name"),
				),
			},
			{
				Config: testAccDcsV1InstanceWhitelistUpdate(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceInstanceName, instance),
					resource.TestCheckResourceAttr(resourceInstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceInstanceName, "engine", "Redis"),
				),
			},
			{
				Config: testAccDcsV1InstanceWhitelistSecondUpdate(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceInstanceName, instance),
					resource.TestCheckResourceAttr(resourceInstanceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceInstanceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceInstanceName, "enable_whitelist", "true"),
				),
			},
		},
	})
}

func testAccCheckDcsV1InstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DcsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating DCSv1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dcs_instance_v1" {
			continue
		}

		_, err := instances.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("DCS instance still exists")
		}
	}
	return nil
}

func testAccCheckDcsV1InstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		dcsClient, err := config.DcsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating DCSv1 client: %w", err)
		}

		v, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("error getting instance (%s): %w", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("DCS instance not found")
		}
		instance = *v
		return nil
	}
}

func testAccDcsV1InstanceBasic(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
  code = "%s"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "redis.ha.xu1.tiny.r2.128"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name            = "%s"
  engine_version  = "5.0"
  password        = "Hungarian_rapsody"
  engine          = "Redis"
  capacity        = 0.125
  vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id      = data.opentelekomcloud_dcs_product_v1.product_1.id
  backup_policy {
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  configuration {
    parameter_id    = "1"
    parameter_name  = "timeout"
    parameter_value = "100"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName)
}
func testAccDcsV1InstanceUpdated(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
  code = "%s"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "redis.ha.xu1.tiny.r2.128"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name            = "%s"
  engine_version  = "5.0"
  password        = "Hungarian_rapsody"
  engine          = "Redis"
  capacity        = 0.125
  vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id      = data.opentelekomcloud_dcs_product_v1.product_1.id
  backup_policy {
    backup_type = "manual"
    begin_at    = "01:00-02:00"
    period_type = "weekly"
    backup_at   = [1, 2, 4]
    save_days   = 2
  }

  configuration {
    parameter_id    = "1"
    parameter_name  = "timeout"
    parameter_value = "200"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName)
}

func testAccDcsV1InstanceSingle(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
  code = "%s"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "redis.single.xu1.tiny.128"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name            = "%s"
  engine_version  = "4.0"
  password        = "Hungarian_rapsody"
  engine          = "Redis"
  capacity        = 0.125
  vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id      = data.opentelekomcloud_dcs_product_v1.product_1.id
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName)
}

func testAccDcsV1InstanceValidation(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
  code = "%s"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "redis.single.xu1.tiny.128"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name            = "%s"
  engine_version  = "3.0"
  password        = "Hungarian_rapsody"
  engine          = "Redis"
  capacity        = 0.125
  vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id      = data.opentelekomcloud_dcs_product_v1.product_1.id
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName)
}

func testAccDcsV1InstanceWhitelistValidation(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
  code = "%s"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "redis.single.xu1.tiny.128"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name            = "%s"
  engine_version  = "3.0"
  password        = "Hungarian_rapsody"
  engine          = "Redis"
  capacity        = 0.125
  vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id      = data.opentelekomcloud_dcs_product_v1.product_1.id

  enable_whitelist = true
  whitelist {
    group_name = "test-group-name"
    ip_list    = ["10.10.10.1", "10.10.10.2"]
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName)
}

func testAccDcsV1InstanceEngineV3(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
  code = "%s"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "dcs.master_standby"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name              = "%s"
  engine_version    = "3.0"
  password          = "Hungarian_rapsody"
  engine            = "Redis"
  capacity          = 2
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  available_zones   = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id        = data.opentelekomcloud_dcs_product_v1.product_1.id
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName)
}

func testAccDcsV1InstanceWhitelist(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
  code = "%s"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "redis.ha.xu1.tiny.r2.128"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name            = "%s"
  engine_version  = "5.0"
  password        = "Hungarian_rapsody"
  engine          = "Redis"
  capacity        = 0.125
  vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id      = data.opentelekomcloud_dcs_product_v1.product_1.id
  backup_policy {
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  configuration {
    parameter_id    = "1"
    parameter_name  = "timeout"
    parameter_value = "100"
  }

  enable_whitelist = false
  whitelist {
    group_name = "test-group-name"
    ip_list    = ["10.10.10.1", "10.10.10.2"]
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName)
}

func testAccDcsV1InstanceWhitelistUpdate(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
  code = "%s"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "redis.ha.xu1.tiny.r2.128"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name            = "%s"
  engine_version  = "5.0"
  password        = "Hungarian_rapsody"
  engine          = "Redis"
  capacity        = 0.125
  vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id      = data.opentelekomcloud_dcs_product_v1.product_1.id
  backup_policy {
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName)
}

func testAccDcsV1InstanceWhitelistSecondUpdate(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
  code = "%s"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "redis.ha.xu1.tiny.r2.128"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name            = "%s"
  engine_version  = "5.0"
  password        = "Hungarian_rapsody"
  engine          = "Redis"
  capacity        = 0.125
  vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id      = data.opentelekomcloud_dcs_product_v1.product_1.id
  backup_policy {
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  configuration {
    parameter_id    = "1"
    parameter_name  = "timeout"
    parameter_value = "100"
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
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName)
}
