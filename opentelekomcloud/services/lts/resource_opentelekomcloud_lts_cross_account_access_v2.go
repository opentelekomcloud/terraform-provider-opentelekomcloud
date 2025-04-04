package lts

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	ac "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/access-config"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCrossAccountAccessV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCrossAccountAccessV2Create,
		UpdateContext: resourceCrossAccountAccessV2Update,
		ReadContext:   resourceCrossAccountAccessV2Read,
		DeleteContext: resourceHostAccessConfigV3Delete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_agency_stream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_agency_stream_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_agency_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_agency_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"agency_project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"agency_domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_config_type": {
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

func resourceCrossAccountAccessV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV20, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV20Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV20Client, err)
	}

	access, err := ac.CrossAccess(client, ac.CreateCrossOpts{
		PreviewAgencyList: []ac.PreviewAgencyLogAccess{
			{
				Type: "AGENCYACCESS",
				Name: d.Get("name").(string),

				StreamName: d.Get("log_stream_name").(string),
				StreamId:   d.Get("log_stream_id").(string),
				GroupName:  d.Get("log_group_name").(string),
				GroupId:    d.Get("log_group_id").(string),

				AgencyStreamName: d.Get("log_agency_stream_name").(string),
				AgencyStreamId:   d.Get("log_agency_stream_id").(string),
				AgencyGroupName:  d.Get("log_agency_group_name").(string),
				AgencyGroupId:    d.Get("log_agency_group_id").(string),

				ProjectId:        client.ProjectID,
				AgencyProjectId:  d.Get("agency_project_id").(string),
				AgencyDomainName: d.Get("agency_domain_name").(string),
				AgencyName:       d.Get("agency_name").(string),
			},
		},
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS v2 cross account access: %s", err)
	}
	d.SetId(access[0].ID)

	if _, ok := d.GetOk("tags"); ok {
		errTags := updateTags(d, meta, "ltsAccessConfig", d.Id())
		if errTags != nil {
			return diag.Errorf("error creating LTS cross account access tags: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceCrossAccountAccessV2Read(clientCtx, d, meta)
}

func resourceCrossAccountAccessV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var configResult ac.AccessConfigInfo
	for _, acc := range requestResp.Result {
		if acc.Name == d.Get("name").(string) {
			configResult = acc
		}
	}
	if configResult.ID == "" {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find OpenTelekomCloud LTS v3 cross account access config by its ID (%s)", d.Id()))
	}
	tagsMap := make(map[string]string)
	for _, tag := range configResult.Tags {
		tagsMap[tag.Key] = tag.Value
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", configResult.Name),
		d.Set("created_at", common.FormatTimeStampRFC3339(configResult.CreatedAt/1000, false)),
		d.Set("tags", tagsMap),
		d.Set("access_config_type", configResult.Type),
		d.Set("log_group_id", configResult.LogInfo.LogGroupId),
		d.Set("log_group_name", configResult.LogInfo.LogGroupName),
		d.Set("log_stream_id", configResult.LogInfo.LogStreamId),
		d.Set("log_stream_name", configResult.LogInfo.LogStreamName),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCrossAccountAccessV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("tags") {
		tagErr := updateTags(d, meta, "ltsAccessConfig", d.Id())
		if tagErr != nil {
			return diag.Errorf("unable to update tags for OpenTelekomCloud LTS v2 cross account access: %s, err:%s", d.Id(), tagErr)
		}
	}
	return resourceCrossAccountAccessV2Read(ctx, d, meta)
}
