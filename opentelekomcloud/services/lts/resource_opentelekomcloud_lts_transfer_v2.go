package lts

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/transfers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLTSNewTransferV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsTransferV2Create,
		UpdateContext: resourceLtsTransferV2Update,
		ReadContext:   resourceLtsTransferV2Read,
		DeleteContext: resourceLtsTransferV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_streams": {
				Type:     schema.TypeList,
				Elem:     ltsTransferLogStreamsSchema(),
				Required: true,
				ForceNew: true,
			},
			"log_transfer_info": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ltsTransferLogTransferInfoSchema(),
				Required: true,
			},
			"log_group_name": {
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

func ltsTransferLogStreamsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func ltsTransferLogTransferInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"log_transfer_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_transfer_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_storage_format": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_transfer_status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_agency_transfer": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ltsTransferLogAgencySchema(),
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"log_transfer_detail": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ltsTransferLogDetailSchema(),
				Required: true,
			},
			"log_created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func ltsTransferLogAgencySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"agency_domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"agency_domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"agency_project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	return &sc
}

func ltsTransferLogDetailSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"obs_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"obs_period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"obs_bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"obs_transfer_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"obs_dir_prefix_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"obs_prefix_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"obs_eps_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"obs_encrypted_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"obs_encrypted_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"obs_time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"obs_time_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func resourceLtsTransferV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	createOpts := transfers.CreateOpts{
		LogGroupId:      d.Get("log_group_id").(string),
		LogStreams:      buildLogStreams(d.Get("log_streams")),
		LogTransferInfo: buildLogTransferInfo(d.Get("log_transfer_info"), config.DomainID, client.ProjectID),
	}
	resp, err := transfers.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS v2 log transfer: %s", err)
	}
	d.SetId(resp.LogTransferId)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceLtsTransferV2Read(clientCtx, d, meta)
}

func buildLogStreams(rawParams interface{}) []transfers.LogStreams {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]transfers.LogStreams, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = transfers.LogStreams{
				LogStreamId:   raw["log_stream_id"].(string),
				LogStreamName: raw["log_stream_name"].(string),
			}
		}
		return rst
	}
	return nil
}

func buildLogTransferInfo(rawParams interface{}, domainID, projectID string) *transfers.LogTransferInfo {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := transfers.LogTransferInfo{
			LogTransferType:   raw["log_transfer_type"].(string),
			LogTransferMode:   raw["log_transfer_mode"].(string),
			LogStorageFormat:  raw["log_storage_format"].(string),
			LogTransferStatus: raw["log_transfer_status"].(string),
			LogAgencyTransfer: buildLogTransferInfoLogAgency(raw["log_agency_transfer"], domainID, projectID),
			LogTransferDetail: buildLogTransferInfoLogTransferDetail(raw["log_transfer_detail"]),
		}
		return &params
	}
	return nil
}

func buildLogTransferInfoLogAgency(rawParams interface{}, domainID, projectID string) *transfers.LogAgencyTransfer {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := transfers.LogAgencyTransfer{
			AgencyDomainId:    raw["agency_domain_id"].(string),
			AgencyDomainName:  raw["agency_domain_name"].(string),
			AgencyName:        raw["agency_name"].(string),
			AgencyProjectId:   raw["agency_project_id"].(string),
			BeAgencyDomainId:  domainID,
			BeAgencyProjectId: projectID,
		}
		return &params
	}
	return nil
}

func buildLogTransferInfoLogTransferDetail(rawParams interface{}) *transfers.TransferDetail {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := transfers.TransferDetail{
			ObsPeriod:           raw["obs_period"].(int),
			ObsPeriodUnit:       raw["obs_period_unit"].(string),
			ObsBucketName:       raw["obs_bucket_name"].(string),
			ObsTransferPath:     raw["obs_transfer_path"].(string),
			ObsDirPreFixName:    raw["obs_dir_prefix_name"].(string),
			ObsPrefixName:       raw["obs_prefix_name"].(string),
			EnterpriseProjectID: raw["obs_eps_id"].(string),
			ObsEncryptedEnable:  raw["obs_encrypted_enable"].(bool),
			ObsEncryptedId:      raw["obs_encrypted_id"].(string),
			ObsTimeZone:         raw["obs_time_zone"].(string),
			ObsTimeZoneId:       raw["obs_time_zone_id"].(string),
			Tags:                common.ExpandToStringList(raw["tags"].([]interface{})),
		}
		return &params
	}
	return nil
}

