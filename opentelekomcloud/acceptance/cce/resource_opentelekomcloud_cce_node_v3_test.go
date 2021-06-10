package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodes"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

var privateIP = "192.168.1.13"

const nodeName = "opentelekomcloud_cce_node_v3.node_1"

func TestAccCCENodesV3Basic(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
					resource.TestCheckResourceAttr(nodeName, "name", "test-node"),
					resource.TestCheckResourceAttr(nodeName, "flavor_id", "s2.xlarge.2"),
					resource.TestCheckResourceAttr(nodeName, "os", "EulerOS 2.5"),
					resource.TestCheckResourceAttr(nodeName, "private_ip", privateIP),
				),
			},
			{
				Config: testAccCCENodeV3Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(nodeName, "name", "test-node2"),
				),
			},
		},
	})
}

func TestAccCCENodesV3Timeout(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
				),
			},
		},
	})
}
func TestAccCCENodesV3OS(t *testing.T) {
	var node nodes.Nodes
	var node2 nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3OS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
					resource.TestCheckResourceAttr(nodeName, "os", "EulerOS 2.5"),
					testAccCheckCCENodeV3Exists("opentelekomcloud_cce_node_v3.node_2", "opentelekomcloud_cce_cluster_v3.cluster_1", &node2),
					resource.TestCheckResourceAttr("opentelekomcloud_cce_node_v3.node_2", "os", "CentOS 7.7"),
				),
			},
		},
	})
}

func TestAccCCENodesV3BandWidthResize(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3Ip,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
					resource.TestCheckResourceAttr(nodeName, "iptype", "5_bgp"),
					resource.TestCheckResourceAttr(nodeName, "sharetype", "PER"),
					resource.TestCheckResourceAttr(nodeName, "bandwidth_charge_mode", "traffic"),
					resource.TestCheckResourceAttr(nodeName, "bandwidth_size", "100"),
				),
			},
			{
				Config: testAccCCENodeV3BandWidthResize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
					resource.TestCheckResourceAttr(nodeName, "bandwidth_size", "10"),
				),
			},
		},
	})
}

// TODO: Need to be tested
func TestAccCCENodesV3_eipIds(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3IpIDs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
				),
			},
			{
				Config: testAccCCENodeV3IpIDsUnset,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
				),
			},
		},
	})
}

func TestAccCCENodesV3IpSetNull(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3Ip,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
					resource.TestCheckResourceAttr(nodeName, "iptype", "5_bgp"),
					resource.TestCheckResourceAttr(nodeName, "sharetype", "PER"),
					resource.TestCheckResourceAttr(nodeName, "bandwidth_charge_mode", "traffic"),
				),
			},
			{
				Config: testAccCCENodeV3IpUnset,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
				),
			},
		},
	})
}

func TestAccCCENodesV3IpCreate(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3IpUnset,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
				),
			},
			{
				Config: testAccCCENodeV3Ip,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
				),
			},
		},
	})
}

func TestAccCCENodesV3IpWithExtendedParameters(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3IpParams,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
					resource.TestCheckResourceAttr(nodeName, "iptype", "5_bgp"),
					resource.TestCheckResourceAttr(nodeName, "sharetype", "PER"),
					resource.TestCheckResourceAttr(nodeName, "bandwidth_charge_mode", "traffic"),
				),
			},
		},
	})
}

func TestAccCCENodesV3IpNulls(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3IpNull,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
				),
			},
		},
	})
}

func TestAccCCENodesV3EncryptedVolume(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3EncryptedVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(nodeName, "opentelekomcloud_cce_cluster_v3.cluster_1", &node),
				),
			},
		},
	})
}

func testAccCheckCCENodeV3Destroy(s *terraform.State) error {
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

		if rs.Type != "opentelekomcloud_cce_node_v3" {
			continue
		}

		_, err := nodes.Get(client, clusterID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("node still exists")
		}
	}

	return nil
}

func testAccCheckCCENodeV3Exists(n string, cluster string, node *nodes.Nodes) resource.TestCheckFunc {
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

		found, err := nodes.Get(client, c.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("node not found")
		}

		*node = *found

		return nil
	}
}

var (
	testAccCCENodeV3OS = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name       = "test-node"
  flavor_id  = "s2.large.2"
  os         = "EulerOS 2.5"

  availability_zone = "%s"
  key_pair          = "%s"

  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}

resource "opentelekomcloud_cce_node_v3" "node_2" {
  cluster_id = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name       = "test-node"
  flavor_id  = "s2.large.2"
  os         = "CentOS 7.7"

  availability_zone = "%[3]s"
  key_pair          = "%[4]s"

  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}

`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodeV3Basic = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name       = "test-node"
  flavor_id  = "s2.xlarge.2"

  availability_zone = "%s"
  key_pair          = "%s"

  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }

  private_ip = "%s"
}`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME, privateIP)

	testAccCCENodeV3Update = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name       = "test-node2"
  flavor_id  = "s2.xlarge.2"

  availability_zone = "%s"
  key_pair          ="%s"

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }

  private_ip = "%s"
}`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME, privateIP)

	testAccCCENodeV3Timeout = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name              = "test-node1"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "%s"
  key_pair          = "%s"

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
  timeouts {
    create = "10m"
    delete = "10m"
  }
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodeV3Ip = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name              = "cce-node-1"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "%s"
  key_pair          = "%s"
  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  bandwidth_size = 100

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodeV3BandWidthResize = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name              = "cce-node-1"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "%s"
  key_pair          = "%s"
  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  bandwidth_size = 10

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodeV3IpUnset = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name              = "cce-node-1"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "%s"
  key_pair          = "%s"
  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodeV3IpParams = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name              = "cce-node-1"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "%s"
  key_pair          = "%s"
  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  bandwidth_size = 100
  sharetype      = "PER"
  iptype         = "5_bgp"

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodeV3IpNull = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name              = "cce-node-1"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "%s"
  key_pair          = "%s"
  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  bandwidth_size = null
  sharetype      = null
  iptype         = null

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodeV3IpIDs = fmt.Sprintf(`
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}
resource "opentelekomcloud_networking_floatingip_v2" "fip_2" {}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce-ids"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name              = "cce-node-1"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "%s"
  key_pair          = "%s"
  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  eip_ids = [opentelekomcloud_networking_floatingip_v2.fip_1.id]

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodeV3IpIDsUnset = fmt.Sprintf(`
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}
resource "opentelekomcloud_networking_floatingip_v2" "fip_2" {}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce-ids"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  flavor_id         = "s2.xlarge.2"
  availability_zone = "%s"
  key_pair          = "%s"
  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  eip_ids = [opentelekomcloud_networking_floatingip_v2.fip_2.id]

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)

	testAccCCENodeV3EncryptedVolume = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name         = "opentelekomcloud-cce-ids"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"

  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  flavor_id         = "s2.xlarge.2"
  availability_zone = "%s"
  key_pair          = "%s"

  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  data_volumes {
    size       = 100
    volumetype = "SATA"
    kms_id     = "%s"
  }
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME, env.OS_KMS_ID)
)
