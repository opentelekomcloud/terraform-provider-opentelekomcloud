package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/function"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getResourceObj(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FuncGraphV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud FunctionGraph client: %s", err)
	}
	return function.GetMetadata(c, state.Primary.ID)
}

func TestAccFgsV2Function_basic(t *testing.T) {
	var (
		f               function.FuncGraph
		randName        = fmt.Sprintf("fgs-acc-api%s", acctest.RandString(5))
		obsObjectConfig = zipFileUploadResourcesConfig()
		resourceName    = "opentelekomcloud_fgs_function_v2.test"
	)

	rc := common.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_basic_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "functiongraph_version", regexp.MustCompile(`v1|v2`)),
					resource.TestCheckResourceAttr(resourceName, "description", "function test"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "code_type", "inline"),
				),
			},
			{
				Config: testAccFgsV2Function_basic_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", "function test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(resourceName, "tags.newkey", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			{
				Config: testAccFgsV2Function_basic_step3(randName, obsObjectConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "code_type", "obs"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"func_code",
					"agency",
					"tags",
				},
			},
		},
	})
}

func TestAccFgsV2Function_text(t *testing.T) {
	var f function.FuncGraph
	randName := fmt.Sprintf("fgs-acc-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_fgs_function_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_text(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"func_code",
				},
			},
		},
	})
}

func TestAccFgsV2Function_createByImage(t *testing.T) {
	var f function.FuncGraph
	randName := fmt.Sprintf("fgs-acc-%s", acctest.RandString(5))
	rName1 := "opentelekomcloud_fgs_function_v2.create_with_vpc_access"
	rName2 := "opentelekomcloud_fgs_function_v2.create_without_vpc_access"

	rc1 := common.InitResourceCheck(
		rName1,
		&f,
		getResourceObj,
	)

	rc2 := common.InitResourceCheck(
		rName2,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckComponentDeployment(t)
			common.TestAccPreCheckImageUrlUpdated(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_createByImage_step_1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", randName+"_1"),
					resource.TestCheckResourceAttr(rName1, "agency", "functiongraph_swr_trust"),
					resource.TestCheckResourceAttr(rName1, "runtime", "Custom Image"),
					resource.TestCheckResourceAttr(rName1, "handler", "-"),
					resource.TestCheckResourceAttr(rName1, "custom_image.0.url", common.OTC_BUILD_IMAGE_URL),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", randName+"_2"),
					resource.TestCheckResourceAttr(rName2, "agency", "functiongraph_swr_trust"),
					resource.TestCheckResourceAttr(rName2, "runtime", "Custom Image"),
					resource.TestCheckResourceAttr(rName2, "handler", "-"),
					resource.TestCheckResourceAttr(rName2, "custom_image.0.url", common.OTC_BUILD_IMAGE_URL),
					resource.TestCheckResourceAttr(rName2, "vpc_id", ""),
					resource.TestCheckResourceAttr(rName2, "network_id", ""),
				),
			},
			{
				Config: testAccFgsV2Function_createByImage_step_2(randName),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "handler", "index.py"),
					resource.TestCheckResourceAttr(rName1, "vpc_id", ""),
					resource.TestCheckResourceAttr(rName1, "network_id", ""),
					resource.TestCheckResourceAttr(rName1, "custom_image.0.url", common.OTC_BUILD_IMAGE_URL_UPDATED),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "handler", "-"),
					resource.TestCheckResourceAttr(rName2, "custom_image.0.url", common.OTC_BUILD_IMAGE_URL_UPDATED),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"xrole",
					"agency",
				},
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"xrole",
					"agency",
				},
			},
		},
	})
}

func TestAccFgsV2Function_logConfig(t *testing.T) {
	var f function.FuncGraph
	randName := fmt.Sprintf("fgs-acc-%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_fgs_function_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_logConfig(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttrSet(resourceName, "log_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "log_topic_id"),
				),
			},
			{
				Config: testAccFgsV2Function_logConfigUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttrSet(resourceName, "log_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "log_topic_id"),
				),
			},
		},
	})
}

func zipFileUploadResourcesConfig() string {
	randName := fmt.Sprintf("fgs-acc-%s", acctest.RandString(5))

	return fmt.Sprintf(`
variable "script_content" {
  type    = string
  default = <<EOT
def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
EOT
}

resource "opentelekomcloud_obs_bucket" "test" {
  bucket = "%[1]s"
  acl    = "private"

  provisioner "local-exec" {
    command = "echo '${var.script_content}' >> test.py\nzip -r test.zip test.py"
  }
  provisioner "local-exec" {
    command = "rm test.zip test.py"
    when    = destroy
  }
}

resource "opentelekomcloud_obs_bucket_object" "test" {
  bucket = opentelekomcloud_obs_bucket.test.bucket
  key    = "test.zip"
  source = abspath("./test.zip")
}`, randName)
}

