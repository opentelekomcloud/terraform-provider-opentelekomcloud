package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/addons"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccCCEAddonV3Basic(t *testing.T) {
	resName := "opentelekomcloud_cce_addon_v3.autoscaler"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3Basic,
				Check: resource.ComposeTestCheckFunc(
					checkScaleDownForAutoscaler(resName, true),
					resource.TestCheckResourceAttr(resName, "values.0.custom.scaleDownDelayAfterDelete", "11"),
				),
			},
			{
				Config: testAccCCEAddonV3Updated,
				Check: resource.ComposeTestCheckFunc(
					checkScaleDownForAutoscaler(resName, false),
					resource.TestCheckResourceAttr(resName, "values.0.custom.scaleDownDelayAfterDelete", "8"),
				),
			},
		},
	})
}

func TestAccCCEAddonV3ForceNewCCE(t *testing.T) {
	resName := "opentelekomcloud_cce_addon_v3.autoscaler"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3Basic,
				Check: resource.ComposeTestCheckFunc(
					checkScaleDownForAutoscaler(resName, true),
					resource.TestCheckResourceAttr(resName, "values.0.custom.scaleDownDelayAfterDelete", "11"),
				),
			},
			{
				Config: testAccCCEAddonV3ForceNew,
				Check: resource.ComposeTestCheckFunc(
					checkScaleDownForAutoscaler(resName, true),
					resource.TestCheckResourceAttr(resName, "values.0.custom.scaleDownDelayAfterDelete", "11"),
				),
			},
		},
	})
}

func TestAccCCEAddonV3EmptyBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3EmptyBasic,
			},
		},
	})
}

func testAccCheckCCEAddonV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	cceClient, err := config.CceV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud CCE client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cce_addon_v3" {
			continue
		}

		_, err := addons.Get(cceClient, rs.Primary.ID, rs.Primary.Attributes["cluster_id"]).Extract()
		if err == nil {
			return fmt.Errorf("addon still exists")
		}
	}

	return nil
}

