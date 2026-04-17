package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataFunctions_basic(t *testing.T) {
	var (
		base = "opentelekomcloud_fgs_function_v2.test"

		all               = "data.opentelekomcloud_fgs_functions_v2.all"
		dcForAllFunctions = common.InitDataSourceCheck(all)

		byPackageName           = "data.opentelekomcloud_fgs_functions_v2.filter_by_package_name"
		dcByPackageName         = common.InitDataSourceCheck(byPackageName)
		byNotFoundPackageName   = "data.opentelekomcloud_fgs_functions_v2.filter_by_not_found_package_name"
		dcByNotFoundPackageName = common.InitDataSourceCheck(byNotFoundPackageName)

		byFunctionUrn           = "data.opentelekomcloud_fgs_functions_v2.filter_by_function_urn"
		dcByFunctionUrn         = common.InitDataSourceCheck(byFunctionUrn)
		byNotFoundFunctionUrn   = "data.opentelekomcloud_fgs_functions_v2.filter_by_not_found_function_urn"
		dcByNotFoundFunctionUrn = common.InitDataSourceCheck(byNotFoundFunctionUrn)

		byFunctionName           = "data.opentelekomcloud_fgs_functions_v2.filter_by_name"
		dcByFunctionName         = common.InitDataSourceCheck(byFunctionName)
		byNotFoundFunctionName   = "data.opentelekomcloud_fgs_functions_v2.filter_by_not_found_name"
		dcByNotFoundFunctionName = common.InitDataSourceCheck(byNotFoundFunctionName)

		byRuntime           = "data.opentelekomcloud_fgs_functions_v2.filter_by_runtime"
		dcByRuntime         = common.InitDataSourceCheck(byRuntime)
		byNotFoundRuntime   = "data.opentelekomcloud_fgs_functions_v2.filter_by_not_found_runtime"
		dcByNotFoundRuntime = common.InitDataSourceCheck(byNotFoundRuntime)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFunctions_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dcForAllFunctions.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "functions.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Filter by package name.
					dcByPackageName.CheckResourceExists(),
					resource.TestCheckOutput("is_package_name_filter_useful", "true"),
					dcByNotFoundPackageName.CheckResourceExists(),
					resource.TestCheckOutput("package_name_not_found_validation_pass", "true"),
					// Filter by function URN.
					dcByFunctionUrn.CheckResourceExists(),
					resource.TestCheckOutput("is_function_urn_filter_useful", "true"),
					dcByNotFoundFunctionUrn.CheckResourceExists(),
					resource.TestCheckOutput("function_urn_not_found_validation_pass", "true"),
					// Filter by function name.
					dcByFunctionName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundFunctionName.CheckResourceExists(),
					resource.TestCheckOutput("name_not_found_validation_pass", "true"),
					// Filter by function runtime.
					dcByRuntime.CheckResourceExists(),
					resource.TestCheckOutput("is_runtime_filter_useful", "true"),
					dcByNotFoundRuntime.CheckResourceExists(),
					resource.TestCheckOutput("runtime_not_found_validation_pass", "true"),
					// Check the attributes.
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.name", base, "name"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.urn", base, "urn"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.package", base, "app"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.runtime", base, "runtime"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.timeout", base, "timeout"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.handler", base, "handler"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.memory_size", base, "memory_size"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.code_type", base, "code_type"),
					resource.TestCheckResourceAttrPair(byFunctionUrn, "functions.0.description", base, "description"),
				),
			},
		},
	})
}

func testAccDataFunctionsV2_base(rName string) string {
	return fmt.Sprintf(`
variable "function_code_content" {
  type    = string
  default = <<EOT
def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
EOT
}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[1]s"
  memory_size = 128
  runtime     = "Python3.9"
  timeout     = 3
  handler     = "index.handler"
  app         = "default"
  description = "fuction test"
  code_type   = "inline"
  func_code   = base64encode(var.function_code_content)
}
`, rName)
}

