---
subcategory: "Cloud Container Instance (CCI)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cci_namespace_v2"
description: |-
  Manages a CCI v2 Namespace resource within OpenTelekomCloud.
---

# opentelekomcloud_cci_namespace_v2

Manages a CCI v2 namespace resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "name" {}

resource "opentelekomcloud_cci_namespace_v2" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, NonUpdatable) Specifies the unique name of the namespace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is also the name of the namespace.

* `api_version` - The API version of the namespace.

* `kind` - The kind of the namespace.

* `annotations` - The annotations of the namespace.

* `labels` - The labels of the namespace.

* `creation_timestamp` - The creation timestamp of the namespace.

* `resource_version` - The resource version of the namespace.

* `uid` - The uid of the namespace.

* `finalizers` - The finalizers of the namespace.

* `status` - The status of the namespace.

* `region` - The region of the namespace

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 3 minutes.

## Import

The CCI v2 namespace can be imported using `name`, e.g.

```bash
$ terraform import opentelekomcloud_cci_namespace_v2.test <name>
```
