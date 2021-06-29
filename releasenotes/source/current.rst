=====================
Current Release Notes
=====================

.. release-notes::

1.24.2
------

New Features
============

* `resource/opentelekomcloud_waf_certificate_v1`: Add possibility to configure TLS ([#1132](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/1132))
* `resource/opentelekomcloud_ecs_instance_v1`: Add possibility to encrypt `data_volumes` ([#1135](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/1135))

Bug Fixes
=========

* `provider`: Store temporary AK/SK in OBS client only ([#1133](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/1133))
* `resource/opentelekomcloud_rds_instance_v3`: Fix not setting `availability_zone` during import/refresh ([#1137](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/1137))
* `resource/opentelekomcloud_obs_bucket`: Fix not creating bucket for `eu-nl` region ([#1139](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/1139))
* `resource/opentelekomcloud_rds_instance_v3`: Fix configuration template applied partially ([#1150](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/1150))

Other Notes
===========

* `resources/opentelekomcloud_nat_gateway_v2`: Note that `router_id` can be a VPC ID ([#1140](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/1140))

1.24.1
------

New Features
============

* **New Resource:** `opentelekomcloud_swr_organization_v2` ([#1120](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/1120))
* **New Resource:** `opentelekomcloud_swr_repository_v2` ([#1123](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/1123))
* **New Resource:** `opentelekomcloud_swr_domain_v2` ([#1127](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/1127))
* **New Resource:** `opentelekomcloud_swr_organization_permissions_v2` ([#1129](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/1129))
* `resource/opentelekomcloud_waf_certificate_v1`: Add name certificate name validation ([#1116](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/1116))
* `resource/opentelekomcloud_cce_node_v3`: Add possibility to encrypt `data_volumes` ([#1117](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/1117))
* `resource/opentelekomcloud_waf_certificate_v1`: Add possibility to import by name ([#1119](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/1119))
* `resource/opentelekomcloud_obs_bucket_v1`: Add possibility to create buckets with encryption ([#1125](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/1125))

Bug Fixes
=========

* `resource/opentelekomcloud_as_configuration_v1`: Fix wrong SHA1 sum in `Read` operation ([#1130](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/1130))

1.24.0
------

Deprecation Notes
=================

* Move to `terraform-plugin-sdk/v2`, only Terraform versions `0.12+` are supported ([#1104](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues/1104))

Other Notes
===========

This release contains no functional differences from `1.23.13`.

