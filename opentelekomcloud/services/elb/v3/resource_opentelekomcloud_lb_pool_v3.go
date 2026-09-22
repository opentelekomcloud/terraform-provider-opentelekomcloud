package v3

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/structs"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/pools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLBPoolV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBPoolV3Create,
		ReadContext:   resourceLBPoolV3Read,
		UpdateContext: resourceLBPoolV3Update,
		DeleteContext: resourceLBPoolV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 0xff),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 0xff),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"lb_algorithm": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP", "QUIC_CID"},
					false,
				),
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"TCP", "UDP", "HTTP", "QUIC", "HTTPS"},
					true,
				),
			},
			"listener_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"listener_id", "loadbalancer_id"},
			},
			"loadbalancer_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"listener_id", "loadbalancer_id"},
			},
			"session_persistence": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(0, 1024),
								validation.StringMatch(regexp.MustCompile(`^[\w-.]+$`),
									"The value can contain only letters, digits, hyphens (-), underscores (_), and periods (.)."),
							),
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"SOURCE_IP", "HTTP_COOKIE", "APP_COOKIE"},
								false,
							),
						},
						"persistence_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"slow_start": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      30,
							ValidateFunc: validation.IntBetween(30, 1200),
						},
					},
				},
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"healthmonitor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"loadbalancer_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"member_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"member_deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"nonProtection", "consoleProtection",
				}, false),
			},
			"protection_reason": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
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

func resourceLBPoolV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}
	deletionProtection := d.Get("member_deletion_protection").(bool)
	opts := pools.CreateOpts{
		LBMethod:                 d.Get("lb_algorithm").(string),
		Protocol:                 d.Get("protocol").(string),
		LoadbalancerID:           d.Get("loadbalancer_id").(string),
		ListenerID:               d.Get("listener_id").(string),
		ProjectID:                d.Get("project_id").(string),
		Name:                     d.Get("name").(string),
		Description:              d.Get("description").(string),
		DeletionProtectionEnable: &deletionProtection,
		Type:                     d.Get("type").(string),
		VpcId:                    d.Get("vpc_id").(string),
		ProtectionStatus:         d.Get("protection_status").(string),
		ProtectionReason:         d.Get("protection_reason").(string),
	}

	if d.Get("session_persistence.#").(int) > 0 {
		persistMap := d.Get("session_persistence.0").(map[string]interface{})
		opts.Persistence = mapToPersistence(persistMap)
	}
	if d.Get("slow_start.#").(int) > 0 {
		opts.SlowStart = expandPoolSlowStart(d.Get("slow_start"))
	}

	pool, err := pools.Create(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating LB Pool v3: %w", err)
	}
	d.SetId(pool.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBPoolV3Read(clientCtx, d, meta)
}

func resourceLBPoolV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	pool, err := pools.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error viewing details of LB Pool v3")
	}

	return setLBPoolV3Fields(d, pool)
}

func setLBPoolV3Fields(d *schema.ResourceData, pool *pools.Pool) diag.Diagnostics {
	mErr := multierror.Append(
		d.Set("name", pool.Name),
		d.Set("description", pool.Description),
		d.Set("lb_algorithm", pool.LBMethod),
		d.Set("project_id", pool.ProjectID),
		d.Set("protocol", pool.Protocol),
		d.Set("session_persistence", expandPersistence(pool.Persistence)),
		d.Set("slow_start", flattenPoolSlowStart(pool.SlowStart)),
		d.Set("ip_version", pool.IpVersion),
		d.Set("healthmonitor_id", pool.MonitorID),
		d.Set("listener_ids", poolRefIDs(pool.Listeners)),
		d.Set("loadbalancer_ids", poolRefIDs(pool.Loadbalancers)),
		d.Set("member_ids", poolRefIDs(pool.Members)),
		d.Set("member_deletion_protection", pool.DeletionProtectionEnable),
		d.Set("type", pool.Type),
		d.Set("vpc_id", pool.VpcId),
		d.Set("protection_status", pool.ProtectionStatus),
		d.Set("protection_reason", pool.ProtectionReason),
		d.Set("created_at", pool.CreatedAt),
		d.Set("updated_at", pool.UpdatedAt),
	)
	if len(pool.Loadbalancers) > 0 {
		mErr = multierror.Append(mErr, d.Set("loadbalancer_id", pool.Loadbalancers[0].ID))
	}
	if len(pool.Listeners) > 0 {
		mErr = multierror.Append(mErr, d.Set("listener_id", pool.Listeners[0].ID))
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting LB Pool v3 fields: %w", err)
	}

	d.SetId(pool.ID)
	return nil
}

