package dns

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/zones"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceDNSZoneV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSZoneV2Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"zone_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"serial": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"transferred_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			"masters": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"links": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags": common.TagsSchema(),
		},
	}
}

type TaggedListOpts struct {
	zones.ListOpts

	Tags string `q:"tags"`
}

func (opts TaggedListOpts) ToZoneListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func tagsAsString(tags map[string]interface{}) string {
	tagList := make([]string, 0, len(tags))
	for tag, value := range tags {
		tagList = append(tagList, fmt.Sprintf("%s,%s", tag, value))
	}
	return strings.Join(tagList, "|")
}

func dataSourceDNSZoneV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	dnsClient, err := config.DnsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := TaggedListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOk("email"); ok {
		listOpts.Email = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	if v, ok := d.GetOk("ttl"); ok {
		listOpts.TTL = v.(int)
	}

	if v, ok := d.GetOk("zone_type"); ok {
		listOpts.Type = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		listOpts.Tags = tagsAsString(v.(map[string]interface{}))
	}

	pages, err := zones.List(dnsClient, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to retrieve zones: %s", err)
	}

	allZones, err := zones.ExtractZones(pages)
	if err != nil {
		return fmterr.Errorf("unable to extract zones: %s", err)
	}

	if len(allZones) < 1 {
		return fmterr.Errorf("your query returned no results." +
			"Please change your search criteria and try again")
	}

	if len(allZones) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	zone := allZones[0]

	log.Printf("[DEBUG] Retrieved DNS Zone %s: %+v", zone.ID, zone)
	d.SetId(zone.ID)

	me := &multierror.Error{}
	// strings
	me = multierror.Append(me,
		d.Set("name", zone.Name),
		d.Set("pool_id", zone.PoolID),
		d.Set("email", zone.Email),
		d.Set("description", zone.Description),
		d.Set("status", zone.Status),
		d.Set("zone_type", zone.ZoneType),

		// ints
		d.Set("ttl", zone.TTL),
		d.Set("version", zone.Version),
		d.Set("serial", zone.Serial),

		// time.Times
		d.Set("created_at", zone.CreatedAt.Format(time.RFC3339)),
		d.Set("updated_at", zone.UpdatedAt.Format(time.RFC3339)),
	)
	if err := me.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// maps
	err = d.Set("links", zone.Links)
	if err != nil {
		log.Printf("[DEBUG] Unable to set links: %s", err)
		return diag.FromErr(err)
	}
	// slices
	err = d.Set("masters", zone.Masters)
	if err != nil {
		log.Printf("[DEBUG] Unable to set masters: %s", err)
		return diag.FromErr(err)
	}

	return nil
}
