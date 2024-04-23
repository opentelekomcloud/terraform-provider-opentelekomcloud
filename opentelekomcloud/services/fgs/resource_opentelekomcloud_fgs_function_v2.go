package fgs

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	aliases "github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/alias"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/function"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/reserved"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func ResourceFgsFunctionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFgsFunctionV2Create,
		ReadContext:   resourceFgsFunctionV2Read,
		UpdateContext: resourceFgsFunctionV2Update,
		DeleteContext: resourceFgsFunctionV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"runtime": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"code_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"handler": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `schema: Required; The entry point of the function.`,
			},
			"functiongraph_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v1", "v2",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code_filename": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypted_user_data": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"agency": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_agency": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"func_code": {
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: hashcode.DecodeHashAndHexEncode,
			},
			"depend_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"initializer_handler": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"initializer_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"network_id"},
			},
			"network_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"vpc_id"},
			},
			"mount_user_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"mount_user_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				RequiredWith: []string{
					"log_topic_id", "log_group_name", "log_topic_name"},
			},
			"log_topic_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"log_group_id"},
			},
			"log_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"log_group_id"},
			},
			"log_topic_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"log_group_id"},
			},
			"func_mounts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mount_resource": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mount_share_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"local_mount_path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"custom_image": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"max_instance_num": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\-?\d+$`),
					`invalid value of maximum instance number, want an integer number or integer string.`),
			},
			"versions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The version name.",
						},
						"aliases": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the version alias.",
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The description of the version alias.",
									},
								},
							},
							Description: "The aliases management for specified version.",
						},
					},
				},
				Description: "The versions management of the function.",
			},
			"tags": common.TagsSchema(),
			"reserved_instances": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"qualifier_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"qualifier_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"idle_mode": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"tactics_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     tacticsConfigsSchema(),
						},
					},
				},
			},
			"concurrency_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"gpu_memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"gpu_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_list": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func tacticsConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cron_configs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cron": {
							Type:     schema.TypeString,
							Required: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"expired_time": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func buildCustomImage(imageConfig []interface{}) *function.CustomImage {
	if len(imageConfig) < 1 {
		return nil
	}

	cfg := imageConfig[0].(map[string]interface{})
	return &function.CustomImage{
		Enabled: pointerto.Bool(true),
		Image:   cfg["url"].(string),
	}
}

func buildFgsFunctionParameters(d *schema.ResourceData) (function.CreateOpts, error) {
	// check app
	app, appOk := d.GetOk("app")
	if !appOk {
		return function.CreateOpts{}, fmt.Errorf("app must be configured")
	}
	packV := ""
	if appOk {
		packV = app.(string)
	}

	agencyV := ""
	if v, ok := d.GetOk("agency"); ok {
		agencyV = v.(string)
	}

	result := function.CreateOpts{
		Name:              d.Get("name").(string),
		Type:              d.Get("functiongraph_version").(string),
		Package:           packV,
		CodeType:          d.Get("code_type").(string),
		CodeURL:           d.Get("code_url").(string),
		Description:       d.Get("description").(string),
		CodeFilename:      d.Get("code_filename").(string),
		Handler:           d.Get("handler").(string),
		MemorySize:        d.Get("memory_size").(int),
		Runtime:           d.Get("runtime").(string),
		Timeout:           d.Get("timeout").(int),
		UserData:          d.Get("user_data").(string),
		EncryptedUserData: d.Get("encrypted_user_data").(string),
		Xrole:             agencyV,
		CustomImage:       buildCustomImage(d.Get("custom_image").([]interface{})),
		GpuMemory:         pointerto.Int(d.Get("gpu_memory").(int)),
	}
	if v, ok := d.GetOk("func_code"); ok {
		funcCode := function.FuncCode{
			File: hashcode.TryBase64EncodeString(v.(string)),
		}
		result.FuncCode = &funcCode
	}
	if v, ok := d.GetOk("log_group_id"); ok {
		logConfig := function.FuncLogConfig{
			GroupID:    v.(string),
			StreamID:   d.Get("log_topic_id").(string),
			GroupName:  d.Get("log_group_name").(string),
			StreamName: d.Get("log_topic_name").(string),
		}
		result.LogConfig = &logConfig
	}
	return result, nil
}

