---
subcategory: "Cloud Connect (CC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cc_central_network_policy_apply_v3"
sidebar_current: "docs-opentelekomcloud-resource-cc-central-network-policy-apply-v3"
description: |-
  Applies a Cloud Connect Central Network Policy within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Cloud Connect Central Network Policy you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-connect/api-ref/api/central_network_policies/index.html)

# opentelekomcloud_cc_central_network_policy_apply_v3

Applies a Cloud Connect (CC) central network policy within OpenTelekomCloud.

A central network can hold several policies, but only one can be active. This resource applies the policy
referenced by `policy_id` to the central network and waits until it becomes effective.

-> Deleting this resource does not delete the policy. Instead, it reverts the central network to its default
policy (the initial policy with version `1`, which has no associated enterprise routers).

## Example Usage

```hcl
variable "central_network_id" {}

resource "opentelekomcloud_cc_central_network_policy_v3" "test" {
  central_network_id = var.central_network_id

  er_instances {
    project_id           = var.project_id
    region_id            = "eu-de"
    enterprise_router_id = var.enterprise_router_id
  }
}

resource "opentelekomcloud_cc_central_network_policy_apply_v3" "test" {
  central_network_id = var.central_network_id
  policy_id          = opentelekomcloud_cc_central_network_policy_v3.test.id
}
```

## Argument Reference

The following arguments are supported:

* `central_network_id` - (Required, String, ForceNew) The ID of the central network the policy is applied to.
  Changing this parameter will create a new resource.

* `policy_id` - (Required, String) The ID of the central network policy to apply.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. Equal to the central network ID.

* `region` - The region in which the central network policy is applied.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The applied central network policy can be imported using the `central_network_id` and `policy_id`, separated by
a slash, e.g.

```sh
$ terraform import opentelekomcloud_cc_central_network_policy_apply_v3.test <central_network_id>/<policy_id>
```
