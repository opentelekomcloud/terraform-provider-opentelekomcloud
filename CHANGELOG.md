## 1.2.0 (Unreleased)

FEATURES:

* **New Data Source:** `opentelekomcloud_deh_host_v1` [GH-98]
* **New Data Source:** `opentelekomcloud_deh_server_v1` [GH-98]
* **New Data Source:** `opentelekomcloud_rts_software_config_v1` [GH-97]
* **New Data Source:** `opentelekomcloud_rts_software_deployment_v1` [GH-97]
* **New Data Source:** `opentelekomcloud_vpc_v1` [GH-87]
* **New Data Source:** `opentelekomcloud_vpc_subnet_v1` [GH-87]
* **New Data Source:** `opentelekomcloud_vpc_subnet_ids_v1` [GH-87]
* **New Data Source:** `opentelekomcloud_vpc_route_v2` [GH-87]
* **New Data Source:** `opentelekomcloud_vpc_route_ids_v2` [GH-87]
* **New Data Source:** `opentelekomcloud_vpc_peering_connection_v2` [GH-87]
* **New Data Source:** `opentelekomcloud_compute_bms_nic_v2` [GH-101]
* **New Data Source:** `opentelekomcloud_compute_bms_keypairs_v2` [GH-101]
* **New Data Source:** `opentelekomcloud_compute_bms_flavors_v2` [GH-101]
* **New Data Source:** `opentelekomcloud_compute_bms_server_v2` [GH-101]
* **New Data Source:** `opentelekomcloud_rts_stack_v1` [GH-95]
* **New Data Source:** `opentelekomcloud_rts_stack_resource_v1` [GH-95]
* **New Data Source:** `opentelekomcloud_sfs_file_system_v2` [GH-92]
* **New Data Source:** `opentelekomcloud_csbs_backup_v1` [GH-117]
* **New Data Source:** `opentelekomcloud_csbs_backup_policy_v1` [GH-117]
* **New Resource:** `opentelekomcloud_deh_host_v1` [GH-98]
* **New Resource:** `opentelekomcloud_rts_software_config_v1` [GH-97]
* **New Resource:** `opentelekomcloud_rts_software_deployment_v1` [GH-97]
* **New Resource:** `opentelekomcloud_vpc_v1` [GH-87]
* **New Resource:** `opentelekomcloud_vpc_subnet_v1` [GH-87]
* **New Resource:** `opentelekomcloud_vpc_route_v2` [GH-87]
* **New Resource:** `opentelekomcloud_vpc_peering_connection_v2` [GH-87]
* **New Resource:** `opentelekomcloud_vpc_peering_connection_accepter_v2` [GH-87]
* **New Resource:** `opentelekomcloud_sfs_file_system_v2` [GH-92]
* **New Resource:** `opentelekomcloud_rts_stack_v1` [GH-95]
* **New Resource:** `opentelekomcloud_nat_gateway_v2` [GH-107]
* **New Resource:** `opentelekomcloud_nat_snat_rule_v2` [GH-107]
* **New Resource:** `opentelekomcloud_as_configuration_v1` [GH-108]
* **New Resource:** `opentelekomcloud_as_group_v1` [GH-108]
* **New Resource:** `opentelekomcloud_as_policy_v1` [GH-108]
* **New Resource:** `opentelekomcloud_dms_queue_v1` [GH-114]
* **New Resource:** `opentelekomcloud_dms_group_v1` [GH-114]
* **New Resource:** `opentelekomcloud_csbs_backup_v1` [GH-117]
* **New Resource:** `opentelekomcloud_csbs_backup_policy_v1` [GH-117]
* **New Resource:** `opentelekomcloud_networking_vip_v2` [GH-119]
* **New Resource:** `opentelekomcloud_networking_vip_associate_v2` [GH-119]

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
