package cc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	gcb "github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/global_connection_bandwidth"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCcGlobalConnectionBandwidthV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCcGlobalConnectionBandwidthV3Create,
		ReadContext:   resourceCcGlobalConnectionBandwidthV3Read,
		UpdateContext: resourceCcGlobalConnectionBandwidthV3Update,
		DeleteContext: resourceCcGlobalConnectionBandwidthV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bordercross": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"sla_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"local_area": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"remote_area": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"spec_code_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"binding_service": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_site_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_site_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frozen": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_share": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"eps_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"directional_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_site_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_site_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCcGlobalConnectionBandwidthV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	createOpts := gcb.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Bordercross:         pointerto.Bool(d.Get("bordercross").(bool)),
		Type:                d.Get("type").(string),
		EnterpriseProjectId: d.Get("enterprise_project_id").(string),
		ChargeMode:          d.Get("charge_mode").(string),
		Size:                d.Get("size").(int),
		SlaLevel:            d.Get("sla_level").(string),
		LocalArea:           d.Get("local_area").(string),
		RemoteArea:          d.Get("remote_area").(string),
		SpecCodeId:          d.Get("spec_code_id").(string),
	}

	created, err := gcb.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud CC global connection bandwidth: %s", err)
	}

	d.SetId(created.ID)

	clientCtx := common.CtxWithClient(ctx, client, ccClientV3)
	return resourceCcGlobalConnectionBandwidthV3Read(clientCtx, d, meta)
}

func resourceCcGlobalConnectionBandwidthV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	bandwidth, err := gcb.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud CC global connection bandwidth")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", bandwidth.Name),
		d.Set("description", bandwidth.Description),
		d.Set("bordercross", bandwidth.Bordercross),
		d.Set("type", bandwidth.Type),
		d.Set("charge_mode", bandwidth.ChargeMode),
		d.Set("size", bandwidth.Size),
		d.Set("enterprise_project_id", bandwidth.EnterpriseProjectId),
		d.Set("sla_level", bandwidth.SlaLevel),
		d.Set("local_area", bandwidth.LocalArea),
		d.Set("remote_area", bandwidth.RemoteArea),
		d.Set("spec_code_id", bandwidth.SpecCodeId),
		d.Set("binding_service", bandwidth.BindingService),
		d.Set("domain_id", bandwidth.DomainId),
		d.Set("local_site_code", bandwidth.LocalSiteCode),
		d.Set("remote_site_code", bandwidth.RemoteSiteCode),
		d.Set("admin_state", bandwidth.AdminState),
		d.Set("frozen", bandwidth.Frozen),
		d.Set("enable_share", bandwidth.EnableShare),
		d.Set("eps_id", bandwidth.EpsId),
		d.Set("created_at", bandwidth.CreatedAt),
		d.Set("updated_at", bandwidth.UpdatedAt),
		d.Set("instances", flattenGcbInstances(bandwidth.Instances)),
		d.Set("directional_connections", flattenGcbDirectionalConnections(bandwidth.DirectionalConnections)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGcbInstances(instances []gcb.AssociatedInstance) []map[string]interface{} {
	if len(instances) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(instances))
	for i, instance := range instances {
		result[i] = map[string]interface{}{
			"id":         instance.ID,
			"type":       instance.Type,
			"region_id":  instance.RegionId,
			"project_id": instance.ProjectId,
		}
	}
	return result
}

func flattenGcbDirectionalConnections(connections []gcb.DirectionalConnection) []map[string]interface{} {
	if len(connections) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(connections))
	for i, connection := range connections {
		result[i] = map[string]interface{}{
			"id":               connection.ID,
			"name":             connection.Name,
			"local_site_code":  connection.LocalSiteCode,
			"remote_site_code": connection.RemoteSiteCode,
		}
	}
	return result
}

func resourceCcGlobalConnectionBandwidthV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	if d.HasChanges("name", "description", "size", "charge_mode", "sla_level", "binding_service", "spec_code_id") {
		updateOpts := gcb.UpdateOpts{
			ID:             d.Id(),
			Name:           d.Get("name").(string),
			Description:    pointerto.String(d.Get("description").(string)),
			Size:           d.Get("size").(int),
			ChargeMode:     d.Get("charge_mode").(string),
			SlaLevel:       d.Get("sla_level").(string),
			BindingService: d.Get("binding_service").(string),
			SpecCodeId:     d.Get("spec_code_id").(string),
		}
		if _, err = gcb.Update(client, updateOpts); err != nil {
			return diag.Errorf("error updating OpenTelekomCloud CC global connection bandwidth (%s): %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, ccClientV3)
	return resourceCcGlobalConnectionBandwidthV3Read(clientCtx, d, meta)
}

func resourceCcGlobalConnectionBandwidthV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	if err = gcb.Delete(client, d.Id()); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud CC global connection bandwidth")
	}

	return nil
}
