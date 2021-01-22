package opentelekomcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/clusters"
)

var clusterName = fmt.Sprintf("cce-%s", acctest.RandString(5))

func TestAccCCEClusterV3_basic(t *testing.T) {
	var cluster clusters.Clusters

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists("opentelekomcloud_cce_cluster_v3.cluster_1", &cluster),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "name", clusterName),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "status", "Available"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "authentication_mode", "x509"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "kube_proxy_mode", "iptables"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "kubernetes_svc_ip_range", "10.247.0.0/16"),
				),
			},
			{
				Config: testAccCCEClusterV3_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "description", "new description"),
				),
			},
		},
	})
}

func TestAccCCEClusterV3_invalidNetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCCEClusterV3_invalidSubnet,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't find subnet.+`),
			},
			{
				Config:      testAccCCEClusterV3_invalidVPC,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't find VPC.+`),
			},
			{
				Config:             testAccCCEClusterV3_computed,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCCEClusterV3_proxyAuth(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_authProxy,
				Check: resource.TestCheckResourceAttr(
					"opentelekomcloud_cce_cluster_v3.cluster_1", "authentication_mode", "authenticating_proxy"),
			},
		},
	})
}

func TestAccCCEClusterV3_timeout(t *testing.T) {
	var cluster clusters.Clusters

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists("opentelekomcloud_cce_cluster_v3.cluster_1", &cluster),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "authentication_mode", "rbac"),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud CCE client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cce_cluster_v3" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
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

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.cceV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud CCE client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_withInvalidVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists("opentelekomcloud_cce_cluster_v3.cluster_1", &cluster),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "name", "opentelekomcloud-cce"),
				),
			},
		},
	})
}

var testAccCCEClusterV3_basic = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}`, clusterName, OS_VPC_ID, OS_NETWORK_ID)

var testAccCCEClusterV3_update = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
  description="new description"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}`, clusterName, OS_VPC_ID, OS_NETWORK_ID)

var testAccCCEClusterV3_timeout = fmt.Sprintf(`
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s2.small"
  cluster_version = "v1.9.2-r2"
  vpc_id="%s"
  subnet_id="%s"
  eip=opentelekomcloud_networking_floatingip_v2.fip_1.address
  container_network_type="overlay_l2"
  authentication_mode = "rbac"
    timeouts {
    create = "20m"
    delete = "10m"
  }
  multi_az = true
}
}`, clusterName, OS_VPC_ID, OS_NETWORK_ID)

var testAccCCEClusterV3_withInvalidVersion = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  cluster_version = "v1.9.2"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
  description="new description"
}`, clusterName, OS_VPC_ID, OS_NETWORK_ID)

var testAccCCEClusterV3_authProxy = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
  authentication_mode = "authenticating_proxy"
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
}`, clusterName, OS_VPC_ID, OS_NETWORK_ID)

var testAccCCEClusterV3_invalidSubnet = fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc" {
  cidr = "192.168.0.0/16"
  name = "cce-test"
}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = opentelekomcloud_vpc_v1.vpc.id
  subnet_id               = "abc"
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}
`, clusterName)

var testAccCCEClusterV3_invalidVPC = fmt.Sprintf(`
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

var testAccCCEClusterV3_computed = fmt.Sprintf(`
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
