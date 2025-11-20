---
subcategory: "Relational Database Service (RDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rds_major_versions_v3"
sidebar_current: "docs-opentelekomcloud-datasource-rds-major-versions-v3"
description: |-
  Get available major upgrade versions for RDS instance upgrade
---

Up-to-date reference of API arguments for RDS major version upgrade you can get at
[documentation portal](https://docs.otc.t-systems.com/relational-database-service/api-ref/api_v3_recommended/upgrading_a_major_version/querying_the_target_version_to_which_a_db_instance_can_be_upgraded_rds_for_postgresql.html)

# opentelekomcloud_rds_major_versions_v3

Use this data source to get available major versions for upgrading an OpenTelekomCloud RDS instance.

## Example Usage
```hcl
data "opentelekomcloud_rds_major_versions_v3" "versions" {
  instance_id = var.instance_id
}
```

## Argument Reference

* `instance_id` - (Required) Specifies the ID of the RDS instance.

## Attributes Reference

In addition, the following attributes are exported:

* `available_versions` - List of available major versions for upgrade.
