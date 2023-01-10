package lts

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLTSGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupV2Create,
		ReadContext:   resourceGroupV2Read,
		UpdateContext: resourceGroupV2Update,
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
				Required: true,
			},
			"creation_time": {
				Type:     schema.TypeInt,
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

	createOpts := groups.CreateOpts{
		LogGroupName: d.Get("group_name").(string),
		TTLInDays:    d.Get("ttl_in_days").(int),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	groupCreate, err := groups.CreateLogGroup(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating log group: %s", err)
	}

	d.SetId(groupCreate)
	return resourceGroupV2Read(ctx, d, meta)
}

func resourceGroupV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	allGroups, err := groups.ListLogGroups(client)
	if err != nil {
		return fmterr.Errorf("error listing OpenTelekomCloud log groups")
	}

	var ltsGroup groups.LogGroup
	for _, group := range allGroups {
		if group.LogGroupId == d.Id() {
			ltsGroup = group
			break
		}
	}

	if ltsGroup.LogGroupId == "" {
		return fmterr.Errorf("OpenTelekomCloud log group %s was not found", d.Id())
	}

	log.Printf("[DEBUG] Retrieved Cluster %s: %#v", d.Id(), ltsGroup)
	d.SetId(ltsGroup.LogGroupId)
	mErr := multierror.Append(
		d.Set("group_name", ltsGroup.LogGroupName),
		d.Set("ttl_in_days", ltsGroup.TTLInDays),
		d.Set("creation_time", ltsGroup.CreationTime),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	_, err = groups.UpdateLogGroup(client, groups.UpdateLogGroupOpts{
		TTLInDays:  int32(d.Get("ttl_in_days").(int)),
		LogGroupId: d.Id(),
	})
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud LTS Log Group: %w", err)
	}

	return resourceGroupV2Read(ctx, d, meta)
}

func resourceGroupV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	err = groups.DeleteLogGroup(client, d.Id())
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault400); ok {
			d.SetId("")
			return nil
		} else {
			return common.CheckDeletedDiag(d, err, "Error deleting log group")
		}
	}

	d.SetId("")
	return nil
}
