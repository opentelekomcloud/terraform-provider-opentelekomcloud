package vpc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/security/group"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/security/rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpcSecGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcSecGroupV3Create,
		ReadContext:   resourceVpcSecGroupV3Read,
		UpdateContext: resourceVpcSecGroupV3Update,
		DeleteContext: resourceVpcSecGroupV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"delete_default_rules": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcSecGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	opts := group.CreateOpts{
		SecurityGroup: group.SecurityGroupOptions{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			EnterpriseProjectId: d.Get("enterprise_project_id").(string),
			Tags:                getVpcSecGroupTags(d),
		},
	}

	securityGroup, err := group.Create(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(securityGroup.ID)

	if d.Get("delete_default_rules").(bool) {
		for _, rule := range securityGroup.SecurityGroupRules {
			err = rules.Delete(client, rule.ID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	log.Printf("[DEBUG] OpenTelekomCloud VPC Security Group `%s` created: %#v", d.Id(), securityGroup)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceVpcSecGroupV3Read(clientCtx, d, meta)
}

func resourceVpcSecGroupV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	securityGroup, err := group.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud Neutron Security group")
	}

	mErr := multierror.Append(nil,
		d.Set("name", securityGroup.Name),
		d.Set("description", securityGroup.Description),
		d.Set("enterprise_project_id", securityGroup.EnterpriseProjectID),
		d.Set("tags", setVpcSecGroupTags(securityGroup.Tags)),
		d.Set("project_id", securityGroup.ProjectID),
		d.Set("created_at", securityGroup.CreatedAt),
		d.Set("updated_at", securityGroup.UpdatedAt),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpcSecGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var groupOpts group.SecurityGroupUpdateOptions

	if d.HasChange("name") {
		groupOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		groupOpts.Description = d.Get("description").(string)
	}

	updateOpts := group.UpdateOpts{
		SecurityGroup: groupOpts,
	}

	log.Printf("[DEBUG] Updating SecGroup %s with options: %#v", d.Id(), updateOpts)
	_, err = group.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud networking SecGroup: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceVpcSecGroupV3Read(clientCtx, d, meta)
}

func resourceVpcSecGroupV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = group.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud networking SecGroup: %s", err)
	}

	d.SetId("")
	return nil
}

func getVpcSecGroupTags(d *schema.ResourceData) []tags.ResourceTag {
	tagsInput := d.Get("tags").(map[string]interface{})
	result := make([]tags.ResourceTag, 0, len(tagsInput))

	for key, value := range tagsInput {
		result = append(result, tags.ResourceTag{
			Key:   key,
			Value: value.(string),
		})
	}
	return result
}

func setVpcSecGroupTags(tagsOutput []tags.ResourceTag) map[string]interface{} {
	result := make(map[string]interface{})
	for _, tag := range tagsOutput {
		result[tag.Key] = tag.Value
	}
	return result
}
