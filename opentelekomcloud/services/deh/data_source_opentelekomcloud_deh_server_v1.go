package deh

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/deh/v1/hosts"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceDEHServersV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDEHServersV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dedicated_host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDEHServersV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	dehClient, err := config.DehV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating DEHv1 client: %w", err)
	}

	listServerOpts := hosts.ListServerOpts{
		ID:     d.Get("server_id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
		UserID: d.Get("user_id").(string),
	}
	pages, err := hosts.ListServer(dehClient, d.Get("dedicated_host_id").(string), listServerOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve deh server: %s", err)
	}

	if len(pages) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again")
	}
	if len(pages) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	DehServer := pages[0]

	log.Printf("[INFO] Retrieved Deh Server using given filter %s: %+v", DehServer.ID, DehServer)
	d.SetId(DehServer.ID)

	mErr := multierror.Append(
		d.Set("server_id", DehServer.ID),
		d.Set("user_id", DehServer.UserID),
		d.Set("name", DehServer.Name),
		d.Set("status", DehServer.Status),
		d.Set("flavor", DehServer.Flavor),
		d.Set("metadata", DehServer.Metadata),
		d.Set("tenant_id", DehServer.TenantID),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	networks, err := flattenInstanceNetwork(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("addresses", networks); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving network to state for OpenTelekomCloud server (%s): %s", d.Id(), err)
	}

	return nil
}
