---
subcategory: "Document Database Service (DDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dds_instance_v3"
sidebar_current: "docs-opentelekomcloud-datasource-dds-instance-v3"
description: |-
  Get DDS instance from OpenTelekomCloud
---

Up-to-date reference of API arguments for DDS instance you can get at
[documentation portal](https://docs.otc.t-systems.com/document-database-service/api-ref/apis_v3.0_recommended/db_instance_management/querying_instances_and_details.html)

# opentelekomcloud_dds_instance_v3

Use this data source to get info of the OpenTelekomCloud DDS instance.

## Example Usage

```hcl
variable "instance_id" {}

data "opentelekomcloud_dds_instance_v3" "instance" {
  instance_id = var.instance_id
}
```


## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional) Specifies the DB instance ID.

* `name` - (Optional) Specifies the DB instance name.

* `datastore_type` - (Optional) Specifies the database type. The value is `DDS-Community`.

* `vpc_id` - (Optional) Specifies the VPC ID. You can log in to the VPC console and
  obtain the ID of the VPC where the DDS instance is located.

* `subnet_id` - (Optional) Specifies the network ID of the subnet. You can log in to
  the VPC console and obtain the network ID of the subnet in the VPC where the DDS
  instance is located.


## Attributes Reference

The following attributes are exported:

* `id` - Indicates the DB instance ID.

* `region` - Indicates the region where the DB instance is deployed.

* `name` - Indicates the DB instance name.

* `availability_zone` - Indicates the AZ.

* `vpc_id` - Indicates the VPC ID.

* `subnet_id` - Indicates the subnet ID.

* `security_group_id` - Indicates the security group ID.

* `disk_encryption_id` - Indicates the disk encryption key ID. This parameter is returned
  only when the instance disk is encrypted.

* `mode` - Indicates the instance type, which is the same as the request parameter.

* `db_username` - Indicates the default username.

* `status` - Indicates the DB instance status.

* `ssl` - Indicates that SSL is enabled or not.

* `datastore/type` - Indicates the DB engine.

* `datastore/version` - Indicates the database version.

* `datastore/storage_engine` - Specifies the storage engine.

* `backup_strategy/start_time` -  Indicates the backup time window. Automated backups will
  be triggered during the backup time window. The current time is the UTC time.

* `backup_strategy/keep_days` - Indicates the number of days to retain the generated backup
  files. The value range is from 0 to 732.

* `nodes/id` - Indicates the node ID.

* `nodes/name` - Indicates the node name.

* `nodes/role` - Indicates the node role.

* `nodes/type` - Indicates the node type.

* `nodes/private_ip` - Indicates the private IP address of a node.

* `nodes/public_ip` - Indicates the EIP that has been bound on a node.

* `nodes/status` - Indicates the node status.
