package cbr

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/vaults"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
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
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr, // see ConfigMode documentation for the reasoning
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
						"extra_info": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
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

func resourceCBRVaultV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	vault, err := vaults.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("error getting vault details: %s", err)
	}

	resourceList := make([]interface{}, len(vault.Resources))
	for i, resource := range vault.Resources {
		data, _ := json.Marshal(resource)
		resMap := make(map[string]interface{})
		err = json.Unmarshal(data, &resMap)
		if err != nil {
			return diag.Errorf("error converting resource list: %s", err)
		}
		resourceList[i] = resMap
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
		d.Set("resource", resourceList),
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
		return diag.Errorf("error setting vault fields: %s", err)
	}

	return nil
}

func resourceCBRVaultV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	resources, err := cbrVaultResourcesCreate(d)
	if err != nil {
		return diag.Errorf("error constructing resources list: %s", err)
	}

	opts := vaults.CreateOpts{
		BackupPolicyID:      d.Get("backup_policy_id").(string),
		Billing:             cbrVaultBillingCreate(d),
		Description:         d.Get("description").(string),
		Name:                d.Get("name").(string),
		Resources:           resources,
		Tags:                cbrVaultTags(d),
		EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		AutoBind:            d.Get("auto_bind").(bool),
		BindRules:           cbrVaultBindRules(d),
		AutoExpand:          d.Get("auto_expand").(bool),
	}

	vault, err := vaults.Create(client, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating vaults: %s", err)
	}
	d.SetId(vault.ID)

	if policy := d.Get("backup_policy_id").(string); policy != "" {
		_, err := vaults.BindPolicy(client, d.Id(), vaults.BindPolicyOpts{PolicyID: policy}).Extract()
		if err != nil {
			return diag.Errorf("error binding policy to vault: %s", err)
		}
	}

	return resourceCBRVaultV3Read(ctx, d, meta)
}

func resFieldName(i int, field string) string {
	return fmt.Sprintf("resource.%d.%s", i, field)
}

func resourceExtraMapToExtra(src map[string]interface{}) (*vaults.ResourceExtraInfo, error) {
	if src == nil {
		return nil, nil
	}
	data, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	extra := new(vaults.ResourceExtraInfo)
	err = json.Unmarshal(data, extra)
	if err != nil {
		return nil, err
	}
	return extra, nil
}

func cbrVaultResourcesCreate(d *schema.ResourceData) (res []vaults.ResourceCreate, err error) {
	resourceCount := d.Get("resource.#").(int)
	res = make([]vaults.ResourceCreate, resourceCount)
	for i := 0; i < resourceCount; i++ {
		rawExtra := d.Get(resFieldName(i, "extra_info")).(map[string]interface{})
		resourceID := d.Get(resFieldName(i, "id")).(string)
		extra, err := resourceExtraMapToExtra(rawExtra)
		if err != nil {
			return nil, fmt.Errorf("error converting \"%s\" resource extra: %v", resourceID, rawExtra)
		}
		res[i] = vaults.ResourceCreate{
			ID:        resourceID,
			Type:      d.Get(resFieldName(i, "type")).(string),
			Name:      d.Get(resFieldName(i, "name")).(string),
			ExtraInfo: extra,
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
	tags := d.Get("bind_rules").([]interface{})
	if len(tags) == 0 {
		return
	}
	rules = new(vaults.VaultBindRules)
	rules.Tags = make([]vaults.Tag, len(tags))
	for i, tag := range tags {
		tagMap := tag.(map[string]interface{})
		rules.Tags[i] = vaults.Tag{
			Key:   tagMap["key"].(string),
			Value: tagMap["value"].(string),
		}
	}
	return rules
}

func cbrVaultTags(d *schema.ResourceData) []vaults.Tag {
	tags := d.Get("tags").(map[string]interface{})
	var tagSlice []vaults.Tag
	for k, v := range tags {
		tagSlice = append(tagSlice, vaults.Tag{Key: k, Value: v.(string)})
	}
	return tagSlice
}

func vaultAddedResources(d *schema.ResourceData) (res []vaults.ResourceCreate, err error) {
	oldR, newR := d.GetChange("resource")
	oldSlice := oldR.([]interface{})
	newSlice := newR.([]interface{})

	oldIDs := common.NewStringSearcher()
	for _, oldOne := range oldSlice {
		oldMap := oldOne.(map[string]interface{})
		oldIDs.AddToIndex(oldMap["id"].(string))
	}

	for _, newOne := range newSlice {
		newMap := newOne.(map[string]interface{})
		newID := newMap["id"].(string)
		if !oldIDs.Contains(newID) {
			newResource := vaults.ResourceCreate{
				ID:   newID,
				Type: newMap["type"].(string),
				Name: newMap["name"].(string),
			}
			extraMap, ok := newMap["extra"].(map[string]interface{})
			if ok {
				extra, err := resourceExtraMapToExtra(extraMap)
				if err != nil {
					return nil, err
				}
				newResource.ExtraInfo = extra
			}
			res = append(res, newResource)
		}
	}
	return
}

func vaultRemovedResources(d *schema.ResourceData) (ids []string) {
	oldR, newR := d.GetChange("resource")

	oldSlice := oldR.([]interface{})
	newSlice := newR.([]interface{})

	newIDs := common.NewStringSearcher()
	for _, newOne := range newSlice {
		newIDs.AddToIndex(newOne.(map[string]interface{})["id"].(string))
	}

	for _, oldOne := range oldSlice {
		oldID := oldOne.(map[string]interface{})["id"].(string)
		if !newIDs.Contains(oldID) {
			ids = append(ids, oldID)
		}
	}
	return
}

func updateResources(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if removedIDs := vaultRemovedResources(d); removedIDs != nil {
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
	if addedResources != nil {
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
		return diag.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
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
			return diag.Errorf("error updating the vault: %s", err)
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

	return resourceCBRVaultV3Read(ctx, d, meta)
}

func resourceCBRVaultV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	if err := vaults.Delete(client, d.Id()).ExtractErr(); err != nil {
		return diag.Errorf("error deleting CBRv3 vault: %s", err)
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

func cbrVaultRequiredFields(ctx context.Context, d *schema.ResourceDiff, _ interface{}) error {
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
