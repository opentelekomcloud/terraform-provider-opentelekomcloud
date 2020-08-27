# Data Source: opentelekomcloud_compute_availability_zones_v2

`opentelekomcloud_compute_availability_zones_v2` is used to get a list of availability zones from OpenTelekomCloud

## Example Usage

```hcl
data "opentelekomcloud_compute_availability_zones_v2" "zones" { }
```

## Argument Reference

* `region` - (Optional) The `region` to fetch availability zones from, defaults to the provider's `region`
* `state` - (Optional) The `state` of the availability zones to match, default ("available").


## Attributes Reference

`id` is set to hash of the returned zone list. In addition, the following attributes
are exported:

* `names` - The names of the availability zones, ordered alphanumerically, that match the queried `state`
