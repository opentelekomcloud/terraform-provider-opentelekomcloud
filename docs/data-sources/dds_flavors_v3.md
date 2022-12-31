---
subcategory: "Document Database Service (DDS)"
---

# opentelekomcloud_dds_flavors_v3

Use this data source to get info of available OpenTelekomCloud DDS flavors.

## Example Usage

```hcl
data "opentelekomcloud_dds_flavors_v3" "flavor" {
  engine_name = "DDS-Community"
  vcpus       = 8
}
```

## Argument Reference

* `engine_name` - (Required) Specifies the engine name of the DDS, `DDS-Community` is supported.

* `type` - (Optional) Specifies the type of the DDS flavor. `mongos`, `shard`, `config` and `replica` are supported.

* `vcpus` - (Optional) Specifies the vCPUs of the DDS flavor.

* `memory` - (Optional) Specifies the RAM of the DDS flavor in GB.


## Attributes Reference

* `region` - See Argument Reference above.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:
  * `spec_code` - The name of the DDS flavor.
  * `type` - See `type` above.
  * `vcpus` - See `vcpus` above.
  * `memory` - See `memory` above.
  * `az_status` - Indicates the status of specifications in an AZ.
