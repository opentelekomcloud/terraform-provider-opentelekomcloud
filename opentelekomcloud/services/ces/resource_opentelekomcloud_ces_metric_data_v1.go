package ces

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metricdata"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCesMetricDataV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCesMetricDataCreate,
		ReadContext:   resourceCesMetricDataRead,
		DeleteContext: resourceCesMetricDataDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"metric": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(3, 32),
								validation.StringMatch(
									regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*\.[a-zA-Z][a-zA-Z0-9_]*$`),
									"Must be of type service.item. service and item each must only have lowercase/uppercase letters, digits, and underscores (_) and must start with a letter.",
								),
								validation.StringDoesNotMatch(
									regexp.MustCompile(`^(SYS\.|AGT\.|SRE\.).*|^SERVICE\.BMS$`),
									"service in namespace cannot be SYS, AGT, or SRE, and namespace cannot be SERVICE.BMS",
								),
							),
						},
						"dimensions": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 3,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 32),
											validation.StringMatch(
												regexp.MustCompile(`^[a-zA-Z].+`),
												"Must start with a letter."),
											validation.StringMatch(
												regexp.MustCompile(`^[\w-]+$`),
												"Only lowercase/uppercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed.",
											),
										),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 256),
											validation.StringMatch(
												regexp.MustCompile(`^[a-zA-Z0-9].+`),
												"Must start with a letter."),
											validation.StringMatch(
												regexp.MustCompile(`^[\w-]+$`),
												"Only lowercase/uppercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed.",
											),
										),
									},
								},
							},
						},
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 604800),
			},
			"collect_time": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeFloat,
				Required: true,
				ForceNew: true,
			},
			"unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 32),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"int", "float",
				}, false),
			},
		},
	}
}

func getMetricDataMetric(d *schema.ResourceData) metricdata.MetricInfo {
	metricListRaw := d.Get("metric").([]interface{})
	metricElement := metricListRaw[0].(map[string]interface{})

	metricDimensions := metricElement["dimensions"].([]interface{})
	dimensionOpts := make([]metricdata.MetricsDimension, len(metricDimensions))
	for i, dimensionElement := range metricDimensions {
		dimension := dimensionElement.(map[string]interface{})
		dimensionOpts[i] = metricdata.MetricsDimension{
			Name:  dimension["name"].(string),
			Value: dimension["value"].(string),
		}
	}
	return metricdata.MetricInfo{
		Namespace:  metricElement["namespace"].(string),
		MetricName: metricElement["metric_name"].(string),
		Dimensions: dimensionOpts,
	}
}

func resourceCesMetricDataCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	createOpts := metricdata.MetricDataItem{
		Metric:      getMetricDataMetric(d),
		Ttl:         d.Get("ttl").(int),
		CollectTime: InterfaceToInt64(d.Get("collect_time")),
		Value:       d.Get("value").(float64),
		Unit:        d.Get("unit").(string),
		Type:        d.Get("type").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	err = metricdata.CreateMetricData(client, []metricdata.MetricDataItem{
		createOpts,
	})
	if err != nil {
		return fmterr.Errorf("error creating metric_data: %w", err)
	}
	log.Printf("[DEBUG] Created custom metric data: %#v", createOpts)
	d.SetId(createOpts.Metric.MetricName)

	clientCtx := common.CtxWithClient(ctx, client, cesClientV1)
	return resourceCesMetricDataRead(clientCtx, d, meta)
}

func resourceCesMetricDataRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceCesMetricDataDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The action resource is only removed from the state the task remains in the cloud.",
		},
	}
}
