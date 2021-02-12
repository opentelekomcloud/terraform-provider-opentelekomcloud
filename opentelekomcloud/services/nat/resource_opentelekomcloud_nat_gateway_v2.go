package nat

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/natgateways"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceNatGatewayV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNatGatewayV2Create,
		Read:   resourceNatGatewayV2Read,
		Update: resourceNatGatewayV2Update,
		Delete: resourceNatGatewayV2Delete,

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

func resourceNatGatewayV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	NatV2Client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud nat client: %s", err)
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
	natGateway, err := natgateways.Create(NatV2Client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creatting Nat Gateway: %s", err)
	}

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud Nat Gateway (%s) to become available.", natGateway.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForNatGatewayActive(NatV2Client, natGateway.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Nat Gateway: %s", err)
	}

	d.SetId(natGateway.ID)

	return resourceNatGatewayV2Read(d, meta)
}

func resourceNatGatewayV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	NatV2Client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud nat client: %s", err)
	}

	natGateway, err := natgateways.Get(NatV2Client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Nat Gateway")
	}

	d.Set("name", natGateway.Name)
	d.Set("description", natGateway.Description)
	d.Set("spec", natGateway.Spec)
	d.Set("router_id", natGateway.RouterID)
	d.Set("internal_network_id", natGateway.InternalNetworkID)
	d.Set("tenant_id", natGateway.TenantID)

	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceNatGatewayV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	NatV2Client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud nat client: %s", err)
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

	_, err = natgateways.Update(NatV2Client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating Nat Gateway: %s", err)
	}

	return resourceNatGatewayV2Read(d, meta)
}

func resourceNatGatewayV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	NatV2Client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud nat client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNatGatewayDelete(NatV2Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud Nat Gateway: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForNatGatewayActive(NatV2Client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := natgateways.Get(NatV2Client, nId).Extract()
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

func waitForNatGatewayDelete(NatV2Client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud Nat Gateway %s.\n", nId)

		n, err := natgateways.Get(NatV2Client, nId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Nat gateway %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = natgateways.Delete(NatV2Client, nId).ExtractErr()
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