func resourceLtsTransferV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	requestResp, err := transfers.List(client, transfers.ListTransfersOpts{})
	var transferResult transfers.Transfer

	for _, tr := range requestResp {
		if tr.LogTransferId == d.Id() {
			transferResult = tr
			break
		}
	}
	if transferResult.LogTransferId == "" {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find OpenTelekomCloud LTS v2 log transfer by its ID (%s)", d.Id()))
	}
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("log_group_id", transferResult.LogGroupId),
		d.Set("log_group_name", transferResult.LogGroupName),
		d.Set("log_streams", flattenGetTransferResponseBodyLogStreams(transferResult.LogStreams)),
		d.Set("log_transfer_info", flattenGetTransferResponseBodyLogTransferInfo(&transferResult.LogTransferInfo)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetTransferResponseBodyLogStreams(resp []transfers.LogStreamsResponse) []interface{} {
	if len(resp) < 1 {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"log_stream_id":   v.LogStreamId,
			"log_stream_name": v.LogStreamName,
		})
	}
	return rst
}

func flattenGetTransferResponseBodyLogTransferInfo(resp *transfers.LogTransferInfoResponse) []interface{} {
	if resp == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"log_transfer_type":   resp.LogTransferType,
			"log_transfer_mode":   resp.LogTransferMode,
			"log_storage_format":  resp.LogStorageFormat,
			"log_transfer_status": resp.LogTransferStatus,
			"log_agency_transfer": flattenLogTransferInfoLogAgency(resp.LogAgencyTransfer),
			"log_transfer_detail": flattenLogTransferInfoLogTransferDetail(resp.LogTransferDetail),
		},
	}
	return rst
}

func flattenLogTransferInfoLogAgency(resp *transfers.LogAgencyTransferResponse) []interface{} {
	if resp == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"agency_domain_id":   resp.AgencyDomainId,
			"agency_domain_name": resp.AgencyDomainName,
			"agency_name":        resp.AgencyName,
			"agency_project_id":  resp.AgencyProjectId,
		},
	}
	return rst
}

func flattenLogTransferInfoLogTransferDetail(resp *transfers.TransferDetailResponse) []interface{} {
	if resp == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"obs_period":           resp.ObsPeriod,
			"obs_period_unit":      resp.ObsPeriodUnit,
			"obs_bucket_name":      resp.ObsBucketName,
			"obs_transfer_path":    resp.ObsTransferPath,
			"obs_dir_prefix_name":  resp.ObsDirPreFixName,
			"obs_prefix_name":      resp.ObsPrefixName,
			"obs_eps_id":           resp.EnterpriseProjectID,
			"obs_encrypted_enable": resp.ObsEncryptedEnable,
			"obs_encrypted_id":     resp.ObsEncryptedId,
			"obs_time_zone":        resp.ObsTimeZone,
			"obs_time_zone_id":     resp.ObsTimeZoneId,
			"tags":                 resp.Tags,
		},
	}
	return rst
}

func resourceLtsTransferV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	updateTransferChanges := []string{
		"log_transfer_info",
	}

	if d.HasChanges(updateTransferChanges...) {
		_, err = transfers.Update(client, transfers.UpdateTransferOpts{
			TransferId:   d.Id(),
			TransferInfo: buildTransferInfoUpdate(d.Get("log_transfer_info")),
		})
		if err != nil {
			return diag.Errorf("error updating OpenTelekomCloud LTS v2 transfer: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceLtsTransferV2Read(clientCtx, d, meta)
}

func buildTransferInfoUpdate(rawParams interface{}) *transfers.TransferInfoUpdate {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := transfers.TransferInfoUpdate{
			StorageFormat:  raw["log_storage_format"].(string),
			TransferStatus: raw["log_transfer_status"].(string),
			TransferDetail: buildLogTransferInfoLogTransferDetail(raw["log_transfer_detail"]),
		}
		return &params
	}
	return nil
}

func resourceLtsTransferV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = transfers.Delete(client, d.Id())
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault400); ok {
			d.SetId("")
			return nil
		} else {
			return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud LTS v2 log transfer")
		}
	}

	return nil
}
