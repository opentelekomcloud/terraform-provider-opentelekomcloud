---
subcategory: "Domain Name Service (DNS)"
---

# opentelekomcloud_dns_zone_v2

Use this data source to get the ID of an available OpenStack DNS zone.

## Example Usage

```hcl
data "opentelekomcloud_dns_zone_v2" "zone_1" {
  name = "example.com"
}
```

## Argument Reference

* `zone_type` - (Optional) The type of the zone: `private` or `public`.

* `name` - (Optional) The name of the zone.

* `description` - (Optional) A description of the zone.

* `email` - (Optional) The email contact for the zone record.

* `status` - (Optional) The zone's status.

* `ttl` - (Optional) The time to live (TTL) of the zone.

* `tags` - (Optional) Tags map to be matched.
  An exact match will be performed. If the value starts with an
  asterisk (*), the string following the asterisk is fuzzy matched.

## Attributes Reference

`id` is set to the ID of the found zone. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `email` - See Argument Reference above.
* `zone_type` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `description` - See Argument Reference above.
* `status` - See Argument Reference above.
* `masters` - An array of master DNS servers.
* `created_at` - The time the zone was created.
* `updated_at` - The time the zone was last updated.
* `version` - The version of the zone.
* `serial` - The serial number of the zone.
* `pool_id` - The ID of the pool hosting the zone.
* `project_id` - The project ID that owns the zone.
