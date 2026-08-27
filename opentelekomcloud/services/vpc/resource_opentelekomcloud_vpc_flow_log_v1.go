package vpc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/flow_logs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpcFlowLogV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcFlowLogV1Create,
		ReadContext:   resourceVpcFlowLogV1Read,
		UpdateContext: resourceVpcFlowLogV1Update,
		DeleteContext: resourceVpcFlowLogV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: common.ValidateName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"port", "vpc", "network",
				}, true),
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"traffic_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"all", "accept", "reject",
				}, true),
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"index_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
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
		},
	}
}

func resourceVpcFlowLogV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud vpc client: %s", err)
	}

	createOpts := flow_logs.CreateOpts{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		ResourceType: d.Get("resource_type").(string),
		ResourceID:   d.Get("resource_id").(string),
		TrafficType:  d.Get("traffic_type").(string),
		LogGroupID:   d.Get("log_group_id").(string),
		LogTopicID:   d.Get("log_topic_id").(string),
	}
	createOpts.IndexEnabled = pointerto.Bool(d.Get("index_enabled").(bool))

	log.Printf("[DEBUG] Create VPC Flow Log Options: %#v", createOpts)
	fl, err := flow_logs.Create(vpcClient, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC flow log: %s", err)
	}

	d.SetId(fl.ID)
	_, err = flow_logs.Update(vpcClient, d.Id(), flow_logs.UpdateOpts{
		AdminState: pointerto.Bool(d.Get("enabled").(bool)),
	})
	if err != nil {
		return fmterr.Errorf("error setting OpenTelekomCloud VPC flow log admin state: %s", err)
	}
	return resourceVpcFlowLogV1Read(ctx, d, config)
}

func resourceVpcFlowLogV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud vpc client: %s", err)
	}

	fl, err := flow_logs.Get(vpcClient, d.Id())
	if err != nil {
		// ignore ErrDefault404
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud flowlog: %s", err)
	}

	mErr := multierror.Append(
		d.Set("name", fl.Name),
		d.Set("description", fl.Description),
		d.Set("resource_type", fl.ResourceType),
		d.Set("resource_id", fl.ResourceID),
		d.Set("traffic_type", fl.TrafficType),
		d.Set("log_group_id", fl.LogGroupID),
		d.Set("log_topic_id", fl.LogTopicID),
		d.Set("enabled", fl.AdminState),
		d.Set("status", fl.Status),
		d.Set("tenant_id", fl.TenantID),
		d.Set("created_at", fl.CreatedAt),
		d.Set("updated_at", fl.UpdatedAt),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpcFlowLogV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud vpc client: %s", err)
	}

	var updateOpts flow_logs.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = pointerto.String(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateOpts.Description = pointerto.String(d.Get("description").(string))
	}
	if d.HasChange("enabled") {
		updateOpts.AdminState = pointerto.Bool(d.Get("enabled").(bool))
	}

	_, err = flow_logs.Update(vpcClient, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud VPC flow log: %s", err)
	}

	return resourceVpcFlowLogV1Read(ctx, d, meta)
}

func resourceVpcFlowLogV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud vpc client: %s", err)
	}

	err = flow_logs.Delete(vpcClient, d.Id())
	if err != nil {
		// ignore ErrDefault404
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[INFO] Successfully deleted OpenTelekomCloud vpc flow log %s", d.Id())
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
