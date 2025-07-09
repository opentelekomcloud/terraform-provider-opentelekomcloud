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

func DataSourceCesMultipleMetricDataV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesMultipleMetricDataRead,

		Schema: map[string]*schema.Schema{
			"metrics": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 10,
				Elem: &schema.Resource{
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
						"dimensions": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 4,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"unit": {
							Type:     schema.TypeString,
							Computed: true,
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
								},
							},
						},
					},
				},
			},
			"from": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"to": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"period": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"1", "300", "1200", "3600", "14400", "86400",
				}, true),
			},
			"filter": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"average", "max", "min", "sum", "variance",
				}, false),
			},
		},
	}
}

func getMetricDataBatchMetrics(d *schema.ResourceData) []metricdata.Metric {
	metricListRaw := d.Get("metrics").([]interface{})
	var metricsList []metricdata.Metric

	for _, v := range metricListRaw {
		metricRaw := v.(map[string]interface{})
		var metricDimensions []metricdata.MetricsDimension
		dimensionsRaw := metricRaw["dimensions"].([]interface{})
		for _, vd := range dimensionsRaw {
			dimensionRaw := vd.(map[string]interface{})
			dimensions := metricdata.MetricsDimension{
				Name:  dimensionRaw["name"].(string),
				Value: dimensionRaw["value"].(string),
			}
			metricDimensions = append(metricDimensions, dimensions)
		}
		metric := metricdata.Metric{
			Namespace:  metricRaw["namespace"].(string),
			MetricName: metricRaw["metric_name"].(string),
			Dimensions: metricDimensions,
		}
		metricsList = append(metricsList, metric)
	}

	return metricsList
}

func dataSourceCesMultipleMetricDataRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	metricDataBatchOpts := metricdata.BatchListMetricDataOpts{
		Metrics: getMetricDataBatchMetrics(d),
		From:    InterfaceToInt64(d.Get("from")),
		To:      InterfaceToInt64(d.Get("to")),
		Period:  d.Get("period").(string),
		Filter:  d.Get("filter").(string),
	}

	metricData, err := metricdata.BatchListMetricData(client, metricDataBatchOpts)
	if err != nil {
		return fmterr.Errorf("error fetching metric data : %w", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("metrics", setBatchMetrics(metricData)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setBatchMetrics(metricsListInResp []metricdata.BatchMetricData) []map[string]interface{} {
	var metricsList []map[string]interface{}
	for _, metricsInResp := range metricsListInResp {
		metrics := map[string]interface{}{
			"namespace":   metricsInResp.Namespace,
			"metric_name": metricsInResp.MetricName,
			"dimensions":  setMultipleDimensions(metricsInResp.Dimensions),
			"unit":        metricsInResp.Unit,
			"datapoints":  setMultipleDatapoints(metricsInResp.Datapoints),
		}

		metricsList = append(metricsList, metrics)
	}
	return metricsList
}

func setMultipleDimensions(metricsDimensionsInResp []metricdata.MetricsDimension) []map[string]interface{} {
	var metricsDimensions []map[string]interface{}
	for _, dimensionInResp := range metricsDimensionsInResp {
		dimension := map[string]interface{}{
			"name":  dimensionInResp.Name,
			"value": dimensionInResp.Value,
		}

		metricsDimensions = append(metricsDimensions, dimension)
	}
	return metricsDimensions
}

func setMultipleDatapoints(datapointsInResp []metricdata.DatapointForBatchMetric) []map[string]interface{} {
	var datapoints []map[string]interface{}
	for _, datapointInResp := range datapointsInResp {
		datapoint := map[string]interface{}{
			"average":   datapointInResp.Average,
			"max":       datapointInResp.Max,
			"min":       datapointInResp.Min,
			"sum":       datapointInResp.Sum,
			"variance":  datapointInResp.Variance,
			"timestamp": datapointInResp.Timestamp,
		}

		datapoints = append(datapoints, datapoint)
	}
	return datapoints
}
