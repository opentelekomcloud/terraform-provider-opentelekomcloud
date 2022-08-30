package lts

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/loggroups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLTSGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupV2Create,
		ReadContext:   resourceGroupV2Read,
		DeleteContext: resourceGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	createOpts := &loggroups.CreateOpts{
		LogGroupName: d.Get("group_name").(string),
		TTL:          d.Get("ttl_in_days").(int),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	groupCreate, err := loggroups.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating log group: %s", err)
	}

	d.SetId(groupCreate.ID)
	return resourceGroupV2Read(ctx, d, meta)
}

func resourceGroupV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	group, err := loggroups.Get(client, d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error getting OpenTelekomCloud log group %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Retrieved Cluster %s: %#v", d.Id(), group)
	d.SetId(group.ID)
	mErr := multierror.Append(
		d.Set("group_name", group.Name),
		d.Set("ttl_in_days", group.TTLinDays),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGroupV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	err = loggroups.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error deleting log group")
	}

	d.SetId("")
	return nil
}
