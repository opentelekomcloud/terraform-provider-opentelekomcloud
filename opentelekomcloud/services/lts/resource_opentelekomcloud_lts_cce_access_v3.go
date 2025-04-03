package lts

import (
	"context"
	"fmt"

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

func ResourceCceAccessV3Config() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCceAccessConfigV3Create,
		UpdateContext: resourceCceAccessConfigV3Update,
		ReadContext:   resourceCceAccessConfigV3Read,
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
				Elem:     cceAccessConfigDetailSchema(),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_group_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"binary_collect": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"log_split": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
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
			"created_at": {
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

func cceAccessConfigDetailSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"path_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"paths": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"black_paths": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"windows_log_info": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     windowsLogInfoSchema(),
				Optional: true,
			},
			"single_log_format": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				Elem:         logFormatSchema(),
				ExactlyOneOf: []string{"access_config.0.multi_log_format"},
			},
			"multi_log_format": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     logFormatSchema(),
			},
			"stdout": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"stderr": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"name_space_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pod_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"include_labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"exclude_labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"log_envs": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"include_envs": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"exclude_envs": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"log_k8s": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"include_k8s_labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"exclude_k8s_labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func resourceCceAccessConfigV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		Type:    "K8S_CCE",
		Details: buildCceDetailRequestBody(d.Get("access_config")),
		LogInfo: &ac.LogInfo{
			LogGroupId:  d.Get("log_group_id").(string),
			LogStreamId: d.Get("log_stream_id").(string),
		},
		HostGroupInfo: &ac.HostGroupInfo{
			HostGroupIds: &hosts,
		},
		Tags:          ltsTags(d),
		ClusterId:     d.Get("cluster_id").(string),
		BinaryCollect: pointerto.Bool(d.Get("binary_collect").(bool)),
		LogSplit:      pointerto.Bool(d.Get("log_split").(bool)),
	}
	access, err := ac.Create(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS v3 cce access: %s", err)
	}
	d.SetId(access.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceCceAccessConfigV3Read(clientCtx, d, meta)
}

func buildCceDetailRequestBody(rawParams interface{}) *ac.AccessConfigDetails {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) < 1 {
		return nil
	}

	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	params := ac.AccessConfigDetails{
		PathType:           raw["path_type"].(string),
		Paths:              common.ExpandToStringList(raw["paths"].(*schema.Set).List()),
		BlackPaths:         common.ExpandToStringList(raw["black_paths"].(*schema.Set).List()),
		Format:             buildFormatBody(raw),
		WindowsLogInfo:     buildWindowsLogInfoBody(raw["windows_log_info"]),
		Stdout:             pointerto.Bool(raw["stdout"].(bool)),
		Stderr:             pointerto.Bool(raw["stderr"].(bool)),
		NamespaceRegex:     raw["name_space_regex"].(string),
		PodNameRegex:       raw["pod_name_regex"].(string),
		ContainerNameRegex: raw["container_name_regex"].(string),
		LogLabels:          common.ConvertToMapString(raw["log_labels"]),
		IncludeLabels:      common.ConvertToMapString(raw["include_labels"]),
		ExcludeLabels:      common.ConvertToMapString(raw["exclude_labels"]),
		LogEnvs:            common.ConvertToMapString(raw["log_envs"]),
		IncludeEnvs:        common.ConvertToMapString(raw["include_envs"]),
		ExcludeEnvs:        common.ConvertToMapString(raw["exclude_envs"]),
		LogK8s:             common.ConvertToMapString(raw["log_k8s"]),
		IncludeK8sLabels:   common.ConvertToMapString(raw["include_k8s_labels"]),
		ExcludeK8sLabels:   common.ConvertToMapString(raw["exclude_k8s_labels"]),
	}
	return &params
}