func resourceFgsFunctionV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := config.FuncGraphV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	createOpts, err := buildFgsFunctionParameters(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	f, err := function.Create(fgsClient, createOpts)
	if err != nil {
		return diag.Errorf("error creating function: %s", err)
	}

	d.SetId(f.FuncURN)
	urn := resourceFgsFunctionUrn(d.Id())
	if d.HasChanges("vpc_id", "func_mounts", "app_agency", "initializer_handler", "initializer_timeout", "concurrency_num") {
		err := resourceFgsFunctionMetadataUpdate(fgsClient, urn, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("depend_list") {
		err := resourceFgsFunctionCodeUpdate(fgsClient, urn, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if strNum, ok := d.GetOk("max_instance_num"); ok {
		// The integer string of the maximum instance number has been already checked in the schema validation.
		maxInstanceNum, _ := strconv.Atoi(strNum.(string))

		_, err = function.UpdateMaxInstances(fgsClient, function.UpdateFuncInstancesOpts{
			MaxInstanceNum: maxInstanceNum,
			FuncUrn:        urn,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if tagList, ok := d.GetOk("tags"); ok {
		opts := tags.TagsActionOpts{
			Tags:   common.ExpandResourceTags(tagList.(map[string]interface{})),
			Id:     d.Id(),
			Action: "create",
		}
		if err := tags.CreateResourceTag(fgsClient, opts); err != nil {
			return diag.Errorf("failed to add tags to FunctionGraph function (%s): %s", d.Id(), err)
		}
	}

	if err = createFunctionVersions(fgsClient, urn, d.Get("versions").(*schema.Set)); err != nil {
		return diag.Errorf("error creating function versions: %s", err)
	}

	if d.HasChanges("reserved_instances") {
		if err = updateReservedInstanceConfig(fgsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFgsFunctionV2Read(ctx, d, meta)
}

func createFunctionVersions(client *golangsdk.ServiceClient, functionUrn string, versionSet *schema.Set) error {
	for _, v := range versionSet.List() {
		version := v.(map[string]interface{})
		versionNum := version["name"].(string)
		aliasCfg := version["aliases"].([]interface{})
		if len(aliasCfg) < 1 {
			continue
		}
		alias := aliasCfg[0].(map[string]interface{})
		opt := aliases.CreateAliasOpts{
			FuncUrn:     functionUrn,
			Name:        alias["name"].(string),
			Version:     versionNum,
			Description: alias["description"].(string),
		}
		_, err := aliases.CreateAlias(client, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

func setFgsFunctionApp(d *schema.ResourceData, app string) error {
	if _, ok := d.GetOk("app"); ok {
		return d.Set("app", app)
	}
	return nil
}

func setFgsFunctionVpcAccess(d *schema.ResourceData, funcVpc function.FuncVpc) error {
	mErr := multierror.Append(
		d.Set("vpc_id", funcVpc.VpcID),
		d.Set("network_id", funcVpc.SubnetID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting vault fields: %s", err)
	}
	return nil
}

func setFunctionMountConfig(d *schema.ResourceData, mountConfig function.MountConfig) error {
	// set mount_config
	if mountConfig.MountUser != (function.MountUser{}) {
		funcMounts := make([]map[string]string, 0, len(mountConfig.FuncMounts))
		for _, v := range mountConfig.FuncMounts {
			funcMount := map[string]string{
				"mount_type":       v.MountType,
				"mount_resource":   v.MountResource,
				"mount_share_path": v.MountSharePath,
				"local_mount_path": v.LocalMountPath,
			}
			funcMounts = append(funcMounts, funcMount)
		}
		mErr := multierror.Append(
			d.Set("func_mounts", funcMounts),
			d.Set("mount_user_id", mountConfig.MountUser.UserID),
			d.Set("mount_user_group_id", mountConfig.MountUser.UserGroupID),
		)
		if err := mErr.ErrorOrNil(); err != nil {
			return fmt.Errorf("error setting vault fields: %s", err)
		}
	}
	return nil
}

func flattenFgsCustomImage(imageConfig function.CustomImage) []map[string]interface{} {
	if (imageConfig != function.CustomImage{}) {
		return []map[string]interface{}{
			{
				"url": imageConfig.Image,
			},
		}
	}
	return nil
}

func queryFunctionVersions(client *golangsdk.ServiceClient, functionUrn string) ([]string, error) {
	queryOpts := aliases.ListVersionOpts{
		FuncUrn: functionUrn,
	}
	versionList, err := aliases.ListVersion(client, queryOpts)
	if err != nil {
		return nil, fmt.Errorf("error querying version list for the specified function URN: %s", err)
	}
	// The length of the function version list is at least 1 (when creating a function, a version named latest is
	// created by default).
	result := make([]string, len(versionList.Functions))
	for i, version := range versionList.Functions {
		result[i] = version.Version
	}
	return result, nil
}

func queryFunctionAliases(client *golangsdk.ServiceClient, functionUrn string) (map[string][]interface{}, error) {
	aliasList, err := aliases.ListAlias(client, functionUrn)
	if err != nil {
		return nil, fmt.Errorf("error querying alias list for the specified function URN: %s", err)
	}

	// Multiple version aliases may exist in the future.
	result := make(map[string][]interface{})
	for _, v := range aliasList {
		result[v.Version] = append(result[v.Version], map[string]interface{}{
			"name":        v.Name,
			"description": v.Description,
		})
	}
	return result, nil
}

func parseFunctionVersions(client *golangsdk.ServiceClient, functionUrn string) ([]map[string]interface{}, error) {
	versionList, err := queryFunctionVersions(client, functionUrn)
	if err != nil {
		return nil, err
	}
	aliasesConfig, err := queryFunctionAliases(client, functionUrn)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(versionList))
	for _, versionNum := range versionList {
		version := map[string]interface{}{
			"name": versionNum, // The version name, also name as the version number.
		}
		if v, ok := aliasesConfig[versionNum]; ok {
			version["aliases"] = v
		}
		result = append(result, version)
	}

	return result, nil
}

func flattenTacticsConfigs(policyConfig reserved.TacticsConfig) []map[string]interface{} {
	if len(policyConfig.CronConfigs) == 0 {
		return nil
	}

	cronConfigRst := make([]map[string]interface{}, len(policyConfig.CronConfigs))
	for i, v := range policyConfig.CronConfigs {
		cronConfigRst[i] = map[string]interface{}{
			"name":         v.Name,
			"cron":         v.Cron,
			"count":        v.Count,
			"start_time":   v.StartTime,
			"expired_time": v.ExpiredTime,
		}
	}

	return []map[string]interface{}{
		{
			"cron_configs": cronConfigRst,
		},
	}
}

func getReservedInstanceConfig(c *golangsdk.ServiceClient, d *schema.ResourceData) ([]map[string]interface{}, error) {
	opts := reserved.ListConfigOpts{
		FuncUrn: d.Id(),
	}
	reservedInstances, err := reserved.ListReservedInstConfigs(c, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting list of the function reserved instance config: %s", err)
	}

	result := make([]map[string]interface{}, len(reservedInstances.ReservedInstances))
	for i, v := range reservedInstances.ReservedInstances {
		result[i] = map[string]interface{}{
			"count":          v.MinCount,
			"idle_mode":      v.IdleMode,
			"qualifier_name": v.QualifierName,
			"qualifier_type": v.QualifierType,
			"tactics_config": flattenTacticsConfigs(v.TacticsConfig),
		}
	}
	return result, nil
}

func getConcurrencyNum(concurrencyNum *int) int {
	return *concurrencyNum
}

func resourceFgsFunctionV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := config.FuncGraphV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	functionUrn := resourceFgsFunctionUrn(d.Id())
	f, err := function.GetMetadata(fgsClient, functionUrn)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "FunctionGraph function")
	}

	versionConfig, err := parseFunctionVersions(fgsClient, functionUrn)
	if err != nil {
		// Not all regions support the version related API calls.
		log.Printf("[ERROR] Unable to parsing the function versions: %s", err)
	}
	log.Printf("[DEBUG] Retrieved Function %s: %+v", functionUrn, f)
	mErr := multierror.Append(
		d.Set("name", f.FuncName),
		d.Set("code_type", f.CodeType),
		d.Set("code_url", f.CodeURL),
		d.Set("description", f.Description),
		d.Set("code_filename", f.CodeFilename),
		d.Set("handler", f.Handler),
		d.Set("memory_size", f.MemorySize),
		d.Set("runtime", f.Runtime),
		d.Set("timeout", f.Timeout),
		d.Set("user_data", f.UserData),
		d.Set("encrypted_user_data", f.EncryptedUserData),
		d.Set("version", f.Version),
		d.Set("urn", functionUrn),
		d.Set("app_agency", f.AppXrole),
		d.Set("depend_list", f.DependVersionList),
		d.Set("initializer_handler", f.InitHandler),
		d.Set("initializer_timeout", f.InitTimeout),
		d.Set("functiongraph_version", f.Type),
		d.Set("custom_image", flattenFgsCustomImage(f.CustomImage)),
		d.Set("max_instance_num", strconv.Itoa(f.StrategyConfig.Concurrency)),
		d.Set("dns_list", f.DomainNames),
		d.Set("log_group_id", f.LogGroupID),
		d.Set("log_topic_id", f.LogStreamID),
		d.Set("agency", f.Xrole),
		setFgsFunctionApp(d, f.Package),
		setFgsFunctionVpcAccess(d, f.FuncVpc),
		setFunctionMountConfig(d, f.MountConfig),
		d.Set("concurrency_num", getConcurrencyNum(pointerto.Int(f.StrategyConfig.ConcurrentNum))),
		d.Set("versions", versionConfig),
		d.Set("gpu_memory", f.GpuMemory),
		d.Set("gpu_type", f.GpuType),
	)

	reservedInstances, err := getReservedInstanceConfig(fgsClient, d)
	if err != nil {
		return diag.Errorf("error retrieving function reserved instance: %s", err)
	}

	mErr = multierror.Append(mErr,
		d.Set("reserved_instances", reservedInstances),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting function fields: %s", err)
	}

	return nil
}

func updateFunctionTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		oRaw, nRaw  = d.GetChange("tags")
		oMap        = oRaw.(map[string]interface{})
		nMap        = nRaw.(map[string]interface{})
		functionUrn = d.Id()
	)

	if len(oMap) > 0 {
		opts := tags.TagsActionOpts{
			Tags:   common.ExpandResourceTags(oMap),
			Id:     functionUrn,
			Action: "delete",
		}
		if err := tags.DeleteResourceTag(client, opts); err != nil {
			return fmt.Errorf("failed to delete tags from FunctionGraph function (%s): %s", functionUrn, err)
		}
	}

	if len(nMap) > 0 {
		opts := tags.TagsActionOpts{
			Tags:   common.ExpandResourceTags(nMap),
			Id:     functionUrn,
			Action: "create",
		}
		if err := tags.CreateResourceTag(client, opts); err != nil {
			return fmt.Errorf("failed to add tags to FunctionGraph function (%s): %s", functionUrn, err)
		}
	}
	return nil
}

func deleteFunctionVersions(client *golangsdk.ServiceClient, functionUrn string, versionSet *schema.Set) error {
	for _, v := range versionSet.List() {
		version := v.(map[string]interface{})
		aliasCfg := version["aliases"].([]interface{})
		if len(aliasCfg) > 0 {
			alias := aliasCfg[0].(map[string]interface{})
			err := aliases.Delete(client, functionUrn, alias["name"].(string))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateFunctionVersions(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		functionUrn = resourceFgsFunctionUrn(d.Id())

		oldSet, newSet = d.GetChange("versions")
		decrease       = oldSet.(*schema.Set).Difference(newSet.(*schema.Set))
		increase       = newSet.(*schema.Set).Difference(oldSet.(*schema.Set))
	)

	err := deleteFunctionVersions(client, functionUrn, decrease)
	if err != nil {
		return fmt.Errorf("error deleting function versions: %s", err)
	}

	err = createFunctionVersions(client, functionUrn, increase)
	if err != nil {
		return fmt.Errorf("error creating function versions: %s", err)
	}

	return nil
}

func buildCronConfigs(cronConfigs []interface{}) []reserved.CronConfig {
	if len(cronConfigs) < 1 {
		return nil
	}

	result := make([]reserved.CronConfig, len(cronConfigs))
	for i, v := range cronConfigs {
		cronConfig := v.(map[string]interface{})
		result[i] = reserved.CronConfig{
			Name:        cronConfig["name"].(string),
			Cron:        cronConfig["cron"].(string),
			Count:       cronConfig["count"].(int),
			StartTime:   cronConfig["start_time"].(int),
			ExpiredTime: cronConfig["expired_time"].(int),
		}
	}
	return result
}

func buildTacticsConfigs(tacticsConfigs []interface{}) *reserved.TacticsConfig {
	if len(tacticsConfigs) < 1 {
		return nil
	}

	tacticsConfig := tacticsConfigs[0].(map[string]interface{})
	result := reserved.TacticsConfig{
		CronConfigs: buildCronConfigs(tacticsConfig["cron_configs"].([]interface{})),
	}
	return &result
}

func getVersionUrn(client *golangsdk.ServiceClient, functionUrn string, qualifierName string) (string, error) {
	queryOpts := aliases.ListVersionOpts{
		FuncUrn: functionUrn,
	}
	versionList, err := aliases.ListVersion(client, queryOpts)
	if err != nil {
		return "", fmt.Errorf("error querying version list for the specified function URN: %s", err)
	}

	for _, val := range versionList.Functions {
		if val.Version == qualifierName {
			return val.FuncURN, nil
		}
	}

	return "", nil
}

func getReservedInstanceUrn(client *golangsdk.ServiceClient, functionUrn string, policy map[string]interface{}) (string, error) {
	qualifierName := policy["qualifier_name"].(string)
	if policy["qualifier_type"].(string) == "version" {
		urn, err := getVersionUrn(client, functionUrn, qualifierName)
		if err != nil {
			return "", err
		}
		return urn, nil
	}

	aliasList, err := aliases.ListAlias(client, functionUrn)
	if err != nil {
		return "", fmt.Errorf("error querying alias list for the specified function URN: %s", err)
	}
	for _, val := range aliasList {
		if val.Name == qualifierName {
			return val.AliasUrn, nil
		}
	}

	return "", nil
}

func removeReservedInstances(client *golangsdk.ServiceClient, functionUrn string, policies []interface{}) error {
	for _, v := range policies {
		policy := v.(map[string]interface{})
		urn, err := getReservedInstanceUrn(client, functionUrn, policy)
		if err != nil {
			return err
		}
		// Deleting the alias will also delete the corresponding reserved instance.
		if urn == "" {
			return nil
		}
		opts := reserved.UpdateOpts{
			FuncUrn:  urn,
			Count:    pointerto.Int(0),
			IdleMode: pointerto.Bool(false),
		}
		_, err = reserved.Update(client, opts)
		if err != nil {
			return fmt.Errorf("error removing function reversed instance: %s", err)
		}
	}

	return nil
}

func addReservedInstances(client *golangsdk.ServiceClient, functionUrn string, addPolicies []interface{}) error {
	for _, v := range addPolicies {
		addPolicy := v.(map[string]interface{})
		urn, err := getReservedInstanceUrn(client, functionUrn, addPolicy)

		if err != nil {
			return err
		}

		opts := reserved.UpdateOpts{
			FuncUrn:       urn,
			Count:         pointerto.Int(addPolicy["count"].(int)),
			IdleMode:      pointerto.Bool(addPolicy["idle_mode"].(bool)),
			TacticsConfig: buildTacticsConfigs(addPolicy["tactics_config"].([]interface{})),
		}
		_, err = reserved.Update(client, opts)
		if err != nil {
			return fmt.Errorf("error updating function reversed instance: %s", err)
		}
	}

	return nil
}

func updateReservedInstanceConfig(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldRaw, newRaw := d.GetChange("reserved_instances")
	addRaw := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
	removeRaw := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
	functionUrn := resourceFgsFunctionUrn(d.Id())
	if removeRaw.Len() > 0 {
		if err := removeReservedInstances(client, functionUrn, removeRaw.List()); err != nil {
			return err
		}
	}

	if addRaw.Len() > 0 {
		if err := addReservedInstances(client, functionUrn, addRaw.List()); err != nil {
			return err
		}
	}

	return nil
}

func resourceFgsFunctionV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := config.FuncGraphV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	urn := resourceFgsFunctionUrn(d.Id())

	// lintignore:R019
	if d.HasChanges("code_type", "code_url", "code_filename", "depend_list", "func_code") {
		err := resourceFgsFunctionCodeUpdate(fgsClient, urn, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// lintignore:R019
	if d.HasChanges("app", "handler", "memory_size", "timeout", "encrypted_user_data",
		"user_data", "agency", "app_agency", "description", "initializer_handler", "initializer_timeout",
		"vpc_id", "network_id", "mount_user_id", "mount_user_group_id", "func_mounts", "custom_image",
		"log_group_id", "log_topic_id", "log_group_name", "log_topic_name", "concurrency_num", "gpu_memory", "gpu_type") {
		err := resourceFgsFunctionMetadataUpdate(fgsClient, urn, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("max_instance_num") {
		// The integer string of the maximum instance number has been already checked in the schema validation.
		maxInstanceNum, _ := strconv.Atoi(d.Get("max_instance_num").(string))

		_, err = function.UpdateMaxInstances(fgsClient, function.UpdateFuncInstancesOpts{
			FuncUrn:        urn,
			MaxInstanceNum: maxInstanceNum,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		if err = updateFunctionTags(fgsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("versions") {
		if err = updateFunctionVersions(fgsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("reserved_instances") {
		if err = updateReservedInstanceConfig(fgsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFgsFunctionV2Read(ctx, d, meta)
}

func resourceFgsFunctionV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := config.FuncGraphV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	urn := resourceFgsFunctionUrn(d.Id())

	err = function.Delete(fgsClient, urn)
	if err != nil {
		return diag.Errorf("error deleting function: %s", err)
	}
	return nil
}

func resourceFgsFunctionMetadataUpdate(fgsClient *golangsdk.ServiceClient, urn string, d *schema.ResourceData) error {
	// check app
	app, appOk := d.GetOk("app")
	if !appOk {
		return fmt.Errorf("app must be configured")
	}
	packV := ""
	if appOk {
		packV = app.(string)
	}

	agencyV := ""
	if v, ok := d.GetOk("agency"); ok {
		agencyV = v.(string)
	}

	updateMetadateOpts := function.UpdateFuncMetadataOpts{
		Name:              d.Get("name").(string),
		Handler:           d.Get("handler").(string),
		MemorySize:        d.Get("memory_size").(int),
		Timeout:           d.Get("timeout").(int),
		Runtime:           d.Get("runtime").(string),
		Package:           packV,
		Description:       d.Get("description").(string),
		UserData:          d.Get("user_data").(string),
		EncryptedUserData: d.Get("encrypted_user_data").(string),
		Xrole:             agencyV,
		AppXrole:          d.Get("app_agency").(string),
		InitHandler:       d.Get("initializer_handler").(string),
		InitTimeout:       pointerto.Int(d.Get("initializer_timeout").(int)),
		CustomImage:       buildCustomImage(d.Get("custom_image").([]interface{})),
		GpuMemory:         pointerto.Int(d.Get("gpu_memory").(int)),
	}

	if _, ok := d.GetOk("vpc_id"); ok {
		updateMetadateOpts.FuncVpc = resourceFgsFunctionFuncVpc(d)
	}

	if _, ok := d.GetOk("func_mounts"); ok {
		updateMetadateOpts.MountConfig = resourceFgsFunctionMountConfig(d)
	}

	// check name here as it will only save to sate if specified before
	if v, ok := d.GetOk("log_group_name"); ok {
		logConfig := function.FuncLogConfig{
			GroupID:    d.Get("log_group_id").(string),
			StreamID:   d.Get("log_topic_id").(string),
			GroupName:  v.(string),
			StreamName: d.Get("log_topic_name").(string),
		}
		updateMetadateOpts.LogConfig = &logConfig
	}

	if v, ok := d.GetOk("concurrency_num"); ok {
		strategyConfig := function.StrategyConfig{
			ConcurrentNum: v.(int),
		}
		updateMetadateOpts.StrategyConfig = &strategyConfig
	}

	log.Printf("[DEBUG] Metaddata Update Options: %#v", updateMetadateOpts)
	updateMetadateOpts.FuncUrn = urn

	_, err := function.UpdateFuncMetadata(fgsClient, updateMetadateOpts)
	if err != nil {
		return fmt.Errorf("error updating metadata of function: %s", err)
	}

	return nil
}

func resourceFgsFunctionCodeUpdate(fgsClient *golangsdk.ServiceClient, urn string, d *schema.ResourceData) error {
	updateCodeOpts := function.UpdateFuncCodeOpts{
		CodeType:     d.Get("code_type").(string),
		CodeURL:      d.Get("code_url").(string),
		CodeFilename: d.Get("code_filename").(string),
	}

	if v, ok := d.GetOk("depend_list"); ok {
		dependListRaw := v.(*schema.Set)
		dependList := make([]string, 0, dependListRaw.Len())
		for _, depend := range dependListRaw.List() {
			dependList = append(dependList, depend.(string))
		}
		updateCodeOpts.DependVersionList = dependList
	}

	if v, ok := d.GetOk("func_code"); ok {
		funcCode := function.FuncCode{
			File: hashcode.TryBase64EncodeString(v.(string)),
		}
		updateCodeOpts.FuncCode = &funcCode
	}

	updateCodeOpts.FuncUrn = urn
	log.Printf("[DEBUG] Code Update Options: %#v", updateCodeOpts)
	_, err := function.UpdateFuncCode(fgsClient, updateCodeOpts)
	if err != nil {
		return fmt.Errorf("error updating code of function: %s", err)
	}

	return nil
}

func resourceFgsFunctionFuncVpc(d *schema.ResourceData) *function.FuncVpc {
	var funcVpc function.FuncVpc
	funcVpc.VpcID = d.Get("vpc_id").(string)
	funcVpc.SubnetID = d.Get("network_id").(string)
	return &funcVpc
}

func resourceFgsFunctionMountConfig(d *schema.ResourceData) *function.MountConfig {
	var mountConfig function.MountConfig
	funcMountsRaw := d.Get("func_mounts").([]interface{})
	if len(funcMountsRaw) >= 1 {
		funcMounts := make([]function.FuncMount, 0, len(funcMountsRaw))
		for _, funcMountRaw := range funcMountsRaw {
			var funcMount function.FuncMount
			funcMountMap := funcMountRaw.(map[string]interface{})
			funcMount.MountType = funcMountMap["mount_type"].(string)
			funcMount.MountResource = funcMountMap["mount_resource"].(string)
			funcMount.MountSharePath = funcMountMap["mount_share_path"].(string)
			funcMount.LocalMountPath = funcMountMap["local_mount_path"].(string)

			funcMounts = append(funcMounts, funcMount)
		}

		mountConfig.FuncMounts = funcMounts

		mountUser := function.MountUser{
			UserID:      strconv.Itoa(-1),
			UserGroupID: strconv.Itoa(-1),
		}

		if v, ok := d.GetOk("mount_user_id"); ok {
			mountUser.UserID = v.(string)
		}

		if v, ok := d.GetOk("mount_user_group_id"); ok {
			mountUser.UserGroupID = v.(string)
		}

		mountConfig.MountUser = mountUser
	}
	return &mountConfig
}

/*
 * Parse urn according from fun_urn.
 * If the separator is not ":" then return to the original value.
 */
func resourceFgsFunctionUrn(urn string) string {
	index := strings.LastIndex(urn, ":")
	if index != -1 {
		urn = urn[0:index]
	}
	return urn
}
