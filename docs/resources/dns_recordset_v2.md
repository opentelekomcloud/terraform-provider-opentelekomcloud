---
subcategory: "Domain Name Service (DNS)"
---

Up-to-date reference of API arguments for DNS recordset you can get at
`https://docs.otc.t-systems.com/domain-name-service/api-ref/apis/record_set_management`.

# opentelekomcloud_dns_recordset_v2

Manages a DNS record set in the OpenTelekomCloud DNS Service.

## Example Usage

### Automatically detect the correct network

```hcl
resource "opentelekomcloud_dns_zone_v2" "example_zone" {
  name        = "example.com."
  email       = "email2@example.com"
  description = "a zone"
  ttl         = 6000
  type        = "public"
}

resource "opentelekomcloud_dns_recordset_v2" "rs_example_com" {
  zone_id     = opentelekomcloud_dns_zone_v2.example_zone.id
  name        = "rs.example.com."
  description = "An example record set"
  ttl         = 3000
  type        = "A"
  records     = ["10.0.0.1"]
}

resource "opentelekomcloud_dns_recordset_v2" "rs_txt_example" {
  zone_id     = opentelekomcloud_dns_zone_v2.zone_1.id
  name        = "%[1]s"
  type        = "TXT"
  description = "a record set"
  ttl         = 300
  records     = ["v=spf1 include:my.example.try.com -all"]
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required) The ID of the zone in which to create the record set.
  Changing this creates a new DNS  record set.

* `name` - (Required) The name of the record set. Changing this creates a new DNS  record set.

-> **Note:** The `.` at the end of the name.

* `type` - (Required) The type of record set. Examples: "A", "MX".
  Changing this creates a new DNS  record set.

* `ttl` - (Optional) The time to live (TTL) of the record set.

* `description` - (Optional) A description of the  record set.

* `records` - (Required) An array of DNS records.

* `tags` - (Optional) The key/value pairs to associate with the zone.

* `value_specs` - (Optional) Map of additional options. Changing this creates a
  new record set.

->
If all `zone_id`, `type`, `name` and `ttl` duplicate the existing DNS record set value,
the new record set won't be managed by the Terraform.
DNS `recordset` resource will be marked as `shared.`

If `type="TXT"` records should pass as plain text without quotation, look at `rs_txt_example`.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `type` - See Argument Reference above.

* `ttl` - See Argument Reference above.

* `records` - See Argument Reference above.

* `description` - See Argument Reference above.

* `tags` - See Argument Reference above.

* `zone_id` - See Argument Reference above.

* `value_specs` - See Argument Reference above.

## Import

This resource can be imported by specifying the zone ID and recordset ID,
separated by a forward slash.

```sh
terraform import opentelekomcloud_dns_recordset_v2.recordset_1 <zone_id>/<recordset_id>
```

Imported key pairs are considered to be not shared.
