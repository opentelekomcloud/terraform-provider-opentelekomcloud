---
subcategory: "Cloud Connect (CC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cc_central_network_policy_v3"
sidebar_current: "docs-opentelekomcloud-resource-cc-central-network-policy-v3"
description: |-
  Manages a Cloud Connect Central Network Policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Cloud Connect Central Network Policy you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-connect/api-ref/api/central_network_policies/index.html)

# opentelekomcloud_cc_central_network_policy_v3

Manages a Cloud Connect (CC) central network policy resource within OpenTelekomCloud.

A central network policy describes which enterprise routers belong to a central network and how their route
tables are associated. Creating a policy only stages it; use `opentelekomcloud_cc_central_network_policy_apply_v3`
to make it the active policy of the central network.

-> Central network policies are immutable. Any change to the arguments forces a new policy to be created.

## Example Usage

```hcl
variable "central_network_id" {}
variable "project_id" {}
variable "enterprise_router_a_id" {}
variable "enterprise_router_a_table_id" {}
variable "enterprise_router_b_id" {}
variable "enterprise_router_b_table_id" {}

resource "opentelekomcloud_cc_central_network_policy_v3" "test" {
  central_network_id = var.central_network_id

  er_instances {
    project_id           = var.project_id
    region_id            = "eu-de"
    enterprise_router_id = var.enterprise_router_a_id
  }

  er_instances {
    project_id           = var.project_id
    region_id            = "eu-nl"
    enterprise_router_id = var.enterprise_router_b_id
  }

  planes {
    associate_er_tables {
      project_id                 = var.project_id
      region_id                  = "eu-de"
      enterprise_router_id       = var.enterprise_router_a_id
      enterprise_router_table_id = var.enterprise_router_a_table_id
    }

    associate_er_tables {
      project_id                 = var.project_id
      region_id                  = "eu-nl"
      enterprise_router_id       = var.enterprise_router_b_id
      enterprise_router_table_id = var.enterprise_router_b_table_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `central_network_id` - (Required, String, ForceNew) The ID of the central network the policy belongs to.
  Changing this parameter will create a new resource.

* `er_instances` - (Optional, List, ForceNew) The list of the enterprise routers on the central network policy.
  Changing this parameter will create a new resource.
  The [er_instances](#policy_er_instances) structure is documented below.

* `planes` - (Optional, List, ForceNew) The list of the central network policy planes.
  Changing this parameter will create a new resource.
  The [planes](#policy_planes) structure is documented below.

<a name="policy_er_instances"></a>
The `er_instances` block supports:

* `project_id` - (Required, String, ForceNew) The project ID of the enterprise router.

* `region_id` - (Required, String, ForceNew) The region ID of the enterprise router.

* `enterprise_router_id` - (Required, String, ForceNew) The ID of the enterprise router.

<a name="policy_planes"></a>
The `planes` block supports:

* `associate_er_tables` - (Optional, List, ForceNew) The list of route tables associated with the central
  network policy. The [associate_er_tables](#policy_associate_er_tables) structure is documented below.

* `exclude_er_connections` - (Optional, List, ForceNew) The list of the enterprise router connections excluded
  from the central network policy. The [exclude_er_connections](#policy_exclude_er_connections) structure is
  documented below.

<a name="policy_associate_er_tables"></a>
The `associate_er_tables` block supports:

* `project_id` - (Required, String, ForceNew) The project ID of the enterprise router.

* `region_id` - (Required, String, ForceNew) The region ID of the enterprise router.

* `enterprise_router_id` - (Required, String, ForceNew) The ID of the enterprise router.

* `enterprise_router_table_id` - (Required, String, ForceNew) The ID of the enterprise router route table.

<a name="policy_exclude_er_connections"></a>
The `exclude_er_connections` block supports:

* `exclude_er_instances` - (Required, List, ForceNew) The list of enterprise routers that will not establish a
  connection with each other. The [exclude_er_instances](#policy_er_instances) structure is the same as the
  `er_instances` structure documented above.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format. Equal to the central network policy ID.

* `document_template_version` - The version of the central network policy document template.

* `is_applied` - Whether the central network policy is applied.

* `version` - The version of the central network policy.

* `state` - The status of the central network policy. The value can be **AVAILABLE**, **CANCELING**,
  **APPLYING**, **FAILED** or **DELETED**.

* `region` - The region in which the central network policy is managed.

## Import

The central network policy can be imported using the `central_network_id` and `id` (policy ID), separated by a
slash, e.g.

```sh
$ terraform import opentelekomcloud_cc_central_network_policy_v3.test <central_network_id>/<id>
```
