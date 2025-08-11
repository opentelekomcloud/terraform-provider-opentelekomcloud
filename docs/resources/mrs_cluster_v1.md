---
subcategory: "MapReduce Service (MRS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_mrs_cluster_v1"
sidebar_current: "docs-opentelekomcloud-resource-mrs-cluster-v1"
description: |-
  Manages a MRS Cluster resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for MRS cluster you can get at
[documentation portal](https://docs.otc.t-systems.com/mapreduce-service/api-ref/apis/cluster_management_apis_v1)

# opentelekomcloud_mrs_cluster_v1

Manages resource cluster within OpenTelekomCloud MRS.

## Example Usage

```hcl
variable "az" {}
variable "vpc_id" {}
variable "network_id" {}
variable "keyname" {}

resource "opentelekomcloud_mrs_cluster_v1" "this" {
  cluster_name          = "mrs-cluster"
  billing_type          = 12
  master_node_num       = 2
  core_node_num         = 3
  master_node_size      = "c4.4xlarge.4.linux.mrs"
  core_node_size        = "c4.4xlarge.4.linux.mrs"
  available_zone_id     = var.az
  vpc_id                = var.vpc_id
  subnet_id             = var.network_id
  cluster_version       = "MRS 3.1.2-LTS.6"
  volume_type           = "SAS"
  volume_size           = 100
  cluster_type          = 0
  safe_mode             = 1
  node_public_cert_name = var.keyname
  cluster_admin_secret  = "Qwerty!123"
  component_list {
    component_name = "Hadoop"
  }
  component_list {
    component_name = "Spark"
  }
  component_list {
    component_name = "HBase"
  }
  component_list {
    component_name = "Hive"
  }
  component_list {
    component_name = "Flink"
  }

  bootstrap_scripts {
    name       = "Modify os config"
    uri        = "s3a://bootstrap/modify_os_config.sh"
    parameters = "param1 param2"
    nodes = [
      "master",
      "core",
      "task",
    ]
    active_master          = true
    before_component_start = true
    fail_action            = "continue"
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `billing_type` - (Required, Integer, ForceNew) The value is `12`, indicating on-demand payment.

* `master_node_num` - (Required, Integer, ForceNew) Number of Master nodes.

* `master_node_size` - (Required, String, ForceNew) Instance specifications of the Master node, for example, `c6.4xlarge.4linux.mrs`. MRS supports host specifications determined by CPU, memory, and disk space. For details about instance specifications, see [ECS Specifications Used by MRS](https://docs.otc.t-systems.com/mapreduce-service/api-ref/appendix/ecs_specifications_used_by_mrs.html#mrs-01-9005).

* `core_node_num` - (Required, Integer, ForceNew) Number of Core nodes Value range: `1` to `500`. A maximum of `500` Core nodes are supported by default. If more than `500` Core nodes are required, contact technical support engineers or invoke background APIs to modify the database.

* `core_node_size` - (Required, String, ForceNew) Instance specification of a Core node Configuration method of this parameter is identical to that of `master_node_size`.

* `available_zone_id` - (Required, String, ForceNew) ID of an available zone. Obtain the value from Regions and Endpoints.

* `cluster_name` - (Required, String, ForceNew) Cluster name, which is globally unique and contains only `1` to `64` letters, digits, hyphens (-), and underscores (_).

* `vpc_id` - (Required, String, ForceNew) ID of the VPC where the subnet locates Obtain the VPC ID from the management console as follows: Register an account and log in to the management console. Click Virtual Private Cloud and select Virtual Private Cloud from the left list. On the Virtual Private Cloud page, obtain the VPC ID from the list.

* `subnet_id` - (Required, String, ForceNew) Subnet ID Obtain the subnet ID from the management console as follows: Register an account and log in to the management console. Click Virtual Private Cloud and select Virtual Private Cloud from the left list. On the Virtual Private Cloud page, obtain the subnet ID from the list.

* `cluster_version` - (Required, String, ForceNew) Version of the clusters. Please refer to `Table 1` in the [API document](https://docs.otc.t-systems.com/mapreduce-service/api-ref/apis/cluster_management_apis_v1/creating_a_cluster_and_running_a_job.html) for supported versions.

* `cluster_type` - (Optional, Integer, ForceNew) Type of clusters. `0`: analysis cluster, `1`: streaming cluster The default value is `0`.

* `volume_type` - (Optional, String, ForceNew) Type of disks. Supported values: `SAS` (High I/O), `SSD` (Ultra-high I/O).

* `volume_size` - (Optional, Integer, ForceNew) Data disk storage space of a Core node Users can add disks to expand storage capacity when creating a cluster. There are the following scenarios: Separation of data storage and computing: Data is stored in the OBS system. Costs of clusters are relatively low but computing performance is poor. The clusters can be deleted at any time. It is recommended when data computing is not frequently performed. Integration of data storage and computing: Data is stored in the HDFS system. Costs of clusters are relatively high but computing performance is good. The clusters cannot be deleted in a short term. It is recommended when data computing is frequently performed. Value range: `100` GB to `32000` GB.

* `master_data_volume_type` - (Optional, String, ForceNew) Data disk storage type of the Master node. Supported values: `SAS` (High I/O), `SSD` (Ultra-high I/O).

* `master_data_volume_size` - (Optional, Integer, ForceNew) Data disk size of the Master node. Value range: `100` GB to `32000` GB.

* `master_data_volume_count` - (Optional, Integer, ForceNew) Number of data disks of the Master node. The value can be set to `1` only.

* `core_data_volume_type` - (Optional, String, ForceNew) Data disk storage type of the Core node.  Supported values: `SAS` (High I/O), `SSD` (Ultra-high I/O).

* `core_data_volume_size` - (Optional, Integer, ForceNew) Data disk size of the Core node. Value range: `100` GB to `32000` GB.

* `core_data_volume_count` - (Optional, Integer, ForceNew) Number of data disks of the Core node. Value range: `1` to `10`.

* `node_public_cert_name` - (Required, String, ForceNew) Name of a key pair You can use a key to log in to the Master node in the cluster.

* `safe_mode` - (Required, Integer, ForceNew) MRS cluster running mode `0`: common mode. The value indicates that the Kerberos authentication is disabled. Users can use all functions provided by the cluster. `1`: safe mode. The value indicates that the Kerberos authentication is enabled. Common users cannot use the file management or job management functions of an MRS cluster and cannot view cluster resource usage or the job records of Hadoop and Spark. To use these functions, the users must obtain the relevant permissions from the MRS Manager administrator. The request has the `cluster_admin_secret` parameter only when `safe_mode` is set to `1`.

* `cluster_admin_secret` - (Optional, String, ForceNew) Indicates the password of the MRS Manager administrator. The password must contain `8` to `32` characters. Must contain at least two types of the following: Lowercase letters, Uppercase letters, Digits, Special characters `~!@#$%^&*()-_=+\|[{}];:'",<.>/?` and spaces.

* `log_collection` - (Optional, Integer, ForceNew) Indicates whether logs are collected when cluster installation fails. `0`: not collected. `1`: collected. The default value is `0`. If `log_collection` is set to `1`, OBS buckets will be created to collect the MRS logs. These buckets will be charged.

* `component_list` - (Required, List, ForceNew) Service component list.
  * `component_name` - (Required, String, ForceNew) Component name.

* `add_jobs` - (Optional, List, ForceNew) You can submit a job when you create a cluster to save time and use MRS easily. Only one job can be added.
  * `job_type` - (Required, Integer, ForceNew) Type. `1`: MapReduce, `2`: Spark, `3`: Hive Script, `4`: HiveQL (not supported currently), `5`: DistCp, importing and exporting data (not supported in this API currently), `6`: Spark Script, `7`: Spark SQL, submitting Spark SQL statements (not supported in this API currently).
  
  * `job_name` - (Required, String, ForceNew) It contains only `1` to `64` letters, digits, hyphens (-), and underscores (_).
  
  * `jar_path` - (Required, String, ForceNew) Path of the `.jar` file or `.sql` file for program execution. The parameter must meet the following requirements: Contains a maximum of `1023` characters, excluding special characters such as `;|&><'$`. The address cannot be empty or full of spaces. Starts with `/` or `s3a://`. Spark Script must end with `.sql` while `MapReduce` and `Spark Jar` must end with `.jar`. `sql` and `jar` are case-insensitive.
  
  * `arguments` - (Optional, String, ForceNew) Key parameter for program execution. The parameter is specified by the function of the user's program. MRS is only responsible for loading the parameter. The parameter contains a maximum of `2047` characters, excluding special characters such as `;|&>'<$`, and can be empty.
  
  * `input` - (Optional, String, ForceNew) Path for inputting data, which must start with `/` or `s3a://`. A correct OBS path is required. The parameter contains a maximum of `1023` characters, excluding special characters such as `;|&>'<$`, and can be empty. 
  
  * `output` - (Optional, String, ForceNew) Path for outputting data, which must start with `/` or `s3a://`. A correct OBS path is required. If the path does not exist, the system automatically creates it. The parameter contains a maximum of `1023` characters, excluding special characters such as `;|&>'<$`, and can be empty.
  
  * `job_log` - (Optional, String, ForceNew) Path for storing job logs that record job running status. This path must start with `/` or `s3a://`. A correct OBS path is required. The parameter contains a maximum of `1023` characters, excluding special characters such as `;|&>'<$`, and can be empty.
  
  * `shutdown_cluster` - (Optional, Boolean, ForceNew) Whether to delete the cluster after the jobs are complete. `true`: Yes, `false`: No.
  
  * `file_action` - (Optional, String, ForceNew) Data import and export import export
  
  * `submit_job_once_cluster_run` - (Required, Boolean, ForceNew) Possible values are: `true` a job is submitted when a cluster is created and `false` a job is submitted separately.
  
  * `hql` - (Optional, String, ForceNew) HiveQL statement.
  
  * `hive_script_path` - (Optional, String, ForceNew) SQL program path. This parameter is needed by Spark Script and Hive Script jobs only and must meet the following requirements: Contains a maximum of `1023` characters, excluding special characters such as `;|&><'$`. The address cannot be empty or full of spaces. Starts with `/` or `s3a://`. Ends with `.sql`. `sql` is case-insensitive.

* `bootstrap_scripts` - (Optional, List, ForceNew) Bootstrap action scripts. For details, see bootstrap_scripts block below.
  * `name` - (Required, String, ForceNew) Name of a bootstrap action script. It must be unique in a cluster. The value can contain only digits, letters, spaces, hyphens (-), and underscores (_) and cannot start with a space.The value can contain `1` to `64` characters.
  
  * `uri` - (Required, String, ForceNew) Path of a bootstrap action script. Set this parameter to an OBS bucket path or a local VM path.
    **OBS bucket path**: Enter a script path manually. For example, enter the path of the public sample script provided by MRS. Example: s3a://bootstrap/presto/presto-install.sh. If dualroles is installed, the parameter of the presto-install.sh script is dualroles. If worker is installed, the parameter of the presto-install.sh script is worker. Based on the Presto usage habit, you are advised to install dualroles on the active Master nodes and worker on the Core nodes.
    **Local VM path**: Enter a script path. The script path must start with a slash (/) and end with .sh.
  
  * `parameters` - (Optional, String, ForceNew) Bootstrap action script parameters.
  
  * `nodes` - (Required, List of Strings, ForceNew) Type of node where the bootstrap action script is executed, including `master`, `core`, and `task`.
  
  * `active_master` - (Optional, Boolean, ForceNew) Whether the bootstrap action script runs only on active Master nodes.
  
  * `before_component_start` - (Optional, Boolean, ForceNew) Time when the bootstrap action script is executed. Currently, the script can be executed before and after the component is started. The default value is `false`, indicating that the bootstrap action script is executed after the component is started.
  
  * `fail_action` - (Required, String, ForceNew) Whether to continue to execute subsequent scripts and create a cluster after the bootstrap action script fails to be executed. `continue`: Continue to execute subsequent scripts. `errorout`: Stop the action.

  -> **NOTE:**
  Please refer to [user guide](https://docs.otc.t-systems.com/mapreduce-service/umn/managing_clusters/bootstrap_actions/index.html) for configuring bootstrap actions.

* `tags` - (Optional, Map) Tags key/value pairs to associate with the cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `order_id` - Order ID for creating clusters.

* `cluster_id` - Cluster ID.

* `available_zone_name` - Name of an availability zone.

* `instance_id` - Instance ID.

* `hadoop_version` - Hadoop version.

* `master_node_ip` - IP address of a Master node.

* `externalIp` - Internal IP address.

* `private_ip_first` - Primary private IP address.

* `external_ip` - External IP address.

* `slave_security_groups_id` - Standby security group ID.

* `security_groups_id` - Security group ID.

* `external_alternate_ip` - Backup external IP address.

* `master_node_spec_id` - Specification ID of a Master node.

* `core_node_spec_id` - Specification ID of a Core node.

* `master_node_product_id` - Product ID of a Master node.

* `core_node_product_id` - Product ID of a Core node.

* `duration` - Cluster subscription duration.

* `vnc` - URI address for remote login of the elastic cloud server.

* `fee` - Cluster creation fee, which is automatically calculated.

* `deployment_id` - Deployment ID of a cluster.

* `cluster_state` - Cluster status. Valid values include: existing history: `starting`,
  `running`, `terminated`, `failed`, `abnormal`, `terminating`, `rebooting`,
  `shutdown`, `frozen`, `scaling-out`, `scaling-in`, `scaling-error`.

* `tenant_id` - Project ID.

* `create_at` - Cluster creation time.

* `update_at` - Cluster update time.

* `error_info` - Error information.

* `charging_start_time` - Time when charging starts.

* `remark` - Remarks of a cluster.

The `component_list` attributes:

* `component_id` - Component ID.

* `component_version` - Component version.

* `componen_desc` - Component description.

## Timeouts

This resource provides the following timeouts configuration options:

- `create` - Default is 30 minutes.

- `delete` - Default is 5 minutes.

## Import

Cluster can be imported using the `cluster_id`, e.g.

```shell
terraform import opentelekomcloud_mrs_cluster_v1.cluster_1 4729ab1c-7c1a-4411-a02e-93dfc361b32d
```
