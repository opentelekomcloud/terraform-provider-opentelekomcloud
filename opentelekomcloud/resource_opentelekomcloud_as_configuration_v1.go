package opentelekomcloud

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/configurations"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/groups"
)

func resourceASConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceASConfigurationCreate,
		Read:   resourceASConfigurationRead,
		Update: nil,
		Delete: resourceASConfigurationDelete,

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
				ValidateFunc: validateName,
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
						},
						"flavor": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  getDefaultFlavor(),
						},
						"image": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							// just stash the hash for state & diff comparisons
							StateFunc: func(v interface{}) string {
								switch v.(type) {
								case string:
									hash := sha1.Sum([]byte(v.(string)))
									return hex.EncodeToString(hash[:])
								default:
									return ""
								}
							},
						},
						"disk": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 32768),
									},
									"volume_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"SATA", "SAS", "SSD", "co-p1", "uh-l1",
										}, false),
									},
									"disk_type": {
										Type:     schema.TypeString,
										Required: true,
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"public_ip": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eip": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_type": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														"5_bgp", "5_mailbgp",
													}, false),
												},
												"bandwidth": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"size": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: validation.IntBetween(1, 500),
															},
															"share_type": {
																Type:     schema.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	"PER",
																}, false),
															},
															"charging_mode": {
																Type:     schema.TypeString,
																Required: true,
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
						},
					},
				},
			},
		},
	}
}

func getDefaultFlavor() string {
	flavorID := os.Getenv("OS_FLAVOR_ID")
	if flavorID != "" {
		return flavorID
	}
	return ""
}

func getDisk(diskMeta []interface{}) (*[]configurations.DiskOpts, error) {
	var diskOptsList []configurations.DiskOpts

	for _, v := range diskMeta {
		disk := v.(map[string]interface{})
		size := disk["size"].(int)
		volumeType := disk["volume_type"].(string)
		diskType := disk["disk_type"].(string)
		if diskType == "SYS" {
			if size < 1 || size > 1024 {
				return nil, fmt.Errorf("for system disk size should be [1, 1024]")
			}
		}
		if diskType == "DATA" {
			if size < 10 || size > 32768 {
				return nil, fmt.Errorf("for data disk size should be [10, 32768]")
			}
		}
		diskOpts := configurations.DiskOpts{
			Size:       size,
			VolumeType: volumeType,
			DiskType:   diskType,
		}
		kmsID := disk["kms_id"].(string)
		if kmsID != "" {
			meta := make(map[string]string)
			meta["__system__cmkid"] = kmsID
			meta["__system__encrypted"] = "1"
			diskOpts.Metadata = meta
		}
		diskOptsList = append(diskOptsList, diskOpts)
	}

	return &diskOptsList, nil
}

func getPersonality(personalityMeta []interface{}) []configurations.PersonalityOpts {
	var personalityOptsList []configurations.PersonalityOpts

	for _, v := range personalityMeta {
		personality := v.(map[string]interface{})
		personalityOpts := configurations.PersonalityOpts{
			Path:    personality["path"].(string),
			Content: personality["content"].(string),
		}
		personalityOptsList = append(personalityOptsList, personalityOpts)
	}

	return personalityOptsList
}

func getPublicIps(publicIpMeta map[string]interface{}) configurations.PublicIpOpts {
	eipMap := publicIpMeta["eip"].([]interface{})[0].(map[string]interface{})
	bandWidthMap := eipMap["bandwidth"].([]interface{})[0].(map[string]interface{})
	bandWidthOpts := configurations.BandwidthOpts{
		Size:         bandWidthMap["size"].(int),
		ShareType:    bandWidthMap["share_type"].(string),
		ChargingMode: bandWidthMap["charging_mode"].(string),
	}

	eipOpts := configurations.EipOpts{
		IpType:    eipMap["ip_type"].(string),
		Bandwidth: bandWidthOpts,
	}

	publicIpOpts := configurations.PublicIpOpts{
		Eip: eipOpts,
	}

	return publicIpOpts
}

