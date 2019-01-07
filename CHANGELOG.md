## 1.6.0 (Unreleased)
## 1.5.0 (January 07, 2019)

FEATURES:

* **New Data Source:** `opentelekomcloud_dcs_az_v1` ([#154](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/154))
* **New Data Source:** `opentelekomcloud_dcs_maintainwindow_v1` ([#154](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/154))
* **New Data Source:** `opentelekomcloud_dcs_product_v1` ([#154](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/154))
* **New Resource:** `opentelekomcloud_networking_floatingip_associate_v2` ([#153](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/153))
* **New Resource:** `opentelekomcloud_dcs_instance_v1` ([#154](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/154))

BUG FIXES:

* `resource/opentelekomcloud_vpc_subnet_v1`: Remove UNKNOWN status to avoid error ([#158](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/158))
* `resource/opentelekomcloud_rds_instance_v1`: Suppress rds name change ([#161](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/161))
* `resource/opentelekomcloud_kms_key_v1`: Add default value of pending_days ([#163](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/163))
* `all resources`: Expose real error message of BadRequest error ([#164](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/164))
* `resource/opentelekomcloud_sfs_file_system_v2`: Suppress sfs system metadata ([#168](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/168))

ENHANCEMENTS:

* Add AKSK authentication support ([#157](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/157))
* `data/opentelekomcloud_images_image_v2`: Add properties filter support for images data source ([#165](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/165))
* `resource/opentelekomcloud_compute_instance_v2`: Add key/value tag support ([#169](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/169))
* `data/opentelekomcloud_vpc_subnet_v1`: Sort vpc subnet ids by network ip availabilities ([#171](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/171))

## 1.4.0 (December 10, 2018)

FEATURES:

* **New Data Source:** `opentelekomcloud_cts_tracker_v1` ([#135](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/135))
* **New Data Source:** `opentelekomcloud_antiddos_v1` ([#138](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/138))
* **New Data Source:** `opentelekomcloud_cce_node_v3` ([#140](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/140))
* **New Data Source:** `opentelekomcloud_cce_cluster_v3` ([#140](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/140))
* **New Resource:** `opentelekomcloud_compute_bms_server_v2` ([#132](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/132))
* **New Resource:** `opentelekomcloud_cts_tracker_v1` ([#135](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/135))
* **New Resource:** `opentelekomcloud_antiddos_v1` ([#138](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/138))
* **New Resource:** `opentelekomcloud_cce_node_v3` ([#140](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/140))
* **New Resource:** `opentelekomcloud_cce_cluster_v3` ([#140](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/140))
* **New Resource:** `opentelekomcloud_maas_task_v1` ([#142](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/142))

## 1.3.0 (November 05, 2018)

FEATURES:

* **New Data Source:** `opentelekomcloud_vbs_backup_policy_v2` ([#121](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/121))
* **New Data Source:** `opentelekomcloud_vbs_backup_v2` ([#121](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/121))
* **New Resource:** `opentelekomcloud_vbs_backup_policy_v2` ([#121](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/121))
* **New Resource:** `opentelekomcloud_vbs_backup_v2` ([#121](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/121))
* **New Resource:** `opentelekomcloud_vbs_backup_share_v2` ([#121](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/121))
* **New Resource:** `opentelekomcloud_mrs_cluster_v1` ([#126](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/126))
* **New Resource:** `opentelekomcloud_mrs_job_v1` ([#126](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/126))

BUG FIXES:

* `resource/opentelekomcloud_elb_loadbalancer`: Fix ELB client error ([#129](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/129))

## 1.2.0 (October 01, 2018)

FEATURES:

* **New Data Source:** `opentelekomcloud_deh_host_v1` ([#98](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/98))
* **New Data Source:** `opentelekomcloud_deh_server_v1` ([#98](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/98))
* **New Data Source:** `opentelekomcloud_rts_software_config_v1` ([#97](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/97))
* **New Data Source:** `opentelekomcloud_rts_software_deployment_v1` ([#97](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/97))
* **New Data Source:** `opentelekomcloud_vpc_v1` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Data Source:** `opentelekomcloud_vpc_subnet_v1` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Data Source:** `opentelekomcloud_vpc_subnet_ids_v1` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Data Source:** `opentelekomcloud_vpc_route_v2` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Data Source:** `opentelekomcloud_vpc_route_ids_v2` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Data Source:** `opentelekomcloud_vpc_peering_connection_v2` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Data Source:** `opentelekomcloud_compute_bms_nic_v2` ([#101](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/101))
* **New Data Source:** `opentelekomcloud_compute_bms_keypairs_v2` ([#101](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/101))
* **New Data Source:** `opentelekomcloud_compute_bms_flavors_v2` ([#101](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/101))
* **New Data Source:** `opentelekomcloud_compute_bms_server_v2` ([#101](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/101))
* **New Data Source:** `opentelekomcloud_rts_stack_v1` ([#95](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/95))
* **New Data Source:** `opentelekomcloud_rts_stack_resource_v1` ([#95](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/95))
* **New Data Source:** `opentelekomcloud_sfs_file_system_v2` ([#92](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/92))
* **New Data Source:** `opentelekomcloud_csbs_backup_v1` ([#117](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/117))
* **New Data Source:** `opentelekomcloud_csbs_backup_policy_v1` ([#117](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/117))
* **New Resource:** `opentelekomcloud_deh_host_v1` ([#98](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/98))
* **New Resource:** `opentelekomcloud_rts_software_config_v1` ([#97](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/97))
* **New Resource:** `opentelekomcloud_rts_software_deployment_v1` ([#97](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/97))
* **New Resource:** `opentelekomcloud_vpc_v1` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Resource:** `opentelekomcloud_vpc_subnet_v1` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Resource:** `opentelekomcloud_vpc_route_v2` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Resource:** `opentelekomcloud_vpc_peering_connection_v2` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Resource:** `opentelekomcloud_vpc_peering_connection_accepter_v2` ([#87](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/87))
* **New Resource:** `opentelekomcloud_sfs_file_system_v2` ([#92](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/92))
* **New Resource:** `opentelekomcloud_rts_stack_v1` ([#95](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/95))
* **New Resource:** `opentelekomcloud_nat_gateway_v2` ([#107](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/107))
* **New Resource:** `opentelekomcloud_nat_snat_rule_v2` ([#107](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/107))
* **New Resource:** `opentelekomcloud_as_configuration_v1` ([#108](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/108))
* **New Resource:** `opentelekomcloud_as_group_v1` ([#108](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/108))
* **New Resource:** `opentelekomcloud_as_policy_v1` ([#108](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/108))
* **New Resource:** `opentelekomcloud_dms_queue_v1` ([#114](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/114))
* **New Resource:** `opentelekomcloud_dms_group_v1` ([#114](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/114))
* **New Resource:** `opentelekomcloud_csbs_backup_v1` ([#117](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/117))
* **New Resource:** `opentelekomcloud_csbs_backup_policy_v1` ([#117](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/117))
* **New Resource:** `opentelekomcloud_networking_vip_v2` ([#119](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/119))
* **New Resource:** `opentelekomcloud_networking_vip_associate_v2` ([#119](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/119))

## 1.1.0 (May 26, 2018)

FEATURES:

* **New Data Source:** `opentelekomcloud_kms_key_v1` ([#14](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/14))
* **New Data Source:** `opentelekomcloud_kms_data_key_v1` ([#14](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/14))
* **New Data Source:** `opentelekomcloud_rds_flavors_v1` ([#15](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/15))
* **New Resource:** `opentelekomcloud_kms_key_v1` ([#14](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/14))
* **New Resource:** `opentelekomcloud_rds_instance_v1` ([#15](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/15))
* **New Resource:** `opentelekomcloud_vpc_eip_v1` ([#48](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/48))

ENHANCEMENTS:
* resource/opentelekomcloud_compute_instance_v2: Add `auto_recovery` argument ([#20](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/20))
* resource/opentelekomcloud_networking_router_v2: Add `enable_snat` argument ([#53](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/53))

## 1.0.0 (December 08, 2017)

Initial release of the OpenTelekom Cloud Provider
