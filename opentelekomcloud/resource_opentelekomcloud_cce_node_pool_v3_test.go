package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodepools"
)

func TestAccCCENodePoolsV3_basic(t *testing.T) {
	var nodePool nodepools.NodePool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCCEKeyPairPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodePoolV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3Exists("opentelekomcloud_cce_node_pool_v3.node_pool", "opentelekomcloud_cce_cluster_v3.cluster", &nodePool),
					resource.TestCheckResourceAttr("opentelekomcloud_cce_node_pool_v3.node_pool", "name", "opentelekomcloud-cce-node-pool"),
					resource.TestCheckResourceAttr("opentelekomcloud_cce_node_pool_v3.node_pool", "flavor_id", "s6.large.2"),
					resource.TestCheckResourceAttr("opentelekomcloud_cce_node_pool_v3.node_pool", "os", "EulerOS 2.5"),
				),
			},
			{
				Config: testAccCCENodePoolV3_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_cce_node_v3.node_pool", "initial_node_count", "2"),
				),
			},
		},
	})
}

func testAccCheckCCENodePoolV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud CCE client: %s", err)
	}

	var clusterId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "opentelekomcloud_cce_cluster_v3" {
			clusterId = rs.Primary.ID
		}

		if rs.Type != "opentelekomcloud_cce_node_pool_v3" {
			continue
		}

		_, err := nodepools.Get(cceClient, clusterId, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("node pool still exists")
		}
	}

	return nil
}

func testAccCheckCCENodePoolV3Exists(n string, cluster string, nodepool *nodepools.NodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		c, ok := s.RootModule().Resources[cluster]
		if !ok {
			return fmt.Errorf("cluster not found: %s", c)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		if c.Primary.ID == "" {
			return fmt.Errorf("cluster id is not set")
		}

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.cceV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud CCE client: %s", err)
		}

		found, err := nodepools.Get(cceClient, c.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("node pool not found")
		}

		*nodepool = *found

		return nil
	}
}

var testAccCCENodePoolV3_basic = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster" {
  name         = "opentelekomcloud-cce-np"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_pool_v3" "node_pool" {
  cluster_id         = opentelekomcloud_cce_cluster_v3.cluster.id
  name               = "opentelekomcloud-cce-node-pool"
  os                 = "EulerOS 2.5"
  flavor             = "s2.xlarge.2"
  initial_node_count = 1
  availability_zone  = "%s"
  key_pair           = "%s"

  scale_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}`, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE, OS_KEYPAIR_NAME)

var testAccCCENodePoolV3_update = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster" {
  name         = "opentelekomcloud-cce-np"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_pool_v3" "node_pool" {
  cluster_id         = opentelekomcloud_cce_cluster_v3.cluster.id
  name               = "opentelekomcloud-cce-node-pool"
  os                 = "EulerOS 2.5"
  flavor             = "s2.xlarge.2"
  initial_node_count = 2
  availability_zone  = "%s"
  key_pair           = "%s"

  scale_enable             = true
  min_node_count           = 2
  max_node_count           = 9
  scale_down_cooldown_time = 100
  priority                 = 1

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}`, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE, OS_KEYPAIR_NAME)
