package cbr

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/vaults"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func ResourceCBRVaultV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCBRVaultV3Create,
		ReadContext:   resourceCBRVaultV3Read,
		UpdateContext: resourceCBRVaultV3Update,
		DeleteContext: resourceCBRVaultV3Delete,

		CustomizeDiff: common.MultipleCustomizeDiffs(cbrVaultRequiredFields),

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
				ForceNew:     true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"resource": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr, // see ConfigMode documentation for the reasoning
				Set:        hashID,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(0, 255),
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"OS::Nova::Server", "OS::Cinder::Volume",
							}, false),
						},
						"exclude_volumes": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"include_volumes": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backup_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backup_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"billing": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"public", "hybrid",
							}, false),
							Default: "public",
						},
						"consistent_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "crash_consistent",
						},
						"object_type": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"protect_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"backup", "replication",
							}, false),
						},
						"size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 10485760),
						},
						"charging_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"post_paid", "pre_paid",
							}, false),
							Default: "post_paid",
						},
						"period_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"year", "month",
							}, false),
							Default: "month",
						},
						"period_num": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"is_auto_renew": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"is_auto_pay": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"console_url": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						"extra_info": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"order_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allocated": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frozen_scene": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"backup_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_bind": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"bind_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"auto_expand": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCBRVaultV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	vault, err := vaults.Get(client, d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error getting vault details: %s", err)
	}

	var resourceInfo []map[string]interface{}
	for _, resource := range vault.Resources {
		resourceMap := map[string]interface{}{
			"id":              resource.ID,
			"name":            resource.Name,
			"type":            resource.Type,
			"exclude_volumes": resource.ExtraInfo.ExcludeVolumes,
			"include_volumes": resource.ExtraInfo.IncludeVolumes,
		}
		resourceInfo = append(resourceInfo, resourceMap)
	}

	tagsMap := make(map[string]string)
	for _, tag := range vault.Tags {
		tagsMap[tag.Key] = tag.Value
	}

	bindRules := make([]interface{}, len(vault.BindRules.Tags))
	for i, rule := range vault.BindRules.Tags {
		bindRules[i] = map[string]interface{}{
			"key":   rule.Key,
			"value": rule.Value,
		}
	}

	mErr := multierror.Append(
		d.Set("description", vault.Description),
		d.Set("name", vault.Name),
		d.Set("project_id", vault.ProjectID),
		d.Set("provider_id", vault.ProviderID),
		d.Set("resource", resourceInfo),
		d.Set("tags", tagsMap),
		d.Set("enterprise_project_id", vault.EnterpriseProjectID),
		d.Set("auto_bind", vault.AutoBind),
		d.Set("auto_expand", vault.AutoExpand),
		d.Set("bind_rules", bindRules),
		d.Set("user_id", vault.UserID),
		d.Set("created_at", vault.CreatedAt),

		setVaultBilling(d, &vault.Billing),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting vault fields: %s", err)
	}

	return nil
}

func resourceCBRVaultV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	resources, err := cbrVaultResourcesCreate(d)
	if err != nil {
		return fmterr.Errorf("error constructing resources list: %s", err)
	}

	opts := vaults.CreateOpts{
		BackupPolicyID:      d.Get("backup_policy_id").(string),
		Billing:             cbrVaultBillingCreate(d),
		Description:         d.Get("description").(string),
		Name:                d.Get("name").(string),
		Resources:           resources,
		EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		AutoBind:            d.Get("auto_bind").(bool),
		BindRules:           cbrVaultBindRules(d),
		AutoExpand:          d.Get("auto_expand").(bool),
	}

	vault, err := vaults.Create(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating vaults: %s", err)
	}
	d.SetId(vault.ID)

	if policy := d.Get("backup_policy_id").(string); policy != "" {
		_, err := vaults.BindPolicy(client, d.Id(), vaults.BindPolicyOpts{PolicyID: policy}).Extract()
		if err != nil {
			return fmterr.Errorf("error binding policy to vault: %s", err)
		}
	}

	if err := common.UpdateResourceTags(client, d, "vault", d.Id()); err != nil {
		return diag.Errorf("error setting tags of CBR vault: %s", err)
	}

	return resourceCBRVaultV3Read(ctx, d, meta)
}

