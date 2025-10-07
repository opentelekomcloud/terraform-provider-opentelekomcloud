---
subcategory: "FunctionGraph"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_fgs_functions_v2"
sidebar_current: "docs-opentelekomcloud-datasource-fgs-functions-v2"
description: |-
  Get details about FGS functions from OpenTelekomCloud
---

# opentelekomcloud_fgs_functions_v2

Use this data source to filter FGS functions within OpenTelekomCloud.

## Example Usage

### Obtain all public functions

```hcl
data "opentelekomcloud_fgs_functions_v2" "test" {}
```

## Argument Reference

The following arguments are supported:

* `package_name` - (Optional, String) Specifies the package name used to query the functions.

* `urn` - (Optional, String) Specifies the function URN used to query the specified function.

* `name` - (Optional, String) Specifies the function name used to query the specified function.

* `runtime` - (Optional, String) Specifies the dependency package runtime used to query the functions.
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
  + **Custom**
  + **PHP7.3**
  + **http**
  + **Custom Image**
  + **Cangjie1.0**

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the functions belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `functions` - All functions that match the filter parameters.
  The [functions](#fgs_functions) structure is documented below.

<a name="fgs_functions"></a>
The `functions` block supports:

* `name` - The function name.

* `urn` - The function URN.

* `package` - The package name that the function used.

* `runtime` - The dependency package runtime of the function.

* `timeout` - The timeout interval of the function.

* `handler` - The entry point of the function.

* `memory_size` - The memory size allocated to the function, the unit is MB.

* `code_type` - The function code type.
  + **inline**: inline code.
  + **zip**: ZIP file.
  + **jar**: JAR file or java functions.
  + **obs**: function code stored in an OBS bucket.

* `code_url` - The code URL.

* `code_filename` - The name of the function file.

* `user_data` - The custom user data (key/value pairs) defined for the function.

* `encrypted_user_data` - The custom user data (key/value pairs) defined to be encrypted for the function.

* `version` - The function version.

* `agency` - The IAM agency name for the function configuration.

  -> The configuration agency name that used to create a trigger to access the relevant service, such as DMS and DIS.

* `app_agency` - The IAM agency name for the function execution.

  -> The execution agency name that used to obtain the Token or AK/SK for accessing other cloud services.

* `description` - The description of the function.

* `vpc_id` - The VPC ID to which the function belongs.

* `network_id` - The network ID of subnet to which the function belongs.

* `max_instance_num` - The maximum number of instances for a single function.

* `initializer_handler` - The initializer of the function.

* `initializer_timeout` - The maximum duration the function can be initialized.

* `enterprise_project_id` - The enterprise project ID to which the function belongs.

* `log_group_id` - The LTS log group ID.

* `log_stream_id` - The LTS log stream ID.

* `functiongraph_version` - The functionGraph version.

* `region` - The functionGraph region.