func expandPersistence(p *pools.SessionPersistence) []interface{} {
	if p == nil {
		return nil
	}
	persistenceMap := map[string]interface{}{
		"cookie_name":         p.CookieName,
		"type":                p.Type,
		"persistence_timeout": p.PersistenceTimeout,
	}
	return []interface{}{persistenceMap}
}

func mapToPersistence(src map[string]interface{}) *pools.SessionPersistence {
	return &pools.SessionPersistence{
		Type:               src["type"].(string),
		CookieName:         src["cookie_name"].(string),
		PersistenceTimeout: src["persistence_timeout"].(int),
	}
}

func expandPoolSlowStart(raw interface{}) *pools.SlowStart {
	items := raw.([]interface{})
	if len(items) == 0 || items[0] == nil {
		return nil
	}
	item := items[0].(map[string]interface{})
	return &pools.SlowStart{
		Enable:   item["enable"].(bool),
		Duration: item["duration"].(int),
	}
}

func flattenPoolSlowStart(slowStart *pools.SlowStart) []interface{} {
	if slowStart == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"enable":   slowStart.Enable,
		"duration": slowStart.Duration,
	}}
}

func poolRefIDs(refs []structs.ResourceRef) []string {
	ids := make([]string, len(refs))
	for i, ref := range refs {
		ids[i] = ref.ID
	}
	return ids
}

func resourceLBPoolV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	opts := pools.UpdateOpts{}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
	}
	if d.HasChange("lb_algorithm") {
		opts.LBMethod = d.Get("lb_algorithm").(string)
	}
	if d.HasChange("session_persistence") {
		if d.Get("session_persistence.#").(int) > 0 {
			persistMap := d.Get("session_persistence.0").(map[string]interface{})
			opts.Persistence = mapToPersistence(persistMap)
		}
	}
	if d.HasChange("slow_start") {
		opts.SlowStart = expandPoolSlowStart(d.Get("slow_start"))
		if opts.SlowStart == nil {
			opts.SlowStart = &pools.SlowStart{Enable: false, Duration: 30}
		}
	}
	if d.HasChange("member_deletion_protection") {
		memberDeletionProtection := d.Get("member_deletion_protection").(bool)
		opts.DeletionProtectionEnable = &memberDeletionProtection
	}
	if d.HasChange("protection_status") {
		opts.ProtectionStatus = d.Get("protection_status").(string)
	}
	if d.HasChange("protection_reason") {
		opts.ProtectionReason = d.Get("protection_reason").(string)
	}
	// https://jira.tsi-dev.otc-service.com/browse/BM-1642
	// if d.HasChange("type") {
	// 	opts.Type = d.Get("type").(string)
	// }
	// if d.HasChange("vpc_id") {
	// 	opts.VpcId = d.Get("vpc_id").(string)
	// }

	_, err = pools.Update(client, d.Id(), opts)
	if err != nil {
		return fmterr.Errorf("error updating LB Pool v3: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBPoolV3Read(clientCtx, d, meta)
}

func resourceLBPoolV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	if err := pools.Delete(client, d.Id()); err != nil {
		return fmterr.Errorf("error deleting LB Pool v3: %w", err)
	}

	return nil
}
