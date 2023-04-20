package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodepools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/cce/shared"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	ecs "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/ecs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const nodePoolResourceName = "opentelekomcloud_cce_node_pool_v3.node_pool"

func TestAccCCENodePoolsV3_basic(t *testing.T) {
	var nodePool nodepools.NodePool

	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.Server, Count: 2},
		{Q: quotas.Volume, Count: 2},
		{Q: quotas.VolumeSize, Count: 40 + 100},
	}
	qts = append(qts, ecs.QuotasForFlavor("s2.large.2")...)
	quotas.BookMany(t, qts)
	shared.BookCluster(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3Exists(nodePoolResourceName, shared.DataSourceClusterName, &nodePool),
					resource.TestCheckResourceAttr(nodePoolResourceName, "name", "opentelekomcloud-cce-node-pool"),
					resource.TestCheckResourceAttr(nodePoolResourceName, "flavor", "s2.large.2"),
					resource.TestCheckResourceAttr(nodePoolResourceName, "os", "EulerOS 2.5"),
					resource.TestCheckResourceAttr(nodePoolResourceName, "k8s_tags.kubelet.kubernetes.io/namespace", "muh"),
					resource.TestCheckResourceAttr(nodePoolResourceName, "data_volumes.0.extend_params.useType", "docker"),
				),
			},
			{
				Config: testAccCCENodePoolV3Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(nodePoolResourceName, "initial_node_count", "2"),
					resource.TestCheckResourceAttr(nodePoolResourceName, "k8s_tags.kubelet.kubernetes.io/namespace", "kuh"),
					resource.TestCheckResourceAttr(nodePoolResourceName, "data_volumes.0.extend_params.useType", "docker"),
				),
			},
		},
	})
}

func TestAccCCENodePoolV3ImportBasic(t *testing.T) {
	t.Parallel()
	shared.BookCluster(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3Basic,
			},
			{
				ResourceName:      nodePoolResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENodePoolV3ImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"max_node_count", "min_node_count", "priority", "scale_down_cooldown_time", "initial_node_count",
				},
			},
		},
	})
}

func testAccCCENodePoolV3ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterID string
		var nodePoolID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_cce_cluster_v3" {
				clusterID = rs.Primary.ID
			} else if rs.Type == "opentelekomcloud_cce_node_pool_v3" {
				nodePoolID = rs.Primary.ID
			}
		}
		if clusterID == "" || nodePoolID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", clusterID, nodePoolID)
		}
		return fmt.Sprintf("%s/%s", clusterID, nodePoolID), nil
	}
}

func TestAccCCENodePoolsV3_randomAZ(t *testing.T) {
	var nodePool nodepools.NodePool

	t.Parallel()
	shared.BookCluster(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3RandomAZ,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3Exists(nodePoolResourceName, shared.DataSourceClusterName, &nodePool),
					resource.TestCheckResourceAttr(nodePoolResourceName, "availability_zone", "random"),
				),
			},
		},
	})
}

func TestAccCCENodePoolsV3EncryptedVolume(t *testing.T) {
	var nodePool nodepools.NodePool

	t.Parallel()
	shared.BookCluster(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3Encrypted,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3Exists(nodePoolResourceName, shared.DataSourceClusterName, &nodePool),
					resource.TestCheckResourceAttr(nodePoolResourceName, "data_volumes.0.kms_id", env.OS_KMS_ID),
					resource.TestCheckResourceAttr(nodePoolResourceName, "root_volume.0.kms_id", env.OS_KMS_ID),
				),
			},
		},
	})
}

func TestAccCCENodePoolsV3ExtendParams(t *testing.T) {
	var nodePool nodepools.NodePool

	t.Parallel()
	shared.BookCluster(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3ExtendParams,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3Exists(nodePoolResourceName, shared.DataSourceClusterName, &nodePool),
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
			return fmt.Errorf("cluster ID is not set")
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

var testAccCCENodePoolV3Basic = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_node_pool_v3" "node_pool" {
  cluster_id         = data.opentelekomcloud_cce_cluster_v3.cluster.id
  name               = "opentelekomcloud-cce-node-pool"
  os                 = "EulerOS 2.9"
  flavor             = "s2.large.2"
  initial_node_count = 1
  availability_zone  = "%s"
  key_pair           = "%s"
  runtime            = "containerd"

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
    extend_params = {
      "useType" = "docker"
    }
  }

  k8s_tags = {
    "kubelet.kubernetes.io/namespace" = "muh"
  }
}`, shared.DataSourceCluster, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

var testAccCCENodePoolV3Update = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_node_pool_v3" "node_pool" {
  cluster_id         = data.opentelekomcloud_cce_cluster_v3.cluster.id
  name               = "opentelekomcloud-cce-node-pool"
  os                 = "EulerOS 2.5"
  flavor             = "s2.large.2"
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
    extend_params = {
      "useType" = "docker"
    }
  }

  k8s_tags = {
    "kubelet.kubernetes.io/namespace" = "kuh"
  }
}`, shared.DataSourceCluster, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

var testAccCCENodePoolV3RandomAZ = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_node_pool_v3" "node_pool" {
  cluster_id         = data.opentelekomcloud_cce_cluster_v3.cluster.id
  name               = "opentelekomcloud-cce-node-pool"
  os                 = "EulerOS 2.5"
  flavor             = "s3.medium.1"
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
}`, shared.DataSourceCluster, env.OS_KEYPAIR_NAME)

var testAccCCENodePoolV3Encrypted = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_node_pool_v3" "node_pool" {
  cluster_id         = data.opentelekomcloud_cce_cluster_v3.cluster.id
  name               = "opentelekomcloud-cce-node-pool"
  os                 = "EulerOS 2.5"
  flavor             = "s3.medium.1"
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
    kms_id     = "%s"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_id     = "%s"
  }
}`, shared.DataSourceCluster, env.OS_KEYPAIR_NAME, env.OS_KMS_ID, env.OS_KMS_ID)

var testAccCCENodePoolV3ExtendParams = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_node_pool_v3" "node_pool" {
  cluster_id         = data.opentelekomcloud_cce_cluster_v3.cluster.id
  name               = "opentelekomcloud-cce-node-pool"
  os                 = "EulerOS 2.5"
  flavor             = "s3.medium.1"
  initial_node_count = 1
  key_pair           = "%s"
  availability_zone  = "random"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  max_pods         = 16
  docker_base_size = 32
}`, shared.DataSourceCluster, env.OS_KEYPAIR_NAME)
