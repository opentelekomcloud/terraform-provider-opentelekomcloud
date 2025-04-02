package lts

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	ac "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/access-config"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceHostAccessConfigV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostAccessConfigV3Create,
		UpdateContext: resourceHostAccessConfigV3Update,
		ReadContext:   resourceHostAccessConfigV3Read,
		DeleteContext: resourceHostAccessConfigV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem:     detailSchema("access_config.0"),
			},
			"host_group_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"log_split": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"binary_collect": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"access_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_stream_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func detailSchema(parent string) *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"paths": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"black_paths": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"single_log_format": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     logFormatSchema(),
				ExactlyOneOf: []string{
					fmt.Sprintf("%s.multi_log_format", parent),
				},
			},
			"multi_log_format": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     logFormatSchema(),
			},
			"windows_log_info": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     windowsLogInfoSchema(),
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func logFormatSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func windowsLogInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"categories": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"event_level": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"time_offset": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"time_offset_unit": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func resourceHostAccessConfigV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	hosts := common.ExpandToStringList(d.Get("host_group_ids").([]interface{}))
	opts := ac.CreateOpts{
		Name:    d.Get("name").(string),
		Type:    "AGENT",
		Details: buildDetailRequestBody(d.Get("access_config")),
		LogInfo: &ac.LogInfo{
			LogGroupId:  d.Get("log_group_id").(string),
			LogStreamId: d.Get("log_stream_id").(string),
		},
		HostGroupInfo: &ac.HostGroupInfo{
			HostGroupIds: &hosts,
		},
		Tags:          ltsTags(d),
		BinaryCollect: pointerto.Bool(d.Get("binary_collect").(bool)),
		LogSplit:      pointerto.Bool(d.Get("log_split").(bool)),
	}
	access, err := ac.Create(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS v3 host group: %s", err)
	}
	d.SetId(access.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceHostAccessConfigV3Read(clientCtx, d, meta)
}

func buildDetailRequestBody(rawParams interface{}) *ac.AccessConfigDetails {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := ac.AccessConfigDetails{
			Paths:          common.ExpandToStringList(raw["paths"].(*schema.Set).List()),
			BlackPaths:     common.ExpandToStringList(raw["black_paths"].(*schema.Set).List()),
			Format:         buildFormatBody(raw),
			WindowsLogInfo: buildWindowsLogInfoBody(raw["windows_log_info"]),
		}
		return &params
	}
	return nil
}

func buildFormatBody(rawParams map[string]interface{}) *ac.AccessConfigFormat {
	log.Printf("[DEBUG] single_log_format: %#v", rawParams["single_log_format"])
	log.Printf("[DEBUG] multi_log_format: %#v", rawParams["multi_log_format"])

	if singleRaw, ok := rawParams["single_log_format"]; ok {
		if rawArray, ok := singleRaw.([]interface{}); ok {
			if len(rawArray) > 0 {
				raw, ok := rawArray[0].(map[string]interface{})
				if ok && len(raw) > 0 {
					return &ac.AccessConfigFormat{
						Single: buildFormatModeBody(raw),
					}
				}
			}
		}
	}

	if multiRaw, ok := rawParams["multi_log_format"]; ok {
		if rawArray, ok := multiRaw.([]interface{}); ok {
			if len(rawArray) > 0 {
				raw, ok := rawArray[0].(map[string]interface{})
				if ok && len(raw) > 0 {
					return &ac.AccessConfigFormat{
						Multi: buildFormatModeBody(raw),
					}
				}
			}
		}
	}

	return nil
}

func buildFormatModeBody(rawParams map[string]interface{}) *ac.AccessConfigFormatBody {
	mode := rawParams["mode"].(string)
	value := rawParams["value"].(string)

	if mode == "system" && value == "" {
		value = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}
	return &ac.AccessConfigFormatBody{
		Mode:  mode,
		Value: value,
	}
}

func buildWindowsLogInfoBody(rawParams interface{}) *ac.AccessConfigWindowsLogInfo {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		rawCategories := raw["categories"].([]interface{})
		categories := make([]string, len(rawCategories))
		for i, r := range rawCategories {
			categories[i] = r.(string)
		}

		rawEventLevels := raw["event_level"].([]interface{})
		eventLevels := make([]string, len(rawEventLevels))
		for i, r := range rawEventLevels {
			eventLevels[i] = r.(string)
		}

		timeOffsetOpts := ac.AccessConfigTimeOffset{
			Offset: int64(raw["time_offset"].(int)),
			Unit:   raw["time_offset_unit"].(string),
		}
		params := ac.AccessConfigWindowsLogInfo{
			Categories: categories,
			EventLevel: eventLevels,
			TimeOffset: &timeOffsetOpts,
		}
		return &params
	}
	return nil
}

