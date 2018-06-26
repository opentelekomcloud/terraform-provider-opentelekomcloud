package opentelekomcloud

import (
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

// Provider returns a schema.Provider for OpenTelekomCloud.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
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

			"auth_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", nil),
				Description: descriptions["auth_url"],
			},

			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["region"],
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", ""),
			},

			"user_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: descriptions["user_name"],
			},

			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: descriptions["user_name"],
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},

			"tenant_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
				Description: descriptions["tenant_name"],
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: descriptions["password"],
			},

			"token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_TOKEN", ""),
				Description: descriptions["token"],
			},

			"domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_ID",
					"OS_PROJECT_DOMAIN_ID",
					"OS_DOMAIN_ID",
				}, ""),
				Description: descriptions["domain_id"],
			},

			"domain_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_NAME",
					"OS_PROJECT_DOMAIN_NAME",
					"OS_DOMAIN_NAME",
					"OS_DEFAULT_DOMAIN",
				}, ""),
				Description: descriptions["domain_name"],
			},

			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", ""),
				Description: descriptions["insecure"],
			},

			"endpoint_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ENDPOINT_TYPE", ""),
			},

			"cacert_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"swauth": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SWAUTH", ""),
				Description: descriptions["swauth"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"opentelekomcloud_images_image_v2":        dataSourceImagesImageV2(),
			"opentelekomcloud_networking_network_v2":  dataSourceNetworkingNetworkV2(),
			"opentelekomcloud_networking_secgroup_v2": dataSourceNetworkingSecGroupV2(),
			"opentelekomcloud_s3_bucket_object":       dataSourceS3BucketObject(),
			"opentelekomcloud_kms_key_v1":             dataSourceKmsKeyV1(),
			"opentelekomcloud_kms_data_key_v1":        dataSourceKmsDataKeyV1(),
			"opentelekomcloud_rds_flavors_v1":         dataSourceRdsFlavorV1(),
			"opentelekomcloud_vpc_v1":                 dataSourceVirtualPrivateCloudVpcV1(),
			"opentelekomcloud_vpc_subnet_v1":          dataSourceVpcSubnetV1(),
			"opentelekomcloud_vpc_subnet_ids_v1":      dataSourceVpcSubnetIdsV1(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"opentelekomcloud_blockstorage_volume_v2":          resourceBlockStorageVolumeV2(),
			"opentelekomcloud_compute_instance_v2":             resourceComputeInstanceV2(),
			"opentelekomcloud_compute_keypair_v2":              resourceComputeKeypairV2(),
			"opentelekomcloud_compute_secgroup_v2":             resourceComputeSecGroupV2(),
			"opentelekomcloud_compute_servergroup_v2":          resourceComputeServerGroupV2(),
			"opentelekomcloud_compute_floatingip_v2":           resourceComputeFloatingIPV2(),
			"opentelekomcloud_compute_floatingip_associate_v2": resourceComputeFloatingIPAssociateV2(),
			"opentelekomcloud_compute_volume_attach_v2":        resourceComputeVolumeAttachV2(),
			"opentelekomcloud_dns_recordset_v2":                resourceDNSRecordSetV2(),
			"opentelekomcloud_dns_zone_v2":                     resourceDNSZoneV2(),
			"opentelekomcloud_fw_firewall_group_v2":            resourceFWFirewallGroupV2(),
			"opentelekomcloud_fw_policy_v2":                    resourceFWPolicyV2(),
			"opentelekomcloud_fw_rule_v2":                      resourceFWRuleV2(),
			"opentelekomcloud_images_image_v2":                 resourceImagesImageV2(),
			"opentelekomcloud_kms_key_v1":                      resourceKmsKeyV1(),
			"opentelekomcloud_lb_loadbalancer_v2":              resourceLoadBalancerV2(),
			"opentelekomcloud_lb_listener_v2":                  resourceListenerV2(),
			"opentelekomcloud_lb_pool_v2":                      resourcePoolV2(),
			"opentelekomcloud_lb_member_v2":                    resourceMemberV2(),
			"opentelekomcloud_lb_monitor_v2":                   resourceMonitorV2(),
			"opentelekomcloud_networking_network_v2":           resourceNetworkingNetworkV2(),
			"opentelekomcloud_networking_subnet_v2":            resourceNetworkingSubnetV2(),
			"opentelekomcloud_networking_floatingip_v2":        resourceNetworkingFloatingIPV2(),
			"opentelekomcloud_networking_port_v2":              resourceNetworkingPortV2(),
			"opentelekomcloud_networking_router_v2":            resourceNetworkingRouterV2(),
			"opentelekomcloud_networking_router_interface_v2":  resourceNetworkingRouterInterfaceV2(),
			"opentelekomcloud_networking_router_route_v2":      resourceNetworkingRouterRouteV2(),
			"opentelekomcloud_networking_secgroup_v2":          resourceNetworkingSecGroupV2(),
			"opentelekomcloud_networking_secgroup_rule_v2":     resourceNetworkingSecGroupRuleV2(),
			"opentelekomcloud_s3_bucket":                       resourceS3Bucket(),
			"opentelekomcloud_s3_bucket_policy":                resourceS3BucketPolicy(),
			"opentelekomcloud_s3_bucket_object":                resourceS3BucketObject(),
			"opentelekomcloud_elb_loadbalancer":                resourceELoadBalancer(),
			"opentelekomcloud_elb_listener":                    resourceEListener(),
			"opentelekomcloud_elb_backend":                     resourceBackend(),
			"opentelekomcloud_elb_health":                      resourceHealth(),
			"opentelekomcloud_ces_alarmrule":                   resourceAlarmRule(),
			"opentelekomcloud_smn_topic_v2":                    resourceTopic(),
			"opentelekomcloud_smn_subscription_v2":             resourceSubscription(),
			"opentelekomcloud_rds_instance_v1":                 resourceRdsInstance(),
			"opentelekomcloud_vpc_eip_v1":                      resourceVpcEIPV1(),
			"opentelekomcloud_vpc_v1":                          resourceVirtualPrivateCloudV1(),
			"opentelekomcloud_vpc_subnet_v1":                   resourceVpcSubnetV1(),
		},

		ConfigureFunc: configureProvider,
	}
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

		"domain_id": "The ID of the Domain to scope to (Identity v3).",

		"domain_name": "The name of the Domain to scope to (Identity v3).",

		"insecure": "Trust self-signed certificates.",

		"cacert_file": "A Custom CA certificate.",

		"endpoint_type": "The catalog endpoint type to use.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",

		"swauth": "Use Swift's authentication system instead of Keystone. Only used for\n" +
			"interaction with Swift.",
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKey:        d.Get("access_key").(string),
		SecretKey:        d.Get("secret_key").(string),
		CACertFile:       d.Get("cacert_file").(string),
		ClientCertFile:   d.Get("cert").(string),
		ClientKeyFile:    d.Get("key").(string),
		DomainID:         d.Get("domain_id").(string),
		DomainName:       d.Get("domain_name").(string),
		EndpointType:     d.Get("endpoint_type").(string),
		IdentityEndpoint: d.Get("auth_url").(string),
		Insecure:         d.Get("insecure").(bool),
		Password:         d.Get("password").(string),
		Region:           d.Get("region").(string),
		Swauth:           d.Get("swauth").(bool),
		Token:            d.Get("token").(string),
		TenantID:         d.Get("tenant_id").(string),
		TenantName:       d.Get("tenant_name").(string),
		Username:         d.Get("user_name").(string),
		UserID:           d.Get("user_id").(string),
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
