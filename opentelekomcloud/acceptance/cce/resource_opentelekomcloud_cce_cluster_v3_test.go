package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/clusters"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

var clusterName = fmt.Sprintf("cce-%s", acctest.RandString(5))

const resourceName = "opentelekomcloud_cce_cluster_v3.cluster_1"

func TestAccCCEClusterV3_basic(t *testing.T) {
	var cluster clusters.Clusters

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "x509"),
					resource.TestCheckResourceAttr(resourceName, "kube_proxy_mode", "iptables"),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_svc_ip_range", "10.247.0.0/16"),
					resource.TestCheckResourceAttrSet(resourceName, "security_group_control"),
					resource.TestCheckResourceAttrSet(resourceName, "security_group_node"),
				),
			},
			{
				Config: testAccCCEClusterV3Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "new description"),
				),
			},
		},
	})
}

func TestAccCCEClusterV3_invalidNetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCCEClusterV3InvalidSubnet,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't find subnet.+`),
			},
			{
				Config:      testAccCCEClusterV3InvalidVPC,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't find VPC.+`),
			},
			{
				Config:             testAccCCEClusterV3Computed,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCCEClusterV3_proxyAuth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3AuthProxy,
				Check:  resource.TestCheckResourceAttr(resourceName, "authentication_mode", "authenticating_proxy"),
			},
		},
	})
}

func TestAccCCEClusterV3_timeout(t *testing.T) {
	var cluster clusters.Clusters

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
		},
	})
}

func TestAccCCEClusterV3NoAddons(t *testing.T) {
	var cluster clusters.Clusters

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3NoAddons,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "installed_addons.#", "0"),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.CceV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud CCE client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cce_cluster_v3" {
			continue
		}

		_, err := clusters.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("cluster still exists")
		}
	}

	return nil
}

func testAccCheckCCEClusterV3Exists(n string, cluster *clusters.Clusters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CceV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud CCE client: %s", err)
		}

		found, err := clusters.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func TestAccCCEClusterV3_withVersionDiff(t *testing.T) {
	var cluster clusters.Clusters

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3WithInvalidVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists("opentelekomcloud_cce_cluster_v3.cluster_1", &cluster),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "name", clusterName),
				),
			},
		},
	})
}

var (
	testAccCCEClusterV3Basic = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}`, common.DataSourceSubnet, clusterName)

	testAccCCEClusterV3Update = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  container_network_type  = "overlay_l2"
  description             = "new description"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}`, common.DataSourceSubnet, clusterName)

	testAccCCEClusterV3Timeout = fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s2.small"
  vpc_id                 = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  eip                    = opentelekomcloud_networking_floatingip_v2.fip_1.address
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  timeouts {
    create = "20m"
    delete = "10m"
  }

  multi_az = true
}
`, common.DataSourceSubnet, clusterName)

	testAccCCEClusterV3WithInvalidVersion = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  cluster_version        = "v1.9.2"
  vpc_id                 = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  container_network_type = "overlay_l2"
  description            = "new description"
}`, common.DataSourceSubnet, clusterName)

	testAccCCEClusterV3AuthProxy = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
  authentication_mode     = "authenticating_proxy"
  authenticating_proxy_ca = <<EOT
