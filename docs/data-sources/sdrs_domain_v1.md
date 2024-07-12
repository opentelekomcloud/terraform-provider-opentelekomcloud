---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_sdrs_domain_v1"
sidebar_current: "docs-opentelekomcloud-datasource-sdrs-domain-v1"
description: |-
  Get the ID of an available SDRS domain from OpenTelekomCloud
---

Up-to-date reference of API arguments for SDRS domain you can get at
[documentation portal](https://docs.otc.t-systems.com/storage-disaster-recovery-service/api-ref/sdrs_apis/active-active_domain/querying_an_active-active_domain.html#sdrs-05-0301)

# opentelekomcloud_sdrs_domain_v1

Use this data source to get the ID of an available OpenTelekomcloud SDRS domain.

~>
    OTC supports a single ``active-active domain`` with default name ``domain_001``.

## Example Usage

~>
  **Result of both examples will be the same.**

### Querying ``active-active domain`` with ``name`` parameter.

```hcl
data "opentelekomcloud_sdrs_domain_v1" "dom_1" {
  name = "domain_001"
}

```

### Querying ``active-active domain`` without ``name`` parameter.


```hcl
data "opentelekomcloud_sdrs_domain_v1" "dom_1" {}

```

## Argument Reference

* `name` - (Optional) Specifies the name of an active-active domain.
  This parameter serves as filter for querying ``active-active`` domains and can be skipped in current version.

## Attributes Reference

`id` is set to the ID of the active-active domain. In addition, the following attributes are exported:

* `name` - See Argument Reference above.

* `description` - Specifies the description of an active-active domain.
