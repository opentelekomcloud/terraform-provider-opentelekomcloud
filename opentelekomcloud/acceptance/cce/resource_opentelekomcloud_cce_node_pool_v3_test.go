package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodepools"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccCCENodePoolsV3_basic(t *testing.T) {
	var nodePool nodepools.NodePool
	nodePoolName := "opentelekomcloud_cce_node_pool_v3.node_pool"
	clusterName := "opentelekomcloud_cce_cluster_v3.cluster"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3Exists(nodePoolName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(nodePoolName, "name", "opentelekomcloud-cce-node-pool"),
					resource.TestCheckResourceAttr(nodePoolName, "flavor", "s2.xlarge.2"),
					resource.TestCheckResourceAttr(nodePoolName, "os", "EulerOS 2.5"),
					resource.TestCheckResourceAttr(nodePoolName, "k8s_tags.kubelet.kubernetes.io/namespace", "muh"),
				),
			},
			{
				Config: testAccCCENodePoolV3Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(nodePoolName, "initial_node_count", "2"),
					resource.TestCheckResourceAttr(nodePoolName, "k8s_tags.kubelet.kubernetes.io/namespace", "kuh"),
				),
			},
		},
	})
}

func TestAccCCENodePoolsV3_randomAZ(t *testing.T) {
	var nodePool nodepools.NodePool
	nodePoolName := "opentelekomcloud_cce_node_pool_v3.node_pool"
	clusterName := "opentelekomcloud_cce_cluster_v3.cluster"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3RandomAZ,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3Exists(nodePoolName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(nodePoolName, "availability_zone", "random"),
				),
			},
		},
	})
}

func TestAccCCENodePoolsV3EncryptedVolume(t *testing.T) {
	var nodePool nodepools.NodePool
	nodePoolName := "opentelekomcloud_cce_node_pool_v3.node_pool"
	clusterName := "opentelekomcloud_cce_cluster_v3.cluster"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3Encrypted,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3Exists(nodePoolName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(nodePoolName, "data_volumes.0.kms_id", env.OS_KMS_ID),
				),
			},
		},
	})
}

func testAccCheckCCENodePoolV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.CceV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud CCE client: %s", err)
	}

	var clusterID string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "opentelekomcloud_cce_cluster_v3" {
			clusterID = rs.Primary.ID
		}

		if rs.Type != "opentelekomcloud_cce_node_pool_v3" {
			continue
		}

		_, err := nodepools.Get(client, clusterID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("node pool still exists")
		}
	}

	return nil
}

func testAccCheckCCENodePoolV3Exists(n string, cluster string, nodePool *nodepools.NodePool) resource.TestCheckFunc {
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

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CceV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud CCE client: %s", err)
		}

		found, err := nodepools.Get(client, c.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("node pool not found")
		}

		*nodePool = *found

		return nil
	}
}

var (
	testAccCCENodePoolV3Basic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_cce_cluster_v3" "cluster" {
  name         = "opentelekomcloud-cce-np"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = data.opentelekomcloud_vpc_v1.shared_vpc.id
  subnet_id    = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id

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
  min_node_count           = 1
  max_node_count           = 3
  scale_down_cooldown_time = 6
  priority                 = 1

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  k8s_tags = {
    "kubelet.kubernetes.io/namespace" = "muh"
  }
}`, common.DataSourceVPC, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodePoolV3Update = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_cce_cluster_v3" "cluster" {
  name         = "opentelekomcloud-cce-np"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = data.opentelekomcloud_vpc_v1.shared_vpc.id
  subnet_id    = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id

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

  k8s_tags = {
    "kubelet.kubernetes.io/namespace" = "kuh"
  }
}`, common.DataSourceVPC, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodePoolV3RandomAZ = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_cce_cluster_v3" "cluster" {
  name         = "opentelekomcloud-cce-np"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = data.opentelekomcloud_vpc_v1.shared_vpc.id
  subnet_id    = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_pool_v3" "node_pool" {
  cluster_id         = opentelekomcloud_cce_cluster_v3.cluster.id
  name               = "opentelekomcloud-cce-node-pool"
  os                 = "EulerOS 2.5"
  flavor             = "s2.xlarge.2"
  initial_node_count = 1
  key_pair           = "%s"
  availability_zone  = "random"

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
}`, common.DataSourceVPC, common.DataSourceSubnet, env.OS_KEYPAIR_NAME)

	testAccCCENodePoolV3Encrypted = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_cce_cluster_v3" "cluster" {
  name         = "opentelekomcloud-cce-np"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = data.opentelekomcloud_vpc_v1.shared_vpc.id
  subnet_id    = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_pool_v3" "node_pool" {
  cluster_id         = opentelekomcloud_cce_cluster_v3.cluster.id
  name               = "opentelekomcloud-cce-node-pool"
  os                 = "EulerOS 2.5"
  flavor             = "s2.xlarge.2"
  initial_node_count = 1
  key_pair           = "%s"
  availability_zone  = "random"

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
    kms_id     = "%s"
  }
}`, common.DataSourceVPC, common.DataSourceSubnet, env.OS_KEYPAIR_NAME, env.OS_KMS_ID)
)
