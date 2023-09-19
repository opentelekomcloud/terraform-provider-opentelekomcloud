package opentelekomcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/antiddos"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/as"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/bms"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/cbr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/cce"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/ces"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/csbs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/css"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/cts"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/dcaas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/dcs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/dds"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/deh"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/dis"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/dms"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/dns"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/drs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/dws"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/ecs"
	elbv2 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v2"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/evs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/fw"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/iam"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/ims"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/kms"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/lts"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/mrs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/nat"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/obs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/rds"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/rts"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/s3"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/sdrs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/sfs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/smn"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/swr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/tms"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/vbs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/vpc"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/vpcep"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/vpn"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/waf"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/version"
)

// Provider returns a schema.Provider for OpenTelekomCloud.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ACCESS_KEY", ""),
				Description: common.Descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SECRET_KEY", ""),
				Description: common.Descriptions["secret_key"],
			},
			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", ""),
				Description: common.Descriptions["auth_url"],
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: common.Descriptions["region"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_REGION_NAME",
					"OS_REGION",
				}, ""),
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: common.Descriptions["user_name"],
			},
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: common.Descriptions["user_id"],
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: common.Descriptions["tenant_id"],
			},
			"tenant_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
				Description: common.Descriptions["tenant_name"],
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: common.Descriptions["password"],
			},
			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TOKEN",
					"OS_AUTH_TOKEN",
				}, ""),
				Description: common.Descriptions["token"],
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SECURITY_TOKEN", ""),
				Description: common.Descriptions["security_token"],
			},
			"passcode": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSCODE", ""),
				Description: common.Descriptions["passcode"],
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_ID",
					"OS_PROJECT_DOMAIN_ID",
					"OS_DOMAIN_ID",
				}, ""),
				Description: common.Descriptions["domain_id"],
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_NAME",
					"OS_PROJECT_DOMAIN_NAME",
					"OS_DOMAIN_NAME",
					"OS_DEFAULT_DOMAIN",
				}, ""),
				Description: common.Descriptions["domain_name"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", false),
				Description: common.Descriptions["insecure"],
			},
			"endpoint_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ENDPOINT_TYPE", ""),
			},
			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: common.Descriptions["cacert_file"],
			},
			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: common.Descriptions["cert"],
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: common.Descriptions["key"],
			},
			"swauth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SWAUTH", false),
				Description: common.Descriptions["swauth"],
			},
			"agency_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AGENCY_NAME", ""),
				Description: common.Descriptions["agency_name"],
			},
			"agency_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AGENCY_DOMAIN_NAME", ""),
				Description: common.Descriptions["agency_domain_name"],
			},
			"delegated_project": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DELEGATED_PROJECT", ""),
				Description: common.Descriptions["delegated_project"],
			},
			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CLOUD", ""),
				Description: common.Descriptions["cloud"],
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: common.Descriptions["max_retries"],
			},
			"max_backoff_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_MAX_BACKOFF_RETRIES", 5),
				Description: common.Descriptions["max_backoff_retries"],
			},
			"backoff_retry_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_BACKOFF_RETRY_TIMEOUT", 60),
				Description: common.Descriptions["backoff_retry_timeout"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"opentelekomcloud_antiddos_v1":                     antiddos.DataSourceAntiDdosV1(),
			"opentelekomcloud_cbr_backup_v3":                   cbr.DataSourceCBRBackupsV3(),
			"opentelekomcloud_cbr_backup_ids_v3":               cbr.DataSourceCBRBackupsIdsV3(),
			"opentelekomcloud_cce_cluster_v3":                  cce.DataSourceCCEClusterV3(),
			"opentelekomcloud_cce_cluster_kubeconfig_v3":       cce.DataSourceCCEClusterKubeConfigV3(),
			"opentelekomcloud_cce_addon_template_v3":           cce.DataSourceCceAddonTemplateV3(),
			"opentelekomcloud_cce_addon_templates_v3":          cce.DataSourceCceAddonTemplatesV3(),
			"opentelekomcloud_cce_node_ids_v3":                 cce.DataSourceCceNodeIdsV3(),
			"opentelekomcloud_cce_node_v3":                     cce.DataSourceCceNodesV3(),
			"opentelekomcloud_compute_availability_zones_v2":   ecs.DataSourceComputeAvailabilityZonesV2(),
			"opentelekomcloud_compute_bms_flavors_v2":          bms.DataSourceBMSFlavorV2(),
			"opentelekomcloud_compute_bms_keypairs_v2":         bms.DataSourceBMSKeyPairV2(),
			"opentelekomcloud_compute_bms_nic_v2":              bms.DataSourceBMSNicV2(),
			"opentelekomcloud_compute_bms_server_v2":           bms.DataSourceBMSServersV2(),
			"opentelekomcloud_compute_flavor_v2":               ecs.DataSourceComputeFlavorV2(),
			"opentelekomcloud_compute_instance_v2":             ecs.DataSourceComputeInstanceV2(),
			"opentelekomcloud_compute_instances_v2":            ecs.DataSourceComputeInstancesV2(),
			"opentelekomcloud_csbs_backup_v1":                  csbs.DataSourceCSBSBackupV1(),
			"opentelekomcloud_csbs_backup_policy_v1":           csbs.DataSourceCSBSBackupPolicyV1(),
			"opentelekomcloud_css_flavor_v1":                   css.DataSourceCSSFlavorV1(),
			"opentelekomcloud_cts_tracker_v1":                  cts.DataSourceCTSTrackerV1(),
			"opentelekomcloud_direct_connect_v2":               dcaas.DataSourceDirectConnectV2(),
			"opentelekomcloud_dcs_az_v1":                       dcs.DataSourceDcsAZV1(),
			"opentelekomcloud_dcs_maintainwindow_v1":           dcs.DataSourceDcsMaintainWindowV1(),
			"opentelekomcloud_dcs_product_v1":                  dcs.DataSourceDcsProductV1(),
			"opentelekomcloud_deh_host_v1":                     deh.DataSourceDEHHostV1(),
			"opentelekomcloud_deh_server_v1":                   deh.DataSourceDEHServersV1(),
			"opentelekomcloud_dds_flavors_v3":                  dds.DataSourceDdsFlavorV3(),
			"opentelekomcloud_dds_instance_v3":                 dds.DataSourceDdsInstanceV3(),
			"opentelekomcloud_dms_az_v1":                       dms.DataSourceDmsAZV1(),
			"opentelekomcloud_dms_product_v1":                  dms.DataSourceDmsProductV1(),
			"opentelekomcloud_dms_maintainwindow_v1":           dms.DataSourceDmsMaintainWindowV1(),
			"opentelekomcloud_dns_nameservers_v2":              dns.DataSourceDNSNameserversV2(),
			"opentelekomcloud_dns_zone_v2":                     dns.DataSourceDNSZoneV2(),
			"opentelekomcloud_dws_flavors_v2":                  dws.DataSourceDwsFlavorsV2(),
			"opentelekomcloud_identity_agency_v3":              iam.DataSourceIdentityAgencyV3(),
			"opentelekomcloud_identity_auth_scope_v3":          iam.DataSourceIdentityAuthScopeV3(),
			"opentelekomcloud_identity_credential_v3":          iam.DataSourceIdentityCredentialV3(),
			"opentelekomcloud_identity_group_v3":               iam.DataSourceIdentityGroupV3(),
			"opentelekomcloud_identity_project_v3":             iam.DataSourceIdentityProjectV3(),
			"opentelekomcloud_identity_projects_v3":            iam.DataSourceIdentityProjectsV3(),
			"opentelekomcloud_identity_role_v3":                iam.DataSourceIdentityRoleV3(),
			"opentelekomcloud_identity_role_custom_v3":         iam.DataSourceIdentityRoleCustomV3(),
			"opentelekomcloud_identity_user_v3":                iam.DataSourceIdentityUserV3(),
			"opentelekomcloud_images_image_v2":                 ims.DataSourceImagesImageV2(),
			"opentelekomcloud_kms_key_v1":                      kms.DataSourceKmsKeyV1(),
			"opentelekomcloud_kms_data_key_v1":                 kms.DataSourceKmsDataKeyV1(),
			"opentelekomcloud_lb_certificate_v3":               elbv3.DataSourceCertificateV3(),
			"opentelekomcloud_lb_flavor_v3":                    elbv3.DataSourceLBFlavorV3(),
			"opentelekomcloud_lb_flavors_v3":                   elbv3.DataSourceLBFlavorsV3(),
			"opentelekomcloud_lb_loadbalancer_v3":              elbv3.DataSourceLoadBalancerV3(),
			"opentelekomcloud_lb_listener_v3":                  elbv3.DataSourceListenerV3(),
			"opentelekomcloud_lb_member_ids_v2":                elbv2.DataSourceLBMemberIDsV2(),
			"opentelekomcloud_nat_gateway_v2":                  nat.DataSourceNatGatewayV2(),
			"opentelekomcloud_networking_network_v2":           vpc.DataSourceNetworkingNetworkV2(),
			"opentelekomcloud_networking_port_v2":              vpc.DataSourceNetworkingPortV2(),
			"opentelekomcloud_networking_secgroup_v2":          vpc.DataSourceNetworkingSecGroupV2(),
			"opentelekomcloud_networking_secgroup_rule_ids_v2": vpc.DataSourceNetworkingSecGroupRuleIdsV2(),
			"opentelekomcloud_obs_bucket":                      obs.DataSourceObsBucket(),
			"opentelekomcloud_obs_bucket_object":               obs.DataSourceObsBucketObject(),
			"opentelekomcloud_rds_instance_v3":                 rds.DataSourceRdsInstanceV3(),
			"opentelekomcloud_rds_backup_v3":                   rds.DataSourceRDSv3Backup(),
			"opentelekomcloud_rds_flavors_v1":                  rds.DataSourceRdsFlavorV1(),
			"opentelekomcloud_rds_flavors_v3":                  rds.DataSourceRdsFlavorV3(),
			"opentelekomcloud_rds_versions_v3":                 rds.DataSourceRdsVersionsV3(),
			"opentelekomcloud_rts_software_deployment_v1":      rts.DataSourceRtsSoftwareDeploymentV1(),
			"opentelekomcloud_rts_software_config_v1":          rts.DataSourceRtsSoftwareConfigV1(),
			"opentelekomcloud_rts_stack_resource_v1":           rts.DataSourceRTSStackResourcesV1(),
			"opentelekomcloud_rts_stack_v1":                    rts.DataSourceRTSStackV1(),
			"opentelekomcloud_s3_bucket_object":                s3.DataSourceS3BucketObject(),
			"opentelekomcloud_sfs_file_system_v2":              sfs.DataSourceSFSFileSystemV2(),
			"opentelekomcloud_sfs_turbo_share_v1":              sfs.DataSourceSFSTurboShareV1(),
			"opentelekomcloud_sdrs_domain_v1":                  sdrs.DataSourceSdrsDomainV1(),
			"opentelekomcloud_vpc_eip_v1":                      vpc.DataSourceVPCEipV1(),
			"opentelekomcloud_vpc_v1":                          vpc.DataSourceVirtualPrivateCloudVpcV1(),
			"opentelekomcloud_vpc_bandwidth":                   vpc.DataSourceBandWidth(),
			"opentelekomcloud_vpc_bandwidth_v2":                vpc.DataSourceBandWidthV2(),
			"opentelekomcloud_vpcep_public_service_v1":         vpcep.DataSourceVPCEPPublicServiceV1(),
			"opentelekomcloud_vpcep_service_v1":                vpcep.DataSourceVPCEPServiceV1(),
			"opentelekomcloud_vbs_backup_v2":                   vbs.DataSourceVBSBackupV2(),
			"opentelekomcloud_vbs_backup_policy_v2":            vbs.DataSourceVBSBackupPolicyV2(),
			"opentelekomcloud_vpc_peering_connection_v2":       vpc.DataSourceVpcPeeringConnectionV2(),
			"opentelekomcloud_vpc_route_v2":                    vpc.DataSourceVPCRouteV2(),
			"opentelekomcloud_vpc_route_ids_v2":                vpc.DataSourceVPCRouteIdsV2(),
			"opentelekomcloud_vpc_subnet_v1":                   vpc.DataSourceVpcSubnetV1(),
			"opentelekomcloud_vpc_subnet_ids_v1":               vpc.DataSourceVpcSubnetIdsV1(),
			"opentelekomcloud_vpnaas_service_v2":               vpn.DataSourceVpnServiceV2(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"opentelekomcloud_antiddos_v1":                               antiddos.ResourceAntiDdosV1(),
			"opentelekomcloud_as_configuration_v1":                       as.ResourceASConfiguration(),
			"opentelekomcloud_as_group_v1":                               as.ResourceASGroup(),
			"opentelekomcloud_as_policy_v1":                              as.ResourceASPolicy(),
			"opentelekomcloud_as_policy_v2":                              as.ResourceASPolicyV2(),
			"opentelekomcloud_blockstorage_volume_v2":                    evs.ResourceBlockStorageVolumeV2(),
			"opentelekomcloud_cbr_policy_v3":                             cbr.ResourceCBRPolicyV3(),
			"opentelekomcloud_cbr_vault_v3":                              cbr.ResourceCBRVaultV3(),
			"opentelekomcloud_cce_addon_v3":                              cce.ResourceCCEAddonV3(),
			"opentelekomcloud_cce_cluster_v3":                            cce.ResourceCCEClusterV3(),
			"opentelekomcloud_cce_node_v3":                               cce.ResourceCCENodeV3(),
			"opentelekomcloud_cce_node_pool_v3":                          cce.ResourceCCENodePoolV3(),
			"opentelekomcloud_ces_alarmrule":                             ces.ResourceAlarmRule(),
			"opentelekomcloud_compute_bms_server_v2":                     bms.ResourceComputeBMSInstanceV2(),
			"opentelekomcloud_compute_bms_tags_v2":                       bms.ResourceBMSTagsV2(),
			"opentelekomcloud_compute_secgroup_v2":                       ecs.ResourceComputeSecGroupV2(),
			"opentelekomcloud_compute_servergroup_v2":                    ecs.ResourceComputeServerGroupV2(),
			"opentelekomcloud_compute_floatingip_v2":                     ecs.ResourceComputeFloatingIPV2(),
			"opentelekomcloud_compute_floatingip_associate_v2":           ecs.ResourceComputeFloatingIPAssociateV2(),
			"opentelekomcloud_compute_instance_v2":                       ecs.ResourceComputeInstanceV2(),
			"opentelekomcloud_compute_keypair_v2":                        ecs.ResourceComputeKeypairV2(),
			"opentelekomcloud_compute_volume_attach_v2":                  ecs.ResourceComputeVolumeAttachV2(),
			"opentelekomcloud_csbs_backup_v1":                            csbs.ResourceCSBSBackupV1(),
			"opentelekomcloud_csbs_backup_policy_v1":                     csbs.ResourceCSBSBackupPolicyV1(),
			"opentelekomcloud_cts_event_notification_v3":                 cts.ResourceCTSEventNotificationV3(),
			"opentelekomcloud_cts_tracker_v1":                            cts.ResourceCTSTrackerV1(),
			"opentelekomcloud_cts_tracker_v3":                            cts.ResourceCTSTrackerV3(),
			"opentelekomcloud_css_cluster_v1":                            css.ResourceCssClusterV1(),
			"opentelekomcloud_css_snapshot_configuration_v1":             css.ResourceCssSnapshotConfigurationV1(),
			"opentelekomcloud_direct_connect_v2":                         dcaas.ResourceDirectConnectV2(),
			"opentelekomcloud_dcs_instance_v1":                           dcs.ResourceDcsInstanceV1(),
			"opentelekomcloud_dds_instance_v3":                           dds.ResourceDdsInstanceV3(),
			"opentelekomcloud_deh_host_v1":                               deh.ResourceDeHHostV1(),
			"opentelekomcloud_dis_stream_v2":                             dis.ResourceDisStreamV2(),
			"opentelekomcloud_dis_app_v2":                                dis.ResourceDisAppV2(),
			"opentelekomcloud_dis_checkpoint_v2":                         dis.ResourceDisCheckpointV2(),
			"opentelekomcloud_dis_dump_task_v2":                          dis.ResourceDisDumpV2(),
			"opentelekomcloud_dns_ptrrecord_v2":                          dns.ResourceDNSPtrRecordV2(),
			"opentelekomcloud_dns_recordset_v2":                          dns.ResourceDNSRecordSetV2(),
			"opentelekomcloud_dns_zone_v2":                               dns.ResourceDNSZoneV2(),
			"opentelekomcloud_dms_instance_v1":                           dms.ResourceDmsInstancesV1(),
			"opentelekomcloud_dms_instance_v2":                           dms.ResourceDmsInstancesV2(),
			"opentelekomcloud_dms_topic_v1":                              dms.ResourceDmsTopicsV1(),
			"opentelekomcloud_dms_user_v2":                               dms.ResourceDmsUsersV2(),
			"opentelekomcloud_dms_user_permission_v1":                    dms.ResourceDmsUsersPermissionV1(),
			"opentelekomcloud_drs_task_v3":                               drs.ResourceDrsTaskV3(),
			"opentelekomcloud_dws_cluster_v1":                            dws.ResourceDcsInstanceV1(),
			"opentelekomcloud_ecs_instance_v1":                           ecs.ResourceEcsInstanceV1(),
			"opentelekomcloud_evs_volume_v3":                             evs.ResourceEvsStorageVolumeV3(),
			"opentelekomcloud_fw_firewall_group_v2":                      fw.ResourceFWFirewallGroupV2(),
			"opentelekomcloud_fw_policy_v2":                              fw.ResourceFWPolicyV2(),
			"opentelekomcloud_fw_rule_v2":                                fw.ResourceFWRuleV2(),
			"opentelekomcloud_identity_agency_v3":                        iam.ResourceIdentityAgencyV3(),
			"opentelekomcloud_identity_credential_v3":                    iam.ResourceIdentityCredentialV3(),
			"opentelekomcloud_identity_group_v3":                         iam.ResourceIdentityGroupV3(),
			"opentelekomcloud_identity_group_membership_v3":              iam.ResourceIdentityGroupMembershipV3(),
			"opentelekomcloud_identity_mapping_v3":                       iam.ResourceIdentityMappingV3(),
			"opentelekomcloud_identity_project_v3":                       iam.ResourceIdentityProjectV3(),
			"opentelekomcloud_identity_protocol_v3":                      iam.ResourceIdentityProtocolV3(),
			"opentelekomcloud_identity_provider_v3":                      iam.ResourceIdentityProviderV3(),
			"opentelekomcloud_identity_provider":                         iam.ResourceIdentityProvider(),
			"opentelekomcloud_identity_role_v3":                          iam.ResourceIdentityRoleV3(),
			"opentelekomcloud_identity_role_assignment_v3":               iam.ResourceIdentityRoleAssignmentV3(),
			"opentelekomcloud_identity_user_group_membership_v3":         iam.ResourceIdentityUserGroupMembershipV3(),
			"opentelekomcloud_identity_user_v3":                          iam.ResourceIdentityUserV3(),
			"opentelekomcloud_images_image_access_accept_v2":             ims.ResourceImagesImageAccessAcceptV2(),
			"opentelekomcloud_images_image_access_v2":                    ims.ResourceImagesImageAccessV2(),
			"opentelekomcloud_images_image_v2":                           ims.ResourceImagesImageV2(),
			"opentelekomcloud_ims_data_image_v2":                         ims.ResourceImsDataImageV2(),
			"opentelekomcloud_ims_image_v2":                              ims.ResourceImsImageV2(),
			"opentelekomcloud_kms_grant_v1":                              kms.ResourceKmsGrantV1(),
			"opentelekomcloud_kms_key_v1":                                kms.ResourceKmsKeyV1(),
			"opentelekomcloud_lb_certificate_v2":                         elbv2.ResourceCertificateV2(),
			"opentelekomcloud_lb_certificate_v3":                         elbv3.ResourceCertificateV3(),
			"opentelekomcloud_lb_ipgroup_v3":                             elbv3.ResourceIpGroupV3(),
			"opentelekomcloud_lb_l7policy_v2":                            elbv2.ResourceL7PolicyV2(),
			"opentelekomcloud_lb_l7rule_v2":                              elbv2.ResourceL7RuleV2(),
			"opentelekomcloud_lb_loadbalancer_v2":                        elbv2.ResourceLoadBalancerV2(),
			"opentelekomcloud_lb_loadbalancer_v3":                        elbv3.ResourceLoadBalancerV3(),
			"opentelekomcloud_lb_listener_v2":                            elbv2.ResourceListenerV2(),
			"opentelekomcloud_lb_listener_v3":                            elbv3.ResourceListenerV3(),
			"opentelekomcloud_lb_member_v2":                              elbv2.ResourceMemberV2(),
			"opentelekomcloud_lb_member_v3":                              elbv3.ResourceLBMemberV3(),
			"opentelekomcloud_lb_monitor_v2":                             elbv2.ResourceMonitorV2(),
			"opentelekomcloud_lb_monitor_v3":                             elbv3.ResourceMonitorV3(),
			"opentelekomcloud_lb_policy_v3":                              elbv3.ResourceLBPolicyV3(),
			"opentelekomcloud_lb_pool_v2":                                elbv2.ResourceLBPoolV2(),
			"opentelekomcloud_lb_pool_v3":                                elbv3.ResourceLBPoolV3(),
			"opentelekomcloud_lb_rule_v3":                                elbv3.ResourceLBRuleV3(),
			"opentelekomcloud_lb_security_policy_v3":                     elbv3.ResourceLBSecurityPolicyV3(),
			"opentelekomcloud_lb_whitelist_v2":                           elbv2.ResourceWhitelistV2(),
			"opentelekomcloud_logtank_group_v2":                          lts.ResourceLTSGroupV2(),
			"opentelekomcloud_logtank_topic_v2":                          lts.ResourceLTSTopicV2(),
			"opentelekomcloud_logtank_transfer_v2":                       lts.ResourceLTSTransferV2(),
			"opentelekomcloud_mrs_cluster_v1":                            mrs.ResourceMRSClusterV1(),
			"opentelekomcloud_mrs_job_v1":                                mrs.ResourceMRSJobV1(),
			"opentelekomcloud_nat_gateway_v2":                            nat.ResourceNatGatewayV2(),
			"opentelekomcloud_nat_dnat_rule_v2":                          nat.ResourceNatDnatRuleV2(),
			"opentelekomcloud_nat_snat_rule_v2":                          nat.ResourceNatSnatRuleV2(),
			"opentelekomcloud_networking_floatingip_v2":                  vpc.ResourceNetworkingFloatingIPV2(),
			"opentelekomcloud_networking_floatingip_associate_v2":        vpc.ResourceNetworkingFloatingIPAssociateV2(),
			"opentelekomcloud_networking_network_v2":                     vpc.ResourceNetworkingNetworkV2(),
			"opentelekomcloud_networking_port_v2":                        vpc.ResourceNetworkingPortV2(),
			"opentelekomcloud_networking_router_v2":                      vpc.ResourceNetworkingRouterV2(),
			"opentelekomcloud_networking_router_interface_v2":            vpc.ResourceNetworkingRouterInterfaceV2(),
			"opentelekomcloud_networking_router_route_v2":                vpc.ResourceNetworkingRouterRouteV2(),
			"opentelekomcloud_networking_secgroup_v2":                    vpc.ResourceNetworkingSecGroupV2(),
			"opentelekomcloud_networking_secgroup_rule_v2":               vpc.ResourceNetworkingSecGroupRuleV2(),
			"opentelekomcloud_networking_subnet_v2":                      vpc.ResourceNetworkingSubnetV2(),
			"opentelekomcloud_networking_vip_v2":                         vpc.ResourceNetworkingVIPV2(),
			"opentelekomcloud_networking_vip_associate_v2":               vpc.ResourceNetworkingVIPAssociateV2(),
			"opentelekomcloud_obs_bucket":                                obs.ResourceObsBucket(),
			"opentelekomcloud_obs_bucket_object":                         obs.ResourceObsBucketObject(),
			"opentelekomcloud_obs_bucket_policy":                         obs.ResourceObsBucketPolicy(),
			"opentelekomcloud_obs_bucket_replication":                    obs.ResourceObsBucketReplication(),
			"opentelekomcloud_rds_backup_v3":                             rds.ResourceRdsBackupV3(),
			"opentelekomcloud_rds_instance_v1":                           rds.ResourceRdsInstance(),
			"opentelekomcloud_rds_instance_v3":                           rds.ResourceRdsInstanceV3(),
			"opentelekomcloud_rds_parametergroup_v3":                     rds.ResourceRdsConfigurationV3(),
			"opentelekomcloud_rds_read_replica_v3":                       rds.ResourceRdsReadReplicaV3(),
			"opentelekomcloud_rts_software_deployment_v1":                rts.ResourceRtsSoftwareDeploymentV1(),
			"opentelekomcloud_rts_software_config_v1":                    rts.ResourceSoftwareConfigV1(),
			"opentelekomcloud_rts_stack_v1":                              rts.ResourceRTSStackV1(),
			"opentelekomcloud_s3_bucket":                                 s3.ResourceS3Bucket(),
			"opentelekomcloud_s3_bucket_policy":                          s3.ResourceS3BucketPolicy(),
			"opentelekomcloud_s3_bucket_object":                          s3.ResourceS3BucketObject(),
			"opentelekomcloud_sfs_file_system_v2":                        sfs.ResourceSFSFileSystemV2(),
			"opentelekomcloud_sfs_share_access_rules_v2":                 sfs.ResourceSFSShareAccessRulesV2(),
			"opentelekomcloud_sfs_turbo_share_v1":                        sfs.ResourceSFSTurboShareV1(),
			"opentelekomcloud_smn_topic_v2":                              smn.ResourceTopic(),
			"opentelekomcloud_smn_topic_attribute_v2":                    smn.ResourceSMNTopicAttributeV2(),
			"opentelekomcloud_smn_subscription_v2":                       smn.ResourceSubscription(),
			"opentelekomcloud_swr_domain_v2":                             swr.ResourceSwrDomainV2(),
			"opentelekomcloud_swr_organization_permissions_v2":           swr.ResourceSwrOrganizationPermissionsV2(),
			"opentelekomcloud_swr_organization_v2":                       swr.ResourceSwrOrganizationV2(),
			"opentelekomcloud_swr_repository_v2":                         swr.ResourceSwrRepositoryV2(),
			"opentelekomcloud_tms_tags_v1":                               tms.ResourceTmsTagV1(),
			"opentelekomcloud_vpc_bandwidth_associate_v2":                vpc.ResourceBandwidthAssociateV2(),
			"opentelekomcloud_vpc_bandwidth_v2":                          vpc.ResourceBandwidthV2(),
			"opentelekomcloud_vpc_eip_v1":                                vpc.ResourceVpcEIPV1(),
			"opentelekomcloud_vpc_v1":                                    vpc.ResourceVirtualPrivateCloudV1(),
			"opentelekomcloud_vpc_peering_connection_v2":                 vpc.ResourceVpcPeeringConnectionV2(),
			"opentelekomcloud_vpc_peering_connection_accepter_v2":        vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"opentelekomcloud_vpcep_endpoint_v1":                         vpcep.ResourceVPCEPEndpointV1(),
			"opentelekomcloud_vpcep_service_v1":                          vpcep.ResourceVPCEPServiceV1(),
			"opentelekomcloud_vpc_route_v2":                              vpc.ResourceVPCRouteV2(),
			"opentelekomcloud_vpc_subnet_v1":                             vpc.ResourceVpcSubnetV1(),
			"opentelekomcloud_vpc_flow_log_v1":                           vpc.ResourceVpcFlowLogV1(),
			"opentelekomcloud_vbs_backup_policy_v2":                      vbs.ResourceVBSBackupPolicyV2(),
			"opentelekomcloud_vbs_backup_v2":                             vbs.ResourceVBSBackupV2(),
			"opentelekomcloud_vbs_backup_share_v2":                       vbs.ResourceVBSBackupShareV2(),
			"opentelekomcloud_sdrs_protected_instance_v1":                sdrs.ResourceSdrsProtectedInstanceV1(),
			"opentelekomcloud_sdrs_protectiongroup_v1":                   sdrs.ResourceSdrsProtectiongroupV1(),
			"opentelekomcloud_vpnaas_ipsec_policy_v2":                    vpn.ResourceVpnIPSecPolicyV2(),
			"opentelekomcloud_vpnaas_service_v2":                         vpn.ResourceVpnServiceV2(),
			"opentelekomcloud_vpnaas_ike_policy_v2":                      vpn.ResourceVpnIKEPolicyV2(),
			"opentelekomcloud_vpnaas_endpoint_group_v2":                  vpn.ResourceVpnEndpointGroupV2(),
			"opentelekomcloud_vpnaas_site_connection_v2":                 vpn.ResourceVpnSiteConnectionV2(),
			"opentelekomcloud_waf_alarm_notification_v1":                 waf.ResourceWafAlarmNotificationV1(),
			"opentelekomcloud_waf_certificate_v1":                        waf.ResourceWafCertificateV1(),
			"opentelekomcloud_waf_domain_v1":                             waf.ResourceWafDomainV1(),
			"opentelekomcloud_waf_policy_v1":                             waf.ResourceWafPolicyV1(),
			"opentelekomcloud_waf_whiteblackip_rule_v1":                  waf.ResourceWafWhiteBlackIpRuleV1(),
			"opentelekomcloud_waf_datamasking_rule_v1":                   waf.ResourceWafDataMaskingRuleV1(),
			"opentelekomcloud_waf_falsealarmmasking_rule_v1":             waf.ResourceWafFalseAlarmMaskingRuleV1(),
			"opentelekomcloud_waf_ccattackprotection_rule_v1":            waf.ResourceWafCcAttackProtectionRuleV1(),
			"opentelekomcloud_waf_preciseprotection_rule_v1":             waf.ResourceWafPreciseProtectionRuleV1(),
			"opentelekomcloud_waf_webtamperprotection_rule_v1":           waf.ResourceWafWebTamperProtectionRuleV1(),
			"opentelekomcloud_waf_dedicated_instance_v1":                 waf.ResourceWafDedicatedInstance(),
			"opentelekomcloud_waf_dedicated_domain_v1":                   waf.ResourceWafDedicatedDomain(),
			"opentelekomcloud_waf_dedicated_policy_v1":                   waf.ResourceWafDedicatedPolicy(),
			"opentelekomcloud_waf_dedicated_certificate_v1":              waf.ResourceWafDedicatedCertificateV1(),
			"opentelekomcloud_waf_dedicated_cc_rule_v1":                  waf.ResourceWafDedicatedCcRuleV1(),
			"opentelekomcloud_waf_dedicated_anti_crawler_rule_v1":        waf.ResourceWafDedicatedAntiCrawlerRuleV1(),
			"opentelekomcloud_waf_dedicated_data_masking_rule_v1":        waf.ResourceWafDedicatedDataMaskingRuleV1(),
			"opentelekomcloud_waf_dedicated_known_attack_source_rule_v1": waf.ResourceWafDedicatedKnownAttackSourceRuleV1(),
			"opentelekomcloud_waf_dedicated_web_tamper_rule_v1":          waf.ResourceWafDedicatedWebTamperRuleV1(),
			"opentelekomcloud_waf_dedicated_anti_leakage_rule_v1":        waf.ResourceWafDedicatedAntiLeakageRuleV1(),
			"opentelekomcloud_waf_dedicated_alarm_masking_rule_v1":       waf.ResourceWafDedicatedAlarmMaskingRuleV1(),
			"opentelekomcloud_waf_dedicated_geo_ip_rule_v1":              waf.ResourceWafDedicatedGeoIpRuleV1(),
			"opentelekomcloud_waf_dedicated_blacklist_rule_v1":           waf.ResourceWafDedicatedBlacklistRuleV1(),
			"opentelekomcloud_waf_dedicated_precise_protection_rule_v1":  waf.ResourceWafDedicatedPreciseProtectionRuleV1(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, d, provider)
	}

	return provider
}

