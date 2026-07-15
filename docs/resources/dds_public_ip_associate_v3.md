---
subcategory: "Document Database Service (DDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dds_public_ip_associate_v3"
sidebar_current: "docs-opentelekomcloud-resource-dds-public-ip-associate-v3"
description: |-
  Manages a DDS Public IP association v3 resource within T Cloud Public (formerly OpenTelekomCloud).
---

Up-to-date reference of API arguments for DDS public ip association you can get at
[documentation portal](https://docs.otc.t-systems.com/document-database-service/api-ref/apis_v3.0_recommended/db_instance_management/binding_an_eip.html)

# opentelekomcloud_dds_public_ip_associate_v3

Associates a public IP to a DDS node.

## Example Usage

```hcl
variable node_id {}

resource "opentelekomcloud_dds_public_ip_associate_v3" "public_ip" {
  node_id      = var.node_id
  public_ip    = opentelekomcloud_compute_floatingip_v2.eip.address
  public_ip_id = opentelekomcloud_compute_floatingip_v2.eip.id
}

resource "opentelekomcloud_compute_floatingip_v2" "eip" {}
```

## Argument Reference

The following arguments are supported:

* `node_id` - (Required, String, ForceNew) Specifies the DDS node ID.

* `public_ip` - (Required, String) Specifies the EIP to be bound. The value must be in the standard IP address format.

* `public_ip_id` - (Required, String) Specifies the EIP ID. The value must be in the standard UUID format.
