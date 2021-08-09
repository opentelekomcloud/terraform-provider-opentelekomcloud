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
			"vpc_peering_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVPCPeeringAccepterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	peeringClient, err := config.NetworkingV2Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Peering client: %s", err)
	}

	id := d.Get("vpc_peering_connection_id").(string)

	n, err := peerings.Get(peeringClient, id).Extract()
	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud Vpc Peering Connection: %s", err)
	}

	if n.Status != "PENDING_ACCEPTANCE" {
		return fmterr.Errorf("VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.")
	}

	var expectedStatus string

	if _, ok := d.GetOk("accept"); ok {
		expectedStatus = "ACTIVE"
		_, err := peerings.Accept(peeringClient, id).ExtractResult()

		if err != nil {
			return fmterr.Errorf("unable to accept VPC Peering Connection: %w", err)
		}
	} else {
		expectedStatus = "REJECTED"
		_, err := peerings.Reject(peeringClient, id).ExtractResult()
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
		return fmterr.Errorf("error waiting for perring connection to become %s: %w", expectedStatus, err)
	}

	d.SetId(n.ID)
	log.Printf("[INFO] VPC Peering Connection status: %s", expectedStatus)

	return resourceVpcPeeringAccepterRead(ctx, d, meta)
}

func resourceVpcPeeringAccepterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	peeringclient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud peering client: %s", err)
	}

	n, err := peerings.Get(peeringclient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Vpc Peering Connection: %s", err)
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
		return diag.FromErr(err)
	}

	return nil
}

func resourceVPCPeeringAccepterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("accept") {
		return fmterr.Errorf("VPC peering action not permitted: Can not accept/reject peering request not in pending_acceptance state.'")
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
		n, err := peerings.Get(peeringClient, peeringId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == expectedStatus {
			return n, expectedStatus, nil
		}

		return n, "PENDING", nil
	}
}