func resourceExtraMapToExtra(src map[string]interface{}) (*vaults.ResourceExtraInfo, error) {
	var resExtra *vaults.ResourceExtraInfo
	if src == nil {
		return nil, nil
	}

	rawExclude := src["exclude_volumes"].(*schema.Set).List()
	exclude := make([]string, len(rawExclude))
	for i, raw := range rawExclude {
		exclude[i] = raw.(string)
	}

	rawInclude := src["include_volumes"].(*schema.Set).List()
	include := make([]vaults.ResourceExtraInfoIncludeVolumes, len(rawInclude))
	for i, raw := range rawInclude {
		include[i] = vaults.ResourceExtraInfoIncludeVolumes{
			ID: raw.(string),
		}
	}
	if len(exclude) != 0 || len(include) != 0 {
		resExtra = &vaults.ResourceExtraInfo{
			ExcludeVolumes: exclude,
			IncludeVolumes: include,
		}
	}
	return resExtra, nil
}

func cbrVaultResourcesCreate(d *schema.ResourceData) (res []vaults.ResourceCreate, err error) {
	resources := d.Get("resource").(*schema.Set)
	res = make([]vaults.ResourceCreate, resources.Len())
	for i, v := range resources.List() {
		resource := v.(map[string]interface{})
		resourceID := resource["id"].(string)
		resExtra, err := resourceExtraMapToExtra(resource)
		if err != nil {
			return nil, err
		}
		res[i] = vaults.ResourceCreate{
			ID:        resourceID,
			Type:      resource["type"].(string),
			Name:      resource["name"].(string),
			ExtraInfo: resExtra,
		}
	}
	return
}

func setVaultBilling(d *schema.ResourceData, billing *vaults.Billing) error {
	created := d.Get("billing.0").(map[string]interface{})
	created["allocated"] = billing.Allocated
	created["charging_mode"] = billing.ChargingMode
	created["cloud_type"] = billing.CloudType
	created["consistent_level"] = billing.ConsistentLevel
	created["object_type"] = billing.ObjectType
	created["order_id"] = billing.OrderID
	created["product_id"] = billing.ProductID
	created["protect_type"] = billing.ProtectType
	created["size"] = billing.Size
	created["spec_code"] = billing.SpecCode
	created["status"] = billing.Status
	created["storage_unit"] = billing.StorageUnit
	created["used"] = billing.Used
	created["frozen_scene"] = billing.FrozenScene
	return d.Set("billing", []interface{}{created})
}

func cbrVaultBillingCreate(d *schema.ResourceData) *vaults.BillingCreate {
	var billingExtra *vaults.BillingCreateExtraInfo

	if extra, ok := d.GetOk("billing.0.extra_info"); ok {
		extraMap := extra.(map[string]interface{})
		billingExtra = &vaults.BillingCreateExtraInfo{
			CombinedOrderID:     extraMap["combined_order_id"].(string),
			CombinedOrderECSNum: extraMap["combined_order_ecs_num"].(int),
		}
	}

	billing := &vaults.BillingCreate{
		CloudType:       d.Get("billing.0.cloud_type").(string),
		ConsistentLevel: d.Get("billing.0.consistent_level").(string),
		ObjectType:      d.Get("billing.0.object_type").(string),
		ProtectType:     d.Get("billing.0.protect_type").(string),
		Size:            d.Get("billing.0.size").(int),
		ChargingMode:    d.Get("billing.0.charging_mode").(string),
		PeriodType:      d.Get("billing.0.period_type").(string),
		PeriodNum:       d.Get("billing.0.period_num").(int),
		IsAutoRenew:     d.Get("billing.0.is_auto_renew").(bool),
		IsAutoPay:       d.Get("billing.0.is_auto_pay").(bool),
		ConsoleURL:      d.Get("billing.0.console_url").(string),
		ExtraInfo:       billingExtra,
	}

	return billing
}

func cbrVaultBindRules(d *schema.ResourceData) (rules *vaults.VaultBindRules) {
	bingTags := d.Get("bind_rules").([]interface{})
	if len(bingTags) == 0 {
		return
	}
	rules = new(vaults.VaultBindRules)
	rules.Tags = make([]tags.ResourceTag, len(bingTags))
	for i, tag := range bingTags {
		tagMap := tag.(map[string]interface{})
		rules.Tags[i] = tags.ResourceTag{
			Key:   tagMap["key"].(string),
			Value: tagMap["value"].(string),
		}
	}
	return rules
}

