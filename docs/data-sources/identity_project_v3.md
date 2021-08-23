---
subcategory: "Identity and Access Management (IAM)"
---

# opentelekomcloud_identity_project_v3

Use this data source to get the ID of an OpenTelekomCloud project.

## Example Usage

```hcl
data "opentelekomcloud_identity_project_v3" "project_1" {
  name = "demo"
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional) The domain this project belongs to.

* `enabled` - (Optional) Whether the project is enabled or disabled. Valid values are `true` and `false`.

* `is_domain` - (Optional) Whether this project is a domain. Valid values are `true` and `false`.

* `name` - (Optional) The name of the project.

* `parent_id` - (Optional) The parent of this project.

## Attributes Reference

`id` is set to the ID of the found project. In addition, the following attributes are exported:

* `description` - The description of the project.

* `region` - Indicates the region where the project is present.
