package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ac "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/access-config"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCceAccessConfig_containerFile(t *testing.T) {
	clusterID := os.Getenv("OS_LTS_CCE_CLUSTER_ID")
	if clusterID == "" {
		t.Skip("The OS_LTS_CCE_CLUSTER_ID must be set for the acceptance test")
	}
	var (
		access ac.AccessConfigInfo
		rName  = "opentelekomcloud_lts_cce_access_v3.container_file"
		name   = fmt.Sprintf("lts_cce_access%s", acctest.RandString(3))
		rc     = common.InitResourceCheck(rName, &access, getHostAccessConfigResourceFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCceAccessConfigContainerFile(name, "", clusterID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "/var"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.0", "/var/a.log"),
					resource.TestCheckResourceAttr(rName, "access_config.0.path_type", "container_file"),
					resource.TestCheckResourceAttr(rName, "access_config.0.name_space_regex", "namespace"),
					resource.TestCheckResourceAttr(rName, "access_config.0.pod_name_regex", "podname"),
					resource.TestCheckResourceAttr(rName, "access_config.0.container_name_regex", "containername"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey1", "bar1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey2", "bar"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey1", "incval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey2", "incval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey1", "excval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey2", "excval"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey1", "envval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey2", "envval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey1", "incval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey2", "incval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey1", "excval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey2", "excval"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey1", "k8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey2", "k8sval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey1", "ink8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey2", "ink8sval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey1", "exk8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey2", "exk8sval"),
				),
			},
			{
				Config: testCceAccessConfigContainerFile(name, "-update", clusterID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar-update"),
					resource.TestCheckResourceAttr(rName, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttr(rName, "access_config.0.name_space_regex", "namespace-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.pod_name_regex", "podname-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.container_name_regex", "containername-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey2", "bar-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey2", "incval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey2", "excval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey2", "envval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey2", "incval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey2", "excval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey2", "k8sval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey2", "ink8sval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey2", "exk8sval-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCceAccessConfig_containerStdout(t *testing.T) {
	clusterID := os.Getenv("OS_LTS_CCE_CLUSTER_ID")
	if clusterID == "" {
		t.Skip("The OS_LTS_CCE_CLUSTER_ID must be set for the acceptance test")
	}
	var (
		access ac.AccessConfigInfo
		rName  = "opentelekomcloud_lts_cce_access_v3.container_stdout"
		name   = fmt.Sprintf("lts_cce_access%s", acctest.RandString(3))
		rc     = common.InitResourceCheck(rName, &access, getHostAccessConfigResourceFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCceAccessConfigContainerStdout(name, "", clusterID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttr(rName, "access_config.0.path_type", "container_stdout"),
					resource.TestCheckResourceAttr(rName, "access_config.0.stdout", "true"),
					resource.TestCheckResourceAttr(rName, "access_config.0.name_space_regex", "namespace"),
					resource.TestCheckResourceAttr(rName, "access_config.0.pod_name_regex", "podname"),
					resource.TestCheckResourceAttr(rName, "access_config.0.container_name_regex", "containername"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey1", "bar1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey2", "bar"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey1", "incval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey2", "incval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey1", "excval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey2", "excval"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey1", "envval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey2", "envval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey1", "incval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey2", "incval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey1", "excval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey2", "excval"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey1", "k8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey2", "k8sval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey1", "ink8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey2", "ink8sval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey1", "exk8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey2", "exk8sval"),
				),
			},
			{
				Config: testCceAccessConfigContainerStdout(name, "-update", clusterID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.name_space_regex", "namespace-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.pod_name_regex", "podname-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.container_name_regex", "containername-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey2", "bar-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey2", "incval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey2", "excval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey2", "envval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey2", "incval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey2", "excval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey2", "k8sval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey2", "ink8sval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey2", "exk8sval-update"),
				),
			},
		},
	})
}

func TestAccCceAccessConfig_hostFile(t *testing.T) {
	clusterID := os.Getenv("OS_LTS_CCE_CLUSTER_ID")
	if clusterID == "" {
		t.Skip("The OS_LTS_CCE_CLUSTER_ID must be set for the acceptance test")
	}
	var (
		access ac.AccessConfigInfo
		rName  = "opentelekomcloud_lts_cce_access_v3.host_file"
		name   = fmt.Sprintf("lts_cce_access%s", acctest.RandString(3))
		rc     = common.InitResourceCheck(rName, &access, getHostAccessConfigResourceFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCceAccessConfigHostFile(name, clusterID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.#", "2"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.#", "2"),
					resource.TestCheckResourceAttr(rName, "access_config.0.path_type", "host_file"),
				),
			},
			{
				Config: testCceAccessConfigHostFileUpdate(name, clusterID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.key", "value-updated"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(rName, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.#", "1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "/var/logs"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.#", "0"),
				),
			},
		},
	})
}

func testCceAccessConfigContainerFile(name, suffix, clusterId string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_host_group_v3" "test" {
  name = "%[1]s"
  type = "linux"
}

