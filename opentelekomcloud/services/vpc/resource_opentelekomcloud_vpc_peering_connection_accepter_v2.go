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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v2/peerings"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpcPeeringConnectionAccepterV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCPeeringAccepterV2Create,
		ReadContext:   resourceVpcPeeringAccepterRead,
		UpdateContext: resourceVPCPeeringAccepterUpdate,
		DeleteContext: resourceVPCPeeringAccepterDelete,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_peering_connection_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"accept": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_tenant_id": {
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

func resourceVPCPeeringAccepterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	peeringClient, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}

	id := d.Get("vpc_peering_connection_id").(string)

	n, err := peerings.Get(peeringClient, id)
	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud Vpc Peering Connection: %s", err)
	}

	expectedStatus := "REJECTED"
	if d.Get("accept").(bool) {
		expectedStatus = "ACTIVE"
	}
	if n.Status == expectedStatus {
		d.SetId(n.ID)
		clientCtx := common.CtxWithClient(ctx, peeringClient, keyClientV2)
		return resourceVpcPeeringAccepterRead(clientCtx, d, meta)
	}
	if n.Status != "PENDING_ACCEPTANCE" {
		return fmterr.Errorf("VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.")
	}
	d.SetId(n.ID)

	if d.Get("accept").(bool) {
		_, err = peerings.Accept(peeringClient, id)
		if err != nil {
			return fmterr.Errorf("unable to accept VPC Peering Connection: %w", err)
		}
	} else {
		_, err = peerings.Reject(peeringClient, id)
		if err != nil {
			return fmterr.Errorf("unable to reject VPC Peering Connection: %w", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{expectedStatus},
		Refresh:    waitForVpcPeeringConnStatus(peeringClient, n.ID, expectedStatus),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for peering connection to become %s: %w", expectedStatus, err)
	}

	log.Printf("[INFO] VPC Peering Connection status: %s", expectedStatus)

	clientCtx := common.CtxWithClient(ctx, peeringClient, keyClientV2)
	return resourceVpcPeeringAccepterRead(clientCtx, d, meta)
}

func resourceVpcPeeringAccepterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	peeringClient, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}

	n, err := peerings.Get(peeringClient, d.Id())
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Vpc Peering Connection: %s", err)
	}

	if err := d.Set("vpc_peering_connection_id", n.ID); err != nil {
		return fmterr.Errorf("error setting VPC peering connection ID: %w", err)
	}
	if err := d.Set("accept", n.Status == "ACTIVE"); err != nil {
		return fmterr.Errorf("error setting VPC peering acceptance state: %w", err)
	}
	if err := setVpcPeeringV2Fields(d, n, config.GetRegion(d)); err != nil {
		return fmterr.Errorf("error setting VPC peering attributes: %w", err)
	}

	return nil
}

func resourceVPCPeeringAccepterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("accept") {
		return fmterr.Errorf("VPC peering action not permitted: cannot accept or reject a request that is not pending acceptance")
	}

	return resourceVpcPeeringAccepterRead(ctx, d, meta)
}

func resourceVPCPeeringAccepterDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("[WARN] Will not delete VPC peering connection. Terraform will remove this resource from the state file, however resources may remain.")
	d.SetId("")
	return nil
}

func waitForVpcPeeringConnStatus(peeringClient *golangsdk.ServiceClient, peeringId, expectedStatus string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := peerings.Get(peeringClient, peeringId)
		if err != nil {
			return nil, "", err
		}

		if n.Status == expectedStatus {
			return n, expectedStatus, nil
		}
		if n.Status != "PENDING_ACCEPTANCE" {
			return n, n.Status, fmt.Errorf(
				"VPC peering connection entered status %s while waiting for %s", n.Status, expectedStatus,
			)
		}

		return n, "PENDING", nil
	}
}