func vaultAddedResources(d *schema.ResourceData) ([]vaults.ResourceCreate, error) {
	oldR, newR := d.GetChange("resource")
	addedSet := newR.(*schema.Set).Difference(oldR.(*schema.Set))
	res := make([]vaults.ResourceCreate, addedSet.Len())
	for i, v := range addedSet.List() {
		newMap := v.(map[string]interface{})
		newResource := vaults.ResourceCreate{
			ID:   newMap["id"].(string),
			Type: newMap["type"].(string),
			Name: newMap["name"].(string),
		}

		extra, err := resourceExtraMapToExtra(newMap)
		if err != nil {
			return nil, err
		}

		newResource.ExtraInfo = extra
		res[i] = newResource
	}
	return res, nil
}

func vaultRemovedResources(d *schema.ResourceData) []string {
	oldR, newR := d.GetChange("resource")
	removedSet := oldR.(*schema.Set).Difference(newR.(*schema.Set))

	ids := make([]string, removedSet.Len())
	for i, v := range removedSet.List() {
		removed := v.(map[string]interface{})
		ids[i] = removed["id"].(string)
	}
	return ids
}

func updateResources(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	removedIDs := vaultRemovedResources(d)
	if len(removedIDs) > 0 {
		_, err := vaults.DissociateResources(client, d.Id(), vaults.DissociateResourcesOpts{
			ResourceIDs: removedIDs,
		}).Extract()
		if err != nil {
			return fmt.Errorf("error unbinding resources: %s", err)
		}
	}

	addedResources, err := vaultAddedResources(d)
	if err != nil {
		return err
	}
	if len(addedResources) > 0 {
		_, err := vaults.AssociateResources(client, d.Id(), vaults.AssociateResourcesOpts{
			Resources: addedResources,
		}).Extract()
		if err != nil {
			return fmt.Errorf("error binding resources: %s", err)
		}
	}

	return nil
}

func updatePolicy(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldP, newP := d.GetChange("backup_policy_id")
	if oldP != "" {
		_, err := vaults.UnbindPolicy(client, d.Id(), vaults.BindPolicyOpts{
			PolicyID: oldP.(string),
		}).Extract()
		if err != nil {
			return fmt.Errorf("error unbinding policy from vault: %s", err)
		}
	}
	if newP != "" {
		_, err := vaults.BindPolicy(client, d.Id(), vaults.BindPolicyOpts{
			PolicyID: newP.(string),
		}).Extract()
		if err != nil {
			return fmt.Errorf("error binding policy to vault: %s", err)
		}
	}
	return nil
}

func resourceCBRVaultV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	opts := vaults.UpdateOpts{}
	needsUpdate := false

	if d.HasChange("billing.0.size") {
		opts.Billing = &vaults.BillingUpdate{
			Size: d.Get("billing.0.size").(int),
		}
		needsUpdate = true
	}

	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
		needsUpdate = true
	}

	if d.HasChange("auto_bind") {
		ab := d.Get("auto_bind").(bool)
		opts.AutoBind = &ab
		needsUpdate = true
	}

	if d.HasChange("auto_expand") {
		ae := d.Get("auto_expand").(bool)
		opts.AutoExpand = &ae
		needsUpdate = true
	}

	if d.HasChange("bind_rules") {
		rules := cbrVaultBindRules(d)
		opts.BindRules = rules
		needsUpdate = true
	}

	if needsUpdate {
		_, err := vaults.Update(client, d.Id(), opts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating the vault: %s", err)
		}
	}

	if d.HasChange("resource") {
		if err := updateResources(d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("backup_policy_id") {
		if err := updatePolicy(d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		if err = common.UpdateResourceTags(client, d, "vault", d.Id()); err != nil {
			return diag.Errorf("failed to update CBR tags: %s", err)
		}
	}

	return resourceCBRVaultV3Read(ctx, d, meta)
}

func resourceCBRVaultV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	if err := vaults.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting CBRv3 vault: %s", err)
	}

	d.SetId("")

	return nil
}

func requiredForPrepaid(d *schema.ResourceDiff, field string) error {
	fieldFull := fmt.Sprintf("billing.0.%s", field)
	if _, ok := d.GetOk(fieldFull); !ok {
		return fmt.Errorf("argument \"%s\" is required if \"charging_mode\" is set to \"pre_paid\"", field)
	}
	return nil
}

func cbrVaultRequiredFields(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if d.Get("billing.0.charging_mode") == "pre_paid" {
		mErr := multierror.Append(
			requiredForPrepaid(d, "period_type"),
			requiredForPrepaid(d, "period_num"),
		)
		if err := mErr.ErrorOrNil(); err != nil {
			return err
		}
	}
	return nil
}

func hashID(v interface{}) int {
	res := v.(map[string]interface{})
	return hashcode.String(res["id"].(string))
}
