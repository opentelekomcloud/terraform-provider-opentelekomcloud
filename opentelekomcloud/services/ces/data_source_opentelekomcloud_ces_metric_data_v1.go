package ces

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metricdata"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCesMetricDataV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesMetricDataRead,

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 32),
					validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*\.[a-zA-Z][a-zA-Z0-9_]*$`),
						"Must be of type service.item. service and item each must only have lowercase/uppercase letters, digits, and underscores (_) and must start with a letter.",
					),
				),
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"to": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dim0": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dim1": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datapoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"average": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"max": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"min": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"sum": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"variance": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCesMetricDataRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	metricDataOpts := metricdata.ShowMetricDataOpts{
		Namespace:  d.Get("namespace").(string),
		MetricName: d.Get("metric_name").(string),
		From:       d.Get("from").(string),
		To:         d.Get("to").(string),
		Period:     d.Get("period").(int),
		Filter:     d.Get("filter").(string),
		Dim0:       d.Get("dim0").(string),
		Dim1:       d.Get("dim1").(string),
	}

	metricData, err := metricdata.ShowMetricData(client, metricDataOpts)
	if err != nil {
		return fmterr.Errorf("error fetching metric data : %w", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("datapoints", setDatapoints(metricData.Datapoints)),
		d.Set("metric_name", metricData.MetricName),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setDatapoints(datapointsInResp []metricdata.Datapoint) []map[string]interface{} {
	var datapoints []map[string]interface{}
	for _, datapointInResp := range datapointsInResp {
		datapoint := map[string]interface{}{
			"average":   datapointInResp.Average,
			"max":       datapointInResp.Max,
			"min":       datapointInResp.Min,
			"sum":       datapointInResp.Sum,
			"variance":  datapointInResp.Variance,
			"timestamp": datapointInResp.Timestamp,
			"unit":      datapointInResp.Unit,
		}

		datapoints = append(datapoints, datapoint)
	}
	return datapoints
}