func resourceCceAccessConfigV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find OpenTelekomCloud LTS v2 cce access config by its ID (%s)", d.Id()))
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
		d.Set("access_config", flattenCceAccessConfigDetail(configResult.AccessConfigDetail)),
		d.Set("cluster_id", configResult.ClusterId),
		d.Set("binary_collect", configResult.BinaryCollect),
		d.Set("log_split", configResult.LogSplit),
		d.Set("created_at", common.FormatTimeStampRFC3339(configResult.CreatedAt/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCceAccessConfigDetail(resp *ac.AccessConfigDetailResponse) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"path_type":            resp.PathType,
			"paths":                resp.Paths,
			"black_paths":          resp.BlackPaths,
			"single_log_format":    flattenHostAccessConfigLogFormat(resp.Format.Single),
			"multi_log_format":     flattenHostAccessConfigLogFormat(resp.Format.Multi),
			"windows_log_info":     flattenHostAccessConfigWindowsLogInfo(resp.WindowsLogInfo),
			"stdout":               resp.Stdout,
			"stderr":               resp.Stderr,
			"name_space_regex":     resp.NamespaceRegex,
			"pod_name_regex":       resp.PodNameRegex,
			"container_name_regex": resp.ContainerNameRegex,
			"log_labels":           resp.LogLabels,
			"include_labels":       resp.IncludeLabels,
			"exclude_labels":       resp.ExcludeLabels,
			"log_envs":             resp.LogEnvs,
			"include_envs":         resp.IncludeEnvs,
			"exclude_envs":         resp.ExcludeEnvs,
			"log_k8s":              resp.LogK8s,
			"include_k8s_labels":   resp.IncludeK8sLabels,
			"exclude_k8s_labels":   resp.ExcludeK8sLabels,
		},
	}
}

func resourceCceAccessConfigV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	updateCceAccessConfigChanges := []string{
		"access_config",
		"host_group_ids",
		"tags",
		"binary_collect",
		"log_split",
	}

	if d.HasChanges(updateCceAccessConfigChanges...) {
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
			Details:       buildCceDetailUpdateRequestBody(d.Get("access_config")),
			BinaryCollect: pointerto.Bool(d.Get("binary_collect").(bool)),
			LogSplit:      pointerto.Bool(d.Get("log_split").(bool)),
			ClusterId:     d.Get("cluster_id").(string),
			Tags:          &tagSlice,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceCceAccessConfigV3Read(clientCtx, d, meta)
}

func buildCceDetailUpdateRequestBody(rawParams interface{}) *ac.AccessConfigDetailsUpdate {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) < 1 {
		return nil
	}

	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}
	logInfo := buildWindowsLogInfoUpdateBody(raw["windows_log_info"])
	params := ac.AccessConfigDetailsUpdate{
		PathType:           raw["path_type"].(string),
		Paths:              common.ExpandToStringList(raw["paths"].(*schema.Set).List()),
		BlackPaths:         common.ExpandToStringList(raw["black_paths"].(*schema.Set).List()),
		Format:             buildFormatBody(raw),
		WindowsLogInfo:     logInfo,
		Stdout:             raw["stdout"].(bool),
		Stderr:             raw["stderr"].(bool),
		NamespaceRegex:     raw["name_space_regex"].(string),
		PodNameRegex:       raw["pod_name_regex"].(string),
		ContainerNameRegex: raw["container_name_regex"].(string),
		LogLabels:          common.ConvertToMapString(raw["log_labels"]),
		IncludeLabels:      common.ConvertToMapString(raw["include_labels"]),
		ExcludeLabels:      common.ConvertToMapString(raw["exclude_labels"]),
		LogEnvs:            common.ConvertToMapString(raw["log_envs"]),
		IncludeEnvs:        common.ConvertToMapString(raw["include_envs"]),
		ExcludeEnvs:        common.ConvertToMapString(raw["exclude_envs"]),
		LogK8s:             common.ConvertToMapString(raw["log_k8s"]),
		IncludeK8sLabels:   common.ConvertToMapString(raw["include_k8s_labels"]),
		ExcludeK8sLabels:   common.ConvertToMapString(raw["exclude_k8s_labels"]),
	}
	return &params
}
