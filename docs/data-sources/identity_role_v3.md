---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_role_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-role-v3"
description: |-
Get a IAM role information from OpenTelekomCloud
---

Up-to-date reference of API arguments for IAM role you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/permission_management/querying_a_role_list.html#en-us-topic-0057845591)

# opentelekomcloud_identity_role_v3

Use this data source to get the ID of an OpenTelekomCloud role.

-> For custom user roles usage please refer to `opentelekomcloud_identity_role_custom_v3`

The Role in Terraform is the same as Policy on the console. however,
The policy name is the display name of Role, the Role name cannot
be found on Console. please refer to the following table to configuration
Role:

EU-DE/EU-NL REGIONS

| Role Name       | Policy Name                     |
|-----------------|---------------------------------|
| system_all_1001 | Full Access                     |
| system_all_0    | ECS Admin                       |
| system_all_1    | ECS User                        |
| system_all_2    | ECS Viewer                      |
| system_all_3    | EVS Viewer                      |
| system_all_4    | EVS Admin                       |
| system_all_5    | VPC Admin                       |
| system_all_6    | VPC Viewer                      |
| system_all_7    | CCE Viewer                      |
| system_all_8    | CCE Admin                       |
| system_all_9    | EPS Viewer                      |
| system_all_10   | EPS Admin                       |
| system_all_11   | RDS ManageAccess                |
| system_all_12   | RDS FullAccess                  |
| system_all_13   | RDS ReadOnlyAccess              |
| system_all_14   | DDS ManageAccess                |
| system_all_15   | DDS FullAccess                  |
| system_all_16   | DDS ReadOnlyAccess              |
| system_all_17   | DAS FullAccess                  |
| system_all_18   | GaussDB ReadOnlyAccess          |
| system_all_19   | GaussDB FullAccess              |
| system_all_20   | GaussDB NoSQL FullAccess        |
| system_all_21   | GaussDB NoSQL ReadOnlyAccess    |
| system_all_22   | DRS FullAccess                  |
| system_all_23   | DRS ReadOnlyAccess              |
| system_all_24   | IAM ReadOnlyAccess              |
| system_all_25   | ModelArts FullAccess            |
| system_all_26   | ModelArts CommonOperations      |
| system_all_27   | CCE FullAccess                  |
| system_all_28   | CCE ReadOnlyAccess              |
| system_all_29   | OBS Administrator               |
| system_all_30   | OBS ReadOnlyAccess              |
| system_all_31   | OBS OperateAccess               |
| system_all_32   | EVS FullAccess                  |
| system_all_33   | EVS ReadOnlyAccess              |
| system_all_34   | OBS Administrator               |
| system_all_35   | OBS Administrator               |
| system_all_36   | OBS ReadOnlyAccess              |
| system_all_37   | OBS OperateAccess               |
| system_all_38   | KMS CMKFullAccess               |
| system_all_39   | KMS CMKReadOnlyAccess           |
| system_all_40   | CDM FullAccess                  |
| system_all_41   | CDM FullAccessExceptEIPUpdating |
| system_all_42   | CDM ReadOnlyAccess              |
| system_all_43   | CDM CommonOperations            |
| system_all_44   | CSS FullAccess                  |
| system_all_45   | CSS ReadOnlyAccess              |
| system_all_46   | VPC FullAccess                  |
| system_all_47   | VPC ReadOnlyAccess              |
| system_all_48   | RMS ReadOnlyAccess              |
| system_all_49   | ECS CommonOperations            |
| system_all_50   | ECS ReadOnlyAccess              |
| system_all_51   | ECS FullAccess                  |
| system_all_52   | IMS FullAccess                  |
| system_all_53   | IMS ReadOnlyAccess              |
| system_all_54   | BMS FullAccess                  |
| system_all_55   | BMS ReadOnlyAccess              |
| system_all_56   | BMS CommonOperations            |
| system_all_57   | DeH FullAccess                  |
| system_all_58   | DeH ReadOnlyAccess              |
| system_all_59   | DeH CommonOperations            |
| system_all_60   | CES ReadOnlyAccess              |
| system_all_61   | CES FullAccess                  |
| system_all_62   | SMN FullAccess                  |
| system_all_63   | SMN ReadOnlyAccess              |
| system_all_64   | WAF ReadOnlyAccess              |
| system_all_65   | WAF FullAccess                  |
| system_all_66   | CBR User                        |
| system_all_67   | CBR Administrator               |
| system_all_68   | CBR Viewer                      |
| system_all_69   | DWS FullAccess                  |
| system_all_70   | DWS ReadOnlyAccess              |
| system_all_71   | MRS ReadOnlyAccess              |
| system_all_72   | MRS FullAccess                  |
| system_all_73   | MRS CommonOperations            |
| apm_adm         | APM Administrator               |
| as_adm          | AutoScaling Administrator       |
| cce_adm         | CCE Administrator               |
| ces_adm         | CES Administrator               |
| csbs_adm        | CSBS Administrator              |
| css_adm         | CSS Administrator               |
| das_admin       | DAS Administrator               |
| dayu_adm        | DGC Administrator               |
| dayu_user       | DGC User                        |
| dcs_admin       | DCS Administrator               |
| ddos_adm        | Anti-DDoS Administrator         |
| dds_adm         | DDS Administrator               |
| dis_adm         | DIS Administrator               |
| dms_adm         | DMS Administrator               |
| dns_adm         | DNS Administrator               |
| dws_adm         | DWS Administrator               |
| dws_db_acc      | DWS Database Access             |
| elb_adm         | ELB Administrator               |
| ims_adm         | IMS Administrator               |
| kms_adm         | KMS Administrator               |
| lts_admin       | LTS Administrator               |
| mrs_adm         | MRS Administrator               |
| nat_adm         | NAT Gateway Administrator       |
| obs_b_list      | OBS Buckets Viewer              |
| plas_adm        | Config Plas Connector           |
| rds_adm         | RDS Administrator               |
| readonly        | Tenant Guest                    |
| rts_adm         | RTS Administrator               |
| sdrs_adm        | SDRS Administrator              |
| secu_admin      | Security Administrator          |
| server_adm      | Server Administrator            |
| sfs_adm         | SFS Administrator               |
| smn_adm         | SMN Administrator               |
| swr_adm         | SWR Administrator               |
| te_admin        | Tenant Administrator            |
| te_agency       | Agent Operator                  |
| tms_adm         | TMS Administrator               |
| vbs_adm         | VBS Administrator               |
| vpc_netadm      | VPC Administrator               |
| vpcep_adm       | VPCEndpoint Administrator       |
| vpn_adm         | VPN Administrator               |
| waf_adm         | WAF Administrator               |

