package as

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/configurations"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceASConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASConfigurationCreate,
		ReadContext:   resourceASConfigurationRead,
		DeleteContext: resourceASConfigurationDelete,

		CustomizeDiff: validateDiskSize,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_configuration_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateName,
			},
			"instance_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"flavor": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							DefaultFunc: schema.EnvDefaultFunc("OS_FLAVOR_ID", nil),
						},
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							DefaultFunc: schema.EnvDefaultFunc("OS_IMAGE_ID", nil),
						},
						"key_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"user_data": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: common.SuppressEmptyStringSHA,
							// just stash the hash for state & diff comparisons
							StateFunc: func(v interface{}) string {
								switch v := v.(type) {
								case string:
									return common.InstallScriptHashSum(v)
								default:
									return ""
								}
							},
						},
						"disk": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											"SATA", "SAS", "SSD", "co-p1", "uh-l1",
										}, false),
									},
									"disk_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											"DATA", "SYS",
										}, false),
									},
									"kms_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"personality": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 5,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"content": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"public_ip": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eip": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_type": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
													ValidateFunc: validation.StringInSlice([]string{
														"5_bgp", "5_mailbgp",
													}, false),
												},
												"bandwidth": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Required: true,
													ForceNew: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"size": {
																Type:         schema.TypeInt,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.IntBetween(1, 500),
															},
															"share_type": {
																Type:     schema.TypeString,
																Required: true,
																ForceNew: true,
																ValidateFunc: validation.StringInSlice([]string{
																	"PER",
																}, false),
															},
															"charging_mode": {
																Type:     schema.TypeString,
																Required: true,
																ForceNew: true,
																ValidateFunc: validation.StringInSlice([]string{
																	"traffic",
																}, false),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
						},
						"security_groups": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
			},
		},
	}
}

func getDisk(diskMeta []interface{}) []configurations.Disk {
	var diskOptsList []configurations.Disk

	for _, v := range diskMeta {
		disk := v.(map[string]interface{})
		size := disk["size"].(int)
		volumeType := disk["volume_type"].(string)
		diskType := disk["disk_type"].(string)
		diskOpts := configurations.Disk{
			Size:       size,
			VolumeType: volumeType,
			DiskType:   diskType,
		}
		kmsID := disk["kms_id"].(string)
		if kmsID != "" {
			meta := make(map[string]interface{})
			meta["__system__cmkid"] = kmsID
			meta["__system__encrypted"] = "1"
			diskOpts.Metadata = meta
		}
		diskOptsList = append(diskOptsList, diskOpts)
	}

	return diskOptsList
}

func getPersonality(personalityMeta []interface{}) []configurations.Personality {
	var personalityOptsList []configurations.Personality

	for _, v := range personalityMeta {
		personality := v.(map[string]interface{})
		personalityOpts := configurations.Personality{
			Path:    personality["path"].(string),
			Content: personality["content"].(string),
		}
		personalityOptsList = append(personalityOptsList, personalityOpts)
	}

	return personalityOptsList
}

func getPublicIps(publicIpMeta map[string]interface{}) *configurations.PublicIp {
	eipMap := publicIpMeta["eip"].([]interface{})[0].(map[string]interface{})
	bandWidthMap := eipMap["bandwidth"].([]interface{})[0].(map[string]interface{})

	publicIpOpts := &configurations.PublicIp{
		Eip: configurations.Eip{
			Type: eipMap["ip_type"].(string),
			Bandwidth: configurations.Bandwidth{
				Size:         bandWidthMap["size"].(int),
				ShareType:    bandWidthMap["share_type"].(string),
				ChargingMode: bandWidthMap["charging_mode"].(string),
			},
		},
	}

	return publicIpOpts
}

func getSecurityGroups(d *schema.ResourceData) []configurations.SecurityGroup {
	rawSecGroups := d.Get("instance_config.0.security_groups").(*schema.Set).List()
	secGroups := make([]configurations.SecurityGroup, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secGroups[i] = configurations.SecurityGroup{
			ID: raw.(string),
		}
	}
	return secGroups
}

