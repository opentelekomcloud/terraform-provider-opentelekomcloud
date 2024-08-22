package deh

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/deh/v1/hosts"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDeHHostV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeHHostV1Create,
		ReadContext:   resourceDeHHostV1Read,
		UpdateContext: resourceDeHHostV1Update,
		DeleteContext: resourceDeHHostV1Delete,
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
				Required: true,
			},
			"auto_placement": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"available_vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"available_memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"instance_total": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"instance_uuids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_type_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"cores": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"sockets": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"available_instance_capacities": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceDeHHostV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DehV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud DeH Client: %s", err)
	}

	allocateOpts := hosts.AllocateOpts{
		Name:          d.Get("name").(string),
		HostType:      d.Get("host_type").(string),
		AutoPlacement: d.Get("auto_placement").(string),
		Az:            d.Get("availability_zone").(string),
		Quantity:      1,
	}

	allocate, err := hosts.Allocate(client, allocateOpts).ExtractHost()

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Dedicated Host: %s", err)
	}
	d.SetId(allocate.AllocatedHostIds[0])

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud Dedicated Host (%s) to become available", allocate.AllocatedHostIds[0])

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available", "fault"},
		Refresh:    waitForDeHActive(client, allocate.AllocatedHostIds[0]),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Dedicated Host : %s", stateErr)
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := common.ExpandResourceTags(tagRaw)
		if err := tags.Create(client, "dedicated-host-tags", allocate.AllocatedHostIds[0], tagList).ExtractErr(); err != nil {
			return fmterr.Errorf("error setting tags of DeH Host: %s", err)
		}
	}

	return resourceDeHHostV1Read(ctx, d, meta)
}

func resourceDeHHostV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DehV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud DeH client: %s", err)
	}
	n, err := hosts.Get(client, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmterr.Errorf("error retrieving OpenTelekomCloud Dedicated Host: %s", err)
	}
	mErr := multierror.Append(
		d.Set("name", n.Name),
		d.Set("status", n.State),
		d.Set("auto_placement", n.AutoPlacement),
		d.Set("availability_zone", n.Az),
		d.Set("available_vcpus", n.AvailableVcpus),
		d.Set("available_memory", n.AvailableMemory),
		d.Set("instance_total", n.InstanceTotal),
		d.Set("instance_uuids", n.InstanceUuids),
		d.Set("host_type", n.HostProperties.HostType),
		d.Set("host_type_name", n.HostProperties.HostTypeName),
		d.Set("vcpus", n.HostProperties.Vcpus),
		d.Set("cores", n.HostProperties.Cores),
		d.Set("sockets", n.HostProperties.Sockets),
		d.Set("memory", n.HostProperties.Memory),
		d.Set("available_instance_capacities", getInstanceProperties(n)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	resourceTags, err := tags.Get(client, "dedicated-host-tags", d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud DeH Host tags: %s", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmterr.Errorf("error saving tags for OpenTelekomCloud DeH Host: %s", err)
	}

	return nil
}

func resourceDeHHostV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DehV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud DeH Client: %s", err)
	}
	var updateOpts hosts.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("auto_placement") {
		updateOpts.AutoPlacement = d.Get("auto_placement").(string)
	}

	_, err = hosts.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud Dedicated Host: %s", err)
	}

	if d.HasChange("tags") {
		if err := common.UpdateResourceTags(client, d, "dedicated-host-tags", d.Id()); err != nil {
			return fmterr.Errorf("error updating tags of DeH Host %s: %s", d.Id(), err)
		}
	}

	return resourceDeHHostV1Read(ctx, d, meta)
}

func resourceDeHHostV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DehV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud DeH client: %s", err)
	}

	result := hosts.Delete(client, d.Id())
	if result.Err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Dedicated Host: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "released", "fault", "ERROR"},
		Target:     []string{"deleted"},
		Refresh:    waitForDeHDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Dedicated Host : %s", err)
	}
	d.SetId("")
	return nil
}

func waitForDeHActive(client *golangsdk.ServiceClient, dehID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := hosts.Get(client, dehID).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.State == "creating" {
			return n, "creating", nil
		}

		return n, n.State, nil
	}
}

func waitForDeHDelete(client *golangsdk.ServiceClient, dehID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud Dedicated Host %s.\n", dehID)

		r, err := hosts.Get(client, dehID).Extract()

		log.Printf("[DEBUG] Value after extract: %#v", r)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Dedicated Host %s", dehID)
				return r, "deleted", nil
			}
			return nil, "", err
		}
		if r.State == "deleting" {
			return r, "deleting", nil
		}
		log.Printf("[DEBUG] OpenTelekomCloud Dedicated Host %s still available.\n", dehID)
		return r, r.State, nil
	}
}
func getInstanceProperties(n *hosts.Host) []map[string]interface{} {
	var v []map[string]interface{}
	for _, val := range n.HostProperties.InstanceCapacities {
		mapping := map[string]interface{}{
			"flavor": val.Flavor,
		}
		v = append(v, mapping)
	}
	return v
}
