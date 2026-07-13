---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_enterprise_project_action_v1"
sidebar_current: "docs-opentelekomcloud-resource-enterprise-project-action-v1"
description: |-
  Use this resource to operate the enterprise project within T-Cloud Public (former OpenTelekomCloud).
---

An up-to-date reference of the EPS API is available at the
[documentation portal](https://docs.otc.t-systems.com/enterprise-project-service/api-ref/enterprise_project_management_apis/enterprise_project_management/index.html#en-us-topic-0133321992)

# opentelekomcloud_enterprise_project_action_v1

Use this resource to operate the enterprise project within T-Cloud Public (former OpenTelekomCloud).

-> This is a one-time action resource. Destroying it does not reverse the action or change the enterprise project; it
   only removes the resource from Terraform state.

## Example Usage

```hcl
variable "enterprise_project_id" {}

resource "opentelekomcloud_enterprise_project_action_v1" "test" {
  enterprise_project_id = var.enterprise_project_id
  action                = "disable"
}
```

## Argument Reference

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the ID of enterprise project to be operated.

* `action` - (Required, String, NonUpdatable) Specifies the action type.
  The valid values are as follows:
  + **enable**
  + **disable**

* `enable_force_new` - (Optional, String) Controls how changes to `enterprise_project_id` or `action` are handled.
  Set this argument to `"true"` to replace the resource and execute the new action. If set to `"false"`, Terraform
  returns an error when either non-updatable argument changes. If omitted, the provider-level `enable_force_new`
  setting is used. Valid values are `"true"` and `"false"`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A generated UUID used to track the action resource in Terraform state.
