---
subcategory: "Tag Management Service (TMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_tms_resource_instances_v1"
sidebar_current: "docs-opentelekomcloud-datasource-tms-resource-instances-v1"
description: |-
  Get TMS resource instances resource, used to query resources by tag, within OpenTelekomCloud.
---

Up-to-date reference of API arguments for TMS resource instances you can get at
[documentation portal](https://docs.otc.t-systems.com/tag-management-service/api-ref/api_description/resource_tags/querying_resources_by_tag.html)

# opentelekomcloud_tms_resource_instances_v1

Use this data source to get the list of resource instances filtered by tag within OpenTelekomCloud.

## Example Usage

```hcl
variable "project_id" {}

data "opentelekomcloud_tms_resource_instances_v1" "tags_ds_1" {
  resource_types = ["disk", "ecs"]
  tags {
    key    = "test"
    values = ["test-tf-acc"]
  }
  project_id = var.project_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_types` - (Required, List[String]) Specifies the resource type. This parameter is case-sensitive. Supported resource types can be provided as ecs,scaling_group, images, disk,vpcs,security-groups, shared_bandwidth,eip, cdn.

* `project_id` - (Optional, String) Specifies the Project ID. This parameter is mandatory when resource_type is a region-specific service.

* `without_any_tag` - (Optional, Boolean) Specifies whether to query only untagged resources. If this parameter is set to `true`, only untagged resources are queried.

* `tags` - (Required, List) Specifies the list of tags. The structure is documented below.
    * `key` - (Required, String) Specifies the key. It can contain up to 36 characters. The key cannot be empty. Only digits, letters, hyphens (-), at signs (@), and underscores (_) are allowed.
    * `values` - (Required, List[String]) Specifies tag values. Each value contains a maximum of 43 Unicode characters and can be an empty string. Only digits, letters, hyphens (-), at signs (@), and underscores (_) are allowed.



## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - Indicates the list of resources. The structure is documented below.
    * `project_id` - Indicates the project ID.
    * `project_name` - Indicates the  project name.
    * `resource_id` - Indicates the resource ID.
    * `resource_name` - Indicates the resource name.
    * `resource_type` - Indicates the resource type.
    * `tags` - Indicates the resource tags.