func resourceHostAccessConfigV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	requestResp, err := ac.List(client, ac.ListOpts{})
	if err != nil {
		return diag.FromErr(err)
	}
	var configResult *ac.AccessConfigInfo
	for _, acc := range requestResp.Result {
		if acc.ID == d.Id() {
			configResult = &acc
		}
	}
	if configResult == nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find OpenTelekomCloud LTS v2 host access config by its ID (%s)", d.Id()))
	}
	tagsMap := make(map[string]string)
	for _, tag := range configResult.Tags {
		tagsMap[tag.Key] = tag.Value
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", configResult.Name),
		d.Set("access_type", configResult.Type),
		d.Set("log_group_id", configResult.LogInfo.LogGroupId),
		d.Set("log_stream_id", configResult.LogInfo.LogStreamId),
		d.Set("log_group_name", configResult.LogInfo.LogGroupName),
		d.Set("log_stream_name", configResult.LogInfo.LogStreamName),
		d.Set("host_group_ids", getHostGroupIDs(configResult.HostGroupInfo)),
		d.Set("tags", tagsMap),
		d.Set("access_config", flattenHostAccessConfigDetail(configResult.AccessConfigDetail)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getHostGroupIDs(configResult *ac.AccessConfigHostGroupIdsResponse) []string {
	if configResult == nil || configResult.HostGroupIds == nil {
		return []string{}
	}
	return configResult.HostGroupIds
}

func flattenHostAccessConfigDetail(resp *ac.AccessConfigDetailResponse) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"paths":             resp.Paths,
			"black_paths":       resp.BlackPaths,
			"single_log_format": flattenHostAccessConfigLogFormat(resp.Format.Single),
			"multi_log_format":  flattenHostAccessConfigLogFormat(resp.Format.Multi),
			"windows_log_info":  flattenHostAccessConfigWindowsLogInfo(resp.WindowsLogInfo),
		},
	}
}

func flattenHostAccessConfigLogFormat(resp *ac.AccessConfigFormatBody) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"mode":  resp.Mode,
			"value": resp.Value,
		},
	}
}

func flattenHostAccessConfigWindowsLogInfo(resp *ac.AccessConfigWindowsLogInfoResponse) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"categories":       resp.Categories,
			"event_level":      resp.EventLevel,
			"time_offset":      resp.TimeOffset.Offset,
			"time_offset_unit": resp.TimeOffset.Unit,
		},
	}
}

func resourceHostAccessConfigV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	updateHostAccessConfigChanges := []string{
		"access_config",
		"host_group_ids",
		"tags",
		"log_split",
		"binary_collect",
	}

	if d.HasChanges(updateHostAccessConfigChanges...) {
		hosts := common.ExpandToStringList(d.Get("host_group_ids").([]interface{}))
		tagSlice := ltsTags(d)
		if tagSlice == nil {
			tagSlice = []tags.ResourceTag{}
		}
		_, err = ac.Update(client, ac.UpdateOpts{
			ID: d.Id(),
			HostGroupInfo: &ac.HostGroupInfo{
				HostGroupIds: &hosts,
			},
			Details:       buildUpdateDetailRequestBody(d.Get("access_config")),
			BinaryCollect: pointerto.Bool(d.Get("binary_collect").(bool)),
			LogSplit:      pointerto.Bool(d.Get("log_split").(bool)),
			Tags:          &tagSlice,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceHostAccessConfigV3Read(clientCtx, d, meta)
}

func buildUpdateDetailRequestBody(rawParams interface{}) *ac.AccessConfigDetailsUpdate {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}
		logInfo := buildWindowsLogInfoUpdateBody(raw["windows_log_info"])
		params := ac.AccessConfigDetailsUpdate{
			Paths:          common.ExpandToStringList(raw["paths"].(*schema.Set).List()),
			BlackPaths:     common.ExpandToStringList(raw["black_paths"].(*schema.Set).List()),
			Format:         buildFormatBody(raw),
			WindowsLogInfo: logInfo,
		}
		return &params
	}
	return nil
}

func buildWindowsLogInfoUpdateBody(rawParams interface{}) *ac.AccessConfigWindowsLogInfoUpdate {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		timeOffsetOpts := ac.AccessConfigTimeOffset{
			Offset: int64(raw["time_offset"].(int)),
			Unit:   raw["time_offset_unit"].(string),
		}
		params := ac.AccessConfigWindowsLogInfoUpdate{
			Categories: common.ExpandToStringList(raw["categories"].(*schema.Set).List()),
			EventLevel: common.ExpandToStringList(raw["event_level"].(*schema.Set).List()),
			TimeOffset: &timeOffsetOpts,
		}
		return &params
	}
	return nil
}

func resourceHostAccessConfigV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	_, err = ac.Delete(client, ac.DeleteOpts{AccessConfigIds: []string{d.Id()}})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud LTS v2 host access")
	}

	return nil
}
