---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_bms_keypairs_v2"
sidebar_current: "docs-opentelekomcloud-compute-bms-keypairs-v2"
description: |-
   Used to query SSH key pairs
---

# Data Source: opentelekomcloud_compute_bms_keypairs_v2

`opentelekomcloud_compute_bms_keypairs_v2` used to query SSH key pairs.


## Example Usage

```hcl

    variable "keypair_name" {}

    data "opentelekomcloud_compute_bms_keypairs_v2" "Query_BMS_keypair" 
    {
        name = "${var.keypair_name}"
    }
    
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `name` - (Required) - It is the key pair name.

## Attributes Reference

All of the argument attributes are also exported as result attributes. 

* `public_key` - It gives the information about the public key in the key pair.

* `fingerprint` - It is the fingerprint information about the key pair.
