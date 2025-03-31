package lts

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	hg "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/host-groups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceHostGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostGroupV3Create,
		UpdateContext: resourceHostGroupV3Update,
		ReadContext:   resourceHostGroupV3Read,
		DeleteContext: resourceHostGroupV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"agent_access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceHostGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	if v, ok := d.GetOk("host_ids"); ok {
		stateConf := &resource.StateChangeConf{
			Pending: []string{"pending"},
			Target:  []string{"running"},
			Refresh: waitForHostActive(client, common.ExpandToStringListBySet(v.(*schema.Set))),
			Timeout: d.Timeout(schema.TimeoutCreate),
			Delay:   10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf("error waiting for OpenTelekomCloud ECS host icagent to be active: %w", err)
		}
	}

	group, err := hg.Create(client, hg.CreateOpts{
		Name:            d.Get("name").(string),
		Type:            d.Get("type").(string),
		HostIdList:      common.ExpandToStringListBySet(d.Get("host_ids").(*schema.Set)),
		AgentAccessType: d.Get("agent_access_type").(string),
		Labels:          common.ExpandToStringListBySet(d.Get("labels").(*schema.Set)),
		Tags:            ltsTags(d),
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS v3 host group: %s", err)
	}
	d.SetId(group.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceHostGroupV3Read(clientCtx, d, meta)
}

func waitForHostActive(client *golangsdk.ServiceClient, ids []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := hg.ListHost(client, hg.ListHostOpts{
			HostIdList: ids,
		})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil, "", nil
			}
			return nil, "", err
		}

		if len(resp.Result) == 0 {
			return resp, "pending", nil
		}

		for _, host := range resp.Result {
			if host.HostStatus != "running" {
				return resp, "pending", nil
			}
		}
		return resp, "running", nil
	}
}

func resourceHostGroupV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	requestResp, err := hg.List(client, hg.ListOpts{})
	if err != nil {
		return diag.FromErr(err)
	}
	var groupResult *hg.HostGroupResponse
	for _, gr := range requestResp.Result {
		if gr.ID == d.Id() {
			groupResult = &gr
		}
	}
	if groupResult == nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find OpenTelekomCloud LTS v2 host group by its ID (%s)", d.Id()))
	}
	tagsMap := make(map[string]string)
	for _, tag := range groupResult.Tags {
		tagsMap[tag.Key] = tag.Value
	}
	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", groupResult.Name),
		d.Set("type", groupResult.Type),
		d.Set("host_ids", groupResult.HostIdList),
		d.Set("agent_access_type", groupResult.AgentAccessType),
		d.Set("labels", groupResult.Labels),
		d.Set("tags", tagsMap),
		d.Set("created_at", common.FormatTimeStampRFC3339(groupResult.CreatedAt/1000, false)),
		d.Set("updated_at", common.FormatTimeStampRFC3339(groupResult.UpdatedAt/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceHostGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	updateHostGroupChanges := []string{
		"name",
		"host_ids",
		"tags",
		"labels",
	}

	if d.HasChanges(updateHostGroupChanges...) {
		if v, ok := d.GetOk("host_ids"); ok {
			stateConf := &resource.StateChangeConf{
				Pending: []string{"pending"},
				Target:  []string{"running"},
				Refresh: waitForHostActive(client, common.ExpandToStringListBySet(v.(*schema.Set))),
				Timeout: d.Timeout(schema.TimeoutCreate),
				Delay:   10 * time.Second,
			}

			_, err := stateConf.WaitForStateContext(ctx)
			if err != nil {
				return fmterr.Errorf("error waiting for OpenTelekomCloud ECS host icagent to be active: %w", err)
			}
		}
		h := common.ExpandToStringListBySet(d.Get("host_ids").(*schema.Set))
		tagSlice := ltsTags(d)
		if tagSlice == nil {
			tagSlice = []tags.ResourceTag{}
		}
		l := common.ExpandToStringListBySet(d.Get("labels").(*schema.Set))
		_, err = hg.Update(client, hg.UpdateLogGroupOpts{
			ID:         d.Id(),
			Name:       d.Get("name").(string),
			HostIdList: &h,
			Labels:     &l,
			Tags:       &tagSlice,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceHostGroupV3Read(clientCtx, d, meta)
}

func resourceHostGroupV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	_, err = hg.Delete(client, hg.DeleteOpts{HostGroupIds: []string{d.Id()}})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud LTS v2 host group")
	}

	return nil
}