func testAccFgsV2Function_basic_step1(rName string) string {
	//nolint:revive
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%s"
  app         = "default"
  description = "function test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccFgsV2Function_basic_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[2]s"
  app         = "default"
  description = "function test update"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
  agency      = "function_vpc_trust"
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_id  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id

  tags = {
    foo    = "baar"
    newkey = "value"
  }
}
`, common.DataSourceSubnet, rName)
}

func testAccFgsV2Function_basic_step3(rName, obsConfig string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[3]s"
  app         = "default"
  description = "fuction test update"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "obs"
  code_url    = format("https://%%s/%%s", opentelekomcloud_obs_bucket.test.bucket_domain_name, opentelekomcloud_obs_bucket_object.test.key)
  agency      = "function_vpc_trust"
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_id  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id

  tags = {
    foo    = "baar"
    newkey = "value"
  }
}
`, common.DataSourceSubnet, obsConfig, rName)
}

func testAccFgsV2Function_text(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"

  func_code = <<EOF
# -*- coding:utf-8 -*-
import json
def handler (event, context):
    return {
        "statusCode": 200,
        "isBase64Encoded": False,
        "body": json.dumps(event),
        "headers": {
            "Content-Type": "application/json"
        }
    }
EOF
}
`, rName)
}

func testAccFgsV2Function_createByImage_step_1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_fgs_function_v2" "create_with_vpc_access" {
  name        = "%[2]s_1"
  app         = "default"
  handler     = "index.py"
  memory_size = 128
  runtime     = "Custom Image"
  timeout     = 3
  agency      = "functiongraph_swr_trust"
  code_type   = "Custom-Image-Swr"

  custom_image {
    url = "%[3]s"
  }

  vpc_id     = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}

resource "opentelekomcloud_fgs_function_v2" "create_without_vpc_access" {
  name        = "%[2]s_2"
  app         = "default"
  handler     = "index.py"
  memory_size = 128
  runtime     = "Custom Image"
  timeout     = 3
  agency      = "functiongraph_swr_trust"
  code_type   = "Custom-Image-Swr"

  custom_image {
    url = "%[3]s"
  }
}
`, common.DataSourceSubnet, rName, common.OTC_BUILD_IMAGE_URL)
}

func testAccFgsV2Function_createByImage_step_2(rName string) string {
	return fmt.Sprintf(`
%[1]s

# Closs the VPC access
resource "opentelekomcloud_fgs_function_v2" "create_with_vpc_access" {
  name        = "%[2]s_1"
  app         = "default"
  handler     = "index.py"
  memory_size = 128
  runtime     = "Custom Image"
  timeout     = 3
  agency      = "functiongraph_swr_trust"
  code_type   = "Custom-Image-Swr"

  custom_image {
    url = "%[3]s"
  }
}

# Open the VPC access
resource "opentelekomcloud_fgs_function_v2" "create_without_vpc_access" {
  name        = "%[2]s_2"
  app         = "default"
  handler     = "index.py"
  memory_size = 128
  runtime     = "Custom Image"
  timeout     = 3
  agency      = "functiongraph_swr_trust"
  code_type   = "Custom-Image-Swr"

  custom_image {
    url = "%[3]s"
  }

  vpc_id     = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}
`, common.DataSourceSubnet, rName, common.OTC_BUILD_IMAGE_URL_UPDATED)
}

func TestAccFgsV2Function_strategy(t *testing.T) {
	var (
		f function.FuncGraph

		name         = fmt.Sprintf("fgs-acc-%s", acctest.RandString(5))
		resourceName = "opentelekomcloud_fgs_function_v2.test"
	)

	rc := common.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_strategy_default(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "max_instance_num", "400"),
				),
			},
			{
				Config: testAccFunction_strategy_defined(name, 1000),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "max_instance_num", "1000"),
				),
			},
			// UpdateMaxInstances doesn't work when max_instance is set to 0
			// Response: 400 "no param has changed"
			// {
			// 	Config: testAccFunction_strategy_defined(name, 0),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		rc.CheckResourceExists(),
			// 		resource.TestCheckResourceAttr(resourceName, "max_instance_num", "0"),
			// 	),
			// },
			{
				Config: testAccFunction_strategy_defined(name, -1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "max_instance_num", "-1"),
				),
			},
			{
				Config: testAccFunction_strategy_default(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "max_instance_num", "-1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"func_code",
				},
			},
		},
	})
}

