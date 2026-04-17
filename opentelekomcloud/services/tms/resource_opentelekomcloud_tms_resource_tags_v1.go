package tms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	rt "github.com/opentelekomcloud/gophertelekomcloud/openstack/tms/v1/resource-tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceTmsResourceTagsV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmsResourceTagsV1Create,
		ReadContext:   resourceTmsResourceTagsV1Read,
		UpdateContext: resourceTmsResourceTagsV1Update,
		DeleteContext: resourceTmsResourceTagsV1Delete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func buildResourcesInfo(resources []interface{}) []rt.Resource {
	if len(resources) < 1 {
		return nil
	}

	result := make([]rt.Resource, len(resources))
	for i, val := range resources {
		resource := val.(map[string]interface{})
		result[i] = rt.Resource{
			ResourceType: resource["resource_type"].(string),
			ResourceId:   resource["resource_id"].(string),
		}
	}
	return result
}

func expandResourceTags(tagsInput map[string]interface{}) []rt.ResourceTag {
	result := make([]rt.ResourceTag, 0, len(tagsInput))

	for key, value := range tagsInput {
		result = append(result, rt.ResourceTag{
			Key:   key,
			Value: pointerto.String(value.(string)),
		})
	}
	return result
}

func resourceTmsResourceTagsV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, tmsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.TmsV1Client()
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	opts := rt.BatchOpts{
		ProjectId: d.Get("project_id").(string),
		Resources: buildResourcesInfo(d.Get("resources").([]interface{})),
		Tags:      expandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	failResp, err := rt.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud TMS resource tags: %s", err)
	}
	if len(failResp) > 0 {
		return diag.Errorf("error creating OpenTelekomCloud TMS resource tags: %#v", failResp)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate OpenTelekomCloud TMS resource ID of the TMS tags management: %s", err)
	}
	d.SetId(randUUID)

	return resourceTmsResourceTagsV1Read(ctx, d, meta)
}

func FlattenTagsToMap(tagsResp []rt.ResourceTag) map[string]interface{} {
	result := make(map[string]interface{})
	for _, val := range tagsResp {
		result[val.Key] = *val.Value
	}
	return result
}

func compareTwoTags(localTags, remoteTags map[string]interface{}) (same, diff map[string]interface{}) {
	same = make(map[string]interface{})
	diff = make(map[string]interface{})

	for localKey, localVal := range localTags {
		if remoteVal, ok := remoteTags[localKey]; ok {
			local, isTypeLocalOk := localVal.(string)
			if !isTypeLocalOk {
				log.Printf("[WARN] The type of tag key (%s) in the script is incorrect, want 'string', but got '%T'",
					localKey, localVal)
				continue
			}
			remote, isTypeRemoteOk := remoteVal.(string)
			if !isTypeRemoteOk {
				log.Printf("[WARN] The type of tag key (%s) in the remote response is incorrect, want 'string', but got '%T'",
					localKey, remoteVal)
				continue
			}
			if local == remote {
				same[localKey] = localVal
				continue
			}
		}
		diff[localKey] = localVal
	}
	return
}

func resourceTmsResourceTagsV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, tmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.TmsV2Client()
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	var (
		projectId = d.Get("project_id").(string)
		resources = d.Get("resources").([]interface{})
		resResult = make([]interface{}, 0, len(resources))
		tagsInput = d.Get("tags").(map[string]interface{})
	)

	for _, val := range resources {
		resource := val.(map[string]interface{})
		resourceId := resource["resource_id"].(string)
		opts := rt.ListOpts{
			ResourceType: resource["resource_type"].(string),
			ProjectId:    projectId,
		}
		resp, err := rt.List(client, resourceId, opts)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				continue
			}
			return diag.Errorf("error query resource (%s) tags: %s", resourceId, err)
		}
		actualTags := FlattenTagsToMap(resp)
		same, diff := compareTwoTags(tagsInput, actualTags)
		if len(diff) > 0 {
			log.Printf("[ERROR] The tags of resource (%s) don't contain some tags that are expected to need to be set."+
				" It should contain tags (%#v), but some tags were not set successfully: %#v", resourceId, tagsInput,
				actualTags)
		}
		if len(same) > 0 {
			resResult = append(resResult, val)
		}
	}
	// All tags set failed.
	if len(resResult) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "TMS tags management")
	}
	mErr := multierror.Append(nil,
		d.Set("resources", resResult),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud TMS resources and tags information: %s", err)
	}
	return nil
}

func resourceTmsResourceTagsV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, tmsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.TmsV1Client()
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	var (
		projectId        = d.Get("project_id").(string)
		oldRes, newRes   = d.GetChange("resources")
		oldTags, newTags = d.GetChange("tags")
	)

	deleteOpts := rt.BatchOpts{
		ProjectId: projectId,
		Resources: buildResourcesInfo(oldRes.([]interface{})),
		Tags:      expandResourceTags(oldTags.(map[string]interface{})),
	}
	failResp, err := rt.Delete(client, deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting OpenTelekomCloud TMS resource tags: %s", err)
	}
	if len(failResp) > 0 {
		return diag.Errorf("some tags were not successfully removed: %#v", failResp)
	}

	opts := rt.BatchOpts{
		ProjectId: projectId,
		Resources: buildResourcesInfo(newRes.([]interface{})),
		Tags:      expandResourceTags(newTags.(map[string]interface{})),
	}
	failResp, err = rt.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating resource tags: %s", err)
	}
	if len(failResp) > 0 {
		return diag.Errorf("some tags were not set successfully: %#v", failResp)
	}

	return resourceTmsResourceTagsV1Read(ctx, d, meta)
}

func resourceTmsResourceTagsV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, tmsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.TmsV1Client()
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	deleteOpts := rt.BatchOpts{
		ProjectId: d.Get("project_id").(string),
		Resources: buildResourcesInfo(d.Get("resources").([]interface{})),
		Tags:      expandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	failResp, err := rt.Delete(client, deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud TMS tags management")
	}
	if len(failResp) > 0 {
		return diag.Errorf("some tags were not successfully removed: %#v", failResp)
	}

	return nil
}
