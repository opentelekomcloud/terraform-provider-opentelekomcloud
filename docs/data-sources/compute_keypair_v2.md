---
subcategory: "Elastic Cloud Server (ECS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_keypair_v2"
sidebar_current: "docs-opentelekomcloud-datasource-compute-keypair-v2"
description: |-
Get ECS keypair details from OpenTelekomCloud
---

Up-to-date reference of API arguments for ECS keypair you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-cloud-server/api-ref/native_openstack_nova_apis/key_and_password_management/querying_ssh_key_pairs.html#en-us-topic-0020212676)

# opentelekomcloud_compute_keypair_v2

Use this data source to get details about Compute SSH key pairs from OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_compute_keypair_v2" "kp_1" {
  name       = "key_1"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIALRzbIOR9HUYNwfKtII/et98eGXDJhf8YxHf9BtRdAU"
}

data "opentelekomcloud_compute_keypair_v2" "key_1" {
  name = "key_1"

  depends_on = [opentelekomcloud_compute_keypair_v2.kp_1]
}
```

## Argument Reference

* `name` - (Optional, ForceNew, String) The name of the keypair.

* `name_regex` - (Optional, ForceNew, String) A regex string to apply to the keypairs list.
  This allows more advanced filtering not supported from the OpenTelekomCloud API.
  This filtering is done locally on what OpenTelekomCloud returns.

## Attributes Reference

All the argument attributes are also exported as result attributes.

* `public_key` - It gives the information about the public key in the key pair.

* `fingerprint` - It is the fingerprint information about the key pair.

* `name` - See Argument Reference above.

* `user_id` - The user id of the owner of the key pair. Not filled by API now.
