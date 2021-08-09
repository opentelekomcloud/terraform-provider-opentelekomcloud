package nat

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/natgateways"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNatGatewayV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNatGatewayV2Create,
		ReadContext:   resourceNatGatewayV2Read,
		UpdateContext: resourceNatGatewayV2Update,
		DeleteContext: resourceNatGatewayV2Delete,

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
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"spec": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: resourceNatGatewayV2ValidateSpec,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internal_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNatGatewayV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud nat client: %s", err)
	}

	createOpts := &natgateways.CreateOpts{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Spec:              d.Get("spec").(string),
		TenantID:          d.Get("tenant_id").(string),
		RouterID:          d.Get("router_id").(string),
		InternalNetworkID: d.Get("internal_network_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	natGateway, err := natgateways.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creatting Nat Gateway: %s", err)
	}

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud Nat Gateway (%s) to become available.", natGateway.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForNatGatewayActive(client, natGateway.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Nat Gateway: %s", err)
	}

	d.SetId(natGateway.ID)

	return resourceNatGatewayV2Read(ctx, d, meta)
}

func resourceNatGatewayV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud nat client: %s", err)
	}

	natGateway, err := natgateways.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "Nat Gateway"))
	}

	mErr := multierror.Append(
		d.Set("name", natGateway.Name),
		d.Set("description", natGateway.Description),
		d.Set("spec", natGateway.Spec),
		d.Set("router_id", natGateway.RouterID),
		d.Set("internal_network_id", natGateway.InternalNetworkID),
		d.Set("tenant_id", natGateway.TenantID),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNatGatewayV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud nat client: %s", err)
	}

	var updateOpts natgateways.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("spec") {
		updateOpts.Spec = d.Get("spec").(string)
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = natgateways.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating Nat Gateway: %s", err)
	}

	return resourceNatGatewayV2Read(ctx, d, meta)
}

func resourceNatGatewayV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	NatV2Client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud nat client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNatGatewayDelete(NatV2Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Nat Gateway: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForNatGatewayActive(client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := natgateways.Get(client, nId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Nat Gateway: %+v", n)
		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, "", nil
	}
}

func waitForNatGatewayDelete(client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud Nat Gateway %s.\n", nId)

		n, err := natgateways.Get(client, nId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Nat gateway %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = natgateways.Delete(client, nId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Nat Gateway %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Nat Gateway %s still active.\n", nId)
		return n, "ACTIVE", nil
	}
}

var Specs = [4]string{"1", "2", "3", "4"}

func resourceNatGatewayV2ValidateSpec(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range Specs {
		if value == Specs[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, Specs))
	return
}
