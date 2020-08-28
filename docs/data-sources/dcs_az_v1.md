---
subcategory: "DCS"
---

# opentelekomcloud_dcs_az_v1

Use this data source to get the ID of an available DCS AZ from OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_dcs_az_v1" "az1" {
  name = "AZ1"
  port = "8004"
  code = "sa-chile-1a"
}
```

## Argument Reference

* `name` - (Required) Indicates the name of an AZ.

* `code` - (Optional) Indicates the code of an AZ.

* `port` - (Required) Indicates the port number of an AZ.


## Attributes Reference

`id` is set to the ID of the found AZ. In addition, the following attributes are exported:

* `name` - See Argument Reference above.

* `code` - See Argument Reference above.

* `port` - See Argument Reference above.
