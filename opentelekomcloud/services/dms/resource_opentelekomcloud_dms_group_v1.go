package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsGroupsV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsGroupsV1Create,
		ReadContext:   resourceDmsGroupsV1Read,
		DeleteContext: resourceDmsGroupsV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		DeprecationMessage: "please use `opentelekomcloud_dms_instance_v1` resource instead",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"queue_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumed_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"produced_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"produced_deadletters": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_deadletters": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDmsGroupsV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms group client: %s", err)
	}

	var getGroups []groups.GroupOps

	n := groups.GroupOps{
		Name: d.Get("name").(string),
	}
	getGroups = append(getGroups, n)

	createOpts := &groups.CreateOps{
		Groups: getGroups,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	v, err := groups.Create(DmsV1Client, d.Get("queue_id").(string), createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud group: %s", err)
	}
	log.Printf("[INFO] group Name: %s", v[0].Name)

	// Store the group ID now
	d.SetId(v[0].ID)

	return resourceDmsGroupsV1Read(ctx, d, meta)
}

func resourceDmsGroupsV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms group client: %s", err)
	}

	queueID := d.Get("queue_id").(string)
	page, err := groups.List(DmsV1Client, queueID, false).AllPages()
	if err != nil {
		return fmterr.Errorf("error getting groups in queue %s: %s", queueID, err)
	}

	groupsList, err := groups.ExtractGroups(page)
	if err != nil {
		return fmterr.Errorf("error extracting groups: %w", err)
	}
	if len(groupsList) < 1 {
		return fmterr.Errorf("no matching resource found")
	}

	if len(groupsList) > 1 {
		return fmterr.Errorf("multiple resources matched")
	}

	group := groupsList[0]
	log.Printf("[DEBUG] Dms group %s: %+v", d.Id(), group)

	d.SetId(group.ID)
	mErr := multierror.Append(
		d.Set("name", group.Name),
		d.Set("consumed_messages", group.ConsumedMessages),
		d.Set("available_messages", group.AvailableMessages),
		d.Set("produced_messages", group.ProducedMessages),
		d.Set("produced_deadletters", group.ProducedDeadletters),
		d.Set("available_deadletters", group.AvailableDeadletters),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDmsGroupsV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms group client: %s", err)
	}

	err = groups.Delete(DmsV1Client, d.Get("queue_id").(string), d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud group: %s", err)
	}

	log.Printf("[DEBUG] Dms group %s deactivated.", d.Id())
	d.SetId("")
	return nil
}