var (
	testAccCCEAddonV3Basic = fmt.Sprintf(`
resource opentelekomcloud_cce_cluster_v3 cluster_1 {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = "%s"
  subnet_id               = "%s"
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}

resource "opentelekomcloud_cce_addon_v3" "autoscaler" {
  template_name    = "autoscaler"
  template_version = "1.19.1"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster_1.id

  values {
    basic  = {
      "cceEndpoint": "https://cce.eu-de.otc.t-systems.com",
      "ecsEndpoint": "https://ecs.eu-de.otc.t-systems.com",
      "euleros_version": "2.2.5",
      "region": "eu-de",
      "swr_addr": "100.125.7.25:20202",
      "swr_user": "hwofficial"
    }
    custom = {
      "coresTotal": 32000,
      "maxEmptyBulkDeleteFlag": 10,
      "maxNodesTotal": 1000,
      "memoryTotal": 128000,
      "scaleDownDelayAfterAdd": 11,
      "scaleDownDelayAfterDelete": 11,
      "scaleDownDelayAfterFailure": 3,
      "scaleDownEnabled": true,
      "scaleDownUnneededTime": 10,
      "scaleDownUtilizationThreshold": 0.25,
      "scaleUpCpuUtilizationThreshold": 0.8,
      "scaleUpMemUtilizationThreshold": 0.8,
      "scaleUpUnscheduledPodEnabled": true,
      "scaleUpUtilizationEnabled": true,
      "unremovableNodeRecheckTimeout": 5,
      "tenant_id": "%s"
    }
  }
}
`, clusterName, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_TENANT_ID)

	testAccCCEAddonV3Updated = fmt.Sprintf(`
resource opentelekomcloud_cce_cluster_v3 cluster_1 {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = "%s"
  subnet_id               = "%s"
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}

resource "opentelekomcloud_cce_addon_v3" "autoscaler" {
  template_name    = "autoscaler"
  template_version = "1.19.1"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster_1.id

  values {
    basic  = {
      "cceEndpoint": "https://cce.eu-de.otc.t-systems.com",
      "ecsEndpoint": "https://ecs.eu-de.otc.t-systems.com",
      "euleros_version": "2.2.5",
      "region": "eu-de",
      "swr_addr": "100.125.7.25:20202",
      "swr_user": "hwofficial"
    }
    custom = {
      "coresTotal": 32000,
      "maxEmptyBulkDeleteFlag": 10,
      "maxNodesTotal": 100,
      "memoryTotal": 128000,
      "scaleDownDelayAfterAdd": 9,
      "scaleDownDelayAfterDelete": 8,
      "scaleDownDelayAfterFailure": 4,
      "scaleDownEnabled": false,
      "scaleDownUnneededTime": 10,
      "scaleDownUtilizationThreshold": 0.25,
      "scaleUpCpuUtilizationThreshold": 0.8,
      "scaleUpMemUtilizationThreshold": 0.8,
      "scaleUpUnscheduledPodEnabled": true,
      "scaleUpUtilizationEnabled": true,
      "unremovableNodeRecheckTimeout": 5,
      "tenant_id": "%s"
    }
  }
}
`, clusterName, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_TENANT_ID)

	testAccCCEAddonV3EmptyBasic = fmt.Sprintf(`
resource opentelekomcloud_cce_cluster_v3 cluster {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = "%s"
  subnet_id               = "%s"
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}

resource "opentelekomcloud_cce_addon_v3" "cluster_autoscaler" {
  template_name    = "autoscaler"
  template_version = "1.19.1"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster.id
  values {
    basic  = {}
    custom = {}
  }
}
`, clusterName, env.OS_VPC_ID, env.OS_NETWORK_ID)

	testAccCCEAddonV3ForceNew = fmt.Sprintf(`
resource opentelekomcloud_cce_cluster_v3 cluster_1 {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.medium"
  vpc_id                  = "%s"
  subnet_id               = "%s"
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}

resource "opentelekomcloud_cce_addon_v3" "autoscaler" {
  template_name    = "autoscaler"
  template_version = "1.19.1"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster_1.id

  values {
    basic  = {
      "cceEndpoint": "https://cce.eu-de.otc.t-systems.com",
      "ecsEndpoint": "https://ecs.eu-de.otc.t-systems.com",
      "euleros_version": "2.2.5",
      "region": "eu-de",
      "swr_addr": "100.125.7.25:20202",
      "swr_user": "hwofficial"
    }
    custom = {
      "coresTotal": 32000,
      "maxEmptyBulkDeleteFlag": 10,
      "maxNodesTotal": 1000,
      "memoryTotal": 128000,
      "scaleDownDelayAfterAdd": 11,
      "scaleDownDelayAfterDelete": 11,
      "scaleDownDelayAfterFailure": 3,
      "scaleDownEnabled": true,
      "scaleDownUnneededTime": 10,
      "scaleDownUtilizationThreshold": 0.25,
      "scaleUpCpuUtilizationThreshold": 0.8,
      "scaleUpMemUtilizationThreshold": 0.8,
      "scaleUpUnscheduledPodEnabled": true,
      "scaleUpUtilizationEnabled": true,
      "unremovableNodeRecheckTimeout": 5,
      "tenant_id": "%s"
    }
  }
}
`, clusterName, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_TENANT_ID)
)

func checkScaleDownForAutoscaler(name string, enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CceV3AddonClient(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud CCE client: %w", err)
		}

		found, err := addons.Get(client, rs.Primary.ID, rs.Primary.Attributes["cluster_id"]).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("addon not found")
		}

		if actual := found.Spec.Values.Advanced["scaleDownEnabled"]; actual != enabled {
			return fmt.Errorf("invalid `scaleDownEnabled` value: expected %v, got %v", enabled, actual)
		}

		return nil
	}
}
