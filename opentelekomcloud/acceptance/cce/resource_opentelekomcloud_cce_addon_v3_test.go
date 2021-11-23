package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/addons"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceAddonName = "opentelekomcloud_cce_addon_v3.autoscaler"

func TestAccCCEAddonV3Basic(t *testing.T) {
	clusterName := randClusterName()
	t.Parallel()
	quotas.BookOne(t, quotas.CCEClusterQuota)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3Basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkScaleDownForAutoscaler(resourceAddonName, true),
					resource.TestCheckResourceAttr(resourceAddonName, "values.0.custom.scaleDownDelayAfterDelete", "11"),
				),
			},
			{
				Config: testAccCCEAddonV3Updated(clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkScaleDownForAutoscaler(resourceAddonName, false),
					resource.TestCheckResourceAttr(resourceAddonName, "values.0.custom.scaleDownDelayAfterDelete", "8"),
				),
			},
		},
	})
}

func TestAccCCEAddonV3ImportBasic(t *testing.T) {
	t.Parallel()
	clusterName := randClusterName()
	quotas.BookOne(t, quotas.CCEClusterQuota)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3Basic(clusterName),
			},
			{
				ResourceName:      resourceAddonName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCEAddonV3ImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"values",
				},
			},
		},
	})
}

func testAccCCEAddonV3ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterID string
		var addonID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_cce_cluster_v3" {
				clusterID = rs.Primary.ID
			} else if rs.Type == "opentelekomcloud_cce_addon_v3" {
				addonID = rs.Primary.ID
			}
		}
		if clusterID == "" || addonID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", clusterID, addonID)
		}
		return fmt.Sprintf("%s/%s", clusterID, addonID), nil
	}
}

func TestAccCCEAddonV3ForceNewCCE(t *testing.T) {
	clusterName := randClusterName()
	t.Parallel()
	quotas.BookOne(t, quotas.CCEClusterQuota)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3Basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkScaleDownForAutoscaler(resourceAddonName, true),
					resource.TestCheckResourceAttr(resourceAddonName, "values.0.custom.scaleDownDelayAfterDelete", "11"),
				),
			},
			{
				Config: testAccCCEAddonV3ForceNew(clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkScaleDownForAutoscaler(resourceAddonName, true),
					resource.TestCheckResourceAttr(resourceAddonName, "values.0.custom.scaleDownDelayAfterDelete", "11"),
				),
			},
		},
	})
}

func testAccCheckCCEAddonV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.CceV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud CCEv3 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cce_addon_v3" {
			continue
		}

		_, err := addons.Get(client, rs.Primary.ID, rs.Primary.Attributes["cluster_id"]).Extract()
		if err == nil {
			return fmt.Errorf("addon still exists")
		}
	}

	return nil
}

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

func testAccCCEAddonV3Basic(cName string) string {
	return fmt.Sprintf(`
%s
%s

resource opentelekomcloud_cce_cluster_v3 cluster_1 {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}

resource "opentelekomcloud_cce_addon_v3" "autoscaler" {
  template_name    = "autoscaler"
  template_version = "1.19.1"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster_1.id

  values {
    basic = {
      "cceEndpoint" : "https://cce.eu-de.otc.t-systems.com",
      "ecsEndpoint" : "https://ecs.eu-de.otc.t-systems.com",
      "image_version" : "1.19.1",
      "platform" : "linux-amd64",
      "region" : "eu-de",
      "swr_addr" : "100.125.7.25:20202",
      "swr_user" : "hwofficial"
    }
    custom = {
      "cluster_id" : opentelekomcloud_cce_cluster_v3.cluster_1.id,
      "coresTotal" : 32000,
      "expander" : "priority",
      "logLevel" : 4,
      "maxEmptyBulkDeleteFlag" : 10,
      "maxNodeProvisionTime" : 15,
      "maxNodesTotal" : 1000,
      "memoryTotal" : 128000,
      "scaleDownDelayAfterAdd" : 10,
      "scaleDownDelayAfterDelete" : 11,
      "scaleDownDelayAfterFailure" : 3,
      "scaleDownEnabled" : true,
      "scaleDownUnneededTime" : 10,
      "scaleDownUtilizationThreshold" : 0.5,
      "scaleUpCpuUtilizationThreshold" : 1,
      "scaleUpMemUtilizationThreshold" : 1,
      "scaleUpUnscheduledPodEnabled" : true,
      "scaleUpUtilizationEnabled" : true,
      "tenant_id" : data.opentelekomcloud_identity_project_v3.project.id,
      "unremovableNodeRecheckTimeout" : 5
    }
  }
}
`, common.DataSourceSubnet, common.DataSourceProject, cName)
}

