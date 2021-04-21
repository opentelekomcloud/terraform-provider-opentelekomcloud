## 1.24.0 (Unreleased)

## 1.23.8 (April 21, 2021)

ENHANCEMENTS:
* `resource/opentelekomcloud_as_group_v1`: Make `security_groups` optional ([#991](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/991))

BUG FIXES:
* `resource/opentelekomcloud_lb_certificate_v2`: Fix constantly updating `private_key`, `certificate` and `domain` fields ([#988](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/988))
* `resource/opentelekomcloud_lb_certificate_v2`: Fix deleting certificates used in LB listener ([#987](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/987))
* `resource/opentelekomcloud_vpc_subnet_v1`: Fix subnet creation when `dns_list` is set ([#995](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/995), follow-up of [#977](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/977))

DOCUMENTATION:
* `resource/opentelekomcloud_cce_cluster_v3`: Add note about CCE authorization required ([#998](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/998))

## 1.23.7 (April 15, 2021)

ENHANCEMENTS:
* `resource/opentelekomcloud_networking_subnet_v2`: Add default value for `dns_nameservers` ([#977](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/977))
* `resource/opentelekomcloud_vpc_subnet_v1`: Add default value for `primary_dns`, `secondary_dns` ([#977](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/977))

BUG FIXES:
* `resource/opentelekomcloud_cce_node_v3`: Remove passing empty `private_ip` in create request ([#973](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/973))
* `resource/opentelekomcloud_s3_bucket`: Make unversioned bucket creation possible ([#976](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/976), [#979](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/979))
* `resource/opentelekomcloud_obs_bucket`: Make unversioned bucket creation possible ([#978](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/978))
* `resource/opentelekomcloud_lb_listener_v2`: Fix schema to avoid resource always updating ([#983](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/983), [#984](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/984))

DOCUMENTATION:
* `resource/opentelekomcloud_cce_addon_v3`:  Fix documentation issues ([#969](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/969))

## 1.23.6 (April 08, 2021)

FEATURES:
* **New Resource:** `opentelekomcloud_identity_provider_v3` ([#946](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/946))
* **New Resource:** `opentelekomcloud_identity_mapping_v3` ([#947](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/947))
* **New Resource:** `opentelekomcloud_sfs_share_access_rules_v2` ([#955](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/955))
* **New Resource:** `opentelekomcloud_sdrs_protected_instance_v1` ([#963](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/963))

ENHANCEMENTS:
* `resource/opentelekomcloud_cce_node_v3`: Allow to set a `private_ip` ([#938](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/938))
* `resource/opentelekomcloud_as_configuration_v1`: Allow to set a `security_groups` ([#941](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/941))
* `resource/opentelekomcloud_cce_nodepool_v3`: Increase timeouts ([#945](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/945))
* `resource/opentelekomcloud_sfs_file_share_v2`: Make `access` params optional ([#953](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/953))
* `resource/opentelekomcloud_сompute_instance_v2`: Add possibility to set `power_state` param  ([#956](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/956))

DOCUMENTATION:
* `resource/opentelekomcloud_ecs_instance_v1`: Clarify that `nics` is required ([#951](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/951))
* `resource/opentelekomcloud_dns_recordset_v2`: Clarify that `type` and `records` are required ([#961](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/961))

## 1.23.5 (March 24, 2021)

ENHANCEMENTS:
* `resource/opentelekomcloud_ecs_instance_v1`: Use common `tags` approach in resource ([#919](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/919))
* `provider/opentelekomcloud`: Retry `502` error one time ([#921](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/921))
* `resource/opentelekomcloud_lb_monitor_v2`: Add possibility to set `domain_name` argument ([#925](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/925))
* `resource/opentelekomcloud_css_cluster_v1`: Add possibility to set `datastore` argument ([#926](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/926))
* `resource/opentelekomcloud_compute_instance_v2`: Use common `tags` approach in resource ([#927](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/927))
* `resource/opentelekomcloud_css_cluster_v1`: Add disk size validation during a plan ([#928](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/928))

DOCUMENTATION:
* `resource/opentelekomcloud_compute_instance_v2`: Update `security_groups` description ([#929](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/929))
* `resource/opentelekomcloud_cce_addon_v3`: Add description of addon template input values ([#931](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/931))

## 1.23.4 (March 17, 2021)

FEATURES:
* **New Data Source:** `opentelekomcloud_css_flavor_v1` ([#913](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/913))

ENHANCEMENTS:
* `resource/opentelekomcloud_css_cluster_v1`: Add support of `enable_authority` and `admin_pass` arguments ([#902](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/902))
* `resource/opentelekomcloud_ecs_instance_v1`: Use security group IDs in all operations ([#909](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/909))

BUG FIXES:
* `resource/opentelekomcloud_vpnaas_ipsec_policy_v2`: Missing support of PFS groups ([#906](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/906))
* `resource/opentelekomcloud_as_group_v1`: Limit `security_groups` maximum number to one ([#907](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/907))
* `resource/opentelekomcloud_cce_node_pool_v3`: Fix too strict `k8s_tags` validations ([#911](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/911))
* `resource/opentelekomcloud_cce_node_pool_v3`: Changes in `k8s_tags` and `taints` trigger resource re-creation no more ([#911](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/911))

## 1.23.3 (March 12, 2021)

ENHANCEMENTS:
* `resource/opentelekomcloud_lb_loadbalancer_v2`: Add possibility to set tags ([#890](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/890))
* `resource/opentelekomcloud_lb_listener_v2`: Add possibility to set tags ([#895](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/895))
* `resource/opentelekomcloud_compute_keypair_v2`: Add new keypair creation support ([#896](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/896))

BUG FIXES:
* `resource/opentelekomcloud_cbr_vault_v3`: Fix not unassignable resources ([#897](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/897))

DOCUMENTATION:
* Improve repository `README.md` ([#894](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/894))
* `resource/opentelekomcloud_dns_zone_v2`: Clarify that private zones are not searched by default ([#905](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/905))

## 1.23.2 (March 4, 2021)

ENHANCEMENTS:
* `resource/opentelekomcloud_as_group_v1`: Add possibility to set tags ([#877](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/877))
* `resource/opentelekomcloud_kms_key_v1`: Add possibility to set tags ([#884](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/884))

BUG FIXES:
* `resource/opentelekomcloud_cce_node_v3`: Remove reading empty CCE Node `Spec.ExtendParam` ([#876](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/876))
* `resource/opentelekomcloud_css_cluster_v1`: Fix error with reading cluster without encryption ([#882](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/882))
* `resource/opentelekomcloud_cbr_policy_v3`: Make `timezone` argument required ([#883](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/883))
* `resource/opentelekomcloud_compute_keypair_v2`: Fix raising error on changing existing public key ([#887](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/887))

## 1.23.1 (February 25, 2021)

ENHANCEMENTS:
* `resource/opentelekomcloud_sfs_file_system_v2`: Add possibility to set tags ([#867](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/867))

BUG FIXES:
* `resource/opentelekomcloud_cce_node_pool_v3`: Fix pool not creating with `random` availability zone ([#864](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/864))
* `resource/opentelekomcloud_compute_instance_v2`: Fix ignored `OS_IMAGE_ID` env variable ([#866](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/866))
* `resource/opentelekomcloud_compute_bms_server_v2`: Fix ignored `OS_IMAGE_ID` env variable ([#866](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/866))
* `resource/opentelekomcloud_vbs_backup_policy_v2`: Fix panic on refresh when policy is missing ([#872](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/872))

## 1.23.0 (February 17, 2021)

NOTES/DEPRECATIONS:

Binary build matrix is narrowed ([#858](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/858)).
Binaries for the following OS/architecture combinations are built:

* **Linux**
  * `AMD64`
  * `i386`
  * `ARMv6`
  * `ARMv8` (`ARM64`)

* **Darwin**
  * `AMD64`

* **Windows**
  * `AMD64`
  * `i386`

* **FreeBSD**
    * `AMD64`
    * `i386`

_`Darwin/ARMv8` (new `M1` chip) to be also built in future_

FEATURES:
* **New Resource:** `opentelekomcloud_sfs_turbo_share_v1` ([#852](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/852))

ENHANCEMENTS:
* `resource/opentelekomcloud_dns_zone_v2`: Make DNS resources diff ignore ending dot in name ([#850](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/850))
* `resource/opentelekomcloud_dns_recordset_v2`: Make DNS resources diff ignore ending dot in name ([#850](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/850))
* `resource/opentelekomcloud_dns_ptrrecord_v2`: Make DNS resources diff ignore ending dot in name ([#850](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/850))
* `resource/opentelekomcloud_compute_instance_v2`: Remove `Deprecated` and `Removed` fields from schema ([#859](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/859))

BUG FIXES:
* `resource/opentelekomcloud_dns_recordset_v2`: Fix shared DNS recordset searching ([#848](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/848))
* `resource/opentelekomcloud_vbs_backup_v2`: Fix reading backup description ([#855](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/855))

## 1.22.8 (February 10, 2021)

FEATURES:
* **New Resource:** `opentelekomcloud_cce_node_pool_v3` ([#825](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/825))
* **New Resource:** `opentelekomcloud_cbr_vault_v3` ([#833](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/833))

BUG FIXES:
* `provider/opentelekomcloud`: Fix not loading cloud config ([#828](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/828))
* `resource/opentelekomcloud_vpc_eip_v1`: Fix missing `tags` argument in the documentation ([#830](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/830))
* `resource/opentelekomcloud_dns_prrecord_v2`: Repair tags workflow in resource ([#832](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/832))
* `resource/opentelekomcloud_cce_addon_v3`: Fix crash on empty `basic` addon values ([#836](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/836))

## 1.22.7 (February 03, 2021)

ENHANCEMENTS:
* `resource/opentelekomcloud_ecs_instance_v1`: Implement plan stage network and volume validation ([#820](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/820))
* `resource/opentelekomcloud_as_configuration`: Add possibility to set `5_mailbgp` to `ip_type` ([#821](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/821))
* `resource/opentelekomcloud_evs_volume_v3`: Add volume type validation ([#823](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/823))

NOTES/DEPRECATIONS:
* `resource/opentelekomcloud_compute_instance_v2`: Deprecate `personality` field ([#819](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/819))

## 1.22.6 (January 27, 2021)

BUG FIXES:
* `resource/opentelekomcloud_obs_bucket`: Fix invalid AK/SK signature for OBS ([#811](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/811))
* `resource/opentelekomcloud_obs_bucket_object`: Fix invalid AK/SK signature for OBS ([#811](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/811))
* `resource/opentelekomcloud_obs_bucket_object`: Fix issue with deleting versioned objects ([#812](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/812))

ENHANCEMENTS:
* `provider/opentelekomcloud`: Add provider credentials validation ([#813](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/813))
* `provider/opentelekomcloud`: Mark sensitive fields as `Sensitive` ([#816](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/816))

## 1.22.5 (January 22, 2021)

BUG FIXES:
* `resource/opentelekomcloud_cce_cluster_v3`: Fix `vpc_id` and `subnet_id` validation during plan ([#804](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/804))

## 1.22.4 (January 21, 2021)

FEATURES:
* **New Data Source:** `opentelekomcloud_rds_versions_v3` ([#792](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/792))

ENHANCEMENTS:
* `resource/opentelekomcloud_cce_cluster_v3`: Implement plan stage network validation ([#787](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/787))
* `resource/opentelekomcloud_cce_cluster_v3`: Add timeouts section to CCE documentation ([#788](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/788))
* `resource/opentelekomcloud_cce_node_v3`: Add timeouts section to CCE documentation ([#788](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/788))
* `resource/opentelekomcloud_compute_keypair_v2`: Allow shared ("global") key pairs ([#794](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/794))
* `resource/opentelekomcloud_rds_instance_v3`: Implement plan stage db version validation ([#795](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/795))
* `resource/opentelekomcloud_rds_parametergroup_v3`: Implement plan stage db version validation ([#796](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/796))
* `resource/opentelekomcloud_dns_recordset_v2`: Allow shared ("global") DNS record sets ([#800](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/800))

BUG FIXES:
* `resource/opentelekomcloud_rds_instance_v3`: Add `param_group_id` support ([#784](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/784))
* `resource/opentelekomcloud_rds_parametergroup_v3`: Fix `rds_parametergroup_v3` recreation ([#789](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/789))

## 1.22.3 (December 23, 2020)

FEATURES:
* **New Resource:** `opentelekomcloud_obs_bucket_policy` ([#773](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/773))
* **New Data Source:** `opentelekomcloud_obs_bucket_object` ([#780](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/780))

ENHANCEMENTS:
* `resource/opentelekomcloud_cce_node_v3`: Add `os` argument ([#778](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/778))

BUG FIXES:
* `resource/opentelekomcloud_obs_bucket_object`: Remove unused `credentials` argument ([#781](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/781))

DOCUMENTATION:
* `data_source/opentelekomcloud_s3_bucket_object`:  Move to `"Object Storage Service (S3)"` subcategory ([#772](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/772))
* `resource/opentelekomcloud_s3_bucket`: Move to `"Object Storage Service (S3)"` subcategory ([#772](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/772))
* `resource/opentelekomcloud_s3_bucket_object`: Move to `"Object Storage Service (S3)"` subcategory ([#772](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/772))
* `resource/opentelekomcloud_s3_bucket_policy`: Move to `"Object Storage Service (S3)"` subcategory ([#772](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/772))

## 1.22.2 (December 16, 2020)

FEATURES:
* **New Resource:** `opentelekomcloud_cce_addon_v3` ([#711](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/711))

ENHANCEMENTS:
* `resource/opentelekomcloud_lb_loadbalancer_v2`: Clarify usage `lb_loadbalancer_v2` with `vpc_subnet_v1` in docs ([#766](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/766))
* `resource/opentelekomcloud_cce_cluster_v3`: Add cluster name validation ([#768](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/768))

## 1.22.1 (December 10, 2020)

FEATURES:
* **New Resource:** `opentelekomcloud_cbr_policy_v3`([#758](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/758))

BUG FIXES:
* `resource/opentelekomcloud_identity_credential_v3`: Remove non-existing credential instead returning error ([#753](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/753))
* `resource/opentelekomcloud_lb_pool_v2`: Fix LB protocol to pool protocol mapping description ([#754](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/754))
* `resource/opentelekomcloud_rds_instance_v3`: Fix issue with update volume size ([#755](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/755))
* `data_source/opentelekomcloud_networking_secgroup_v2`: Prevent panic due to unhandled error ([#756](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/756))


## 1.22.0 (December 03, 2020)

FEATURES:
* **New Data Source:** `opetelekomcloud_dds_instance_v3` ([#725](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/725))

ENHANCEMENTS:
* `resource/opentelekomcloud_cce_cluster_v3`: Add new argument `authenticating_proxy_ca` ([#727](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/727))
* `data_source/opentelekomcloud_cce_cluster_v3`: Add new argument `authentication_mode` ([#727](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/727))
* `resource/opentelekomcloud_obs_bucket`: Setting up AK/SK is not required anymore ([#745](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/745))
* `resource/opentelekomcloud_obs_bucket_object`: Setting up AK/SK is not required anymore ([#745](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/745))

BUG FIXES:
* `resource/opentelekomcloud_identity_credential_v3`: Add the missing documentation ([#731](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/731))
* `resource/opentelekomcloud_vpnaas_ike_policy_v2`: Fix hardcoded values for `PFS` and `phase1_negotiation_mode` ([#733](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/733))
* `resource/opentelekomcloud_identity_credential_v3`: Make `user_id` Optional ([#737](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/737))


## 1.21.6 (November 25, 2020)

FEATURES:
* **New Resource:** `opetelekomcloud_dds_instance_v3` ([#717](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/717))
* **New Data Source:** `opetelekomcloud_dds_flavors_v3` ([#718](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/718))
* **New Data Source:** `opentelekomcloud_vpc_bandwidth` ([#719](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/719))

BUG FIXES:
* `resource/opentelekomcloud_as_group_v1`: Fix failing autoscaling group deletion ([#722](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/722))


## 1.21.5 (November 19, 2020)

ENHANCEMENTS:
* `provider/opentelekomcloud`: Add `OS_TOKEN` as alternative env var for `token` ([#706](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/706))
* `resource/opentelekomcloud_lb_monitor_v2`: Add `monitor_port` argument ([#709](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/709))
* `resource/opentelekomcloud_waf_domain_v1`: Rename WAF domain server attributes ([#710](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/710))
* `resource/opentelekomcloud_csbs_backup_policy_v1`: Add fields to CSBS policy ([#714](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/714))


## 1.21.4 (November 12, 2020)

BUG FIXES:
* `provider/opentelekomcloud`: Fix retries for 409 and 503 error codes ([#688](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/688))
* `provider/opentelekomcloud`: Fix region handling ([#697](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/697))
* `resource/opentelekomcloud_s3_bucket`: Fix panic creating `s3_bucket` without `tenant_name` in provider config ([#698](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/698))
* `resource/opentelekomcloud_compute_instance_v2`: Revert changes from [#686](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/686) ([#701](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/701))
* `resource/opentelekomcloud_rds_instance_v3`: Fix RDSv3 instance import ([#704](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/704))

ENHANCEMENTS:
* `resource/opentelekomcloud_lb_listener_v2`: Add new field `type` and make `private_key` as Optional ([#688](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/688))
* `resource/opentelekomcloud_lb_certificate_v2`: Add new fields `http2_enable`, `client_ca_tls_container_ref` and `tls_ciphers_policy` ([#688](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/688))
* `resource/opentelekomcloud_cce_cluster_v3`: Add new fields `kubernetes_svc_ip_range` and `kube_proxy_mode` ([#699](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/699))


## 1.21.3 (November 6, 2020)

BUG FIXES:
* `resource/opentelekomcloud_compute_instance_v2`: Fix diff on every apply when using security group IDs instead of names ([#686](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/686))
* `resource/opentelekomcloud_s3_bucket_policy`: Fix not working policy example in documentation ([#692](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/692))

ENHANCEMENTS:
* `resource/opentelekomcloud_cce_node_v3`: Make `iptype`, `bandwidth_charge_mode`, `sharetype` settable ([#681](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/681))
* `resource/opentelekomcloud_cce_node_v3`: Fix not existing flavor in documentation ([#684](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/684))
* `resource/opentelekomcloud_ecs_instance_v1`: Fix not existing flavor in documentation ([#689](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/689))


## 1.21.2 (October 29, 2020)

BUG FIXES:
* `resource/opentelekomcloud_cce_cluster_v3`: Suppress schema diff in CCE version ([#666](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/666))
* `resource/opentelekomcloud_cce_cluster_v3`: Increase delete timeout to 30m ([#674](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/674))
* `resource/opentelekomcloud_compute_secgroup_v2`: Fix delete group if it's used ([#677](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/677))
* `resource/opentelekomcloud_networking_secgroup_v2`: Fix delete group if it's used ([#676](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/676))

FEATURES:
* **New Data Source:** `opentelekomcloud_identity_auth_scope_v3` ([#669](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/669))

ENHANCEMENTS:
* `resource/opentelekomcloud_identity_user_v3`: Add email field to schema ([668](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/668))


## 1.21.1 (October 23, 2020)

BUG FIXES:
* `resource/opentelekomcloud_rds_instance_v3`: Fix not assigning public IP ([#658](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/658))

ENHANCEMENTS:
* `resource/opentelekomcloud_blockstorage_volume_v2`: Allow expanding volume without re-creation ([#661](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/661))

## 1.21.0 (October 15, 2020)

ENHANCEMENTS:
* Migrate to `opentelekomcloud/gophertelekomcloud` from `huaweicloud/golangsdk`: ([#641](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/641))


## 1.20.3 (October 14, 2020)

BUG FIXES:
* `resource/opentelekomcloud_dcs_instance_v1`: Fix issues with DCS schema ([#643](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/643))
* `data_source/opentelekomcloud_role_v3`: Update role list ([#654](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/654))


## 1.20.2 (September 30, 2020)

BUG FIXES:
* `resource/opentelekomcloud_lb_monitor_v2`: Fix `UDP-CONNECT` in type validation ([#634](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/634))
* `resource/opentelekomcloud_cce_node_v3`: Handle 404 during reading tags for CCE node ([#635](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/635))
* `resource/opentelekomcloud_obs_bucket`: Fix not creating OBS bucket with `security_token` ([#636](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/636))

ENHANCEMENTS:
* `resource/opentelekomcloud_cce_node_v3`: Add k8sTags to CCE node resource ([#621](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/621))
* `resource/opentelekomcloud_csbs_backup_policy_v1`: Add `created_at` attribute ([#628](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/628))
* `provider/opentelekomcloud`: Allow setting security token by env variable ([#627](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/627))


## 1.20.1 (September 24, 2020)

BUG FIXES:
* `resource/opentelekomcloud_cce_node_v3`: `public_key` attribute not setting ([#616](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/616))

FEATURES:
* **New Data Source:** `opentelekomcloud_dns_zone_v2` ([#620](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/620))

ENHANCEMENTS:
* `resource/opentelekomcloud_cce_node_v3`: Only `bandwidth_charge_mode` is now required for EIP creation ([#616](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/616))

## 1.20.0 (September 16, 2020)

BUG FIXES:
* `data_source/opentelekomcloud_cce_cluster_v3`: Update outdated docs ([#614](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/614))
* `resource/opentelekomcloud_cce_cluster_v3`: Update outdated docs ([#614](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/614))
* `resource/opentelekomcloud_lb_listener_v2`: Update outdated docs ([#615](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/615))

FEATURES:
* **New Data Source:** `opentelekomcloud_identity_credential_v3` ([#613](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/613))
* **New Resource:** `opentelekomcloud_identity_credential_v3` ([#613](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/613))

## 1.19.5 (September 4, 2020)

BUG FIXES:
* `resource/opentelekomcloud_blockstorage_volume_v2`: Ignore metadata.policy changes in blockstorage_volume_v2 ([#604](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/604))
* `resource/opentelekomcloud_smn_subscription_v2`: Fix r/smn_subscription_v2 and d/cts_tracker_v1 ([608](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/608))
* `data_source/opentelekomcloud_cts_tracker_v1`: Fix r/smn_subscription_v2 and d/cts_tracker_v1 ([608](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/608))

FEATURES:
* **New Data Source:** `opentelekomcloud_vpnaas_service_v2` ([#605](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/605))

## 1.19.4 (September 1, 2020)

BUG FIXES:
* **Multiple Resources:** Documentation fixes after migration ([#599](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/599))

## 1.19.3 (September 1, 2020)

FEATURES:
* **Removed Resource:** `opentelekomcloud_maas_task_v1` ([#585](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/585))

ENHANCEMENTS:
* `resource/opentelekomcloud_compute_instance_v2`: Fix ECS tags-tag confusion ([#586](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/586))
* `resource/opentelekomcloud_rds_instance_v3`: Add setting public IP for RDS instance v3 ([#596](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/596))

## 1.19.2 (August 24, 2020)

ENHANCEMENTS:
* `data_source/opentelekomcloud_cce_cluster_v3`: Add certificates to cce_cluster_v3 data source ([#581](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/581))
* `resource/opentelekomcloud_vpc_eip_v1`: Add `tags` support ([#570](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/570))

## 1.19.1 (August 21, 2020)

BUG FIXES:
* `resource/opentelekomcloud_rds_instance_v3`: Fix HTTP 415 when retrieving tags after nodes role switch ([#564](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/564))
* `resource/opentelekomcloud_cce_cluster_v3`: Add setting `cluster_version` on resource read ([#568](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/568))

## 1.19.0 (August 08, 2020)

ENHANCEMENTS:
* `resource/opentelekomcloud_as_group_v1`: Add health_periodic_audit_grace_period to as group ([#545](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/545))
* `resource/opentelekomcloud_smn_topic_v2`: Add project_name to SMN topic ([#554](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/554))
* `resource/opentelekomcloud_vpc_eip_v1`: Update documentation ([#550](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/550))
* `resource/opentelekomcloud_cts_tracker_v1`: Add project_name to CTS tracker ([#555](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/555))
* `resource/opentelekomcloud_compute_instance_v2`: Improve getting instance NICs ([#559](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/559))

BUG FIXES:
* `resource/opentelekomcloud_rds_instance_v3`: Fix documentation ([#549](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/549))

FEATURES:
* **New Data Source:** `compute_availability_zones_v2`([#558](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/558))

## 1.18.1 (July 10, 2020)

ENHANCEMENTS:

* `resource/opentelekomcloud_as_group_v1`: Add `current_instance_number` and `status` attributes ([#522](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/522))
* `provider/opentelekomcloud`: Add `max_retries` argument to the provider's options ([#537](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/537))

BUG FIXES:

* `resource/rds_instance_v3`: Fix argument description ([#525](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/525))
* `resource/cce_cluster_v3`: Update subnet_id description of CCE cluster ([#535](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/535))

## 1.18.0 (June 16, 2020)

ENHANCEMENTS:

* `opentelekomcloud_vpc_v1`: Add tag support ([#508](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/508))
* `opentelekomcloud_vpc_subnet_v1`: Add tag support ([#508](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/508))
* `opentelekomcloud_dns_zone_v2`: Add tag support ([#510](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/510))
* `opentelekomcloud_dns_recordset_v2`: Add tag support ([#514](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/514))
* `opentelekomcloud_cce_node_v3`: Add tag support ([#513](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/513))


BUG FIXES:

* `opentelekomcloud_waf_domain_v1`: Fix waf_domain_v1 using old waf API ([#496](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/496))
* `opentelekomcloud_dcs_instance_v1`, `opentelekomcloud_dms_instance_v1`, `opentelekomcloud_rds_instance_v3`: Set sensitive flag for password parameter ([#504](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/504))
* `opentelekomcloud_cts_tracker_v1`: Fix handling of missing tracker ([#518](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/518))

## 1.17.1 (May 07, 2020)

BUG FIXES:

* `resource/opentelekomcloud_vpc_subnet_v1`: Fix VPC subnet delete issue ([#492](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/492))

## 1.17.0 (April 26, 2020)

FEATURES:

* **New Data Source:** `opentelekomcloud_dms_az_v1` ([#485](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/485))
* **New Data Source:** `opentelekomcloud_dms_product_v1` ([#485](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/485))
* **New Data Source:** `opentelekomcloud_dms_maintainwindow_v1` ([#485](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/485))
* **New Resource:** `opentelekomcloud_obs_bucket` ([#467](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/467))
* **New Resource:** `opentelekomcloud_obs_bucket_object` ([#467](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/467))
* **New Resource:** `opentelekomcloud_dns_ptrrecord_v2` ([#480](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/480))
* **New Resource:** `opentelekomcloud_dms_instance_v1` ([#485](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/485))

ENHANCEMENTS:

* `resource/opentelekomcloud_ces_alarmrule`: Add alarm_level argument support ([#481](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/481))
* `resource/opentelekomcloud_vbs_backup_policy_v2`: Add associating volumes support ([#478](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/478))
* `resource/opentelekomcloud_rds_instance_v3`: Clean up ID if the intance couldn't be found ([#479](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/479))
* `resource/opentelekomcloud_vbs_backup_policy_v3`: Add week_frequency and rentention_day support ([#489](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/489))

BUG FIXES:

* `resource/opentelekomcloud_fw_rule_v2`: Fix removing assigned FW rule ([#462](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/462))
* `resource/opentelekomcloud_dns_recordset_v2`: Fix updating only TTL value issue ([#465](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/465))
* `resource/opentelekomcloud_vbs_backup_policy_v2`: Fix missing required `frequency` value ([#469](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/469))
* `resource/opentelekomcloud_mrs_cluster_v1`: Update core nodes number validate func ([#477](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/477))

## 1.16.0 (March 06, 2020)

FEATURES:

* **New Resource:** `opentelekomcloud_nat_dnat_rule_v2` ([#447](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/447))

ENHANCEMENTS:

* `resource/opentelekomcloud_cce_node_v3`: Add preinstall/postinstall script support ([#452](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/452))
* `resource/opentelekomcloud_mrs_cluster_v1`: Add tags parameter support ([#453](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/453))
* `resource/opentelekomcloud_mrs_cluster_v1`: Add bootstrap scripts parameter support ([#455](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/455))

BUG FIXES:

* `resource/opentelekomcloud_elb_loadbalancer`: Increase bandwidth range to 1000 ([#459](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/459))
* `data source/opentelekomcloud_vpc_subnet_v1`: Fix vpc_subnet_v1 retrieval by id ([#460](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/460))

## 1.15.1 (February 11, 2020)

BUG FIXES:

* `resource/opentelekomcloud_rds_instance_v3`: Fix RDS instance node id issue ([#450](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/450))

## 1.15.0 (January 16, 2020)

FEATURES:

* **New Resource:** `opentelekomcloud_logtank_group_v2` ([#435](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/435))
* **New Resource:** `opentelekomcloud_logtank_topic_v2` ([#435](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/435))
* **New Resource:** `opentelekomcloud_lb_certificate_v2` ([#437](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/437))
* **New Resource:** `opentelekomcloud_vpc_flow_log_v1` ([#439](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/439))
* **New Resource:** `opentelekomcloud_lb_l7policy_v2` ([#441](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/441))
* **New Resource:** `opentelekomcloud_lb_l7rule_v2` ([#441](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/441))

ENHANCEMENTS:

* `resource/opentelekomcloud_networking_secgroup_v2`: Add description to secgroup_rule_v2 ([#432](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/432))
* `resource/opentelekomcloud_blockstorage_volume_v2`: Update list of values for volume type ([#433](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/433))
* Add clouds.yaml support ([#434](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/434))

## 1.14.0 (December 02, 2019)

FEATURES:

* **New Data Source:** `opentelekomcloud_cce_node_ids_v3` ([#411](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/411))
* **New Resource:** `opentelekomcloud_vpnaas_endpoint_group_v2` ([#412](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/412))
* **New Resource:** `opentelekomcloud_vpnaas_ike_policy_v2` ([#412](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/412))
* **New Resource:** `opentelekomcloud_vpnaas_ipsec_policy_v2` ([#412](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/412))
* **New Resource:** `opentelekomcloud_vpnaas_service_v2` ([#412](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/412))
* **New Resource:** `opentelekomcloud_vpnaas_site_connection_v2` ([#412](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/412))

ENHANCEMENTS:

* `resource/opentelekomcloud_evs_volume_v3`: Add kms_id parameter support ([#403](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/403))
* `resource/opentelekomcloud_cce_cluster_v3`: Add eip update support ([#410](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/410))
* `resource/opentelekomcloud_compute_instance_v2`: Log fault message when build compute instance failed ([#413](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/413))
* `resource/opentelekomcloud_evs_volume_v3`: Add device_type parameter support ([#419](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/419))
* `resource/opentelekomcloud_evs_volume_v3`: Add wwn attribute support ([#420](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/420))

BUG FIXES:

* `resource/opentelekomcloud_cce_node_v3`: Fix cce node update issue ([#405](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/405))
* `resource/opentelekomcloud_dcs_instance_v1`: Fix ip/port attributes issue ([#408](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/408))
* `resource/opentelekomcloud_mrs_cluster_v1`: Fix MRS region issue ([#409](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/409))
* `resource/opentelekomcloud_compute_bms_server_v2`: Fix BMS boot from volume issue ([#422](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/422))

## 1.13.1 (October 22, 2019)

ENHANCEMENTS:

* `resource/opentelekomcloud_cce_cluster_v3`: Add eip parameter support ([#400](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/400))
* `resource/opentelekomcloud_compute_bms_server_v2`: Add tags parameter support ([#401](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/401))

## 1.13.0 (October 18, 2019)

FEATURES:

* **New Resource:** `opentelekomcloud_evs_volume_v3` ([#380](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/380))
* **New Resource:** `opentelekomcloud_lb_whitelist_v2` ([#390](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/390))
* **New Resource:** `opentelekomcloud_ims_image_v2` ([#391](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/391))
* **New Resource:** `opentelekomcloud_ims_data_image_v2` ([#396](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/396))

ENHANCEMENTS:

* `resource/opentelekomcloud_vpc_subnet_v1`: Add NTP server addresses support ([#369](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/369))
* `resource/opentelekomcloud_rds_instance_v3`: Add tag support ([#373](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/373))
* `resource/opentelekomcloud_rds_instance_v3`: Add flavor update support ([#377](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/377))
* `resource/opentelekomcloud_rds_instance_v3`: Add volume resize support ([#378](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/378))
* `resource/opentelekomcloud_waf_domain_v1`: Add policy_id parameter support ([#381](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/381))
* `resource/opentelekomcloud_as_group_v1`: Add lbaas_listeners parameter support ([#385](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/385))
* `resource/opentelekomcloud_as_configuration_v1`: Add kms_id parameter support ([#389](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/389))

BUG FIXES:

* `resource/opentelekomcloud_rds_instance_v3`: Fix RDS backup_strategy parameter issue ([#367](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/367))
* `data resource/opentelekomcloud_vpc_v1`: Fix id filter issue ([#379](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/379))

## 1.12.0 (August 30, 2019)

FEATURES:

* **New Resource:** `opentelekomcloud_ecs_instance_v1` ([#347](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/347))

ENHANCEMENTS:

* `resource/opentelekomcloud_cce_cluster_v3`: Add CCE cluster certificates ([#349](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/349))
* `resource/opentelekomcloud_cce_cluster_v3`: Add multi-az support for CCE cluster ([#350](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/350))
Add detailed error message for 404 ([#352](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/352))

BUG FIXES:

* `resource/opentelekomcloud_vpc_subnet_v1`: Fix dns_list type issue ([#351](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/351))
* `resource/opentelekomcloud_cce_node_v3`: Fix data_volumes type issue ([#354](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/354))
Fix common user ak/sk authentication issue with domain_name ([#362](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/362))
* `resource/opentelekomcloud_rds_instance_v3`: Fix backup_strategy parameter issue ([#363](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/363))


## 1.11.0 (August 01, 2019)

FEATURES:

* **New Data Source:** `opentelekomcloud_sdrs_domain_v1` ([#328](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/328))
* **New Resource:** `opentelekomcloud_sdrs_protectiongroup_v1` ([#326](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/326))

ENHANCEMENTS:

* `resource/opentelekomcloud_vpc_v1`: Add enable_shared_snat support ([#333](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/333))
* `resource/opentelekomcloud_networking_floatingip_v2`: Add default value for floating_ip pool ([#335](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/335))
* `resource/opentelekomcloud_blockstorage_volume_v2`: Add device_type argument support ([#338](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/338))
* `resource/opentelekomcloud_blockstorage_volume_v2`: Add wwn attribute support ([#339](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/339))

BUG FIXES:

* `resource/opentelekomcloud_sfs_file_system_v2`: Set availability_zone to Computed ([#330](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/330))
* `resource/opentelekomcloud_rds_configuration_v3`: Fix RDS parametergroup acc test ([#331](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/331))

## 1.10.0 (July 01, 2019)

FEATURES:

* **New Resource:** `opentelekomcloud_waf_whiteblackip_rule_v1` ([#313](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/313))
* **New Resource:** `opentelekomcloud_waf_datamasking_rule_v1` ([#315](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/315))
* **New Resource:** `opentelekomcloud_waf_falsealarmmasking_rule_v1` ([#317](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/317))
* **New Resource:** `opentelekomcloud_waf_ccattackprotection_rule_v1` ([#320](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/320))
* **New Resource:** `opentelekomcloud_waf_preciseprotection_rule_v1` ([#322](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/322))
* **New Resource:** `opentelekomcloud_waf_webtamperprotection_rule_v1` ([#324](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/324))

ENHANCEMENTS:

* `resource/opentelekomcloud_mrs_cluster_v1`: Add master/core data volume support to MRS cluster ([#308](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/308))
* `resource/opentelekomcloud_mrs_cluster_v1`: Add SAS volume type support to MRS cluster ([#310](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/310))

BUG FIXES:

* `resource/opentelekomcloud_identity_project_v3`: Fix project creation issue ([#305](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/305))

## 1.9.0 (June 06, 2019)

FEATURES:

* **New Resource:** `opentelekomcloud_waf_certificate_v1` ([#285](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/285))
* **New Resource:** `opentelekomcloud_waf_domain_v1` ([#286](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/286))
* **New Resource:** `opentelekomcloud_waf_policy_v1` ([#293](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/293))
* **New Resource:** `opentelekomcloud_rds_parametergroup_v3` ([#290](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/290))

ENHANCEMENTS:

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.
* `resource/opentelekomcloud_rds_instance_v3`: Add import support ([#274](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/274))
* `resource/opentelekomcloud_cce_node_v3`: Add private_ip attribute ([#280](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/280))

BUG FIXES:

* `resource/opentelekomcloud_cce_node_v3`: Fix eip_count issue ([#279](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/279))

## 1.8.0 (May 06, 2019)

FEATURES:

* **New Data Source:** `opentelekomcloud_networking_port_v2` ([#263](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/263))
* **New Data Source:** `opentelekomcloud_rds_flavors_v3` ([#267](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/267))
* **New Resource:** `opentelekomcloud_identity_role_v3` ([#213](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/213))
* **New Resource:** `opentelekomcloud_css_cluster_v1` ([#255](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/255))
* **New Resource:** `opentelekomcloud_rds_instance_v3` ([#267](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/267))

ENHANCEMENTS:

* `resource/opentelekomcloud_dns_zone_v2`: Add support for attaching multi routers to dns zone ([#261](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/261))
* `resource/opentelekomcloud_cce_cluster_v3`: Add authentication mode option support for CCE cluster ([#262](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/262))
* `provider`: Add security_token option for OBS federated authentication ([#264](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/264))
* `resource/opentelekomcloud_rds_instance_v1`: Add RDS tag support ([#268](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/268))

BUG FIXES:

* `resource/opentelekomcloud_dms_group_v1`: Fix wrong error message ([#260](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/260))
* `data_source/opentelekomcloud_cce_node_v3`: Fix node data source with node_id ([#265](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/265))
* `resource/opentelekomcloud_cce_node_v3`: Remove Abnormal from cce node creating target state ([#266](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/266))

## 1.7.0 (March 20, 2019)

NOTES/DEPRECATIONS:

* provider: The `region`, `tenant_id`, `domain_id`, `user_id` arguments have been deprecated and `tenant_name`, `domain_name` changed to be `required`. Please update your configurations as it might be removed in the future releases.

FEATURES:

* **New Resource:** `opentelekomcloud_identity_agency_v3` ([#232](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/232))

ENHANCEMENTS:

* `provider`: Remove region, tenant_id, domain_id, user_id parameters ([#230](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/230))
* `resource/opentelekomcloud_compute_instance_v2`: Add support of security_groups update ([#234](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/234))
* `resource/opentelekomcloud_nat_snat_rule_v2`: Add `cidr` and `source_type` parameters support ([#237](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/237))

BUG FIXES:

* `resource/opentelekomcloud_identity_role_assignment_v3`: Fix attributes set issue ([#226](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/226))
* `resource/opentelekomcloud_csbs_backup_policy_v1`: Fix csbs policies parameters issue ([#244](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/244))

## 1.6.1 (February 18, 2019)

BUG FIXES:

* `provider authentication`: Fix authentication with tenant ([#216](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/216))
* `resource/opentelekomcloud_dcs_instance_v1`: Update `password` and `engine_version` of dcs instance from Option to Required ([#217](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/217))
* `resource/opentelekomcloud_smn_topic_v2`: Fix some smn topic parameters issue ([#218](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/218))

## 1.6.0 (February 01, 2019)

FEATURES:

* **New Data Source:** `opentelekomcloud_identity_role_v3` ([#167](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/167))
* **New Data Source:** `opentelekomcloud_identity_project_v3` ([#167](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/167))
* **New Data Source:** `opentelekomcloud_identity_user_v3` ([#167](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/167))
* **New Data Source:** `opentelekomcloud_identity_group_v3` ([#167](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/167))
* **New Resource:** `opentelekomcloud_identity_project_v3` ([#167](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/167))
* **New Resource:** `opentelekomcloud_identity_role_v3` ([#167](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/167))
* **New Resource:** `opentelekomcloud_identity_role_assignment_v3` ([#167](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/167))
* **New Resource:** `opentelekomcloud_identity_user_v3` ([#167](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/167))
* **New Resource:** `opentelekomcloud_identity_group_v3` ([#167](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/167))
* **New Resource:** `opentelekomcloud_identity_group_membership_v3` ([#167](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/167))

BUG FIXES:

* `resource/opentelekomcloud_rts_stack_v1`: Re-sign for 302 redirect in ak/sk scenario ([#204](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/204))
* `resource/opentelekomcloud_elb_listener`: Fix elb listener update error for backend_port ([#209](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/209))

## 1.5.2 (January 11, 2019)

BUG FIXES:

* `resource/opentelekomcloud_compute_instance_v2`: Fix instance tag update error ([#178](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/178))
* `resource/opentelekomcloud_dns_recordset_v2`: Fix dns records update error ([#179](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/179))
* `resource/opentelekomcloud_dns_recordset_v2`: Fix dns entries re-sort issue ([#185](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/185))

## 1.5.1 (January 08, 2019)

BUG FIXES:

* Fix ak/sk authentication issue ([#176](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/176))

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
* `data_source/opentelekomcloud_images_image_v2`: Add properties filter support for images data source ([#165](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/165))
* `resource/opentelekomcloud_compute_instance_v2`: Add key/value tag support ([#169](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/169))
* `data_source/opentelekomcloud_vpc_subnet_v1`: Sort vpc subnet ids by network ip availabilities ([#171](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues/171))

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
