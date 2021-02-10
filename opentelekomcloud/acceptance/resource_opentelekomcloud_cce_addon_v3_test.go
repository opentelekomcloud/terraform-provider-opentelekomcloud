package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/addons"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccCCEAddonV3_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3_basic,
			},
		},
	})
}

func TestAccCCEAddonV3_emptyBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3_emptyBasic,
			},
		},
	})
}

func testAccCheckCCEAddonV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*cfg.Config)
	cceClient, err := config.CceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud CCE client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cce_cluster_v3" {
			continue
		}

		_, err := addons.Get(cceClient, rs.Primary.ID, rs.Primary.Attributes["cluster_id"]).Extract()
		if err == nil {
			return fmt.Errorf("cluster still exists")
		}
	}

	return nil
}

var (
	testAccCCEAddonV3_basic = fmt.Sprintf(`
resource opentelekomcloud_cce_cluster_v3 cluster_1 {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = "%s"
  subnet_id               = "%s"
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}

resource opentelekomcloud_cce_addon_v3 addon {
  template_name    = "metrics-server"
  template_version = "1.0.3"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster_1.id

  values {
    basic = {
      euleros_version = "2.5"
      rbac_enabled    = true
      swr_addr        = "100.125.7.25:20202"
      swr_user        = "hwofficial"
    }
  }
}
`, clusterName, OS_VPC_ID, OS_NETWORK_ID)

	testAccCCEAddonV3_emptyBasic = fmt.Sprintf(`
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
  template_version = "1.17.2"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster.id
  values {
    basic = {}
  }
}
`, clusterName, OS_VPC_ID, OS_NETWORK_ID)
)
