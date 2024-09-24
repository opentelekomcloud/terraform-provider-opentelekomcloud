package sdrs

import (
	"context"
	"fmt"
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
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceSdrsProtectiongroupV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	sdrsClient, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
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

	n, err := protectiongroups.Create(sdrsClient, createOpts)
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
		if d.Get("enable").(bool) {
			errEnable := enableProtectionsGroup(d, sdrsClient)
			if errEnable != nil {
				return fmterr.Errorf("error while enabling Protection Group: %s", errEnable)
			}
		}
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
	n, err := protectiongroups.Get(sdrsClient, d.Id())

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
		d.Set("source_availability_zone", n.SourceAvailabilityZone),
		d.Set("target_availability_zone", n.TargetAvailabilityZone),
		d.Set("domain_id", n.DomainID),
		d.Set("source_vpc_id", n.SourceVPCID),
		d.Set("dr_type", n.DRType),
		d.Set("enable", n.Status == "protected"),
		d.Set("created_at", n.CreatedAt),
		d.Set("updated_at", n.UpdatedAt),
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

	if d.HasChange("name") {
		var updateOpts protectiongroups.UpdateOpts
		updateOpts.Name = d.Get("name").(string)

		log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

		_, err = protectiongroups.Update(sdrsClient, d.Id(), updateOpts)
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud SDRS Protectiongroup: %s", err)
		}
	}

	if d.HasChange("enable") {
		var errEnable error
		if d.Get("enable").(bool) {
			errEnable = enableProtectionsGroup(d, sdrsClient)
		} else {
			errEnable = disableProtectionsGroup(d, sdrsClient)
		}
		if errEnable != nil {
			return diag.FromErr(errEnable)
		}
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

	n, err := protectiongroups.Delete(sdrsClient, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud SDRS Protectiongroup: %s", err)
	}

	if err := protectiongroups.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutDelete)/time.Second), n.JobID); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func enableProtectionsGroup(d *schema.ResourceData, sdrsClient *golangsdk.ServiceClient) error {
	enableResponse, err := protectiongroups.Enable(sdrsClient, d.Id())
	if err != nil {
		return fmt.Errorf("error enabling OpenTelekomcomCloud SDRS Protectiongroup: %s", err)
	}
	return protectiongroups.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), enableResponse.JobID)
}

func disableProtectionsGroup(d *schema.ResourceData, sdrsClient *golangsdk.ServiceClient) error {
	disableResponse, err := protectiongroups.Disable(sdrsClient, d.Id())
	if err != nil {
		return fmt.Errorf("error enabling OpenTelekomcomCloud SDRS Protectiongroup: %s", err)
	}
	return protectiongroups.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), disableResponse.JobID)
}