func getInstanceConfig(configDataMap map[string]interface{}) (*configurations.InstanceConfigOpts, error) {
	disksData := configDataMap["disk"].([]interface{})
	disks, err := getDisk(disksData)
	if err != nil {
		return nil, fmt.Errorf("error happened when validating disk size: %s", err)
	}

	personalityData := configDataMap["personality"].([]interface{})
	personalities := getPersonality(personalityData)

	instanceConfigOpts := &configurations.InstanceConfigOpts{
		ID:          configDataMap["instance_id"].(string),
		FlavorRef:   configDataMap["flavor"].(string),
		ImageRef:    configDataMap["image"].(string),
		SSHKey:      configDataMap["key_name"].(string),
		UserData:    []byte(configDataMap["user_data"].(string)),
		Disk:        *disks,
		Personality: personalities,
		Metadata:    configDataMap["metadata"].(map[string]interface{}),
	}

	pubicIpData := configDataMap["public_ip"].([]interface{})
	// user specify public_ip
	if len(pubicIpData) == 1 {
		publicIpMap := pubicIpData[0].(map[string]interface{})
		publicIps := getPublicIps(publicIpMap)
		instanceConfigOpts.PubicIp = publicIps
	}
	log.Printf("[DEBUG] get instanceConfig: %#v", instanceConfigOpts)

	return instanceConfigOpts, nil
}

func resourceASConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud AutoScaling client: %s", err)
	}

	configDataMap := d.Get("instance_config").([]interface{})[0].(map[string]interface{})
	log.Printf("[DEBUG] instance_config is: %#v", configDataMap)
	instanceConfig, err := getInstanceConfig(configDataMap)
	if err != nil {
		return fmt.Errorf("error when getting instance_config info: %s", err)
	}
	createOpts := configurations.CreateOpts{
		Name:           d.Get("scaling_configuration_name").(string),
		InstanceConfig: *instanceConfig,
	}

	log.Printf("[DEBUG] Create AS configuration Options: %#v", createOpts)
	asConfigID, err := configurations.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating ASConfiguration: %s", err)
	}
	d.SetId(asConfigID)
	log.Printf("[DEBUG] Create AS Configuration %q Success!", asConfigID)
	return resourceASConfigurationRead(d, meta)
}

func resourceASConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud AutoScaling client: %s", err)
	}

	asConfig, err := configurations.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "AS Configuration")
	}

	log.Printf("[DEBUG] Retrieved ASConfiguration %q: %+v", d.Id(), asConfig)

	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
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
	instanceConfigInfo["user_data"] = asConfig.InstanceConfig.UserData
	// instanceConfigInfo["disk"] = asConfig.InstanceConfig.Disk               // TODO: Check
	// instanceConfigInfo["personality"] = asConfig.InstanceConfig.Personality // TODO: Check
	instanceConfigInfo["public_ip"] = asConfig.InstanceConfig.PublicIp
	instanceConfigInfo["metadata"] = asConfig.InstanceConfig.Metadata
	instanceConfigList := []interface{}{instanceConfigInfo}

	if err = d.Set("instance_config", instanceConfigList); err != nil {
		return err
	}

	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}

func resourceASConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud AutoScaling client: %s", err)
	}

	asConfigGroups, err := getASGroupsByConfiguration(client, d.Id())
	if err != nil {
		return fmt.Errorf("error getting AS groups by configuration ID %q: %s", d.Id(), err)
	}
	if len(asConfigGroups) > 0 {
		var groupIDs []string
		for _, group := range asConfigGroups {
			groupIDs = append(groupIDs, group.ID)
		}
		return fmt.Errorf("can not delete the configuration %q, it is used by AS groups %s", d.Id(), groupIDs)
	}

	log.Printf("[DEBUG] Begin to delete AS configuration %q", d.Id())
	if err := configurations.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("error deleting AS configuration: %s", err)
	}

	return nil
}

func getASGroupsByConfiguration(client *golangsdk.ServiceClient, configID string) ([]groups.Group, error) {
	var asGroups []groups.Group
	listOpts := groups.ListOpts{
		ConfigurationID: configID,
	}
	page, err := groups.List(client, listOpts).AllPages()
	if err != nil {
		return asGroups, fmt.Errorf("error getting ASGroups by configuration %q: %s", configID, err)
	}
	asGroups, err = page.(groups.GroupPage).Extract()
	return asGroups, err
}
