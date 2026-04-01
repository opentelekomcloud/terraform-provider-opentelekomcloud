package cci

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/pod"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getV2PodResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CciV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CCI v2 client: %s", err)
	}
	return pod.Get(client, state.Primary.Attributes["namespace"], state.Primary.Attributes["name"])
}

func TestAccV2Pod_basic(t *testing.T) {
	cciImage := os.Getenv("OS_CCI_IMAGE")
	if cciImage == "" {
		t.Skip("OS_CCI_IMAGE is not set, skipping CCI pod test")
	}

	var p pod.Pod
	rName := fmt.Sprintf("cci-pod-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_cci_pod_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&p,
		getV2PodResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Pod_basic(rName, cciImage),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "namespace", rName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "containers.0.image", cciImage),
					resource.TestCheckResourceAttr(resourceName, "containers.0.name", "c1"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.limits.cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.limits.memory", "4G"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.requests.cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "containers.0.resources.0.requests.memory", "4G"),
					resource.TestCheckResourceAttr(resourceName, "image_pull_secrets.0.name", "imagepull-secret"),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
				),
			},
			{
				Config: testAccV2Pod_update(rName, cciImage),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "containers.0.image", cciImage),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
				),
			},
		},
	})
}

func testAccV2Pod_basic(rName, image string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cci_pod_v2" "test" {
  depends_on = [opentelekomcloud_cci_network_v2.test]

  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"

  annotations = {
    "description"                    = "test",
    "resource.cci.io/pod-size-specs" = "2.00_4.0",
    "resource.cci.io/instance-type"  = "general-computing",
  }

  containers {
    image = "%[3]s"
    name  = "c1"

    resources {
      limits = {
        cpu    = 2
        memory = "4G"
      }

      requests = {
        cpu    = 2
        memory = "4G"
      }
    }
  }

  image_pull_secrets {
    name = "imagepull-secret"
  }

  lifecycle {
    ignore_changes = [
      annotations, containers,
    ]
  }
}
`, testAccV2Network_basic(rName), rName, image)
}

func testAccV2Pod_update(rName, image string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cci_pod_v2" "test" {
  depends_on = [opentelekomcloud_cci_network_v2.test]

  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"

  annotations = {
    "description"                    = "test-updated",
    "resource.cci.io/pod-size-specs" = "2.00_4.0",
    "resource.cci.io/instance-type"  = "general-computing",
  }

  containers {
    image = "%[3]s"
    name  = "c1"

    resources {
      limits = {
        cpu    = 2
        memory = "4G"
      }

      requests = {
        cpu    = 2
        memory = "4G"
      }
    }
  }

  image_pull_secrets {
    name = "imagepull-secret"
  }

  lifecycle {
    ignore_changes = [
      annotations, containers,
    ]
  }
}
`, testAccV2Network_basic(rName), rName, image)
}
