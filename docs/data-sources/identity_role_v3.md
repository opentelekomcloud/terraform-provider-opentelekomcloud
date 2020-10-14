---
subcategory: "Identity and Access Management (IAM)"
---

# opentelekomcloud_identity_role_v3

Use this data source to get the ID of an OpenTelekomCloud role.

The Role in Terraform is the same as Policy on the console. however,
The policy name is the display name of Role, the Role name cannot
be found on Console. please refer to the following table to configuration
Role:

Role Name | Policy Name
---- | ----
readonly | Tenant Guest
tms_adm | TMS Administrator
cce_adm | CCE Administrator
dcs_admin | DCS Administrator
dis_adm | DIS Administrator
system_all_8 | CCE Admin
system_all_15 | DDS FullAccess
system_all_6 | VPC Viewer
rds_adm | RDS Administrator
system_all_1001 | Full Access
system_all_3 | EVS Viewer
te_agency | Agent Operator
dms_adm | DMS Administrator
ces_adm | CES Administrator
system_all_9 | EPS Viewer
rts_adm | RTS Administrator
obs_b_list | OBS Buckets
system_all_11 | RDS ManageAccess
system_all_5 | VPC Admin
dns_adm | DNS Administrator
system_all_12 | RDS FullAccess
server_adm | Server Administrator
system_all_10 | EPS Admin
sdrs_adm | SDRS Administrator
system_all_14 | DDS ManageAccess
system_all_0 | ECS Admin
wks_adm | Workspace Administrator
te_admin | Tenant Administrator
waf_adm | WAF Administrator
system_all_7 | CCE Viewer
system_all_17 | DAS FullAccess
sfs_adm | SFS Administrator
vpc_netadm | VPC Administrator
css_adm | CSS Administrator
as_adm | AutoScaling Administrator
system_all_16 | DDS ReadOnlyAccess
csbs_adm | CSBS Administrator
swr_adm | SWR Administrator
das_admin | DAS Administrator
system_all_13 | RDS ReadOnlyAccess
secu_admin | Security Administrator
system_all_2 | ECS Viewer
dws_adm | DWS Administrator
mobs_adm | MaaS OBS
vbs_adm | VBS Administrator
ddos_adm | Anti-DDoS Administrator
system_all_4 | EVS Admin
system_all_1 | ECS User
dws_db_acc | DWS Database
kms_adm | KMS Administrator
mrs_adm | MRS Administrator
nat_adm | NAT Gateway
dds_adm | DDS Administrator
ims_adm | IMS Administrator
smn_adm | SMN Administrator
plas_adm | Config Plas
elb_adm | ELB Administrator
vpcep_adm | VPCEndpoint Administrator

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
