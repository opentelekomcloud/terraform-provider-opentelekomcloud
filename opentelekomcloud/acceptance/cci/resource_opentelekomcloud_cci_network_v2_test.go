package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/network"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getV2NetworkResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CciV2NetworkClient(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CCI v2 client: %s", err)
	}
	return network.Get(client, state.Primary.Attributes["namespace"], state.Primary.Attributes["name"])
}

func TestAccV2Network_basic(t *testing.T) {
	var net network.Network
	rName := fmt.Sprintf("cci-net-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_cci_network_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&net,
		getV2NetworkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Network_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "namespace", rName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckOutput("is_warm_pool_size_pass", "true"),
					resource.TestCheckOutput("is_warm_pool_recycle_interval_pass", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "subnets.0.subnet_id",
						"opentelekomcloud_vpc_subnet_v1.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_ids.0",
						"opentelekomcloud_networking_secgroup_v2.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "api_version"),
					resource.TestCheckResourceAttrSet(resourceName, "kind"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "finalizers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
					resource.TestCheckResourceAttrSet(resourceName, "status.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "status.0.conditions.#"),
					resource.TestCheckResourceAttrSet(resourceName, "status.0.subnet_attrs.#"),
					resource.TestCheckResourceAttrPair(resourceName, "status.0.subnet_attrs.0.network_id",
						"opentelekomcloud_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "status.0.subnet_attrs.0.subnet_v4_id",
						"opentelekomcloud_vpc_subnet_v1.test", "subnet_id"),
				),
			},
			{
				Config: testAccV2Network_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckOutput("is_warm_pool_size_pass", "true"),
					resource.TestCheckOutput("is_warm_pool_recycle_interval_pass", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_ids.0",
						"opentelekomcloud_networking_secgroup_v2.test1", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV2NetworkImportStateFunc(resourceName),
			},
		},
	})
}

func testAccV2NetworkImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["namespace"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("the namespace (%s) or name(%s) or ID (%s) is nil",
				rs.Primary.Attributes["namespace"], rs.Primary.Attributes["name"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["namespace"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccV2Network_base(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "test" {
  vpc_id = opentelekomcloud_vpc_v1.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.test.cidr, 4, 1), 1)
}

resource "opentelekomcloud_networking_secgroup_v2" "test" {
  name = "%[1]s"
}
`, rName)
}

func testAccV2Network_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cci_namespace_v2" "test" {
  name = "%[2]s"
}

resource "opentelekomcloud_cci_network_v2" "test" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"

  annotations = {
    "yangtse.io/project-id"                 = opentelekomcloud_cci_namespace_v2.test.annotations["tenant.kubernetes.io/project-id"]
    "yangtse.io/domain-id"                  = opentelekomcloud_cci_namespace_v2.test.annotations["tenant.kubernetes.io/domain-id"]
    "yangtse.io/warm-pool-size"             = "10"
    "yangtse.io/warm-pool-recycle-interval" = "2"
  }

  subnets {
    subnet_id = opentelekomcloud_vpc_subnet_v1.test.subnet_id
  }

  security_group_ids = [opentelekomcloud_networking_secgroup_v2.test.id]
}

output "is_warm_pool_size_pass" {
  value = opentelekomcloud_cci_network_v2.test.annotations["yangtse.io/warm-pool-size"] == "10"
}

output "is_warm_pool_recycle_interval_pass" {
  value = opentelekomcloud_cci_network_v2.test.annotations["yangtse.io/warm-pool-recycle-interval"] == "2"
}
`, testAccV2Network_base(rName), rName)
}

func TestAccV2Network_multipleSecurityGroups(t *testing.T) {
	var net network.Network
	rName := fmt.Sprintf("cci-net-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_cci_network_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&net,
		getV2NetworkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Network_multipleSecurityGroups(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV2NetworkImportStateFunc(resourceName),
			},
		},
	})
}

func testAccV2Network_multipleSecurityGroups(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_networking_secgroup_v2" "sg1" {
  name = "%[2]s-sg1"
}

resource "opentelekomcloud_networking_secgroup_v2" "sg2" {
  name = "%[2]s-sg2"
}

resource "opentelekomcloud_networking_secgroup_v2" "sg3" {
  name = "%[2]s-sg3"
}

resource "opentelekomcloud_cci_namespace_v2" "test" {
  name = "%[2]s"
}

resource "opentelekomcloud_cci_network_v2" "test" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"

  annotations = {
    "yangtse.io/project-id"                 = opentelekomcloud_cci_namespace_v2.test.annotations["tenant.kubernetes.io/project-id"]
    "yangtse.io/domain-id"                  = opentelekomcloud_cci_namespace_v2.test.annotations["tenant.kubernetes.io/domain-id"]
    "yangtse.io/warm-pool-size"             = "10"
    "yangtse.io/warm-pool-recycle-interval" = "2"
  }

  subnets {
    subnet_id = opentelekomcloud_vpc_subnet_v1.test.subnet_id
  }

  security_group_ids = [
    opentelekomcloud_networking_secgroup_v2.sg1.id,
    opentelekomcloud_networking_secgroup_v2.sg2.id,
    opentelekomcloud_networking_secgroup_v2.sg3.id,
  ]
}
`, testAccV2Network_base(rName), rName)
}

func testAccV2Network_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_networking_secgroup_v2" "test1" {
  name = "%[2]s_update"
}

resource "opentelekomcloud_cci_namespace_v2" "test" {
  name = "%[2]s"
}

resource "opentelekomcloud_cci_network_v2" "test" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"

  annotations = {
    "yangtse.io/project-id"                 = opentelekomcloud_cci_namespace_v2.test.annotations["tenant.kubernetes.io/project-id"]
    "yangtse.io/domain-id"                  = opentelekomcloud_cci_namespace_v2.test.annotations["tenant.kubernetes.io/domain-id"]
    "yangtse.io/warm-pool-size"             = "8"
    "yangtse.io/warm-pool-recycle-interval" = "3"
  }

  subnets {
    subnet_id = opentelekomcloud_vpc_subnet_v1.test.subnet_id
  }

  security_group_ids = [opentelekomcloud_networking_secgroup_v2.test1.id]
}

output "is_warm_pool_size_pass" {
  value = opentelekomcloud_cci_network_v2.test.annotations["yangtse.io/warm-pool-size"] == "8"
}

output "is_warm_pool_recycle_interval_pass" {
  value = opentelekomcloud_cci_network_v2.test.annotations["yangtse.io/warm-pool-recycle-interval"] == "3"
}
`, testAccV2Network_base(rName), rName)
}
