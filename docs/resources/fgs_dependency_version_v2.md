---
subcategory: "FunctionGraph"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_fgs_dependency_version_v2"
sidebar_current: "docs-opentelekomcloud-resource-fgs-dependency-version-v2"
description: |-
  Manages a custom dependency version within OpenTelekomCloud.
---

# opentelekomcloud_fgs_dependency_version_v2

Manages a custom dependency version within OpenTelekomCloud.

## Example Usage

### Create a custom dependency version using an OBS bucket path where the ZIP file is located

```hcl
variable "dependency_name" {}
variable "custom_dependency_location" {}

resource "opentelekomcloud_fgs_dependency_version_v2" "test" {
  name    = var.dependency_name
  runtime = "Python3.6"
  link    = var.custom_dependency_location
}
```

### Create a custom dependency version directly from ZIP file

```hcl
variable "dependency_name" {}

resource "opentelekomcloud_fgs_dependency_version_v2" "test" {
  name    = var.dependency_name
  runtime = "Python3.6"
  file    = filebase64("${path.module}/data/requirements.zip")
}
```

## Argument Reference

* `runtime` - (Required, String, ForceNew) Specifies the runtime of the custom dependency version.
  The valid values are as follows:
  + **Java8**
  + **Java11**
  + **Node.js6.10**
  + **Node.js8.10**
  + **Node.js10.16**
  + **Node.js12.13**
  + **Node.js14.18**
  + **Python2.7**
  + **Python3.6**
  + **Python3.9**
  + **Go1.8**
  + **Go1.x**
  + **C#(.NET Core 2.0)**
  + **C#(.NET Core 2.1)**
  + **C#(.NET Core 3.1)**
  + **Custom**
  + **PHP 7.3**
  + **http**

  Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the custom dependency package to which the version
  belongs.
  The name can contain a maximum of `96` characters and must start with a letter and end with a letter or digit.
  Only letters, digits, underscores (_), periods (.), and hyphens (-) are allowed.
  Changing this will create a new resource.

* `link` - (Optional, String, ForceNew) Specifies the OBS bucket path where the dependency package is located.
  The OBS object URL must be in ZIP format, such as
  `https://test-bucket.obs.eu-de.otc.t-systems.com/index.zip`.
  Either `link` or `file` must be specified.
  Changing this will create a new resource.

  -> A link can only be used to create at most one dependency package.

* `file` - (Optional, String, ForceNew) Specifies the file contents in the file stream format and must be a ZIP file encoded using Base64.
  The size of the ZIP file cannot exceed 40 MB. For a larger file, upload it through OBS.
  Either `link` or `file` must be specified.

* `description` - (Optional, String, ForceNew) Specifies the description of the custom dependency version.
  The description can contain a maximum of `512` characters.
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, consists of dependency ID and version number, separated by a slash (/).
  The format is `<name>/<version>`.

* `version` - The dependency package version.

* `file_name` - The dependency file name.

* `download_link` - The temporary download link of a dependency file.

* `owner` - The dependency owner, **public** indicates a public dependency.

* `etag` - The unique ID of the dependency.

* `size` - The dependency size, in bytes.

* `dependency_id` - The ID of the dependency package corresponding to the version.

## Import

Dependency version can be imported using `name` and the `version` number, separated by a slash (/), e.g.

```bash
$ terraform import opentelekomcloud_fgs_dependency_version_v2.test <name>/<version>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `link`, `file`.
It is generally recommended running `terraform plan` after importing a dependency package.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the dependency package. Also you can ignore changes as below.

```hcl
resource "opentelekomcloud_fgs_dependency_version_v2" "test" {

  lifecycle {
    ignore_changes = [
      link,
      file
    ]
  }
}
```