func providerConfigure(_ context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
	config := cfg.Config{
		AccessKey:           d.Get("access_key").(string),
		SecretKey:           d.Get("secret_key").(string),
		CACertFile:          d.Get("cacert_file").(string),
		ClientCertFile:      d.Get("cert").(string),
		ClientKeyFile:       d.Get("key").(string),
		Cloud:               d.Get("cloud").(string),
		DomainID:            d.Get("domain_id").(string),
		DomainName:          d.Get("domain_name").(string),
		EndpointType:        d.Get("endpoint_type").(string),
		IdentityEndpoint:    d.Get("auth_url").(string),
		Insecure:            d.Get("insecure").(bool),
		Password:            d.Get("password").(string),
		Passcode:            d.Get("passcode").(string),
		Region:              d.Get("region").(string),
		Swauth:              d.Get("swauth").(bool),
		Token:               d.Get("token").(string),
		SecurityToken:       d.Get("security_token").(string),
		TenantID:            d.Get("tenant_id").(string),
		TenantName:          d.Get("tenant_name").(string),
		Username:            d.Get("user_name").(string),
		UserID:              d.Get("user_id").(string),
		AgencyName:          d.Get("agency_name").(string),
		AgencyDomainName:    d.Get("agency_domain_name").(string),
		DelegatedProject:    d.Get("delegated_project").(string),
		MaxRetries:          d.Get("max_retries").(int),
		MaxBackoffRetries:   d.Get("max_backoff_retries").(int),
		BackoffRetryTimeout: d.Get("backoff_retry_timeout").(int),
		UserAgent:           p.UserAgent("terraform-provider-opentelekomcloud", version.ProviderVersion),
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, diag.FromErr(err)
	}

	return &config, nil
}
