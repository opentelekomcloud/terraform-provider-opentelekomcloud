package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestResourceCSSSnapshotConfigurationV1_basic(t *testing.T) {
	if osAgency == "" {
		t.Skip("OS_AGENCY is required for the test")
	}

	name := fmt.Sprintf("css-%s", acctest.RandString(10))
	resourceName := "opentelekomcloud_css_snapshot_configuration_v1.config"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.TestAccPreCheckRequiredEnvVars(t)
			acc.TestAccSubnetPreCheck(t)
			acc.TestAccAzPreCheck(t)
		},
		ProviderFactories: acc.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCssClusterV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceCSSSnapshotConfigurationV1Basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "creation_policy.0.prefix", "snap"),
					resource.TestCheckResourceAttrSet(resourceName, "base_path"),
				),
			},
			{
				Config: testResourceCSSSnapshotConfigurationV1Updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "creation_policy.0.prefix", "snapshot"),
					resource.TestCheckResourceAttr(resourceName, "creation_policy.0.keepday", "2"),
				),
			},
		},
	})
}

var osAgency = os.Getenv("OS_AGENCY")

func testResourceCSSSnapshotConfigurationV1Basic(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_css_cluster_v1" "cluster" {
  expect_node_num = 1
  name            = "%s"
  node_config {
    flavor = "css.medium.8"
    network_info {
      security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
      network_id        = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
      vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
    }
    volume {
      volume_type = "COMMON"
      size        = 100
    }

    availability_zone = "%s"
  }

  enable_https     = true
  enable_authority = true
  admin_pass       = "QwertyUI!"
}

resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "tf-snap-testing"
  force_destroy = true
}

resource "opentelekomcloud_css_snapshot_configuration_v1" "config" {
  cluster_id = opentelekomcloud_css_cluster_v1.cluster.id
  configuration {
    bucket = opentelekomcloud_obs_bucket.bucket.bucket
    agency = "%s"
  }
  creation_policy {
    prefix      = "snap"
    period      = "00:00 GMT+03:00"
    keepday     = 1
    enable      = true
    delete_auto = true
  }
}
`, acc.DataSourceSecGroupDefault, acc.DataSourceSubnet, name, env.OS_AVAILABILITY_ZONE, osAgency)
}

func testResourceCSSSnapshotConfigurationV1Updated(name string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_css_cluster_v1" "cluster" {
  expect_node_num = 1
  name            = "%s"
  node_config {
    flavor = "css.medium.8"
    network_info {
      security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
      network_id        = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
      vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
    }
    volume {
      volume_type = "COMMON"
      size        = 100
    }

    availability_zone = "%s"
  }

  enable_https     = true
  enable_authority = true
  admin_pass       = "QwertyUI!"
}

resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "tf-snap-testing"
  force_destroy = true
}

resource "opentelekomcloud_css_snapshot_configuration_v1" "config" {
  cluster_id = opentelekomcloud_css_cluster_v1.cluster.id
  configuration {
    bucket = opentelekomcloud_obs_bucket.bucket.bucket
    agency = "%s"
  }
  creation_policy {
    prefix      = "snapshot"
    period      = "00:00 GMT+03:00"
    keepday     = 2
    enable      = true
    delete_auto = true
  }
}
`, acc.DataSourceSecGroupDefault, acc.DataSourceSubnet, name, env.OS_AVAILABILITY_ZONE, osAgency)
}