func testAccFunction_strategy_default(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  functiongraph_version = "v2"
  name                  = "%[1]s"
  app                   = "default"
  handler               = "index.handler"
  memory_size           = 128
  timeout               = 3
  runtime               = "Python2.7"
  code_type             = "inline"
  func_code             = "dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganN="
}
`, name)
}

func testAccFunction_strategy_defined(name string, maxInstanceNum int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  functiongraph_version = "v2"
  name                  = "%[1]s"
  app                   = "default"
  handler               = "index.handler"
  memory_size           = 128
  timeout               = 3
  runtime               = "Python2.7"
  code_type             = "inline"
  func_code             = "dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganN="
  max_instance_num      = %[2]d
}
`, name, maxInstanceNum)
}

func TestAccFgsV2Function_versions(t *testing.T) {
	var (
		f function.FuncGraph

		name         = fmt.Sprintf("fgs-acc-%s", acctest.RandString(5))
		resourceName = "opentelekomcloud_fgs_function_v2.test"
	)

	rc := common.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_versions_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.0.name", "latest"),
				),
			},
			{
				Config: testAccFunction_versions_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.0.name", "latest"),
					resource.TestCheckResourceAttr(resourceName, "versions.0.aliases.0.name", "demo"),
					resource.TestCheckResourceAttr(resourceName, "versions.0.aliases.0.description",
						"This is a description of the demo alias"),
				),
			},
			{
				Config: testAccFunction_versions_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.0.name", "latest"),
					resource.TestCheckResourceAttr(resourceName, "versions.0.aliases.0.name", "demo_update"),
					resource.TestCheckResourceAttr(resourceName, "versions.0.aliases.0.description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"func_code",
				},
			},
		},
	})
}

func testAccFunction_versions_step1(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  functiongraph_version = "v2"
  name                  = "%[1]s"
  app                   = "default"
  handler               = "index.handler"
  memory_size           = 128
  timeout               = 3
  runtime               = "Python2.7"
  code_type             = "inline"
  func_code             = "dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganN="

  // Test whether 'plan' and 'apply' commands will report an error when only the version number is filled in.
  versions {
    name = "latest"
  }
}
`, name)
}

func testAccFunction_versions_step2(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  functiongraph_version = "v2"
  name                  = "%[1]s"
  app                   = "default"
  handler               = "index.handler"
  memory_size           = 128
  timeout               = 3
  runtime               = "Python2.7"
  code_type             = "inline"
  func_code             = "dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganN="

  versions {
    name = "latest"

    aliases {
      name        = "demo"
      description = "This is a description of the demo alias"
    }
  }
}
`, name)
}

func testAccFunction_versions_step3(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  functiongraph_version = "v2"
  name                  = "%[1]s"
  app                   = "default"
  handler               = "index.handler"
  memory_size           = 128
  timeout               = 3
  runtime               = "Python2.7"
  code_type             = "inline"
  func_code             = "dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganN="

  versions {
    name = "latest"

    aliases {
      name = "demo_update"
    }
  }
}
`, name)
}

func testAccFgsV2Function_logConfig(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_logtank_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_logtank_topic_v2" "test" {
  group_id   = opentelekomcloud_logtank_group_v2.test.id
  topic_name = "%[1]s"
}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[1]s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganN="

  log_group_id   = opentelekomcloud_logtank_group_v2.test.id
  log_topic_id   = opentelekomcloud_logtank_topic_v2.test.id
  log_group_name = opentelekomcloud_logtank_group_v2.test.group_name
  log_topic_name = opentelekomcloud_logtank_topic_v2.test.topic_name
}
`, rName)
}

func testAccFgsV2Function_logConfigUpdate(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_logtank_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_logtank_topic_v2" "test" {
  group_id   = opentelekomcloud_logtank_group_v2.test.id
  topic_name = "%[1]s"
}

resource "opentelekomcloud_logtank_group_v2" "test1" {
  group_name  = "%[1]s-new"
  ttl_in_days = 30
}

