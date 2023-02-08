package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/mrs/v1/cluster"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/mrs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceClusterName = "opentelekomcloud_mrs_cluster_v1.this"

func TestAccMRSV1Cluster_basic(t *testing.T) {
	var mrsCluster cluster.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMRSV1ClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV1ClusterExists(resourceClusterName, &mrsCluster),
					resource.TestCheckResourceAttr(resourceClusterName, "cluster_state", "running"),
				),
			},
		},
	})
}

func testAccCheckMRSV1ClusterDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.MrsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(mrs.ErrCreationClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_mrs_cluster_v1" {
			continue
		}

		clusterGet, err := cluster.Get(client, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("cluster still exists")
		}
		if clusterGet.ClusterState == "terminated" {
			return nil
		}
	}

	return nil
}

func testAccCheckMRSV1ClusterExists(n string, clusterGet *cluster.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set. ")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.MrsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(mrs.ErrCreationClient, err)
		}

		found, err := cluster.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ClusterId != rs.Primary.ID {
			return fmt.Errorf("cluster not found")
		}

		*clusterGet = *found

		return nil
	}
}

var testAccMRSV1ClusterConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_mrs_cluster_v1" "this" {
  cluster_name          = "mrs-cluster-acc"
  billing_type          = 12
  master_node_num       = 2
  core_node_num         = 3
  master_node_size      = "c3.xlarge.4.linux.mrs"
  core_node_size        = "c3.xlarge.4.linux.mrs"
  available_zone_id     = "%s"
  vpc_id                = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  cluster_version       = "MRS 2.1.0"
  volume_type           = "SATA"
  volume_size           = 100
  cluster_type          = 0
  safe_mode             = 1
  node_public_cert_name = "%s"
  cluster_admin_secret  = "SuperQwerty!123"
  component_list {
    component_name = "Presto"
  }
  component_list {
    component_name = "Hadoop"
  }
  component_list {
    component_name = "Spark"
  }
  component_list {
    component_name = "HBase"
  }
  component_list {
    component_name = "Hive"
  }
  component_list {
    component_name = "Hue"
  }
  component_list {
    component_name = "Loader"
  }
  component_list {
    component_name = "Tez"
  }
  component_list {
    component_name = "Flink"
  }

  bootstrap_scripts {
    name       = "Modify os config"
    uri        = "s3a://bootstrap/modify_os_config.sh"
    parameters = "param1 param2"
    nodes = [
      "master",
      "core",
      "task",
    ]
    active_master          = true
    before_component_start = true
    fail_action            = "continue"
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)
