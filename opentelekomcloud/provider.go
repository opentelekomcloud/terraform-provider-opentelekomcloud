package opentelekomcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/jinzhu/copier"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

// Provider returns a schema.Provider for OpenTelekomCloud.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ACCESS_KEY", ""),
				Description: descriptions["access_key"],
			},

			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SECRET_KEY", ""),
				Description: descriptions["secret_key"],
			},

			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", ""),
				Description: descriptions["auth_url"],
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["region"],
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", ""),
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: descriptions["user_name"],
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: descriptions["user_name"],
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},

			"tenant_name": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
				Description: descriptions["tenant_name"],
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: descriptions["password"],
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_TOKEN", ""),
				Description: descriptions["token"],
			},

			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["security_token"],
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_ID",
					"OS_PROJECT_DOMAIN_ID",
					"OS_DOMAIN_ID",
				}, ""),
				Description: descriptions["domain_id"],
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
				Description: descriptions["domain_name"],
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", false),
				Description: descriptions["insecure"],
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
				Description: descriptions["cacert_file"],
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"swauth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SWAUTH", false),
				Description: descriptions["swauth"],
			},
			"agency_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AGENCY_NAME", ""),
				Description: descriptions["agency_name"],
			},

			"agency_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AGENCY_DOMAIN_NAME", ""),
				Description: descriptions["agency_domain_name"],
			},
			"delegated_project": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DELEGATED_PROJECT", ""),
				Description: descriptions["delegated_project"],
			},
			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CLOUD", ""),
				Description: descriptions["cloud"],
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: descriptions["max_retries"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"opentelekomcloud_images_image_v2":               dataSourceImagesImageV2(),
			"opentelekomcloud_networking_network_v2":         dataSourceNetworkingNetworkV2(),
			"opentelekomcloud_networking_port_v2":            dataSourceNetworkingPortV2(),
			"opentelekomcloud_networking_secgroup_v2":        dataSourceNetworkingSecGroupV2(),
			"opentelekomcloud_s3_bucket_object":              dataSourceS3BucketObject(),
			"opentelekomcloud_kms_key_v1":                    dataSourceKmsKeyV1(),
			"opentelekomcloud_kms_data_key_v1":               dataSourceKmsDataKeyV1(),
			"opentelekomcloud_rds_flavors_v1":                dataSourceRdsFlavorV1(),
			"opentelekomcloud_rds_flavors_v3":                dataSourceRdsFlavorV3(),
			"opentelekomcloud_vpc_v1":                        dataSourceVirtualPrivateCloudVpcV1(),
			"opentelekomcloud_vpc_peering_connection_v2":     dataSourceVpcPeeringConnectionV2(),
			"opentelekomcloud_vpc_route_v2":                  dataSourceVPCRouteV2(),
			"opentelekomcloud_vpc_route_ids_v2":              dataSourceVPCRouteIdsV2(),
			"opentelekomcloud_vpc_subnet_v1":                 dataSourceVpcSubnetV1(),
			"opentelekomcloud_vpc_subnet_ids_v1":             dataSourceVpcSubnetIdsV1(),
			"opentelekomcloud_rts_software_deployment_v1":    dataSourceRtsSoftwareDeploymentV1(),
			"opentelekomcloud_rts_software_config_v1":        dataSourceRtsSoftwareConfigV1(),
			"opentelekomcloud_rts_stack_v1":                  dataSourceRTSStackV1(),
			"opentelekomcloud_rts_stack_resource_v1":         dataSourceRTSStackResourcesV1(),
			"opentelekomcloud_sfs_file_system_v2":            dataSourceSFSFileSystemV2(),
			"opentelekomcloud_deh_host_v1":                   dataSourceDEHHostV1(),
			"opentelekomcloud_deh_server_v1":                 dataSourceDEHServersV1(),
			"opentelekomcloud_vbs_backup_policy_v2":          dataSourceVBSBackupPolicyV2(),
			"opentelekomcloud_vbs_backup_v2":                 dataSourceVBSBackupV2(),
			"opentelekomcloud_compute_availability_zones_v2": dataSourceComputeAvailabilityZonesV2(),
			"opentelekomcloud_compute_bms_nic_v2":            dataSourceBMSNicV2(),
			"opentelekomcloud_compute_bms_keypairs_v2":       dataSourceBMSKeyPairV2(),
			"opentelekomcloud_compute_bms_flavors_v2":        dataSourceBMSFlavorV2(),
			"opentelekomcloud_compute_bms_server_v2":         dataSourceBMSServersV2(),
			"opentelekomcloud_csbs_backup_v1":                dataSourceCSBSBackupV1(),
			"opentelekomcloud_csbs_backup_policy_v1":         dataSourceCSBSBackupPolicyV1(),
			"opentelekomcloud_antiddos_v1":                   dataSourceAntiDdosV1(),
			"opentelekomcloud_cts_tracker_v1":                dataSourceCTSTrackerV1(),
			"opentelekomcloud_cce_node_v3":                   dataSourceCceNodesV3(),
			"opentelekomcloud_cce_node_ids_v3":               dataSourceCceNodeIdsV3(),
			"opentelekomcloud_cce_cluster_v3":                dataSourceCCEClusterV3(),
			"opentelekomcloud_dcs_az_v1":                     dataSourceDcsAZV1(),
			"opentelekomcloud_dcs_maintainwindow_v1":         dataSourceDcsMaintainWindowV1(),
			"opentelekomcloud_dcs_product_v1":                dataSourceDcsProductV1(),
			"opentelekomcloud_dms_az_v1":                     dataSourceDmsAZV1(),
			"opentelekomcloud_dms_product_v1":                dataSourceDmsProductV1(),
			"opentelekomcloud_dms_maintainwindow_v1":         dataSourceDmsMaintainWindowV1(),
			"opentelekomcloud_identity_role_v3":              dataSourceIdentityRoleV3(),
			"opentelekomcloud_identity_project_v3":           dataSourceIdentityProjectV3(),
			"opentelekomcloud_identity_user_v3":              dataSourceIdentityUserV3(),
			"opentelekomcloud_identity_group_v3":             dataSourceIdentityGroupV3(),
			"opentelekomcloud_sdrs_domain_v1":                dataSourceSdrsDomainV1(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"opentelekomcloud_blockstorage_volume_v2":             resourceBlockStorageVolumeV2(),
			"opentelekomcloud_evs_volume_v3":                      resourceEvsStorageVolumeV3(),
			"opentelekomcloud_compute_instance_v2":                resourceComputeInstanceV2(),
			"opentelekomcloud_compute_keypair_v2":                 resourceComputeKeypairV2(),
			"opentelekomcloud_compute_secgroup_v2":                resourceComputeSecGroupV2(),
			"opentelekomcloud_compute_servergroup_v2":             resourceComputeServerGroupV2(),
			"opentelekomcloud_compute_floatingip_v2":              resourceComputeFloatingIPV2(),
			"opentelekomcloud_compute_floatingip_associate_v2":    resourceComputeFloatingIPAssociateV2(),
			"opentelekomcloud_compute_volume_attach_v2":           resourceComputeVolumeAttachV2(),
			"opentelekomcloud_dns_recordset_v2":                   resourceDNSRecordSetV2(),
			"opentelekomcloud_dns_zone_v2":                        resourceDNSZoneV2(),
			"opentelekomcloud_dns_ptrrecord_v2":                   resourceDNSPtrRecordV2(),
			"opentelekomcloud_dcs_instance_v1":                    resourceDcsInstanceV1(),
			"opentelekomcloud_ecs_instance_v1":                    resourceEcsInstanceV1(),
			"opentelekomcloud_fw_firewall_group_v2":               resourceFWFirewallGroupV2(),
			"opentelekomcloud_fw_policy_v2":                       resourceFWPolicyV2(),
			"opentelekomcloud_fw_rule_v2":                         resourceFWRuleV2(),
			"opentelekomcloud_identity_agency_v3":                 resourceIdentityAgencyV3(),
			"opentelekomcloud_identity_project_v3":                resourceIdentityProjectV3(),
			"opentelekomcloud_identity_role_v3":                   resourceIdentityRoleV3(),
			"opentelekomcloud_identity_role_assignment_v3":        resourceIdentityRoleAssignmentV3(),
			"opentelekomcloud_identity_user_v3":                   resourceIdentityUserV3(),
			"opentelekomcloud_identity_group_v3":                  resourceIdentityGroupV3(),
			"opentelekomcloud_identity_group_membership_v3":       resourceIdentityGroupMembershipV3(),
			"opentelekomcloud_images_image_v2":                    resourceImagesImageV2(),
			"opentelekomcloud_ims_data_image_v2":                  resourceImsDataImageV2(),
			"opentelekomcloud_ims_image_v2":                       resourceImsImageV2(),
			"opentelekomcloud_kms_key_v1":                         resourceKmsKeyV1(),
			"opentelekomcloud_lb_certificate_v2":                  resourceCertificateV2(),
			"opentelekomcloud_lb_l7policy_v2":                     resourceL7PolicyV2(),
			"opentelekomcloud_lb_l7rule_v2":                       resourceL7RuleV2(),
			"opentelekomcloud_lb_loadbalancer_v2":                 resourceLoadBalancerV2(),
			"opentelekomcloud_lb_listener_v2":                     resourceListenerV2(),
			"opentelekomcloud_lb_pool_v2":                         resourcePoolV2(),
			"opentelekomcloud_lb_member_v2":                       resourceMemberV2(),
			"opentelekomcloud_lb_monitor_v2":                      resourceMonitorV2(),
			"opentelekomcloud_lb_whitelist_v2":                    resourceWhitelistV2(),
			"opentelekomcloud_mrs_cluster_v1":                     resourceMRSClusterV1(),
			"opentelekomcloud_mrs_job_v1":                         resourceMRSJobV1(),
			"opentelekomcloud_nat_gateway_v2":                     resourceNatGatewayV2(),
			"opentelekomcloud_nat_snat_rule_v2":                   resourceNatSnatRuleV2(),
			"opentelekomcloud_nat_dnat_rule_v2":                   resourceNatDnatRuleV2(),
			"opentelekomcloud_networking_network_v2":              resourceNetworkingNetworkV2(),
			"opentelekomcloud_networking_subnet_v2":               resourceNetworkingSubnetV2(),
			"opentelekomcloud_networking_floatingip_v2":           resourceNetworkingFloatingIPV2(),
			"opentelekomcloud_networking_floatingip_associate_v2": resourceNetworkingFloatingIPAssociateV2(),
			"opentelekomcloud_networking_port_v2":                 resourceNetworkingPortV2(),
			"opentelekomcloud_networking_router_v2":               resourceNetworkingRouterV2(),
			"opentelekomcloud_networking_router_interface_v2":     resourceNetworkingRouterInterfaceV2(),
			"opentelekomcloud_networking_router_route_v2":         resourceNetworkingRouterRouteV2(),
			"opentelekomcloud_networking_secgroup_v2":             resourceNetworkingSecGroupV2(),
			"opentelekomcloud_networking_secgroup_rule_v2":        resourceNetworkingSecGroupRuleV2(),
			"opentelekomcloud_s3_bucket":                          resourceS3Bucket(),
			"opentelekomcloud_s3_bucket_policy":                   resourceS3BucketPolicy(),
			"opentelekomcloud_s3_bucket_object":                   resourceS3BucketObject(),
			"opentelekomcloud_obs_bucket":                         resourceObsBucket(),
			"opentelekomcloud_obs_bucket_object":                  resourceObsBucketObject(),
			"opentelekomcloud_elb_loadbalancer":                   resourceELoadBalancer(),
			"opentelekomcloud_elb_listener":                       resourceEListener(),
			"opentelekomcloud_elb_backend":                        resourceBackend(),
			"opentelekomcloud_elb_health":                         resourceHealth(),
			"opentelekomcloud_ces_alarmrule":                      resourceAlarmRule(),
			"opentelekomcloud_smn_topic_v2":                       resourceTopic(),
			"opentelekomcloud_smn_subscription_v2":                resourceSubscription(),
			"opentelekomcloud_rds_instance_v1":                    resourceRdsInstance(),
			"opentelekomcloud_vpc_eip_v1":                         resourceVpcEIPV1(),
			"opentelekomcloud_vpc_v1":                             resourceVirtualPrivateCloudV1(),
			"opentelekomcloud_vpc_peering_connection_v2":          resourceVpcPeeringConnectionV2(),
			"opentelekomcloud_vpc_peering_connection_accepter_v2": resourceVpcPeeringConnectionAccepterV2(),
			"opentelekomcloud_vpc_route_v2":                       resourceVPCRouteV2(),
			"opentelekomcloud_vpc_subnet_v1":                      resourceVpcSubnetV1(),
			"opentelekomcloud_vpc_flow_log_v1":                    resourceVpcFlowLogV1(),
			"opentelekomcloud_rts_software_deployment_v1":         resourceRtsSoftwareDeploymentV1(),
			"opentelekomcloud_rts_software_config_v1":             resourceSoftwareConfigV1(),
			"opentelekomcloud_rts_stack_v1":                       resourceRTSStackV1(),
			"opentelekomcloud_sfs_file_system_v2":                 resourceSFSFileSystemV2(),
			"opentelekomcloud_compute_bms_tags_v2":                resourceBMSTagsV2(),
			"opentelekomcloud_compute_bms_server_v2":              resourceComputeBMSInstanceV2(),
			"opentelekomcloud_as_configuration_v1":                resourceASConfiguration(),
			"opentelekomcloud_as_group_v1":                        resourceASGroup(),
			"opentelekomcloud_as_policy_v1":                       resourceASPolicy(),
			"opentelekomcloud_csbs_backup_v1":                     resourceCSBSBackupV1(),
			"opentelekomcloud_csbs_backup_policy_v1":              resourceCSBSBackupPolicyV1(),
			"opentelekomcloud_deh_host_v1":                        resourceDeHHostV1(),
			"opentelekomcloud_networking_vip_v2":                  resourceNetworkingVIPV2(),
			"opentelekomcloud_networking_vip_associate_v2":        resourceNetworkingVIPAssociateV2(),
			"opentelekomcloud_dms_instance_v1":                    resourceDmsInstancesV1(),
			"opentelekomcloud_dms_queue_v1":                       resourceDmsQueuesV1(),
			"opentelekomcloud_dms_group_v1":                       resourceDmsGroupsV1(),
			"opentelekomcloud_vbs_backup_policy_v2":               resourceVBSBackupPolicyV2(),
			"opentelekomcloud_vbs_backup_v2":                      resourceVBSBackupV2(),
			"opentelekomcloud_vbs_backup_share_v2":                resourceVBSBackupShareV2(),
			"opentelekomcloud_antiddos_v1":                        resourceAntiDdosV1(),
			"opentelekomcloud_cts_tracker_v1":                     resourceCTSTrackerV1(),
			"opentelekomcloud_cce_node_v3":                        resourceCCENodeV3(),
			"opentelekomcloud_cce_cluster_v3":                     resourceCCEClusterV3(),
			"opentelekomcloud_maas_task_v1":                       resourceMaasTaskV1(),
			"opentelekomcloud_css_cluster_v1":                     resourceCssClusterV1(),
			"opentelekomcloud_rds_instance_v3":                    resourceRdsInstanceV3(),
			"opentelekomcloud_rds_parametergroup_v3":              resourceRdsConfigurationV3(),
			"opentelekomcloud_sdrs_protectiongroup_v1":            resourceSdrsProtectiongroupV1(),
			"opentelekomcloud_waf_certificate_v1":                 resourceWafCertificateV1(),
			"opentelekomcloud_waf_domain_v1":                      resourceWafDomainV1(),
			"opentelekomcloud_waf_policy_v1":                      resourceWafPolicyV1(),
			"opentelekomcloud_waf_whiteblackip_rule_v1":           resourceWafWhiteBlackIpRuleV1(),
			"opentelekomcloud_waf_datamasking_rule_v1":            resourceWafDataMaskingRuleV1(),
			"opentelekomcloud_waf_falsealarmmasking_rule_v1":      resourceWafFalseAlarmMaskingRuleV1(),
			"opentelekomcloud_waf_ccattackprotection_rule_v1":     resourceWafCcAttackProtectionRuleV1(),
			"opentelekomcloud_waf_preciseprotection_rule_v1":      resourceWafPreciseProtectionRuleV1(),
			"opentelekomcloud_waf_webtamperprotection_rule_v1":    resourceWafWebTamperProtectionRuleV1(),
			"opentelekomcloud_vpnaas_ipsec_policy_v2":             resourceVpnIPSecPolicyV2(),
			"opentelekomcloud_vpnaas_service_v2":                  resourceVpnServiceV2(),
			"opentelekomcloud_vpnaas_ike_policy_v2":               resourceVpnIKEPolicyV2(),
			"opentelekomcloud_vpnaas_endpoint_group_v2":           resourceVpnEndpointGroupV2(),
			"opentelekomcloud_vpnaas_site_connection_v2":          resourceVpnSiteConnectionV2(),
			"opentelekomcloud_logtank_group_v2":                   resourceLTSGroupV2(),
			"opentelekomcloud_logtank_topic_v2":                   resourceLTSTopicV2(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return configureProvider(d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "The access key for API operations. You can retrieve this\n" +
			"from the 'My Credential' section of the console.",

		"secret_key": "The secret key for API operations. You can retrieve this\n" +
			"from the 'My Credential' section of the console.",

		"auth_url": "The Identity authentication URL.",

		"region": "The OpenTelekomCloud region to connect to.",

		"user_name": "Username to login with.",

		"user_id": "User ID to login with.",

		"tenant_id": "The ID of the Tenant (Identity v2) or Project (Identity v3)\n" +
			"to login with.",

		"tenant_name": "The name of the Tenant (Identity v2) or Project (Identity v3)\n" +
			"to login with.",

		"password": "Password to login with.",

		"token": "Authentication token to use as an alternative to username/password.",

		"security_token": "Security token to use for OBS federated authentication.",

		"domain_id": "The ID of the Domain to scope to (Identity v3).",

		"domain_name": "The name of the Domain to scope to (Identity v3).",

		"insecure": "Trust self-signed certificates.",

		"cacert_file": "A Custom CA certificate.",

		"endpoint_type": "The catalog endpoint type to use.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",

		"swauth": "Use Swift's authentication system instead of Keystone. Only used for\n" +
			"interaction with Swift.",

		"agency_name": "The name of agency",

		"agency_domain_name": "The name of domain who created the agency (Identity v3).",

		"delegated_project": "The name of delegated project (Identity v3).",

		"cloud": "An entry in a `clouds.yaml` file to use.",

		"max_retries": "How many times HTTP connection should be retried until giving up.",
	}
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		AccessKey:        d.Get("access_key").(string),
		SecretKey:        d.Get("secret_key").(string),
		CACertFile:       d.Get("cacert_file").(string),
		ClientCertFile:   d.Get("cert").(string),
		ClientKeyFile:    d.Get("key").(string),
		Cloud:            d.Get("cloud").(string),
		DomainID:         d.Get("domain_id").(string),
		DomainName:       d.Get("domain_name").(string),
		EndpointType:     d.Get("endpoint_type").(string),
		IdentityEndpoint: d.Get("auth_url").(string),
		Insecure:         d.Get("insecure").(bool),
		Password:         d.Get("password").(string),
		Region:           d.Get("region").(string),
		Swauth:           d.Get("swauth").(bool),
		Token:            d.Get("token").(string),
		SecurityToken:    d.Get("security_token").(string),
		TenantID:         d.Get("tenant_id").(string),
		TenantName:       d.Get("tenant_name").(string),
		Username:         d.Get("user_name").(string),
		UserID:           d.Get("user_id").(string),
		AgencyName:       d.Get("agency_name").(string),
		AgencyDomainName: d.Get("agency_domain_name").(string),
		DelegatedProject: d.Get("delegated_project").(string),
		MaxRetries:       d.Get("max_retries").(int),
		terraformVersion: terraformVersion,
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func reconfigProjectName(src *Config, projectName string) (*Config, error) {
	config := &Config{}
	if err := copier.Copy(config, src); err != nil {
		return nil, err
	}
	config.TenantName = projectName
	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}
	return config, nil
}