func getInstanceConfig(d *schema.ResourceData) configurations.InstanceConfigOpts {
	configDataMap := d.Get("instance_config").([]interface{})[0].(map[string]interface{})
	disksData := configDataMap["disk"].([]interface{})
	personalityData := configDataMap["personality"].([]interface{})
	meta := configDataMap["metadata"].(map[string]interface{})
	instanceConfigOpts := configurations.InstanceConfigOpts{
		ID:             configDataMap["instance_id"].(string),
		FlavorRef:      configDataMap["flavor"].(string),
		ImageRef:       configDataMap["image"].(string),
		Disk:           getDisk(disksData),
		SSHKey:         configDataMap["key_name"].(string),
		Personality:    getPersonality(personalityData),
		UserData:       []byte(configDataMap["user_data"].(string)),
		SecurityGroups: getSecurityGroups(d),
	}
	if _, ok := meta["admin_pass"]; ok {
		instanceConfigOpts.Metadata = configurations.AdminPassMetadata{
			AdminPass: meta["admin_pass"].(string),
		}
	}

	pubicIpData := configDataMap["public_ip"].([]interface{})
	// user specify public_ip
	if len(pubicIpData) == 1 {
		publicIpMap := pubicIpData[0].(map[string]interface{})
		instanceConfigOpts.PubicIp = getPublicIps(publicIpMap)
	}

	return instanceConfigOpts
}

func resourceASConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := configurations.CreateOpts{
		Name:           d.Get("scaling_configuration_name").(string),
		InstanceConfig: getInstanceConfig(d),
	}

	log.Printf("[DEBUG] Create AS configuration Options: %#v", createOpts)
	asConfigID, err := configurations.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating ASConfiguration: %s", err)
	}

	d.SetId(asConfigID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceASConfigurationRead(clientCtx, d, meta)
}

func resourceASConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}
	asConfig, err := configurations.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AS Configuration")
	}

	log.Printf("[DEBUG] Retrieved ASConfiguration %q: %+v", d.Id(), asConfig)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("scaling_configuration_name", asConfig.Name),
	)

	instanceConfig := d.Get("instance_config").([]interface{})
	instanceConfigInfo := make(map[string]interface{})
	if len(instanceConfig) != 0 {
		instanceConfigInfo = instanceConfig[0].(map[string]interface{})
	}
	instanceConfigInfo["instance_id"] = asConfig.InstanceConfig.InstanceID
	instanceConfigInfo["flavor"] = asConfig.InstanceConfig.FlavorRef
	instanceConfigInfo["image"] = asConfig.InstanceConfig.ImageRef
	instanceConfigInfo["key_name"] = asConfig.InstanceConfig.SSHKey
	instanceConfigInfo["user_data"] = common.InstallScriptHashSum(asConfig.InstanceConfig.UserData)

	var secGrpIDs []string
	for _, sg := range asConfig.InstanceConfig.SecurityGroups {
		secGrpIDs = append(secGrpIDs, sg.ID)
	}
	instanceConfigInfo["security_groups"] = secGrpIDs
	instanceConfigList := []interface{}{instanceConfigInfo}

	if err := d.Set("instance_config", instanceConfigList); err != nil {
		return diag.FromErr(err)
	}

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceASConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	asConfigGroups, err := getASGroupsByConfiguration(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error getting AS groups by configuration ID %q: %s", d.Id(), err)
	}

	if len(asConfigGroups) > 0 {
		var groupIDs []string
		for _, group := range asConfigGroups {
			groupIDs = append(groupIDs, group.ID)
		}
		return fmterr.Errorf("can not delete the configuration %q, it is used by AS groups %s", d.Id(), groupIDs)
	}

	log.Printf("[DEBUG] Begin to delete AS configuration %q", d.Id())
	if err := configurations.Delete(client, d.Id()); err != nil {
		return fmterr.Errorf("error deleting AS configuration: %s", err)
	}

	return nil
}

func getASGroupsByConfiguration(client *golangsdk.ServiceClient, configID string) ([]groups.Group, error) {
	var asGroups *groups.ListScalingGroupsResponse
	listOpts := groups.ListOpts{
		ConfigurationID: configID,
	}
	asGroups, err := groups.List(client, listOpts)
	if err != nil {
		return asGroups.ScalingGroups, fmt.Errorf("error getting ASGroups by configuration %q: %s", configID, err)
	}
	return asGroups.ScalingGroups, err
}

func validateDiskSize(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	instanceConfigData := d.Get("instance_config").([]interface{})[0].(map[string]interface{})
	disksData := instanceConfigData["disk"].([]interface{})
	mErr := &multierror.Error{}
	for _, v := range disksData {
		disk := v.(map[string]interface{})
		size := disk["size"].(int)
		diskType := disk["disk_type"].(string)
		if diskType == "SYS" {
			if size < 4 || size > 32768 {
				mErr = multierror.Append(mErr, fmt.Errorf("for system disk size should be [4, 32768]"))
			}
		}
		if diskType == "DATA" {
			if size < 10 || size > 32768 {
				mErr = multierror.Append(mErr, fmt.Errorf("for data disk size should be [10, 32768]"))
			}
		}
	}
	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}
