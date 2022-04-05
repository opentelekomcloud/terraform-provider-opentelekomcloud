package vpc

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/provider"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/networks"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ExtractValSFromNid(s string) (string, string) {
	rgs := strings.Split(s, ":")
	if len(rgs) >= 2 {
		log.Printf("[DEBUG] ExtractValSFromNid: %s:%s from (%s)", rgs[0], rgs[1], s)
		return rgs[0], rgs[1]
	}
	log.Printf("[DEBUG] ExtractValSFromNid: true:'%s' from (%s)", s, s)
	return "true", s
}

func ExtractValFromNid(s string) (bool, string) {
	sasu, id := ExtractValSFromNid(s)
	asu, err := strconv.ParseBool(sasu)
	if err != nil {
		return true, id // Should never occur?
	}
	return asu, id
}

func FormatNidFromValS(_ string, id string) string {
	// Causing problems with instance network lookups right now
	// return fmt.Sprintf("%s:%s", asu, id)
	return id
}

func suppressAsuDiff(k, old, new string, d *schema.ResourceData) bool {
	_, idOld := ExtractValSFromNid(old)
	_, idNew := ExtractValSFromNid(new)
	return idOld == idNew
}

func ResourceNetworkingNetworkV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingNetworkV2Create,
		ReadContext:   resourceNetworkingNetworkV2Read,
		UpdateContext: resourceNetworkingNetworkV2Update,
		DeleteContext: resourceNetworkingNetworkV2Delete,
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
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"admin_state_up": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Computed:     true,
				ValidateFunc: common.ValidateTrueOnly,
			},
			"shared": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"segments": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"physical_network": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"segmentation_id": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingNetworkV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	createOpts := NetworkCreateOpts{
		networks.CreateOpts{
			Name:     d.Get("name").(string),
			TenantID: d.Get("tenant_id").(string),
		},
		common.MapValueSpecs(d),
	}

	asuRaw := d.Get("admin_state_up").(string)
	asu := false
	if asuRaw != "" {
		asuT, err := strconv.ParseBool(asuRaw)
		if err != nil {
			return fmterr.Errorf("admin_state_up, if provided, must be either 'true' or 'false'")
		}
		asu = asuT
		// asuFake := true
		// createOpts.AdminStateUp = &asuFake //&asu
	}

	sharedRaw := d.Get("shared").(string)
	if sharedRaw != "" {
		shared, err := strconv.ParseBool(sharedRaw)
		if err != nil {
			return fmterr.Errorf("shared, if provided, must be either 'true' or 'false': %v", err)
		}
		createOpts.Shared = &shared
	}

	segments := resourceNetworkingNetworkV2Segments(d)

	var n *networks.Network
	if len(segments) > 0 {
		providerCreateOpts := provider.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			Segments:          segments,
		}
		log.Printf("[DEBUG] Create Options: %#v", providerCreateOpts)
		n, err = networks.Create(networkingClient, providerCreateOpts).Extract()
	} else {
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		n, err = networks.Create(networkingClient, createOpts).Extract()
	}

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Neutron network: %s", err)
	}
	d.SetId(FormatNidFromValS(strconv.FormatBool(asu), n.ID))

	log.Printf("[INFO] Network ID: %s", n.ID)

	log.Printf("[DEBUG] Waiting for Network (%s) to become available", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForNetworkActive(networkingClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for network to become active: %w", err)
	}

	d.SetId(FormatNidFromValS(strconv.FormatBool(asu), n.ID))

	return resourceNetworkingNetworkV2Read(ctx, d, meta)
}

func resourceNetworkingNetworkV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	asu, id := ExtractValSFromNid(d.Id())
	n, err := networks.Get(networkingClient, id).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "network")
	}

	log.Printf("[DEBUG] Retrieved Network %s: %+v", d.Id(), n)

	d.SetId(FormatNidFromValS(asu, n.ID))
	mErr := multierror.Append(
		d.Set("name", n.Name),
		d.Set("admin_state_up", asu),
		d.Set("shared", strconv.FormatBool(n.Shared)),
		d.Set("tenant_id", n.TenantID),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNetworkingNetworkV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	var updateOpts networks.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	asu := false
	if d.HasChange("admin_state_up") {
		asuRaw := d.Get("admin_state_up").(string)
		if asuRaw != "" {
			asuT, err := strconv.ParseBool(asuRaw) // nolint:staticcheck
			if err != nil {
				return fmterr.Errorf("admin_state_up, if provided, must be either 'true' or 'false'")
			}
			asu = asuT // nolint:ineffassign,staticcheck
			// asuFake := true
			// updateOpts.AdminStateUp = &asuFake //&asu
		}
	}
	if d.HasChange("shared") {
		sharedRaw := d.Get("shared").(string)
		if sharedRaw != "" {
			shared, err := strconv.ParseBool(sharedRaw)
			if err != nil {
				return fmterr.Errorf("shared, if provided, must be either 'true' or 'false': %v", err)
			}
			updateOpts.Shared = &shared
		}
	}
	asu, id := ExtractValFromNid(d.Id())

	log.Printf("[DEBUG] Updating Network %s with options: %+v", id, updateOpts)

	_, err = networks.Update(networkingClient, id, updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud Neutron Network: %s", err)
	}

	d.SetId(FormatNidFromValS(strconv.FormatBool(asu), id))
	return resourceNetworkingNetworkV2Read(ctx, d, meta)
}

func resourceNetworkingNetworkV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	_, id := ExtractValFromNid(d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNetworkDelete(networkingClient, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Neutron Network: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceNetworkingNetworkV2Segments(d *schema.ResourceData) (providerSegments []provider.Segment) {
	segments := d.Get("segments").([]interface{})
	for _, v := range segments {
		var segment provider.Segment
		segmentMap := v.(map[string]interface{})

		if v, ok := segmentMap["physical_network"].(string); ok {
			segment.PhysicalNetwork = v
		}

		if v, ok := segmentMap["network_type"].(string); ok {
			segment.NetworkType = v
		}

		if v, ok := segmentMap["segmentation_id"].(int); ok {
			segment.SegmentationID = v
		}

		providerSegments = append(providerSegments, segment)
	}
	return
}

func waitForNetworkActive(networkingClient *golangsdk.ServiceClient, networkId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := networks.Get(networkingClient, networkId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Neutron Network: %+v", n)
		if n.Status == "DOWN" || n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, n.Status, nil
	}
}

func waitForNetworkDelete(networkingClient *golangsdk.ServiceClient, networkId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud Network %s.\n", networkId)

		n, err := networks.Get(networkingClient, networkId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Network %s", networkId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = networks.Delete(networkingClient, networkId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Network %s", networkId)
				return n, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return n, "ACTIVE", nil
				}
			}
			return n, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Network %s still active.\n", networkId)
		return n, "ACTIVE", nil
	}
}
