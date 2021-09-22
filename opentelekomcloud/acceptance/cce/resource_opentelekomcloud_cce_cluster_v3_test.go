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
		PreCheck:          func() { common.TestAccPreCheck(t) },
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
		PreCheck:          func() { common.TestAccPreCheck(t) },
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
		PreCheck:          func() { common.TestAccPreCheck(t) },
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
		PreCheck:          func() { common.TestAccPreCheck(t) },
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
		PreCheck:          func() { common.TestAccPreCheck(t) },
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
		PreCheck:          func() { common.TestAccPreCheck(t) },
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
  authenticating_proxy {
    ca = <<EOT
-----BEGIN CERTIFICATE-----
MIICZjCCAc+gAwIBAgIUZtMIBg4MdR/h8yPITTx5+B0Xj0swDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMTA5MjExNDEyMThaFw0zMTA5
MTkxNDEyMThaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwgZ8wDQYJKoZIhvcNAQEB
BQADgY0AMIGJAoGBANSP54SyxRLkWjnGkLQcUOSxM6FUfk8k+PvSgq9xF0CeO8BR
NbY+y/+Jr9l1k0XRjAajIe/pdD4Yta24ox9yH5+ay+eKIFjdr2DAsyD5SSnW+A3o
e1eYAwZhWWgGRgq8VG6yik1vm1nsJW8crj+3k/N0kjiLYFGRTelpSThnesb3AgMB
AAGjUzBRMB0GA1UdDgQWBBT2daqVn+fyu7qoWflS7WwQ9EosODAfBgNVHSMEGDAW
gBT2daqVn+fyu7qoWflS7WwQ9EosODAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4GBAEoaBV1HA5wSDEGu1blKUtda13IxBg3ZpHWxwOpDS3gwXYzEOPAL
jgKvjKkLJaJnTp8V+YtO1xBhB272ZgWbr22Cer8TQZchNc16I2qLp+O9AQuPqVYO
15xHZN4yCgCVYcSlUm/HW2tJ3lAmilxkEFvJJcK1uLh7vqMflmcPSLe5
-----END CERTIFICATE-----
EOT
    cert = <<EOT
-----BEGIN CERTIFICATE-----
MIIByDCCATECFAzNYuav3B0dfSfIe7L8pDJ0t0LuMA0GCSqGSIb3DQEBCwUAMEUx
CzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRl
cm5ldCBXaWRnaXRzIFB0eSBMdGQwHhcNMjEwOTIxMTQzMDM3WhcNMzEwOTE5MTQz
MDM3WjBFMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UE
CgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMFwwDQYJKoZIhvcNAQEBBQADSwAw
SAJBALcrRkYvf/pLJQQp21KCPk56AuWh0UxbMd75NOZQWrY1QFXTRJ0xYHAGa/LG
LgposjKO9BELu+AYe+UoIJVwdAkCAwEAATANBgkqhkiG9w0BAQsFAAOBgQBwHUAE
Pzs5sfqJgJe0YDNZSC+V/+NZ67AYTDBzYt35dNEqUxbnsOaoMjuWr5a9/mkcIHPT
aXLd7Y85Iko9GDXrW4Xw58ZAqQO6nLSTvZ3LX5zxpmIY8RKw31bTnZbJWnzlV5UU
AmXWUWcptcu8WYcWfgrzP8Kob5znC1svHsKm0g==
-----END CERTIFICATE-----
EOT
    private_key = <<EOT
-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBALcrRkYvf/pLJQQp21KCPk56AuWh0UxbMd75NOZQWrY1QFXTRJ0x
YHAGa/LGLgposjKO9BELu+AYe+UoIJVwdAkCAwEAAQJBAJXeLHOErdum3DSZ4r+R
nVUsc25bhhpJi3Z6xJOlL3NgoDaAWEapQZ+jGs/XPCu14Skxwy5s9wgXznsfxIav
qWECIQDeVgWmBcvNz2FmQD8V1pIfQoec3hpTH3bVA06Rhg0j7QIhANLnGVpiCI+s
Pgqeqr93J1HojrcD9u5C9kahdt57GgUNAiBI5E7pxVCx4uF90mZcVIKHeRpY1YAv
7ErbP0BM+XPpaQIgNaUu37yb7N+lEFJ3oCgQylbbJlZN0yEZP7IGaGTro2kCIQCc
qYYLFv6yuySapSHrdOaPXnXrhMY4BE0EpzAuh+opxw==
-----END RSA PRIVATE KEY-----
EOT
  }
}
`, common.DataSourceSubnet, clusterName)

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
