package lts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	quick_search "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/quick-search"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceQuickSearchCriteriaV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceQuickSearchCriteriaV1Create,
		ReadContext:   resourceQuickSearchCriteriaV1Read,
		DeleteContext: resourceQuickSearchCriteriaV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSearchCriteriaImportState,
		},

		Schema: map[string]*schema.Schema{
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
			"criteria": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceQuickSearchCriteriaV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV10, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV10Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV10Client, err)
	}

	criteriaId, err := quick_search.Create(client,
		d.Get("log_group_id").(string),
		d.Get("log_stream_id").(string),
		quick_search.CreateOpts{
			Criteria:   d.Get("criteria").(string),
			Name:       d.Get("name").(string),
			SearchType: d.Get("type").(string),
		})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS v1.0 search criteria: %s", err)
	}
	d.SetId(criteriaId)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceQuickSearchCriteriaV1Read(clientCtx, d, meta)
}

func resourceQuickSearchCriteriaV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV10, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV10Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV10Client, err)
	}

	requestResp, err := quick_search.ListCriterias(
		client,
		d.Get("log_group_id").(string),
		d.Get("log_stream_id").(string),
		quick_search.ListOpts{
			SearchType: d.Get("type").(string),
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}
	var searchResult quick_search.SearchCriteria
	for _, sq := range requestResp {
		if sq.ID == d.Id() {
			searchResult = sq
			break
		}
	}
	if searchResult.ID == "" {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find OpenTelekomCloud v1.0 search criteria by its ID (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("criteria", searchResult.Criteria),
		d.Set("name", searchResult.Name),
		d.Set("type", searchResult.SearchType),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceQuickSearchCriteriaV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV10, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV10Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV10Client, err)
	}

	err = quick_search.Delete(client,
		d.Get("log_group_id").(string),
		d.Get("log_stream_id").(string),
		quick_search.DeleteOpts{
			ID: d.Id(),
		})
	if err != nil {
		return diag.Errorf("error deleting search criteria: %s", err)
	}

	return nil
}

func resourceSearchCriteriaImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID format, want '<log_group_id>/<log_stream_id>/<id>', but got '%s'", d.Id())
	}

	groupID := parts[0]
	streamID := parts[1]
	searchCriteriaID := parts[2]

	d.SetId(searchCriteriaID)
	mErr := multierror.Append(nil,
		d.Set("log_group_id", groupID),
		d.Set("log_stream_id", streamID),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