func testAccCCEAddonV3Updated(cName string) string {
	return fmt.Sprintf(`
%s
%s

resource opentelekomcloud_cce_cluster_v3 cluster_1 {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}

resource "opentelekomcloud_cce_addon_v3" "autoscaler" {
  template_name    = "autoscaler"
  template_version = "1.19.1"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster_1.id

  values {
    basic = {
      "cceEndpoint" : "https://cce.eu-de.otc.t-systems.com",
      "ecsEndpoint" : "https://ecs.eu-de.otc.t-systems.com",
      "image_version" : "1.19.1",
      "platform" : "linux-amd64",
      "region" : "eu-de",
      "swr_addr" : "100.125.7.25:20202",
      "swr_user" : "hwofficial"
    }
    custom = {
      "cluster_id" : opentelekomcloud_cce_cluster_v3.cluster_1.id,
      "coresTotal" : 32000,
      "expander" : "priority",
      "logLevel" : 4,
      "maxEmptyBulkDeleteFlag" : 10,
      "maxNodeProvisionTime" : 15,
      "maxNodesTotal" : 1000,
      "memoryTotal" : 128000,
      "scaleDownDelayAfterAdd" : 10,
      "scaleDownDelayAfterDelete" : 8,
      "scaleDownDelayAfterFailure" : 3,
      "scaleDownEnabled" : false,
      "scaleDownUnneededTime" : 10,
      "scaleDownUtilizationThreshold" : 0.5,
      "scaleUpCpuUtilizationThreshold" : 1,
      "scaleUpMemUtilizationThreshold" : 1,
      "scaleUpUnscheduledPodEnabled" : true,
      "scaleUpUtilizationEnabled" : true,
      "tenant_id" : data.opentelekomcloud_identity_project_v3.project.id,
      "unremovableNodeRecheckTimeout" : 5
    }
  }
}
`, common.DataSourceSubnet, common.DataSourceProject, cName)
}

func testAccCCEAddonV3ForceNew(cName string) string {
	return fmt.Sprintf(`
%s
%s

resource opentelekomcloud_cce_cluster_v3 cluster_1 {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.medium"
  vpc_id                  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}

resource "opentelekomcloud_cce_addon_v3" "autoscaler" {
  template_name    = "autoscaler"
  template_version = "1.19.1"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster_1.id

  values {
    basic = {
      "cceEndpoint" : "https://cce.eu-de.otc.t-systems.com",
      "ecsEndpoint" : "https://ecs.eu-de.otc.t-systems.com",
      "image_version" : "1.19.1",
      "platform" : "linux-amd64",
      "region" : "eu-de",
      "swr_addr" : "100.125.7.25:20202",
      "swr_user" : "hwofficial"
    }
    custom = {
      "coresTotal" : 32000,
      "maxEmptyBulkDeleteFlag" : 10,
      "maxNodesTotal" : 1000,
      "memoryTotal" : 128000,
      "scaleDownDelayAfterAdd" : 11,
      "scaleDownDelayAfterDelete" : 11,
      "scaleDownDelayAfterFailure" : 3,
      "scaleDownEnabled" : true,
      "scaleDownUnneededTime" : 10,
      "scaleDownUtilizationThreshold" : 0.25,
      "scaleUpCpuUtilizationThreshold" : 0.8,
      "scaleUpMemUtilizationThreshold" : 0.8,
      "scaleUpUnscheduledPodEnabled" : true,
      "scaleUpUtilizationEnabled" : true,
      "unremovableNodeRecheckTimeout" : 5,
      "tenant_id" : data.opentelekomcloud_identity_project_v3.project.id,
    }
  }
}
`, common.DataSourceSubnet, common.DataSourceProject, cName)
}
