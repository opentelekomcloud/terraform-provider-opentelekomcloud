package waf

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafReferenceTableV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedRefTableV1Create,
		ReadContext:   resourceWafDedicatedRefTableV1Read,
		UpdateContext: resourceWafDedicatedRefTableV1Update,
		DeleteContext: resourceWafDedicatedRefTableV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^([\\w]{1,64})$"),
					"The name can contains of 1 to 64 characters."+
						"Only letters, digits and underscores (_) are allowed."),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"url", "user-agent", "ip", "params", "cookie", "referer", "header",
				}, false),
			},
			"conditions": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 30,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(1, 2048),
				},
				Description: "schema: Required",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceWafDedicatedRefTableV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	opt := rules.CreateReferenceTableOpts{
		Name:   d.Get("name").(string),
		Type:   d.Get("type").(string),
		Values: common.ExpandToStringSlice(d.Get("conditions").([]interface{})),
	}
	log.Printf("[DEBUG] Create OpenTelekomCloud WAF Dedicated reference table options: %#v", opt)

	r, err := rules.CreateReferenceTable(client, opt)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] OpenTelekomCloud Waf Dedicated reference table created: %#v", r)
	d.SetId(r.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedRefTableV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedRefTableV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	r, err := rules.GetReferenceTable(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud WAF dedicated reference table.")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", r.Name),
		d.Set("type", r.Type),
		d.Set("conditions", r.Values),
		d.Set("description", r.Description),
		d.Set("created_at", time.Unix(r.CreatedAt/1000, 0).Format("2006-01-02 15:04:05")),
	)

	if mErr.ErrorOrNil() != nil {
		return fmterr.Errorf("error setting OpenTelekomCloud WAF dedicated reference table fields: %w", err)
	}

	return nil
}

func resourceWafDedicatedRefTableV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	opts := rules.UpdateReferenceTableOpts{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Values:      common.ExpandToStringSlice(d.Get("conditions").([]interface{})),
		Description: d.Get("description").(string),
	}
	log.Printf("[DEBUG] Update OpenTelekomCloud WAF dedicated reference table options: %#v", opts)

	_, err = rules.UpdateReferenceTable(client, d.Id(), opts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud WAF dedicated reference table: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedRefTableV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedRefTableV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	err = rules.DeleteReferenceTable(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF dedicated reference table: %s", err)
	}

	d.SetId("")
	return nil
}
