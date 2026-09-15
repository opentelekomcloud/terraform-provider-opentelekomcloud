package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v2/peerings"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpcPeeringConnectionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCPeeringV2Create,
		ReadContext:   resourceVPCPeeringV2Read,
		UpdateContext: resourceVPCPeeringV2Update,
		DeleteContext: resourceVPCPeeringV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: common.ValidateName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"vpc_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"peer_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func resourceVPCPeeringV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}

	createOpts := peerings.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		RequestVpcInfo: peerings.VpcInfo{
			VpcID: d.Get("vpc_id").(string),
		},
		AcceptVpcInfo: peerings.VpcInfo{
			VpcID:    d.Get("peer_vpc_id").(string),
			TenantID: d.Get("peer_tenant_id").(string),
		},
	}

	n, err := peerings.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Vpc Peering Connection: %s", err)
	}
	d.SetId(n.ID)

	log.Printf("[INFO] Vpc Peering Connection ID: %s", n.ID)

	log.Printf("[INFO] Waiting for OpenTelekomCloud Vpc Peering Connection(%s) to become available", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"PENDING_ACCEPTANCE", "ACTIVE"},
		Refresh:      waitForVpcPeeringActive(client, n.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for VPC peering connection to become available: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVPCPeeringV2Read(clientCtx, d, meta)
}

func resourceVPCPeeringV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}

	n, err := peerings.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "peering")
	}

	if err := setVpcPeeringV2Fields(d, n, config.GetRegion(d)); err != nil {
		return fmterr.Errorf("error setting VPC peering attributes: %w", err)
	}

	return nil
}

func resourceVPCPeeringV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}

	updateOpts := peerings.UpdateOpts{}
	if d.HasChange("name") {
		updateOpts.Name = pointerto.String(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateOpts.Description = pointerto.String(d.Get("description").(string))
	}

	_, err = peerings.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud Vpc Peering Connection: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVPCPeeringV2Read(clientCtx, d, meta)
}

func resourceVPCPeeringV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE"},
		Target:       []string{"DELETED"},
		Refresh:      waitForVpcPeeringDelete(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Vpc Peering Connection: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcPeeringActive(client *golangsdk.ServiceClient, peeringId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := peerings.Get(client, peeringId)
		if err != nil {
			return nil, "", err
		}

		if n.Status == "PENDING_ACCEPTANCE" || n.Status == "ACTIVE" {
			return n, n.Status, nil
		}
		if n.Status == "REJECTED" || n.Status == "EXPIRED" || n.Status == "DELETED" {
			return n, n.Status, fmt.Errorf(
				"VPC peering connection entered terminal status %s while being created", n.Status,
			)
		}

		return n, "CREATING", nil
	}
}

func waitForVpcPeeringDelete(client *golangsdk.ServiceClient, peeringId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := peerings.Get(client, peeringId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud vpc peering connection %s", peeringId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = peerings.Delete(client, peeringId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud vpc peering connection %s", peeringId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}
