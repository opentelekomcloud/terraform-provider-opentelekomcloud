---
subcategory: "Storage Disaster Recovery Service (SDRS)"
---

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
