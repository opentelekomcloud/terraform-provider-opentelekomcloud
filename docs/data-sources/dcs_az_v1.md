---
subcategory: "Distributed Cache Service (DCS)"
---

# opentelekomcloud_dcs_az_v1

Use this data source to get the ID of an available DCS AZ from OpenTelekomCloud.

## Example Usage

### Query AZ `id` by providing `name` argument

```hcl
data "opentelekomcloud_dcs_az_v1" "az1" {
  name = "eu-de-01"
}
```

### Query AZ `id` by providing `port` and `code` arguments

```hcl
data "opentelekomcloud_dcs_az_v1" "az2" {
  port = "8003"
  code = "eu-de-02"
}
```

### Query AZ `id` by providing all arguments

```hcl
data "opentelekomcloud_dcs_az_v1" "az2" {
  name = "eu-de-02"
  port = "8003"
  code = "eu-de-02"
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
