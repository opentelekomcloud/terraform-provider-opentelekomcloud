---
subcategory: "Relational Database Service (RDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rds_public_ip_associate_v3"
sidebar_current: "docs-opentelekomcloud-resource-rds-public-ip-associate-v3"
description: |-
  Manages an RDS Public Ip association v3 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for RDS public ip association you can get at
[documentation portal](https://docs.otc.t-systems.com/relational-database-service/api-ref/api_v3_recommended/db_instance_management/binding_and_unbinding_an_eip.html#rds-05-0009)

# opentelekomcloud_rds_public_ip_associate_v3

Associates a public IP to an RDS instance.

## Example Usage

### Bind

```hcl
variable instance_id {}

resource "opentelekomcloud_rds_public_ip_associate_v3" "public_ip" {
  instance_id  = var.instance_id
  public_ip    = opentelekomcloud_compute_floatingip_v2.eip_2.address
  public_ip_id = opentelekomcloud_compute_floatingip_v2.eip_2.id
}

resource "opentelekomcloud_compute_floatingip_v2" "eip_1" {}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the RDS instance ID.

* `public_ip` - (Required, String) Specifies the EIP to be bound. The value must be in the standard IP address format.

* `public_ip_id` - (Required, String) Specifies the EIP ID. The value must be in the standard UUID format.
