---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_signature_associate_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-signature-associate-v2"
description: |-
Manages a APIGW Signature Associate resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for API Gateway signature associate service you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/binding_unbinding_signature_keys/index.html)

# opentelekomcloud_apigw_signature_associate_v2

Use this resource to bind the APIs to the signature within OpenTelekomCloud.

-> A signature can only create one `opentelekomcloud_apigw_signature_associate_v2` resource.
   And a published ID for API can only bind a signature.

## Example Usage

```hcl
variable "gateway_id" {}
variable "signature_id" {}
variable "api_publish_ids" {
  type = list(string)
}

resource "opentelekomcloud_apigw_signature_associate_v2" "test" {
  instance_id  = var.gateway_id
  signature_id = var.signature_id
  publish_ids  = var.api_publish_ids
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the APIs and the
  signature belong.
  Changing this will create a new resource.

* `signature_id` - (Required, String, ForceNew) Specifies the signature ID for APIs binding.
  Changing this will create a new resource.

* `publish_ids` - (Required, List) Specifies the publish IDs corresponding to the APIs bound by the signature.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID. The format is `<gateway_id>/<signature_id>`.

* `region` - Region where the signature and the APIs are located.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
* `update` - Default is 3 minutes.
* `delete` - Default is 3 minutes.

## Import

Associate resources can be imported using their `signature_id` and the APIGW dedicated gateway ID to which the signature
belongs, separated by a slash, e.g.

```bash
$ terraform import opentelekomcloud_apigw_signature_associate_v2.test <gateway_id>/<signature_id>
```
