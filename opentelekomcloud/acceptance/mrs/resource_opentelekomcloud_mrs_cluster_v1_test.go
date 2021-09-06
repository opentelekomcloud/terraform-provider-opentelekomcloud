package acceptance

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/mrs/v1/cluster"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccMRSV1Cluster_basic(t *testing.T) {
	var clusterGet cluster.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
			common.TestAccAzPreCheck(t)
			common.TestAccKeyPairPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMRSV1ClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV1ClusterExists("opentelekomcloud_mrs_cluster_v1.cluster1", &clusterGet),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_mrs_cluster_v1.cluster1", "cluster_state", "running"),
				),
			},
		},
	})
}

func testAccCheckMRSV1ClusterDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	mrsClient, err := config.MrsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_mrs_cluster_v1" {
			continue
		}

		clusterGet, err := cluster.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("cluster still exists. err : %s", err)
		}
		if clusterGet.Clusterstate == "terminated" {
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
		mrsClient, err := config.MrsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud mrs client: %s ", err)
		}

		found, err := cluster.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Clusterid != rs.Primary.ID {
			return fmt.Errorf("cluster not found. ")
		}

		*clusterGet = *found
		time.Sleep(5 * time.Second)

		return nil
	}
}

var testAccMRSV1ClusterConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_mrs_cluster_v1" "cluster1" {
  cluster_name             = "mrs-cluster-acc"
  billing_type             = 12
  master_node_num          = 1
  core_node_num            = 1
  master_node_size         = "c3.2xlarge.4.linux.mrs"
  core_node_size           = "c3.2xlarge.4.linux.mrs"
  available_zone_id        = "%s"
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  cluster_version          = "MRS 1.7.2"
  master_data_volume_type  = "SAS"
  master_data_volume_size  = 100
  master_data_volume_count = 1
  core_data_volume_type    = "SATA"
  core_data_volume_size    = 100
  core_data_volume_count   = 2
  safe_mode                = 0
  cluster_type             = 0
  node_public_cert_name    = "%s"
  cluster_admin_secret     = "Qwert!123"
  component_list {
    component_name = "Hadoop"
  }
  component_list {
    component_name = "Spark"
  }
  component_list {
    component_name = "Hive"
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
