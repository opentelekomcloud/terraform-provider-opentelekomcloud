package er

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/vpc"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceErVpcAttachmentV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcAttachmentV3Create,
		UpdateContext: resourceVpcAttachmentV3Update,
		ReadContext:   resourceVpcAttachmentV3Read,
		DeleteContext: resourceVpcAttachmentV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceVpcAttachmentV3ImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 255),
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
				),
			},
			"auto_create_vpc_routes": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
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
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcAttachmentV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	instanceId := d.Get("instance_id").(string)

	opts := vpc.CreateOpts{
		RouterID:            instanceId,
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		AutoCreateVpcRoutes: d.Get("auto_create_vpc_routes").(bool),
	}

	resp, err := vpc.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating VPC attachment: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      vpcAttachmentStatusRefreshFunc(client, instanceId, d.Id(), []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, erClientV3)
	return resourceVpcAttachmentV3Read(clientCtx, d, meta)
}

func resourceVpcAttachmentV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		config       = meta.(*cfg.Config)
		instanceId   = d.Get("instance_id").(string)
		attachmentId = d.Id()
		region       = config.GetRegion(d)
	)

	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(region)
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	resp, err := vpc.Get(client, instanceId, attachmentId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ER VPC attachment")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_id", resp.VpcAttachment.VpcId),
		d.Set("subnet_id", resp.VpcAttachment.SubnetId),
		d.Set("name", resp.VpcAttachment.Name),
		d.Set("description", resp.VpcAttachment.Description),
		d.Set("auto_create_vpc_routes", resp.VpcAttachment.AutoCreateVpcRoutes),
		d.Set("status", resp.VpcAttachment.State),
		d.Set("created_at", resp.VpcAttachment.CreatedAt),
		d.Set("updated_at", resp.VpcAttachment.UpdatedAt),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving VPC attachment (%s) fields: %s", d.Id(), mErr)
	}
	return nil
}

func vpcAttachmentStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, attachmentId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := vpc.Get(client, instanceId, attachmentId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}

			return nil, "", err
		}
		log.Printf("[DEBUG] The details of the VPC attachment (%s) is: %#v", attachmentId, resp)

		if common.StrSliceContains([]string{"failed"}, resp.VpcAttachment.State) {
			return resp, "", fmt.Errorf("unexpected status '%s'", resp.VpcAttachment.State)
		}
		if common.StrSliceContains(targets, resp.VpcAttachment.State) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourceVpcAttachmentV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		config       = meta.(*cfg.Config)
		instanceId   = d.Get("instance_id").(string)
		attachmentId = d.Id()
	)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	if d.HasChanges("name", "description") {
		opts := vpc.UpdateOpts{
			RouterID:        instanceId,
			VpcAttachmentID: attachmentId,
			Name:            d.Get("name").(string),
			Description:     pointerto.String(d.Get("description").(string)),
		}

		_, err = vpc.Update(client, opts)
		if err != nil {
			return fmterr.Errorf("error getting VPC attachment (%s) details: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"COMPLETED"},
			Refresh:      vpcAttachmentStatusRefreshFunc(client, instanceId, attachmentId, []string{"available"}),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        5 * time.Second,
			PollInterval: 10 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVpcAttachmentV3Read(ctx, d, meta)
}

func resourceVpcAttachmentV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	instanceId := d.Get("instance_id").(string)
	attachmentId := d.Id()

	err = vpc.Delete(client, instanceId, attachmentId)
	if err != nil {
		return diag.Errorf("error deleting VPC attachment (%s) form the ER instance: %s", attachmentId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      vpcAttachmentStatusRefreshFunc(client, instanceId, attachmentId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceVpcAttachmentV3ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<instance_id>/<attachment_id>', but '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
