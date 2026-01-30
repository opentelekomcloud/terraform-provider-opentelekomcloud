package ces

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/resourcegroups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceResourceGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceGroupV2Create,
		UpdateContext: resourceResourceGroupV2Update,
		ReadContext:   resourceResourceGroupV2Read,
		DeleteContext: resourceResourceGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`),
						"must start with a letter, only letters, digits, underscores (_), and hyphens (-) are allowed",
					),
				),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EPS", "TAG", "Manual",
				}, false),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"associated_eps_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Elem:     resourceGroupResourcesSchema(),
				Optional: true,
			},
			"created_at": {
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

func resourceGroupResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dimensions": {
				Type:     schema.TypeList,
				Elem:     resourceGroupDimensionsSchema(),
				Required: true,
			},
		},
	}
}

func resourceGroupDimensionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceResourceGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	createOpts := resourcegroups.CreateOpts{
		GroupName:           d.Get("name").(string),
		Type:                d.Get("type").(string),
		EnterpriseProjectId: d.Get("enterprise_project_id").(string),
		Tags:                buildResourceGroupTags(d.Get("tags").(map[string]interface{})),
		AssociationEpIds:    expandToStringList(d.Get("associated_eps_ids").([]interface{})),
	}

	log.Printf("[DEBUG] CES Resource Group V2 Create Options: %#v", createOpts)

	groupId, err := resourcegroups.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating CES Resource Group V2: %w", err)
	}

	log.Printf("[DEBUG] CES Resource Group V2 created with ID: %s", groupId)
	d.SetId(groupId)

	// Add resources to the group if specified (for Manual type)
	if v, ok := d.GetOk("resources"); ok {
		createResourcesOpts := resourcegroups.BatchCreateResourcesOpts{
			Resources: buildResourceItems(v.([]interface{})),
		}
		_, err = resourcegroups.BatchCreateResources(client, groupId, createResourcesOpts)
		if err != nil {
			return fmterr.Errorf("error adding resources to CES Resource Group V2: %w", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, cesClientV2)
	return resourceResourceGroupV2Read(clientCtx, d, meta)
}

func resourceResourceGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	group, err := resourcegroups.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CES Resource Group V2")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", group.GroupName),
		d.Set("type", group.Type),
		d.Set("enterprise_project_id", group.EnterpriseProjectId),
		d.Set("tags", flattenResourceGroupTags(group.Tags)),
		d.Set("associated_eps_ids", group.AssociationEpIds),
		d.Set("created_at", group.CreateTime),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceResourceGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	groupId := d.Id()

	// Update name and tags
	if d.HasChanges("name", "tags") {
		updateOpts := resourcegroups.UpdateOpts{
			GroupName: d.Get("name").(string),
			Tags:      buildResourceGroupTags(d.Get("tags").(map[string]interface{})),
		}

		log.Printf("[DEBUG] CES Resource Group V2 Update Options: %#v", updateOpts)

		err = resourcegroups.Update(client, groupId, updateOpts)
		if err != nil {
			return fmterr.Errorf("error updating CES Resource Group V2: %w", err)
		}
	}

	// Update resources (for Manual type)
	if d.HasChange("resources") {
		oldResources, newResources := d.GetChange("resources")

		// Delete old resources
		if len(oldResources.([]interface{})) > 0 {
			deleteResourcesOpts := resourcegroups.BatchDeleteResourcesOpts{
				Resources: buildResourceItems(oldResources.([]interface{})),
			}
			_, err = resourcegroups.BatchDeleteResources(client, groupId, deleteResourcesOpts)
			if err != nil {
				return fmterr.Errorf("error deleting resources from CES Resource Group V2: %w", err)
			}
		}

		// Add new resources
		if len(newResources.([]interface{})) > 0 {
			createResourcesOpts := resourcegroups.BatchCreateResourcesOpts{
				Resources: buildResourceItems(newResources.([]interface{})),
			}
			_, err = resourcegroups.BatchCreateResources(client, groupId, createResourcesOpts)
			if err != nil {
				return fmterr.Errorf("error adding resources to CES Resource Group V2: %w", err)
			}
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, cesClientV2)
	return resourceResourceGroupV2Read(clientCtx, d, meta)
}

func resourceResourceGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	groupId := d.Id()
	log.Printf("[DEBUG] Deleting CES Resource Group V2: %s", groupId)

	deleteOpts := resourcegroups.DeleteOpts{
		GroupIds: []string{groupId},
	}

	_, err = resourcegroups.Delete(client, deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CES Resource Group V2")
	}

	return nil
}

func buildResourceGroupTags(tags map[string]interface{}) []resourcegroups.ResourceGroupTag {
	if len(tags) == 0 {
		return nil
	}

	result := make([]resourcegroups.ResourceGroupTag, 0, len(tags))
	for k, v := range tags {
		result = append(result, resourcegroups.ResourceGroupTag{
			Key:   k,
			Value: v.(string),
		})
	}
	return result
}

func flattenResourceGroupTags(tags []resourcegroups.ResourceGroupTag) map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	result := make(map[string]interface{})
	for _, tag := range tags {
		result[tag.Key] = tag.Value
	}
	return result
}

func expandToStringList(raw []interface{}) []string {
	if len(raw) == 0 {
		return nil
	}

	result := make([]string, len(raw))
	for i, v := range raw {
		result[i] = v.(string)
	}
	return result
}

func buildResourceItems(resources []interface{}) []resourcegroups.ResourceItem {
	if len(resources) == 0 {
		return nil
	}

	result := make([]resourcegroups.ResourceItem, len(resources))
	for i, v := range resources {
		raw := v.(map[string]interface{})
		result[i] = resourcegroups.ResourceItem{
			Namespace:  raw["namespace"].(string),
			Dimensions: buildResourceDimensions(raw["dimensions"].([]interface{})),
		}
	}
	return result
}

func buildResourceDimensions(dimensions []interface{}) []resourcegroups.ResourceDimension {
	if len(dimensions) == 0 {
		return nil
	}

	result := make([]resourcegroups.ResourceDimension, len(dimensions))
	for i, v := range dimensions {
		raw := v.(map[string]interface{})
		result[i] = resourcegroups.ResourceDimension{
			Name:  raw["name"].(string),
			Value: raw["value"].(string),
		}
	}
	return result
}
