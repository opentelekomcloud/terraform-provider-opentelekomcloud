---
subcategory: "Bare Metal Server (BMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_bms_keypairs_v2"
sidebar_current: "docs-opentelekomcloud-datasource-compute-bms-keypairs-v2"
description: |-
Get details about SSH key pairs of BMSs from OpenTelekomCloud
---

Up-to-date reference of API arguments for BMSs SSH key pairs you can get at
[documentation portal](https://docs.otc.t-systems.com/bare-metal-server/api-ref/native_openstack_nova_v2.1_apis/bms_ssh_key_pair_management/querying_ssh_key_pairs_native_openstack_api.html#en-us-topic-0060384658)

# opentelekomcloud_compute_bms_keypairs_v2

Use this data source to get details about SSH key pairs of BMSs from OpenTelekomCloud.

## Example Usage

```hcl
variable "keypair_name" {}

data "opentelekomcloud_compute_bms_keypairs_v2" "query_bms_keypair" {
  name = var.keypair_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `name` - (Required) It is the key pair name.

## Attributes Reference

All of the argument attributes are also exported as result attributes.

* `public_key` - It gives the information about the public key in the key pair.

* `fingerprint` - It is the fingerprint information about the key pair.
