---
subcategory: "Cloud Connect (CC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cc_central_network_v3"
sidebar_current: "docs-opentelekomcloud-resource-cc-central-network-v3"
description: |-
  Manages a Cloud Connect Central Network resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Cloud Connect Central Network you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-connect/api-ref/api/central_networks/index.html)

# opentelekomcloud_cc_central_network_v3

Manages a Cloud Connect (CC) central network resource within OpenTelekomCloud.

A central network lets you connect enterprise routers from different regions and accounts so that they can
communicate over a fully meshed network. It is the top-level container for central network policies and
connections.

-> The central network APIs are account-level (domain-scoped). Make sure the credentials used by the provider
have permission to manage Cloud Connect resources.

## Example Usage

```hcl
resource "opentelekomcloud_cc_central_network_v3" "test" {
  name        = "central-network-demo"
  description = "managed by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the central network.
  The name can contain 1 to 64 characters, only letters, digits, underscores (_) and hyphens (-) are allowed.

* `description` - (Optional, String) The description of the central network.
  The description can contain a maximum of 255 characters, and the angle brackets (`<` and `>`) are not allowed.

* `enterprise_project_id` - (Optional, String, ForceNew) The ID of the enterprise project that the central
  network belongs to. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format. Equal to the central network ID.

* `state` - The status of the central network. The value can be **AVAILABLE**, **CREATING**, **UPDATING**,
  **FAILED**, **DELETING**, **DELETED** or **RESTORING**.

* `default_plane_id` - The ID of the default central network plane.

* `domain_id` - The ID of the account that the central network belongs to.

* `created_at` - The creation time of the central network.

* `updated_at` - The latest update time of the central network.

* `region` - The region in which the central network is managed.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The central network can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_cc_central_network_v3.test 0ce123456a00f2591fabc00385ff1234
```
