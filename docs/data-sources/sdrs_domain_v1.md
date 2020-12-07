---
subcategory: "Storage Disaster Recovery Service (SDRS)"
---

# opentelekomcloud_sdrs_domain_v1

Use this data source to get the ID of an available OpenTelekomcloud SDRS domain.

## Example Usage

```hcl
data "opentelekomcloud_sdrs_domain_v1" "dom_1" {
  name = "domain_001"
}

```

## Argument Reference

* `name` - (Optional) Specifies the name of an active-active domain.

## Attributes Reference

`id` is set to the ID of the active-active domain. In addition, the following attributes are exported:

* `name` - See Argument Reference above.

* `description` - Specifies the description of an active-active domain.
