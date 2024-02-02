package v3

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/members"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLBMemberV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBMemberV3Create,
		ReadContext:   resourceLBMemberV3Read,
		UpdateContext: resourceLBMemberV3Update,
		DeleteContext: resourceLBMemberV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("pool_id", "member_id"),
		},

		Schema: map[string]*schema.Schema{
			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"protocol_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsPortNumber,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"member_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLBMemberV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	opts := members.CreateOpts{
		Address:      d.Get("address").(string),
		ProtocolPort: d.Get("protocol_port").(int),
		Name:         d.Get("name").(string),
		ProjectID:    d.Get("project_id").(string),
		SubnetID:     d.Get("subnet_id").(string),
	}

	weight := common.CheckNull("weight", d)
	if !weight {
		opts.Weight = pointerto.Int(d.Get("weight").(int))
	}

	member, err := members.Create(client, poolID(d), opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating LB pool member v3: %w", err)
	}
	_ = d.Set("member_id", member.ID) // this can't ever return an error

	if err := common.SetComplexID(d, "pool_id", "member_id"); err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBMemberV3Read(clientCtx, d, meta)
}

func resourceLBMemberV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	member, err := members.Get(client, poolID(d), memberID(d)).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error reading LB pool member v3")
	}
	mErr := multierror.Append(
		d.Set("address", member.Address),
		d.Set("name", member.Name),
		d.Set("operating_status", member.OperatingStatus),
		d.Set("project_id", member.ProjectID),
		d.Set("subnet_id", member.SubnetID),
		d.Set("protocol_port", member.ProtocolPort),
		d.Set("weight", member.Weight),
		d.Set("ip_version", member.IpVersion),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting LB pool member fields: %w", err)
	}

	return nil
}

func resourceLBMemberV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	opts := members.UpdateOpts{}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
	}
	if d.HasChange("weight") {
		weight := d.Get("weight").(int)
		opts.Weight = &weight
	}

	_, err = members.Update(client, poolID(d), memberID(d), opts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating LB pool member: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBMemberV3Read(clientCtx, d, meta)
}

func resourceLBMemberV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	err = members.Delete(client, poolID(d), memberID(d)).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting member: %w", err)
	}

	return nil
}

func poolID(d *schema.ResourceData) string {
	return d.Get("pool_id").(string)
}

func memberID(d *schema.ResourceData) string {
	return d.Get("member_id").(string)
}