-----BEGIN CERTIFICATE-----
MIIDpTCCAo2gAwIBAgIJAKdmmOBYnFvoMA0GCSqGSIb3DQEBCwUAMGkxCzAJBgNV
BAYTAnh4MQswCQYDVQQIDAJ4eDELMAkGA1UEBwwCeHgxCzAJBgNVBAoMAnh4MQsw
CQYDVQQLDAJ4eDELMAkGA1UEAwwCeHgxGTAXBgkqhkiG9w0BCQEWCnh4QDE2My5j
b20wHhcNMTcxMjA0MDM0MjQ5WhcNMjAxMjAzMDM0MjQ5WjBpMQswCQYDVQQGEwJ4
eDELMAkGA1UECAwCeHgxCzAJBgNVBAcMAnh4MQswCQYDVQQKDAJ4eDELMAkGA1UE
CwwCeHgxCzAJBgNVBAMMAnh4MRkwFwYJKoZIhvcNAQkBFgp4eEAxNjMuY29tMIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwZ5UJULAjWr7p6FVwGRQRjFN
2s8tZ/6LC3X82fajpVsYqF1xqEuUDndDXVD09E4u83MS6HO6a3bIVQDp6/klnYld
iE6Vp8HH5BSKaCWKVg8lGWg1UM9wZFnlryi14KgmpIFmcu9nA8yV/6MZAe6RSDmb
3iyNBmiZ8aZhGw2pI1YwR+15MVqFFGB+7ExkziROi7L8CFCyCezK2/oOOvQsH1dz
Q8z1JXWdg8/9Zx7Ktvgwu5PQM3cJtSHX6iBPOkMU8Z8TugLlTqQXKZOEgwajwvQ5
mf2DPkVgM08XAgaLJcLigwD513koAdtJd5v+9irw+5LAuO3JclqwTvwy7u/YwwID
AQABo1AwTjAdBgNVHQ4EFgQUo5A2tIu+bcUfvGTD7wmEkhXKFjcwHwYDVR0jBBgw
FoAUo5A2tIu+bcUfvGTD7wmEkhXKFjcwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0B
AQsFAAOCAQEAWJ2rS6Mvlqk3GfEpboezx2J3X7l1z8Sxoqg6ntwB+rezvK3mc9H0
83qcVeUcoH+0A0lSHyFN4FvRQL6X1hEheHarYwJK4agb231vb5erasuGO463eYEG
r4SfTuOm7SyiV2xxbaBKrXJtpBp4WLL/s+LF+nklKjaOxkmxUX0sM4CTA7uFJypY
c8Tdr8lDDNqoUtMD8BrUCJi+7lmMXRcC3Qi3oZJW76ja+kZA5mKVFPd1ATih8TbA
i34R7EQDtFeiSvBdeKRsPp8c0KT8H1B4lXNkkCQs2WX5p4lm99+ZtLD4glw8x6Ic
i1YhgnQbn5E0hz55OLu5jvOkKQjPCW+9Aa==
-----END CERTIFICATE-----
EOT
}`, common.DataSourceSubnet, clusterName)

	testAccCCEClusterV3InvalidSubnet = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = "abc"
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}
`, common.DataSourceSubnet, clusterName)

	testAccCCEClusterV3InvalidVPC = fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc" {
  cidr = "192.168.0.0/16"
  name = "cce-test"
}

locals {
  subnet_cidr  = cidrsubnet(opentelekomcloud_vpc_v1.vpc.cidr, 8, 0)
  subnet_gw_ip = cidrhost(local.subnet_cidr, 1)
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet" {
  cidr       = local.subnet_cidr
  gateway_ip = local.subnet_gw_ip
  name       = "cce-test"
  vpc_id     = opentelekomcloud_vpc_v1.vpc.id
}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = "abc"
  subnet_id               = opentelekomcloud_vpc_subnet_v1.subnet.id
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}
`, clusterName)

	testAccCCEClusterV3Computed = fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc" {
  cidr = "192.168.0.0/16"
  name = "cce-test"
}

locals {
  subnet_cidr  = cidrsubnet(opentelekomcloud_vpc_v1.vpc.cidr, 8, 0)
  subnet_gw_ip = cidrhost(local.subnet_cidr, 1)
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet" {
  cidr       = local.subnet_cidr
  gateway_ip = local.subnet_gw_ip
  name       = "cce-test"
  vpc_id     = opentelekomcloud_vpc_v1.vpc.id
}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = opentelekomcloud_vpc_v1.vpc.id
  subnet_id               = opentelekomcloud_vpc_subnet_v1.subnet.id
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}
`, clusterName)

	testAccCCEClusterV3NoAddons = fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
  no_addons               = true
}`, common.DataSourceSubnet, clusterName)
)