resource "opentelekomcloud_lts_cce_access_v3" "container_file" {
  name           = "%[1]s"
  log_group_id   = opentelekomcloud_lts_group_v2.test.id
  log_stream_id  = opentelekomcloud_lts_stream_v2.test.id
  host_group_ids = [opentelekomcloud_lts_host_group_v3.test.id]
  cluster_id     = "%[2]s"

  access_config {
    path_type            = "container_file"
    paths                = ["/var"]
    black_paths          = ["/var/a.log"]
    name_space_regex     = "namespace%[3]s"
    pod_name_regex       = "podname%[3]s"
    container_name_regex = "containername%[3]s"

    windows_log_info {
      categories       = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    single_log_format {
      mode = "system"
    }

    log_labels = {
      loglabelkey1 = "bar1"
      loglabelkey2 = "bar%[3]s"
    }

    include_labels = {
      includeKey1 = "incval1"
      includeKey2 = "incval%[3]s"
    }

    exclude_labels = {
      excludeKey1 = "excval1"
      excludeKey2 = "excval%[3]s"
    }

    log_envs = {
      envKey1 = "envval1"
      envKey2 = "envval%[3]s"
    }

    include_envs = {
      inEnvKey1 = "incval1"
      inEnvKey2 = "incval%[3]s"
    }

    exclude_envs = {
      exEnvKey1 = "excval1"
      exEnvKey2 = "excval%[3]s"
    }

    log_k8s = {
      k8sKey1 = "k8sval1"
      k8sKey2 = "k8sval%[3]s"
    }

    include_k8s_labels = {
      ink8sKey1 = "ink8sval1"
      ink8sKey2 = "ink8sval%[3]s"
    }

    exclude_k8s_labels = {
      exk8sKey1 = "exk8sval1"
      exk8sKey2 = "exk8sval%[3]s"
    }
  }

  tags = {
    key = "value"
    foo = "bar%[3]s"
  }
}
`, name, clusterId, suffix)
}

func testCceAccessConfigContainerStdout(name, suffix, clusterId string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_host_group_v3" "test" {
  name = "%[1]s"
  type = "linux"
}

resource "opentelekomcloud_lts_cce_access_v3" "container_stdout" {
  name           = "%[1]s"
  log_group_id   = opentelekomcloud_lts_group_v2.test.id
  log_stream_id  = opentelekomcloud_lts_stream_v2.test.id
  host_group_ids = [opentelekomcloud_lts_host_group_v3.test.id]
  cluster_id     = "%[3]s"

  access_config {
    path_type            = "container_stdout"
    stdout               = true
    name_space_regex     = "namespace%[2]s"
    pod_name_regex       = "podname%[2]s"
    container_name_regex = "containername%[2]s"

    windows_log_info {
      categories       = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    single_log_format {
      mode = "system"
    }

    log_labels = {
      loglabelkey1 = "bar1"
      loglabelkey2 = "bar%[2]s"
    }

    include_labels = {
      includeKey1 = "incval1"
      includeKey2 = "incval%[2]s"
    }

    exclude_labels = {
      excludeKey1 = "excval1"
      excludeKey2 = "excval%[2]s"
    }

    log_envs = {
      envKey1 = "envval1"
      envKey2 = "envval%[2]s"
    }

    include_envs = {
      inEnvKey1 = "incval1"
      inEnvKey2 = "incval%[2]s"
    }

    exclude_envs = {
      exEnvKey1 = "excval1"
      exEnvKey2 = "excval%[2]s"
    }

    log_k8s = {
      k8sKey1 = "k8sval1"
      k8sKey2 = "k8sval%[2]s"
    }

    include_k8s_labels = {
      ink8sKey1 = "ink8sval1"
      ink8sKey2 = "ink8sval%[2]s"
    }

    exclude_k8s_labels = {
      exk8sKey1 = "exk8sval1"
      exk8sKey2 = "exk8sval%[2]s"
    }
  }

  tags = {
    key = "value"
    foo = "bar%[2]s"
  }
}
`, name, suffix, clusterId)
}

func testCceAccessConfigHostFile(name, clusterId string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_host_group_v3" "test" {
  name = "%[1]s"
  type = "linux"
}

resource "opentelekomcloud_lts_cce_access_v3" "host_file" {
  name           = "%[1]s"
  log_group_id   = opentelekomcloud_lts_group_v2.test.id
  log_stream_id  = opentelekomcloud_lts_stream_v2.test.id
  host_group_ids = [opentelekomcloud_lts_host_group_v3.test.id]
  cluster_id     = "%[2]s"

  access_config {
    path_type   = "host_file"
    paths       = ["/var", "/temp"]
    black_paths = ["/var/temp.log", "/var/a.log"]

    single_log_format {
      mode = "system"
    }
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, name, clusterId)
}

func testCceAccessConfigHostFileUpdate(name, clusterId string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_host_group_v3" "test" {
  name = "%[1]s"
  type = "linux"
}

resource "opentelekomcloud_lts_cce_access_v3" "host_file" {
  name           = "%[1]s"
  log_group_id   = opentelekomcloud_lts_group_v2.test.id
  log_stream_id  = opentelekomcloud_lts_stream_v2.test.id
  host_group_ids = [opentelekomcloud_lts_host_group_v3.test.id]
  cluster_id     = "%[2]s"

  access_config {
    path_type = "host_file"
    paths     = ["/var/logs"]

    single_log_format {
      mode = "system"
    }
  }

  tags = {
    key   = "value-updated"
    owner = "terraform"
  }
}
`, name, clusterId)
}