func testAccDataFunctions_basic() string {
	randName := fmt.Sprintf("fgs-fns-%s", acctest.RandString(5))
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "opentelekomcloud_fgs_functions_v2" "all" {
  depends_on = [
    opentelekomcloud_fgs_function_v2.test
  ]
}

# Filter by package name.
locals {
  package_name = opentelekomcloud_fgs_function_v2.test.app
}

data "opentelekomcloud_fgs_functions_v2" "filter_by_package_name" {
  depends_on = [
    opentelekomcloud_fgs_function_v2.test,
  ]

  package_name = local.package_name
}

data "opentelekomcloud_fgs_functions_v2" "filter_by_not_found_package_name" {
  depends_on = [
    opentelekomcloud_fgs_function_v2.test,
  ]

  package_name = "package_name_not_found"
}

locals {
  package_name_filter_result = [for v in data.opentelekomcloud_fgs_functions_v2.filter_by_package_name.functions[*].package :
  v == local.package_name]
}

output "is_package_name_filter_useful" {
  value = length(local.package_name_filter_result) > 0 && alltrue(local.package_name_filter_result)
}

output "package_name_not_found_validation_pass" {
  value = length(data.opentelekomcloud_fgs_functions_v2.filter_by_not_found_package_name.functions) == 0
}

# Filter by function URN.
locals {
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
}

data "opentelekomcloud_fgs_functions_v2" "filter_by_function_urn" {
  depends_on = [
    opentelekomcloud_fgs_function_v2.test,
  ]

  urn = local.function_urn
}

data "opentelekomcloud_fgs_functions_v2" "filter_by_not_found_function_urn" {
  depends_on = [
    opentelekomcloud_fgs_function_v2.test,
  ]

  urn = "function_urn_not_found"
}

locals {
  function_urn_filter_result = [for v in data.opentelekomcloud_fgs_functions_v2.filter_by_function_urn.functions[*].urn :
  v == local.function_urn]
}

output "is_function_urn_filter_useful" {
  value = length(local.function_urn_filter_result) > 0 && alltrue(local.function_urn_filter_result)
}

output "function_urn_not_found_validation_pass" {
  value = length(data.opentelekomcloud_fgs_functions_v2.filter_by_not_found_function_urn.functions) == 0
}

# Filter by function name.
locals {
  function_name = opentelekomcloud_fgs_function_v2.test.name
}

data "opentelekomcloud_fgs_functions_v2" "filter_by_name" {
  depends_on = [
    opentelekomcloud_fgs_function_v2.test,
  ]

  name = local.function_name
}

data "opentelekomcloud_fgs_functions_v2" "filter_by_not_found_name" {
  depends_on = [
    opentelekomcloud_fgs_function_v2.test,
  ]

  name = "name_not_found"
}

locals {
  name_filter_result = [for v in data.opentelekomcloud_fgs_functions_v2.filter_by_function_urn.functions[*].name :
  v == local.function_name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

output "name_not_found_validation_pass" {
  value = length(data.opentelekomcloud_fgs_functions_v2.filter_by_not_found_name.functions) == 0
}

locals {
  function_runtime = opentelekomcloud_fgs_function_v2.test.runtime
}

data "opentelekomcloud_fgs_functions_v2" "filter_by_runtime" {
  depends_on = [
    opentelekomcloud_fgs_function_v2.test,
  ]

  runtime = local.function_runtime
}

data "opentelekomcloud_fgs_functions_v2" "filter_by_not_found_runtime" {
  depends_on = [
    opentelekomcloud_fgs_function_v2.test,
  ]

  runtime = "runtime_not_found"
}

locals {
  runtime_filter_result = [for v in data.opentelekomcloud_fgs_functions_v2.filter_by_runtime.functions[*].runtime :
  v == local.function_runtime]
}

output "is_runtime_filter_useful" {
  value = length(local.runtime_filter_result) > 0 && alltrue(local.runtime_filter_result)
}

output "runtime_not_found_validation_pass" {
  value = length(data.opentelekomcloud_fgs_functions_v2.filter_by_not_found_runtime.functions) == 0
}
`, testAccDataFunctionsV2_base(randName))
}
