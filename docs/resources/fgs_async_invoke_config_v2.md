---
subcategory: "FunctionGraph"
---

Up-to-date reference of API arguments for FGS you can get at
`https://docs.otc.t-systems.com/function-graph/api-ref/apis/asynchronous_execution_notification/index.html`.

# opentelekomcloud_fgs_async_invoke_config_v2

Manages the configuration of the asynchronous invocation within OpenTelekomCloud.

-> A function only supports configuring one resource.

## Example Usage

```hcl
variable "function_urn" {}
variable "bucket_name" {}
variable "topic_urn" {}

resource "opentelekomcloud_fgs_async_invoke_config_v2" "test" {
  function_urn                   = var.function_urn
  max_async_event_age_in_seconds = 3500
  max_async_retry_attempts       = 2

  on_success {
    destination = "OBS"
    param = jsonencode({
      bucket  = var.bucket_name
      prefix  = "/success"
      expires = 5
    })
  }

  on_failure {
    destination = "SMN"
    param = jsonencode({
      topic_urn = var.topic_urn
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `function_urn` - (Required, String, ForceNew) Specifies the function URN to which the asynchronous invocation belongs.
  Changing this will create a new resource.

* `max_async_event_age_in_seconds` - (Required, Int) Specifies the maximum validity period of a message.

* `max_async_retry_attempts` - (Required, Int) Specifies the maximum number of retry attempts to be made if
  asynchronous invocation fails.

* `on_success` - (Optional, List) Specifies the target to be invoked when a function is successfully executed.
  The [object](#functiongraph_destination_config) structure is documented below.

* `on_failure` - (Optional, List) Specifies the target to be invoked when a function fails to be executed due to a
  system error or an internal error.
  The [object](#functiongraph_destination_config) structure is documented below.

<a name="functiongraph_destination_config"></a>
The `on_success` and the `on_failure` blocks support:

* `destination` - (Required, String) Specifies the object type.
  The valid values are as follows:
    + **OBS**
    + **SMN**
    + **DIS**
    + **FunctionGraph**

* `param` - (Required, String) Specifies the parameters (map object in JSON format) corresponding to the target service.
    + The **OBS** objects include: `bucket` (bucket name), `prefix` (object directory prefix) and `expires` (object
      expiration time, the valid value ranges from `0` to `365`. If the value is `0`, the object will not expire.).
    + The **SMN** objects include: `topic_urn`.
    + The **DIS** objects include: `stream_name`.
    + The **FunctionGraph** objects include: `func_urn` (function URN).

-> If you enable the destination function, you must ensure that the agent contains the operation authority of the
corresponding service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `region` - Specifies the region in which the asynchronous invocation is configured.

## Import

The configurations can be imported using their related `function_urn`, e.g.

```bash
$ terraform import opentelekomcloud_fgs_async_invoke_config_v2.test <function_urn>
```
