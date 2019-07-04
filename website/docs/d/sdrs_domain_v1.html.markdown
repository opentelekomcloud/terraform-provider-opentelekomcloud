---
layout: "opentelekomcloud"
page_title: "Opentelekomcloud: opentelekomcloud_sdrs_domain_v1"
sidebar_current: "docs-opentelekomcloud-datasource-sdrs-domain-v1"
description: |-
  Get information on an Opentelekomcloud SDRS Active-Active Domain.
---

# opentelekomcloud\_sdrs\_domain_v1

Use this data source to get the ID of an available Opentelekomcloud SDRS domain.

## Example Usage

```hcl

data "opentelekomcloud_sdrs_domain_v1" "dom_1" {
	name = "domain_001"
}

```

## Argument Reference

* `name` - (Optional) Specifies the name of an active-active domain.

## Attributes Reference

`id` is set to the ID of the active-active domain. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `description` - Specifies the description of an active-active domain.
