package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/secret"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getV2SecretResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CciV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CCI v2 client: %s", err)
	}
	return secret.Get(client, state.Primary.Attributes["namespace"], state.Primary.Attributes["name"])
}

func TestAccV2Secret_basic(t *testing.T) {
	var sec secret.Secret
	rName := fmt.Sprintf("cci-secret-%s", acctest.RandString(5))
	nsName := fmt.Sprintf("cci-ns-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_cci_secret_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&sec,
		getV2SecretResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Secret_basic(nsName, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "api_version", "cci/v2"),
					resource.TestCheckResourceAttr(resourceName, "kind", "Secret"),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
					resource.TestCheckResourceAttrSet(resourceName, "data.%"),
					resource.TestCheckOutput("key1_verify", "true"),
				),
			},
			{
				Config: testAccV2Secret_update(nsName, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckOutput("key1_verify", "true"),
					resource.TestCheckOutput("key3_verify", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccV2Secret_basic(nsName, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cci_secret_v2" "test" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"
  type      = "Opaque"

  data = {
    "key1" = "dmFsdWUx"
  }
}

output "key1_verify" {
  value = opentelekomcloud_cci_secret_v2.test.data["key1"] == "dmFsdWUx"
}
`, testAccV2Namespace_basic(nsName), rName)
}

func TestAccV2Secret_dockerRegistry(t *testing.T) {
	var sec secret.Secret
	rName := fmt.Sprintf("cci-secret-%s", acctest.RandString(5))
	nsName := fmt.Sprintf("cci-ns-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_cci_secret_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&sec,
		getV2SecretResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Secret_dockerRegistry(nsName, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "kubernetes.io/dockerconfigjson"),
					resource.TestCheckResourceAttrSet(resourceName, "data..dockerconfigjson"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccV2Secret_dockerRegistry(nsName, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cci_secret_v2" "test" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"
  type      = "kubernetes.io/dockerconfigjson"

  data = {
    ".dockerconfigjson" = base64encode(jsonencode({
      auths = {
        "swr.eu-de.otc.t-systems.com" = {
          username = "testuser"
          password = "testpassword"
          auth     = base64encode("testuser:testpassword")
        }
      }
    }))
  }
}
`, testAccV2Namespace_basic(nsName), rName)
}

func testAccV2Secret_update(nsName, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cci_secret_v2" "test" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"
  type      = "Opaque"

  data = {
    "key1"       = "dXBkYXRlZF92YWx1ZTE="
    "expired.at" = "MjAyNS0wNC0xNlQwNTo1NzowMVo="
  }
}

output "key1_verify" {
  value = opentelekomcloud_cci_secret_v2.test.data["key1"] == "dXBkYXRlZF92YWx1ZTE="
}

output "key3_verify" {
  value = opentelekomcloud_cci_secret_v2.test.data["expired.at"] == "MjAyNS0wNC0xNlQwNTo1NzowMVo="
}
`, testAccV2Namespace_basic(nsName), rName)
}
