package sdrs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/protectiongroups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSdrsProtectiongroupV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSdrsProtectiongroupV1Create,
		ReadContext:   resourceSdrsProtectiongroupV1Read,
		UpdateContext: resourceSdrsProtectiongroupV1Update,
		DeleteContext: resourceSdrsProtectiongroupV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dr_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSdrsProtectiongroupV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	sdrsClient, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud SDRS Client: %s", err)
	}

	createOpts := protectiongroups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		SourceAZ:    d.Get("source_availability_zone").(string),
		TargetAZ:    d.Get("target_availability_zone").(string),
		DomainID:    d.Get("domain_id").(string),
		SourceVpcID: d.Get("source_vpc_id").(string),
		DrType:      d.Get("dr_type").(string),
	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	n, err := protectiongroups.Create(sdrsClient, createOpts).ExtractJobResponse()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud SDRS Protectiongroup: %s", err)
	}

	if err := protectiongroups.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return diag.FromErr(err)
	}

	entity, err := protectiongroups.GetJobEntity(sdrsClient, n.JobID, "server_group_id")
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := entity.(string); ok {
		d.SetId(id)
		clientCtx := common.CtxWithClient(ctx, sdrsClient, sdrsClientV1)
		return resourceSdrsProtectiongroupV1Read(clientCtx, d, meta)
	}

	return fmterr.Errorf("unexpected conversion error in resourceSdrsProtectiongroupV1Create.")
}

func resourceSdrsProtectiongroupV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	sdrsClient, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}
	n, err := protectiongroups.Get(sdrsClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud SDRS Protectiongroup: %s", err)
	}
	mErr := multierror.Append(
		d.Set("name", n.Name),
		d.Set("description", n.Description),
		d.Set("source_availability_zone", n.SourceAZ),
		d.Set("target_availability_zone", n.TargetAZ),
		d.Set("domain_id", n.DomainID),
		d.Set("source_vpc_id", n.SourceVpcID),
		d.Set("dr_type", n.DrType),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSdrsProtectiongroupV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	sdrsClient, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}
	var updateOpts protectiongroups.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	_, err = protectiongroups.Update(sdrsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud SDRS Protectiongroup: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, sdrsClient, sdrsClientV1)
	return resourceSdrsProtectiongroupV1Read(clientCtx, d, meta)
}

func resourceSdrsProtectiongroupV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	sdrsClient, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	n, err := protectiongroups.Delete(sdrsClient, d.Id()).ExtractJobResponse()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud SDRS Protectiongroup: %s", err)
	}

	if err := protectiongroups.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutDelete)/time.Second), n.JobID); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
