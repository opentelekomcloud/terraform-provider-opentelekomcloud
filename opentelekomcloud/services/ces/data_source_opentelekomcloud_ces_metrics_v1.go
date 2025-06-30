package ces

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metrics"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCesMetricsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesMetricsRead,

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dim": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1000,
				ValidateFunc: validation.IntBetween(1, 1000),
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"asc", "desc",
				}, false),
			},
			"metrics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dimensions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"meta_data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"marker": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCesMetricsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	listOpts := metrics.ListMetricsRequest{
		Namespace:  d.Get("namespace").(string),
		MetricName: d.Get("metric_name").(string),
		Dim:        d.Get("dim").(string),
		Start:      d.Get("start").(string),
		Limit:      InterfaceToIntPtr(d.Get("limit")),
		Order:      d.Get("order").(string),
	}

	page := metrics.ListMetrics(client, listOpts)
	pages, err := page.AllPages()
	if err != nil {
		return fmterr.Errorf("error fetching pages : %w", err)
	}
	listMetrics, err := metrics.ExtractAllPagesMetrics(pages)
	if err != nil {
		return fmterr.Errorf("error extracting metrics : %w", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Retrieved metrics list %s: %#v", d.Id(), listMetrics)

	metadata := []map[string]interface{}{
		{
			"count":  listMetrics.MetaData.Count,
			"marker": listMetrics.MetaData.Marker,
			"total":  listMetrics.MetaData.Total,
		},
	}

	mErr := multierror.Append(
		d.Set("metrics", setMetricsInfo(listMetrics.Metrics)),
		d.Set("meta_data", metadata),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setMetricsInfo(dimensionsInResp []metrics.MetricInfoList) []map[string]interface{} {
	var metricsInfoList []map[string]interface{}
	for _, metricsInfoInResp := range dimensionsInResp {
		metricsInfo := map[string]interface{}{
			"namespace":   metricsInfoInResp.Namespace,
			"dimensions":  setMetricsDimension(metricsInfoInResp.Dimensions),
			"metric_name": metricsInfoInResp.MetricName,
			"unit":        metricsInfoInResp.Unit,
		}
		metricsInfoList = append(metricsInfoList, metricsInfo)
	}
	return metricsInfoList
}

func setMetricsDimension(dimensionsListInResp []metrics.MetricsDimension) []map[string]interface{} {
	var dimensionsList []map[string]interface{}
	for _, dimensionsInResp := range dimensionsListInResp {
		dimensions := map[string]interface{}{
			"name":  dimensionsInResp.Name,
			"value": dimensionsInResp.Value,
		}
		dimensionsList = append(dimensionsList, dimensions)
	}
	return dimensionsList
}
