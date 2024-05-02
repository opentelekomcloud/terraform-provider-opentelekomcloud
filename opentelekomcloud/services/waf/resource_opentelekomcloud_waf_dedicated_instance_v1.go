package waf

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/instances"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
)

func ResourceWafDedicatedInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedInstanceV1Create,
		ReadContext:   resourceWafDedicatedInstanceV1Read,
		UpdateContext: resourceWafDedicatedInstanceV1Update,
		DeleteContext: resourceWafDedicatedInstanceV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "x86",
				ForceNew: true,
			},
			"flavor": {
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
			"security_group": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"res_tenant": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"billing_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"upgradable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceWafDedicatedInstanceV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	sg := d.Get("security_group").([]interface{})
	groups := make([]string, 0, len(sg))
	for _, v := range sg {
		groups = append(groups, v.(string))
	}
	log.Printf("[DEBUG] The security_group parameters: %+v.", groups)

	opts := instances.CreateOpts{
		Region:           config.GetRegion(d),
		ChargeMode:       payPerUseMode,
		AvailabilityZone: d.Get("availability_zone").(string),
		Architecture:     d.Get("architecture").(string),
		InstanceName:     d.Get("name").(string),
		Specification:    d.Get("specification").(string),
		Flavor:           d.Get("flavor").(string),
		VpcId:            d.Get("vpc_id").(string),
		SubnetId:         d.Get("subnet_id").(string),
		SecurityGroupsId: groups,
		Count:            defaultCount,
		ResTenant:        pointerto.Bool(d.Get("res_tenant").(bool)),
	}

	r, err := instances.Create(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.Instances[0].Id)

	log.Printf("[DEBUG] Waiting for opentelekomcloud WAF dedicated instance (%s) to become available", r.Instances[0].Id)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Creating"},
		Target:       []string{"Created"},
		Refresh:      waitForInstanceCreated(client, r.Instances[0].Id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF dedicated instance: %s", err)
	}

	if err == nil {
		err = updateInstanceName(client, r.Instances[0].Id, d.Get("name").(string))
	}
	if err != nil {
		log.Printf("[DEBUG] Error while waiting to create Waf dedicated instance. %s : %#v", d.Id(), err)
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedInstanceV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedInstanceV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	instance, err := instances.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud WAF dedicated instance.")
	}
	upgradable := false
	if instance.Upgradable == 1 {
		upgradable = true
	}
	mErr := multierror.Append(nil,
		d.Set("region", instance.Region),
		d.Set("name", instance.Name),
		d.Set("availability_zone", instance.AvailabilityZone),
		d.Set("architecture", instance.Architecture),
		d.Set("flavor", instance.Flavor),
		d.Set("vpc_id", instance.VpcID),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group", instance.SecurityGroups),
		d.Set("server_id", instance.ServerId),
		d.Set("service_ip", instance.ServiceIp),
		d.Set("status", instance.Status),
		d.Set("access_status", instance.AccessStatus),
		d.Set("upgradable", upgradable),
		d.Set("specification", instance.ResourceSpecification),
		d.Set("created_at", instance.CreatedAt),
		d.Set("billing_status", instance.BillingStatus),
	)

	if mErr.ErrorOrNil() != nil {
		return fmterr.Errorf("error setting opentelekomcloud WAF dedicated instance fields: %w", err)
	}
	return nil
}

func resourceWafDedicatedInstanceV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	if d.HasChanges("name") {
		err = updateInstanceName(client, d.Id(), d.Get("name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedInstanceV1Read(clientCtx, d, meta)
}

func waitForInstanceCreated(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := instances.Get(client, id)
		if err != nil {
			return nil, "Error", err
		}
		switch r.Status {
		case runStatusCreating:
			return r, "Creating", nil
		case runStatusRunning:
			return r, "Created", nil
		default:
			err = fmt.Errorf("error while creating opentelekomcloud WAF dedicated instance %s. "+
				"Unexpected status: %v", r.ID, r.Status)
			return r, "Error", err
		}
	}
}

func waitForInstanceDeleted(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := instances.Get(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] The opentelekomcloud Waf dedicated instance has been deleted (ID:%s).", id)
				return r, "Deleted", nil
			}
			return nil, "Error", err
		}
		switch r.Status {
		case runStatusDeleting:
			return r, "Deleting", nil
		case runStatusDeleted:
			return r, "Deleted", nil
		default:
			err = fmt.Errorf("error in delete WAF opentelekomcloud dedicated instance[%s]. "+
				"Unexpected status: %v", r.ID, r.Status)
			return r, "Error", err
		}
	}
}

func resourceWafDedicatedInstanceV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	err = instances.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting opentelekomcloud WAF dedicated instance : %w", err)
	}

	log.Printf("[DEBUG] Waiting for opentelekomcloud WAF dedicated instance to be deleted (ID:%s).", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting"},
		Target:       []string{"Deleted"},
		Refresh:      waitForInstanceDeleted(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		log.Printf("[DEBUG] Error while waiting to delete opentelekomcloud WAF dedicated instance. \n%s : %#v", d.Id(), err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func updateInstanceName(client *golangsdk.ServiceClient, id string, name string) error {
	_, err := instances.Update(client, id, instances.UpdateOpts{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("error updating opentelekomcloud WAF dedicated instance: %w", err)
	}
	return nil
}
