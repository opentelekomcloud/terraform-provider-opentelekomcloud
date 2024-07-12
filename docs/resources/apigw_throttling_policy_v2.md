---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_throttling_policy_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-throttling-policy-v2"
description: |-
  Manages a APIGW Throttling Policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for API Gateway throttling policy service you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/index.html)

# opentelekomcloud_apigw_throttling_policy_v2

API Gateway (APIG) is a high-performance, high-availability, and high-security API hosting service that helps you build,
manage, and deploy APIs at any scale.
With just a few clicks, you can integrate internal systems, and selectively expose capabilities with minimal costs and risks.

## Example Usage

```hcl
resource "opentelekomcloud_apigw_throttling_policy_v2" "policy" {
  instance_id       = opentelekomcloud_apigw_gateway_v2.gateway.id
  name              = "throttling policy"
  type              = "API-shared"
  period            = 10
  period_unit       = "MINUTE"
  max_api_requests  = 70
  max_user_requests = 45
  max_app_requests  = 45
  max_ip_requests   = 45
  description       = "Created by tf"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the throttling
  policy belongs.

* `name` - (Required, String) Specifies the name of the throttling policy.
  The valid length is limited from `3` to `64`, only English letters, digits and underscores (_) are
  allowed.

* `period` - (Required, Int) Specifies the period of time for limiting the number of API calls.
  This parameter applies with each of the API call limits: `max_api_requests`, `max_app_requests`, `max_ip_requests`
  and `max_user_requests`.

* `max_api_requests` - (Required, Int) Specifies the maximum number of times an API can be accessed within a specified
  period. The value of this parameter cannot exceed the default limit `200` TPS.

* `max_app_requests` - (Optional, Int) Specifies the maximum number of times the API can be accessed by an app within
  the same period.
  The value of this parameter must be less than or equal to the value of `max_user_requests`.

* `max_ip_requests` - (Optional, Int) Specifies the maximum number of times the API can be accessed by an IP address
  within the same period.
  The value of this parameter must be less than or equal to the value of `max_api_requests`.

* `max_user_requests` - (Optional, Int) Specifies the maximum number of times the API can be accessed by a user within
  the same period.
  The value of this parameter must be less than or equal to the value of `max_api_requests`.

* `type` - (Optional, String) Specifies the type of the request throttling policy.
  The valid values are as follows:
    + **API-based**: limiting the maximum number of times a single API bound to the policy can be called within the
      specified period.
    + **API-shared**: limiting the maximum number of times all APIs bound to the policy can be called within the specified
      period.

* `description` - (Optional, String) Specifies the description about the API throttling policy.
  The description contain a maximum of `255` characters and the angle brackets (< and >) are not allowed.

* `period_unit` - (Optional, String) Specifies the time unit for limiting the number of API calls.
  The valid values are **SECOND**, **MINUTE**, **HOUR** and **DAY**, defaults to **MINUTE**.

* `user_throttles` - (Optional, List) Specifies the array of one or more special throttling policies for IAM user limit.
  The `user_throttles` object structure is documented below.

* `app_throttles` - (Optional, List) Specifies the array of one or more special throttling policies for APP limit.
  The `app_throttles` object structure is documented below.

The `user_throttles` and `user_throttles` blocks support:

* `max_api_requests` - (Required, Int) Specifies the maximum number of times an API can be accessed within a specified
  period.

* `throttling_object_id` - (Required, String) Specifies the object ID which the special throttling policy belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the API throttling policy.

* `region` - The region where the throttling policy is located.

* `user_throttles` - The array of one or more special throttling policies for IAM user limit.
  The structure is documented below.

* `app_throttles` - The array of one or more special throttling policies for APP limit.
  The structure is documented below.

* `created_at` - The creation time of the throttling policy.

The `user_throttles` and `app_throttles` blocks support:

* `throttling_object_name` - The object name which the special user/application throttling policy belongs.

* `id` - ID of the special user/application throttling policy.

## Import

API Throttling Policies can be imported using their `name` and related dedicated instance ID, separated by a slash, e.g.

```shell
$ terraform import opentelekomcloud_apigw_throttling_policy_v2.test <instance_id>/<name>
```
