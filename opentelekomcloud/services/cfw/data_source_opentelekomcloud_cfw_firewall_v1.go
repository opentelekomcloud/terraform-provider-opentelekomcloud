package cfw

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	cfwmanagementv1 "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/management"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCfwFirewallV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCFWFirewallV1Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
				ValidateFunc: validation.StringInSlice([]string{
					"0",
				}, true),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_code": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"eip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"log_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_eip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_log_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_vpc_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ha_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"engine_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protect_objects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"object_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_old_firewall_instance": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_available_obs": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_support_threat_tags": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_ipv6": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"feature_toggle": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_size_measure_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_url_filtering": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceCFWFirewallV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	d.SetId(d.Get("id").(string))
	serviceType, err := strconv.Atoi(d.Get("service_type").(string))
	if err != nil {
		return fmterr.Errorf("error converting service type to integer: %s", err)
	}
	firewallInstance, err := cfwmanagementv1.Get(client, d.Id(), serviceType)
	if err != nil {
		return fmterr.Errorf("error fetching CFW Firewall instance: %w", err)
	}

	log.Printf("[DEBUG] Retrieved instance %s: %#v", d.Id(), firewallInstance)

	serviceTypeStr := fmt.Sprintf("%d", firewallInstance.ServiceType)
	mErr := multierror.Append(nil,
		d.Set("id", firewallInstance.FwInstanceID),
		d.Set("name", firewallInstance.FwInstanceName),
		d.Set("enterprise_project_id", firewallInstance.EnterpriseProjectID),
		d.Set("ha_type", firewallInstance.HAType),
		d.Set("charge_mode", firewallInstance.ChargeMode),
		d.Set("service_type", serviceTypeStr),
		d.Set("engine_type", firewallInstance.EngineType),
		d.Set("status", firewallInstance.Status),
		d.Set("is_old_firewall_instance", firewallInstance.IsOldFirewallInstance),
		d.Set("is_available_obs", firewallInstance.IsAvailableObs),
		d.Set("is_support_threat_tags", firewallInstance.IsSupportThreatTags),
		d.Set("support_ipv6", firewallInstance.SupportIpv6),
		d.Set("feature_toggle", firewallInstance.FeatureToggle),
		d.Set("resource_id", firewallInstance.ResourceID),
		d.Set("support_url_filtering", firewallInstance.SupportUrlFiltering),
	)

	flavorInResp := firewallInstance.Flavor
	flavor := map[string]interface{}{
		"version_code":        flavorInResp.Version,
		"eip_count":           flavorInResp.EipCount,
		"vpc_count":           flavorInResp.VpcCount,
		"bandwidth":           flavorInResp.Bandwidth,
		"log_storage":         flavorInResp.LogStorage,
		"default_bandwidth":   flavorInResp.DefaultBandwidth,
		"default_eip_count":   flavorInResp.DefaultEipCount,
		"default_log_storage": flavorInResp.DefaultLogStorage,
		"default_vpc_count":   flavorInResp.DefaultVpcCount,
	}
	var flavors []map[string]interface{}
	flavors = append(flavors, flavor)

	mErr = multierror.Append(
		mErr,
		d.Set("flavor", flavors),
	)

	var protectedObjects []map[string]interface{}
	for _, protectedObjectInResp := range firewallInstance.ProtectObjects {
		protectedObject := make(map[string]interface{})
		protectedObject["object_id"] = protectedObjectInResp.ObjectID
		protectedObject["object_name"] = protectedObjectInResp.ObjectName
		protectedObject["type"] = protectedObjectInResp.Type
		protectedObjects = append(protectedObjects, protectedObject)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("protect_objects", protectedObjects),
	)

	var resourcesList []map[string]interface{}
	for _, resourceInResp := range firewallInstance.Resources {
		resource := map[string]interface{}{
			"resource_id":              resourceInResp.ResourceID,
			"cloud_service_type":       resourceInResp.CloudServiceType,
			"resource_type":            resourceInResp.ResourceType,
			"resource_spec_code":       resourceInResp.ResourceSpecCode,
			"resource_size":            resourceInResp.ResourceSize,
			"resource_size_measure_id": resourceInResp.ResourceSizeMeasureID,
		}
		resourcesList = append(resourcesList, resource)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("resources", resourcesList),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
