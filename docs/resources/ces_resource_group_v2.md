---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_resource_group_v2"
sidebar_current: "docs-opentelekomcloud-resource-ces-resource-group-v2"
description: |-
  Manages a CES Resource Group v2 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES resource group you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_v2/resource_groups/index.html)

# opentelekomcloud_ces_resource_group_v2

Manages a CES resource group resource within OpenTelekomCloud.

## Example Usage

### Add resources manually

```hcl
variable "subnet_id" {}

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name        = "ecs-test"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = var.subnet_id
  }
}

resource "opentelekomcloud_ces_resource_group_v2" "test" {
  name = "test"

  resources {
    namespace = "SYS.ECS"
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.vm_1.id
    }
  }

  resources {
    namespace = "SYS.EVS"
    dimensions {
      name  = "disk_name"
      value = "${opentelekomcloud_compute_instance_v2.vm_1.id}-sda"
    }
  }
}
```

### Add resources from enterprise projects

```hcl
variable "eps_id" {}

resource "opentelekomcloud_ces_resource_group_v2" "test" {
  name               = "test"
  type               = "EPS"
  associated_eps_ids = [var.eps_id]
}
```

### Add resources by tags

```hcl
resource "opentelekomcloud_ces_resource_group_v2" "test" {
  name = "test"
  type = "TAG"
  tags = {
    key = "value"
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the resource group name.
  This parameter can contain a maximum of 128 characters, which may consist of letters,
  digits, hyphens (-), and underscores (_). It must start with a letter.

* `type` - (Optional, String, ForceNew) Specifies the resource group type.
  The value can be **EPS**, **TAG**, and **Manual**. If not specified, that means add resources manually.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of the resource group.
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value to match resources.
  It's required if the value of type is **TAG**.

* `associated_eps_ids` - (Optional, List, ForceNew) Specifies the enterprise project IDs where the resources from.
  It's required if the value of type is **EPS**.
  Changing this parameter will create a new resource.

* `resources` - (Optional, List) Specifies the list of resources to add into the group.
  The [resources](#ResourceGroup_resources) structure is documented below.

<a name="ResourceGroup_resources"></a>
The `resources` block supports:

* `namespace` - (Required, String) Specifies the namespace in **service.item** format.
  **service** and **item** each must be a string that starts with a letter and contains only letters, digits, and
  underscores (_). For details,
  see [Services Interconnected with Cloud Eye](https://docs.otc.t-systems.com/cloud-eye/api-ref/appendix/services_interconnected_with_cloud_eye.html).

* `dimensions` - (Required, List) Specifies the list of dimensions.
  The [dimensions](#ResourceGroup_dimensions) structure is documented below.

<a name="ResourceGroup_dimensions"></a>
The `dimensions` block supports:

* `name` - (Required, String) Specifies the dimension name.
  The value can be a string of 1 to 32 characters that must start with a letter
  and contain only letters, digits, and underscores (_).

* `value` - (Required, String) Specifies the dimension value.
  The value can be a string of 1 to 256 characters that must start with a letter or a number
  and contain only letters, digits, underscores (_), hyphens (-), and dots (.).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time.

* `region` - The region in which the resource group is created.

## Import

The resource group can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_ces_resource_group_v2.test 0ce123456a00f2591fabc00385ff1234
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `resources`.
It is generally recommended running `terraform plan` after importing a resource group.
You can then decide if changes should be applied to the resource group, or the resource definition should be updated to
align with the resource group. Also you can ignore changes as below.

```hcl
resource "opentelekomcloud_ces_resource_group_v2" "test" {
  lifecycle {
    ignore_changes = [
      resources,
    ]
  }
}
```
