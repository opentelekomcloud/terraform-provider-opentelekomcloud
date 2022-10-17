---
subcategory: "Identity and Access Management (IAM)"
---

# opentelekomcloud_identity_role_v3

Use this data source to get the ID of an OpenTelekomCloud role.

-> For custom user roles usage please refer to `opentelekomcloud_identity_role_custom_v3`

The Role in Terraform is the same as Policy on the console. however,
The policy name is the display name of Role, the Role name cannot
be found on Console. please refer to the following table to configuration
Role:

Role Name | Policy Name
---- | ----
readonly | Tenant Guest
system_all_1001 | Full Access
system_all_39 | KMS CMKReadOnlyAccess
system_all_38 | KMS CMKFullAccess
system_all_37 | OBS OperateAccess
system_all_36 | OBS ReadOnlyAccess
system_all_35 | OBS Administrator
system_all_34 | OBS Administrator
system_all_33 | EVS ReadOnlyAccess
system_all_32 | EVS FullAccess
system_all_31 | OBS OperateAccess
system_all_30 | OBS ReadOnlyAccess
system_all_29 | OBS Administrator
system_all_28 | CCE ReadOnlyAccess
system_all_27 | CCE FullAccess
system_all_26 | SWR CommonOperations
system_all_25 | ModelArts FullAccess
system_all_24 | IAM ReadOnlyAccess
system_all_23 | DRS ReadOnlyAccess
system_all_22 | DRS FullAccess
system_all_21 | GaussDB NoSQL ReadOnlyAccess
system_all_20 | GaussDB NoSQL FullAccess
system_all_19 | GaussDB FullAccess
system_all_18 | GaussDB ReadOnlyAccess
system_all_17 | DAS FullAccess
system_all_16 | DDS ReadOnlyAccess
system_all_15 | DDS FullAccess
system_all_14 | DDS ManageAccess
system_all_13 | RDS ReadOnlyAccess
system_all_12 | RDS FullAccess
system_all_11 | RDS ManageAccess
system_all_10 | EPS Admin
system_all_9 | EPS Viewer
system_all_8 | CCE Admin
system_all_7 | CCE Viewer
system_all_6 | VPC Viewer
system_all_5 | VPC Admin
system_all_4 | EVS Admin
system_all_3 | EVS Viewer
system_all_2 | ECS Viewer
system_all_1 | ECS User
system_all_0 | ECS Admin
apm_adm | APM Administrator
as_adm | AutoScaling Administrator
cce_adm | CCE Administrator
ces_adm | CES Administrator
csbs_adm | CSBS Administrator
css_adm | CSS Administrator
das_admin | DAS Administrator
dcs_admin | DCS Administrator
dds_adm | DDS Administrator
ddos_adm | Anti-DDoS Administrator
dis_adm | DIS Administrator
dms_adm | DMS Administrator
dns_adm | DNS Administrator
dws_adm | DWS Administrator
dws_db_acc | DWS Database
elb_adm | ELB Administrator
ims_adm | IMS Administrator
kms_adm | KMS Administrator
lts_admin | LTS Administrator
mrs_adm | MRS Administrator
nat_adm | NAT Gateway Administrator
obs_b_list | OBS Buckets Viewer
plas_adm | Config Plas Connector
rds_adm | RDS Administrator
rts_adm | RTS Administrator
sdrs_adm | SDRS Administrator
secu_admin | Security Administrator
server_adm | Server Administrator
sfs_adm | SFS Administrator
smn_adm | SMN Administrator
swr_adm | SWR Administrator
te_admin | Tenant Administrator
te_agency | Agent Operator
tms_adm | TMS Administrator
vbs_adm | VBS Administrator
vpc_netadm | VPC Administrator
vpcep_adm | VPCEndpoint Administrator
vpn_adm | VPN Administrator
waf_adm | WAF Administrator
wks_adm | Workspace Administrator

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
