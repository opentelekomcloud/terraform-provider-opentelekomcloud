---
subcategory: "FunctionGraph"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_fgs_publish_version_v2"
sidebar_current: "docs-opentelekomcloud-resource-fgs-publish-version-v2"
description: |-
  Manages an FGS Publish Version resource within T-Cloud Public (formerly OpenTelekomCloud).
---

Up-to-date reference of API arguments for FGS versions you can get at
[documentation portal](https://docs.otc.t-systems.com/function-graph/api-ref/api/versions_and_aliases/index.html)

# opentelekomcloud_fgs_publish_version_v2

Manages the publishing of a function version within T-Cloud Public (formerly OpenTelekomCloud).

## Example Usage

```hcl
variable "function_urn" {}

resource "opentelekomcloud_fgs_publish_version_v2" "test" {
  function_urn = var.function_urn
  description  = "test version"
}
```

## Argument Reference

The following arguments are supported:

* `function_urn` - (Required, String, ForceNew) Specifies the URN of the function.
  Changing this will create a new resource.

* `digest` - (Optional, String, ForceNew) Specifies the digest of the function code.
  Changing this will create a new resource.

* `version` - (Required, String, ForceNew) Specifies the version of the function.
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in the format of `function_urn/version`.

* `func_name` - Specifies the name of the function.

* `last_modified` - Specifies the last modified time of the version.