EU-CH2 REGION

| Role Name     | Policy Name                    |
|---------------|--------------------------------|
| system_all_0  | CBR BackupsAndVaultsFullAccess |
| system_all_1  | CBR FullAccess                 |
| system_all_2  | CBR ReadOnlyAccess             |
| system_all_3  | VPC FullAccess                 |
| system_all_4  | VPC ReadOnlyAccess             |
| system_all_5  | ELB FullAccess                 |
| system_all_6  | ELB ReadOnlyAccess             |
| system_all_7  | ECS CommonOperations           |
| system_all_8  | ECS FullAccess                 |
| system_all_9  | ECS ReadOnlyAccess             |
| system_all_10 | AutoScaling FullAccess         |
| system_all_11 | AutoScaling ReadOnlyAccess     |
| system_all_12 | DNS ReadOnlyAccess             |
| system_all_13 | DCAAS ReadOnlyAccess           |
| system_all_14 | DNS FullAccess                 |
| system_all_15 | DNS ReadOnlyAccess             |
| system_all_16 | EVS FullAccess                 |
| system_all_17 | EVS ReadOnlyAccess             |
| system_all_18 | BMS FullAccess                 |
| system_all_19 | BMS ReadOnlyAccess             |
| system_all_20 | BMS CommonOperations           |
| system_all_21 | VPCEndpoint FullAccess         |
| system_all_22 | VPCEndpoint ReadOnlyAccess     |
| system_all_23 | SMN FullAccess                 |
| system_all_24 | SMN ReadOnlyAccess             |
| system_all_25 | IMS FullAccess                 |
| system_all_26 | IMS ReadOnlyAccess             |
| system_all_27 | CES FullAccess                 |
| system_all_28 | CES ReadOnlyAccess             |
| system_all_29 | WAF ReadOnlyAccess             |
| system_all_30 | WAF FullAccess                 |
| system_all_31 | OBS Administrator              |
| system_all_32 | OBS ReadOnlyAccess             |
| system_all_33 | OBS OperateAccess              |
| system_all_34 | CTS FullAccess                 |
| system_all_35 | RDS FullAccess                 |
| system_all_36 | RDS ReadOnlyAccess             |
| system_all_37 | RDS ManageAccess               |
| system_all_38 | LTS Administrator              |
| system_all_39 | LTS ReadOnlyAccess             |
| system_all_40 | LTS FullAccess                 |
| system_all_41 | AOM Admin                      |
| system_all_42 | AOM Viewer                     |
| system_all_43 | CTS ReadOnlyAccess             |
| system_all_44 | CCE FullAccess                 |
| system_all_45 | CCE ReadOnlyAccess             |
| apm_adm       | APM Admin                      |
| as_adm        | AutoScaling Administrator      |
| cce_adm       | CCE Administrator              |
| ces_adm       | CES Administrator              |
| dcaas_adm     | Direct Connect Administrator   |
| dns_adm       | DNS Administrator              |
| elb_adm       | ELB Administrator              |
| ims_adm       | IMS Administrator              |
| kms_adm       | KMS Administrator              |
| obs_b_list    | OBS Buckets Viewer             |
| readonly      | Tenant Guest                   |
| sdrs_adm      | SDRS Administrator             |
| secu_admin    | Security Administrator         |
| server_adm    | Server Administrator           |
| smn_adm       | SMN Administrator              |
| svcstg_dev    | SWR Developer                  |
| swr_adm       | SWR Administrator              |
| te_admin      | Tenant Administrator           |
| te_agency     | Agent Operator                 |
| tms_adm       | TMS Administrator              |
| vpc_netadm    | VPC Administrator              |
| vpcep_adm     | VPCEP Administrator            |
| waf_adm       | WAF Administrator              |


```hcl
data "opentelekomcloud_identity_role_v3" "auth_admin" {
  name = "secu_admin"
}
```

## Argument Reference

* `name` - (Required) The name of the role.

* `domain_id` - (Optional) The domain the role belongs to.

## Attributes Reference

`id` is set to the ID of the found role. In addition, the following attributes are exported:

* `name` - See Argument Reference above.

* `domain_id` - See Argument Reference above.