resource "opentelekomcloud_logtank_topic_v2" "test1" {
  group_id   = opentelekomcloud_logtank_group_v2.test1.id
  topic_name = "%[1]s-new"
}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[1]s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganN="

  log_group_id   = opentelekomcloud_logtank_group_v2.test1.id
  log_topic_id   = opentelekomcloud_logtank_topic_v2.test1.id
  log_group_name = opentelekomcloud_logtank_group_v2.test1.group_name
  log_topic_name = opentelekomcloud_logtank_topic_v2.test1.topic_name
}
`, rName)
}

func TestAccFgsV2Function_reservedInstance_version(t *testing.T) {
	var (
		f            function.FuncGraph
		name         = fmt.Sprintf("fgs-acc-%s", acctest.RandString(5))
		resourceName = "opentelekomcloud_fgs_function_v2.test"
	)

	rc := common.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_reservedInstance_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.qualifier_name", "latest"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.qualifier_type", "version"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.count", "1"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.idle_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.tactics_config.0.cron_configs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.tactics_config.0.cron_configs.0.count", "2"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.tactics_config.0.cron_configs.0.cron", "0 */10 * * * ?"),
					resource.TestCheckResourceAttrSet(resourceName, "reserved_instances.0.tactics_config.0.cron_configs.0.start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "reserved_instances.0.tactics_config.0.cron_configs.0.expired_time"),
					resource.TestCheckResourceAttrSet(resourceName, "reserved_instances.0.tactics_config.0.cron_configs.0.name"),
				),
			},
			{
				Config: testAccFgsV2Function_reservedInstance_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.count", "2"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.idle_mode", "false"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.tactics_config.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"func_code",
					"tags",
				},
			},
		},
	})
}

func TestAccFgsV2Function_reservedInstance_alias(t *testing.T) {
	var (
		f               function.FuncGraph
		name            = fmt.Sprintf("fgs-acc-%s", acctest.RandString(5))
		updateAliasName = fmt.Sprintf("fgs-acc-%s-updated", acctest.RandString(5))
		resourceName    = "opentelekomcloud_fgs_function_v2.test"
	)

	rc := common.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_reservedInstance_alias(name, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.qualifier_name", name),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.qualifier_type", "alias"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.count", "1"),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.idle_mode", "false"),
				),
			},
			{
				Config: testAccFgsV2Function_reservedInstance_alias(name, updateAliasName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.qualifier_name", updateAliasName),
					resource.TestCheckResourceAttr(resourceName, "reserved_instances.0.qualifier_type", "alias"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"func_code",
					"tags",
				},
			},
		},
	})
}

func testAccFgsV2Function_reservedInstance_step1(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Node.js16.17"
  code_type   = "inline"

  reserved_instances {
    qualifier_type = "version"
    qualifier_name = "latest"
    count          = 1
    idle_mode      = true

    tactics_config {
      cron_configs {
        name         = "scheme-waekcy"
        cron         = "0 */10 * * * ?"
        start_time   = "1708342889"
        expired_time = "1739878889"
        count        = 2
      }
    }
  }
}
`, rName)
}

func testAccFgsV2Function_reservedInstance_step2(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Node.js16.17"
  code_type   = "inline"

  reserved_instances {
    qualifier_type = "version"
    qualifier_name = "latest"
    count          = 2
    idle_mode      = false
  }
}
`, rName)
}

func testAccFgsV2Function_reservedInstance_alias(rName string, aliasName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Node.js16.17"
  code_type   = "inline"

  versions {
    name = "latest"

    aliases {
      name = "%[2]s"
    }
  }

  reserved_instances {
    qualifier_type = "alias"
    qualifier_name = "%[2]s"
    count          = 1
    idle_mode      = false
  }
}
`, rName, aliasName)
}

func TestAccFgsV2Function_concurrencyNum(t *testing.T) {
	var (
		f function.FuncGraph

		name         = fmt.Sprintf("fgs-acc-%s", acctest.RandString(5))
		resourceName = "opentelekomcloud_fgs_function_v2.test"
	)

	rc := common.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_strategy_default(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "concurrency_num", "1"),
				),
			},
			{
				Config: testAccFunction_concurrencyNum(name, 1000),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "concurrency_num", "1000"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"func_code",
				},
			},
		},
	})
}

func testAccFunction_concurrencyNum(name string, concurrencyNum int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  functiongraph_version = "v2"
  name                  = "%[1]s"
  app                   = "default"
  handler               = "index.handler"
  memory_size           = 128
  timeout               = 3
  runtime               = "Python2.7"
  code_type             = "inline"
  func_code             = "dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganN="
  concurrency_num       = %[2]d
}
`, name, concurrencyNum)
}
