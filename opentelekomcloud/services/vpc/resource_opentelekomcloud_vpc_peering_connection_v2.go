package vpc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/peerings"

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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceVPCPeeringV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	requestvpcinfo := peerings.VpcInfo{
		VpcId: d.Get("vpc_id").(string),
	}

	acceptvpcinfo := peerings.VpcInfo{
		VpcId:    d.Get("peer_vpc_id").(string),
		TenantId: d.Get("peer_tenant_id").(string),
	}

	createOpts := peerings.CreateOpts{
		Name:           d.Get("name").(string),
		RequestVpcInfo: requestvpcinfo,
		AcceptVpcInfo:  acceptvpcinfo,
	}

	n, err := peerings.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Vpc Peering Connection: %s", err)
	}

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
		return fmterr.Errorf("error waiting for VIP to become active: %w", err)
	}

	d.SetId(n.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVPCPeeringV2Read(clientCtx, d, meta)
}

func resourceVPCPeeringV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	n, err := peerings.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "peering")
	}

	mErr := multierror.Append(
		d.Set("name", n.Name),
		d.Set("status", n.Status),
		d.Set("vpc_id", n.RequestVpcInfo.VpcId),
		d.Set("peer_vpc_id", n.AcceptVpcInfo.VpcId),
		d.Set("peer_tenant_id", n.AcceptVpcInfo.TenantId),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting VPC peering attributes: %w", err)
	}

	return nil
}

func resourceVPCPeeringV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	updateOpts := peerings.UpdateOpts{
		Name: d.Get("name").(string),
	}

	_, err = peerings.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud Vpc Peering Connection: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVPCPeeringV2Read(clientCtx, d, meta)
}

func resourceVPCPeeringV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
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
		n, err := peerings.Get(client, peeringId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "PENDING_ACCEPTANCE" || n.Status == "ACTIVE" {
			return n, n.Status, nil
		}

		return n, "CREATING", nil
	}
}

func waitForVpcPeeringDelete(client *golangsdk.ServiceClient, peeringId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := peerings.Get(client, peeringId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud vpc peering connection %s", peeringId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = peerings.Delete(client, peeringId).ExtractErr()
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
