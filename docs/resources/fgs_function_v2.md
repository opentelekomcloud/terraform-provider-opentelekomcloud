---
subcategory: "FunctionGraph"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_fgs_function_v2"
sidebar_current: "docs-opentelekomcloud-resource-fgs-function-v2"
description: |-
  Manages an FGS Function resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for FGS you can get at
[documentation portal](https://docs.otc.t-systems.com/function-graph/api-ref/apis/index.html)

# opentelekomcloud_fgs_function_v2

Manages a V2 function graph resource within OpenTelekomCloud.

## Example Usage

### With base64 func code

```hcl
variable "function_name" {}
variable "function_codes" {}
variable "agency_name" {}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = var.function_name
  app         = "default"
  agency      = var.agency_name
  description = "fuction test"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = base64encode(var.function_codes)
}
```

### Create function using code stored in OBS

```hcl
resource "opentelekomcloud_obs_bucket" "codebucket" {
  bucket = "mycodebucket"
  acl    = "private"
}

resource "opentelekomcloud_obs_bucket_object" "code_object" {
  bucket       = opentelekomcloud_obs_bucket.codebucket.bucket
  key          = "code.zip"
  source       = "code.zip"
  etag         = filemd5("code.zip")
  content_type = "application/zip"
}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "my-function"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python3.9"
  code_type   = "obs"
  code_url    = format("https://%s/%s", opentelekomcloud_obs_bucket.codebucket.bucket_domain_name, opentelekomcloud_obs_bucket_object.code_object.key)

  # Triggers a FunctionGraph code update when the OBS object content changes.
  source_code_hash = opentelekomcloud_obs_bucket_object.code_object.etag
}
```

### With text code

```hcl
variable "function_name" {}
variable "agency_name" {}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = var.function_name
  app         = "default"
  agency      = var.agency_name
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = <<EOF
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
```

### Create function using SWR image

```hcl
variable "function_name" {}
variable "agency_name" {} // The agent name that authorizes FunctionGraph service SWR administrator privilege
variable "image_url" {}

resource "opentelekomcloud_fgs_function_v2" "by_swr_image" {
  name        = var.function_name
  agency      = var.agency_name
  handler     = "-"
  app         = "default"
  runtime     = "Custom Image"
  memory_size = 128
  timeout     = 3

  custom_image {
    url = var.image_url
  }
}
```

### Create function with an alias for latest version

```hcl
variable "function_name" {}
variable "function_codes" {}

resource "opentelekomcloud_fgs_function_v2" "with_alias" {
  name        = var.function_name
  app         = "default"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = base64encode(var.function_codes)

  versions {
    name = "latest"

    aliases {
      name = "demo"
    }
  }
}
```

### Create function with log group and stream

```hcl
variable "function_name" {}
variable "function_codes" {}
variable "log_group_id" {}
variable "log_topic_id" {}
variable "log_group_name" {}
variable "log_topic_name" {}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = var.function_name
  app         = "default"
  agency      = "test"
  description = "fuction test"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = base64encode(var.function_codes)

  log_group_id   = var.log_group_id
  log_topic_id   = var.log_topic_id
  log_group_name = var.log_group_name
  log_topic_name = var.log_topic_name
}
```

### Create function with advanced configuration

```hcl
variable "function_name" {}
variable "function_codes" {}
variable "agency_name" {}
variable "trigger_access_vpc_ids" {
  type = list(string)
}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name                  = var.function_name
  app                   = "default"
  agency                = var.agency_name
  description           = "fuction test"
  handler               = "test.handler"
  memory_size           = 128
  timeout               = 3
  runtime               = "Python2.7"
  code_type             = "inline"
  func_code             = base64encode(var.function_codes)
  functiongraph_version = "v2"

  network_controller {
    disable_public_network = true

    dynamic "trigger_access_vpcs" {
      for_each = var.trigger_access_vpc_ids

      content {
        vpc_id = trigger_access_vpcs.value
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the function.
  Changing this will create a new resource.

* `app` - (Required, String) Specifies the group to which the function belongs.

* `memory_size` - (Required, Int) Specifies the memory size allocated to the function, in MByte (MB).

* `runtime` - (Required, String, ForceNew) Specifies the environment for executing the function.
  The valid values are as follows:
    + **Java8**
    + **Java11**
    + **Node.js6.10**
    + **Node.js8.10**
    + **Node.js10.16**
    + **Node.js12.13**
    + **Node.js14.18**
    + **Node.js16.17**
    + **Node.js18.15**
    + **Python2.7**
    + **Python3.6**
    + **Python3.9**
    + **Python3.10**
    + **Go1.x**
    + **C#(.NET Core 2.1)**
    + **C#(.NET Core 3.1)**
    + **PHP7.3**
    + **Custom**
    + **http**
    + **Custom Image**

* `timeout` - (Required, Int) Specifies the timeout interval of the function, in seconds.
  The value ranges from `3` to `900`.

* `code_type` - (Optional, String) Specifies the function code type, which can be:
    + **inline**: inline code.
    + **zip**: ZIP file.
    + **jar**: JAR file or java functions.
    + **obs**: function code stored in an OBS bucket.
    + **Custom-Image-Swr**: function code comes from the SWR custom image.

* `handler` - (Required, String) Specifies the entry point of the function.

* `functiongraph_version` - (Optional, String, ForceNew) Specifies the FunctionGraph version, default value is **v2**.
  The valid values are as follows:
    + **v1**
    + **v2**

* `func_code` - (Optional, String) Specifies the function code.
  The code value can be encoded using **Base64** or just with the text code.
  Required if the `code_type` is set to **inline**, **zip**, or **jar**.

* `code_url` - (Optional, String) Specifies the code url.
  Required if the `code_type` is set to **obs**.

* `source_code_hash` - (Optional, String) Terraform-only value used to trigger a code update when the code source
  changes but `code_url` stays the same. This is especially useful with `code_type = "obs"`, for example by setting
  it to `opentelekomcloud_obs_bucket_object.example.etag`. If omitted, changes to the OBS object content alone will
  not trigger a FunctionGraph code update.

* `code_filename` - (Optional, String) Specifies the name of a function file.
  Required if the `code_type` is set to **jar** or **zip**.

* `depend_list` - (Optional, List) Specifies the ID list of the dependencies.

* `user_data` - (Optional, String) Specifies the Key/Value information defined for the function.

* `encrypted_user_data` - (Optional, String) Specifies the key/value information defined to be encrypted for the
  function.

* `agency` - (Optional, String) Specifies the agency. This parameter is mandatory if the function needs to access other
  cloud services.

* `app_agency` - (Optional, String) Specifies the execution agency enables you to obtain a token or an AK/SK for
  accessing other cloud services.

* `description` - (Optional, String) Specifies the description of the function.

* `initializer_handler` - (Optional, String) Specifies the initializer of the function.

* `initializer_timeout` - (Optional, Int) Specifies the maximum duration the function can be initialized. Value range:
  1s to 300s.

* `vpc_id` - (Optional, String) Specifies the ID of VPC.

* `network_id` - (Optional, String) Specifies the network ID of subnet.

* `mount_user_id` - (Optional, Int) Specifies the user ID, a non-0 integer from `–1` to `65,534`.
  Defaults to `-1`.

* `mount_user_group_id` - (Optional, Int) Specifies the user group ID, a non-0 integer from `–1` to `65,534`.
  Defaults to `-1`.

* `func_mounts` - (Optional, List) Specifies the file system list. The `func_mounts` object structure is documented
  below.

* `custom_image` - (Optional, List) Specifies the custom image configuration of the function.
  The `custom_image` structure is documented below.
  Required if the parameter `code_type` is **Custom-Image-Swr**.

* `max_instance_num` - (Optional, String) Specifies the maximum number of instances of the function.
  The valid value ranges from `-1` to `1,000`, defaults to `400`.
    + The minimum value is `-1` and means the number of instances is unlimited.
    + `0` means this function is disabled.
    + The empty value means to keep the default (latest updated) value.

  -> This parameter is only supported by the `v2` version of the function.

* `versions` - (Optional, List) Specifies the versions management of the function.
  The `versions` structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the function.

* `log_group_id` - (Optional, String) Specifies the ID of the LTS log group.

* `log_group_name` - (Optional, String) Specifies the name of the LTS log group.

* `log_topic_id` - (Optional, String) Specifies the ID of the LTS log stream.

* `log_topic_name` - (Optional, String) Specifies the name of the LTS stream.

* `reserved_instances` - (Optional, List) Specifies the reserved instance policies of the function.
  The `reserved_instances` structure is documented below.

* `gpu_memory` - (Optional, Int) Specifies the GPU memory size allocated to the function, in MByte (MB).
  The valid value ranges form `1,024` to `16,384`, the value must be a multiple of `1,024`.
  If not specified, the GPU function is disabled.

* `pre_stop_handler` - (Optional, String) Specifies the pre-stop handler of a function. The value must contain a period (.)
 in the format of xx.xx. For example, for Node.js function myfunction.pre_stop_handler, the file name is myfunction.js,
 and the initialization function is pre_stop_handler.

* `pre_stop_timeout` - (Optional, Int) Specifies the maximum duration the function can be initialized. Value range: 1s-90s.

* `network_controller` - (Optional, List) Specifies the network configuration of the function.
  The [network_controller](#function_network_controller) structure is documented below.

* `peering_cidr` - (Optional, String) Specifies the VPC cidr blocks used in the function code to detect whether it
  conflicts with the VPC cidr blocks used by the service.
  The cidr blocks are separated by semicolons and cannot exceed `5`.

* `enable_dynamic_memory` - (Optional, Boolean) Specifies whether to enable dynamic memory allocation.

* `enable_class_isolation` - (Optional, Boolean) Specifies whether to enable class isolation.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the
  function belongs.

* `dns_list` - (Optional, String) Specifies the private DNS configuration of the function network.
  Private DNS list is associated to the function by a string in the following format:
  `[{\"id\":\"ff8080828a07ffea018a17184aa310f5\","domain_name":"functiondebug.example1.com."}]`

  -> Ensure the agency with DNS management permissions specified before using this parameter.

* `enable_lts_log` - (Optional, Bool) Specifies whether to enable the LTS logging feature.
  The valid values are as follows:
  + **null** or omit this parameter definition: Using the default value configured on the FunctionGraph service to
    configure the LTS logging feature.
  + **true**: Explicitly enable the LTS logging feature.
  + **false**: Explicitly disable the LTS logging feature.

* `concurrency_num` - (Optional, Int) Specifies the number of concurrent requests of the function.
  The valid value is range from `1` to `1,000`, the default value is `1`.

  -> 1. This parameter is only supported by the `v2` version of the function.
  <br>2. This parameter is available only when the `runtime` parameter is set to **http** or **Custom Image**.

* `is_stateful_function` - (Optional, Bool) Specifies whether the function is a stateful function.
  Defaults to **false**.

* `enable_auth_in_header` - (Optional, Bool) Specifies whether the authentication in the request header is enabled.
  Defaults to **false**.

* `ephemeral_storage` - (Optional, Int) Specifies the size of the function ephemeral storage.
  The valid values are as follows:
  + **512**
  + **10240**

  Defaults to `512`. Only custom image or http runtime supported.

* `heartbeat_handler` - (Optional, String) Specifies the heartbeat handler of the function.
  The rule is **xx.xx**, such as **com.huawei.demo.TriggerTests.heartBeat**, it must contain periods (.).
  The heartbeat function entry must be in the same file as the function execution entry.

* `lts_custom_tag` - (Optional, Map) Specifies the custom tags configuration that used to filter the LTS logs.
  This parameter is available only when `enable_lts_log` is set to **true**.

  -> This parameter is only supported by the `v2` version of the function.

The `func_mounts` block supports:

* `mount_type` - (Required, String) Specifies the mount type.
    + **sfs**
    + **sfsTurbo**
    + **ecs**

* `mount_resource` - (Required, String) Specifies the ID of the mounted resource (corresponding cloud service).

* `mount_share_path` - (Required, String) Specifies the remote mount path. Example: 192.168.0.12:/data.

* `local_mount_path` - (Required, String) Specifies the function access path.

The `custom_image` block supports:

* `url` - (Required, String) Specifies the URL of SWR image, the URL must start with `swr.`.

* `command` - (Optional, String) Specifies the startup commands of the SWR image.
  Multiple commands are separated by commas (,). e.g. `/bin/sh`.
  If this parameter is not specified, the entrypoint or CMD in the image configuration will be used by default.

* `args` - (Optional, String) Specifies the command line arguments used to start the SWR image.
  If multiple arguments are separated by commas (,). e.g. `-args,value`.
  If this parameter is not specified, the CMD in the image configuration will be used by default.

* `working_dir` - (Optional, String) Specifies the working directory of the SWR image.
  If not specified, the default value is `/`.
  Currently, the folder path can only be set to `/` and it cannot be created or modified.

The `versions` block supports:

* `name` - (Required, String) Specifies the version name.

  -> Currently, only supports the management of the default version (**latest**).

* `aliases` - (Optional, List) Specifies the aliases management for specified version.
  The `aliases` structure is documented below.

The `aliases` block supports:

* `name` - (Required, String) Specifies the name of the version alias.

* `description` - (Optional, String) Specifies the description of the version alias.

The `reserved_instances` block supports:

* `qualifier_type` - (Required, String) Specifies qualifier type of reserved instance. The valid values are as follows:
    + **version**
    + **alias**

  -> Reserved instances cannot be configured for both a function alias and the corresponding version. For example,
  if the alias of the `latest` version is `1.0` and reserved instances have been configured for this version,
  no more instances can be configured for alias `1.0`.

* `qualifier_name` - (Required, String) Specifies the version name or alias name.

* `count` - (Required, Int) Specifies the number of reserved instance.
  The valid value ranges from `0` to `1,000`.
  If this parameter is set to `0`, the reserved instance will not run.

* `idle_mode` - (Optional, Bool) Specifies whether to enable the idle mode. The default value is `false`.
  If this parameter is enabled, reserved instances are initialized and the mode change needs some time to take effect.
  You will still be billed at the price of reserved instances for non-idle mode in this period.

* `tactics_config` - (Optional, List) Specifies the auto scaling policies for reserved instance.
  The `tactics_config` structure is documented below.

The `tactics_config` block supports:

* `cron_configs` - (Optional, List) Specifies the list of scheduled policy configurations.
  The `cron_configs` structure is documented below.

The `cron_configs` block supports:

* `name` - (Required, String) Specifies the name of scheduled policy configuration.
  The valid length is limited from `1` to `60` characters, only letters, digits, hyphens (-), and underscores (_) are allowed.
  The name must start with a letter and ending with a letter or digit.

* `cron` - (Required, String) Specifies the cron expression.

* `count` - (Required, Int) Specifies the number of reserved instance to which the policy belongs.
  The valid value ranges from `0` to `1,000`.

  -> The number of reserved instances must be greater than or equal to the number of reserved instances in the basic configuration.

* `start_time` - (Required, Int) Specifies the effective timestamp of policy. The unit is `s`, e.g. **1740560074**.

* `expired_time` - (Required, Int) Specifies the expiration timestamp of the policy. The unit is `s`, e.g. **1740560074**.

<a name="function_network_controller"></a>
The `network_controller` block supports:

* `trigger_access_vpcs` - (Optional, List) Specifies the configuration of the VPCs that can trigger the function.
  The [trigger_access_vpcs](#function_network_controller_trigger_access_vpcs) structure is documented below.

* `disable_public_network` - (Optional, Bool) Specifies whether to disable the public network access.

<a name="function_network_controller_trigger_access_vpcs"></a>
The `trigger_access_vpcs` block supports:

* `vpc_id` - (Optional, String) Specifies the ID of the VPC that can trigger the function.

* `vpc_name` - (Optional, String) Specifies the name of the VPC that can trigger the function.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, consist of `urn` and current `version`, the format is `<urn>:<version>`.

* `region` - The region in which function graph resource is created.

* `extend_config` - The extended configuration.

* `apig_route_enable` - Whether to configure gateway routing rules.

* `allow_ephemeral_storage` - Indicates whether ephemeral storage can be configured.

* `func_mounts/status` - The status of file system.

* `urn` - Uniform Resource Name.

* `version` - The version of the function.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Functions can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_fgs_function_v2.test <id>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attributes are:
`app`, `func_code`, `agency`, `tags"`.
It is generally recommended running `terraform plan` after importing a function.
You can then decide if changes should be applied to the function, or the resource definition should be updated to align
with the function. Also you can ignore changes as below.

```hcl
resource "opentelekomcloud_fgs_function_v2" "test" {
  lifecycle {
    ignore_changes = [
      app, func_code, agency, tags,
    ]
  }
}
```
