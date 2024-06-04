package apigw

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/authorizer"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

type AuthType string

const (
	AuthTypeFrontend AuthType = "FRONTEND"
	AuthTypeBackend  AuthType = "BACKEND"
)

func ResourceAPICustomAuthorizerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomAuthorizerV2Create,
		ReadContext:   resourceCustomAuthorizerV2Read,
		UpdateContext: resourceCustomAuthorizerV2Update,
		DeleteContext: resourceCustomAuthorizerV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCustomAuthorizerImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z][\w]*$`),
						"Only letters, digits and underscores (_) are allowed, and must start with a letter."),
					validation.StringLenBetween(3, 64),
				),
			},
			"function_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(AuthTypeFrontend),
				ValidateFunc: validation.StringInSlice([]string{
					string(AuthTypeFrontend), string(AuthTypeBackend),
				}, false),
			},
			"is_body_send": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 3600),
			},
			"user_data": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 2048),
			},
			"identity": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"location": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HEADER", "QUERY",
							}, false),
						},
						"validation": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 2048),
						},
					},
				},
			},
			"region": {
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

func buildIdentities(identities *schema.Set) []authorizer.Identity {
	if identities.Len() < 1 {
		return nil
	}
	result := make([]authorizer.Identity, identities.Len())
	for i, val := range identities.List() {
		identity := val.(map[string]interface{})
		result[i] = authorizer.Identity{
			Name:       identity["name"].(string),
			Location:   identity["location"].(string),
			Validation: identity["validation"].(string),
		}
	}
	return result
}

func buildCustomAuthorizerOpts(d *schema.ResourceData) (authorizer.CreateOpts, error) {
	identities := d.Get("identity").(*schema.Set)
	if identities.Len() > 0 && d.Get("type").(string) != string(AuthTypeFrontend) {
		return authorizer.CreateOpts{}, fmt.Errorf("the identities can only be set when the type is 'FRONTEND'")
	}

	return authorizer.CreateOpts{
		GatewayID:      d.Get("gateway_id").(string),
		Name:           d.Get("name").(string),
		Type:           d.Get("type").(string),
		AuthorizerType: "FUNC", // The custom authorizer only support 'FUNC'.
		FunctionUrn:    d.Get("function_urn").(string),
		NeedBody:       pointerto.Bool(d.Get("is_body_send").(bool)),
		Ttl:            pointerto.Int(d.Get("ttl").(int)),
		UserData:       d.Get("user_data").(string),
		Identities:     buildIdentities(identities),
	}, nil
}

func resourceCustomAuthorizerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	opts, err := buildCustomAuthorizerOpts(d)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := authorizer.Create(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "custom authorizer")
	}
	d.SetId(resp.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceCustomAuthorizerV2Read(clientCtx, d, meta)
}

func resourceCustomAuthorizerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	resp, err := authorizer.Get(client, d.Get("gateway_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to get OpenTelekomCloud APIGW custom authorizer: %s", d.Id()))
	}
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("function_urn", resp.FunctionUrn),
		d.Set("type", resp.Type),
		d.Set("is_body_send", resp.NeedBody),
		d.Set("ttl", resp.Ttl),
		d.Set("user_data", resp.UserData),
		d.Set("created_at", resp.CreatedAt),
		d.Set("identity", flattenCustomAuthorizerIdentities(resp.Identities)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud APIGW custom authorizer fields: %s", err)
	}
	return nil
}

func resourceCustomAuthorizerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts, err := buildCustomAuthorizerOpts(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = authorizer.Update(client, d.Id(), opts)
	if err != nil {
		return diag.Errorf("error updating OpenTelekomCloud APIGW custom authorizer (%s): %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceCustomAuthorizerV2Read(clientCtx, d, meta)
}

func resourceCustomAuthorizerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	gatewayId := d.Get("gateway_id").(string)
	err = authorizer.Delete(client, gatewayId, d.Id())
	if err != nil {
		return diag.Errorf("error deleting OpenTelekomCloud APIGW custom authorizer (%s) from gateway (%s): %s",
			d.Id(), gatewayId, err)
	}
	return nil
}

// The ID cannot find on the console, so we need to import by authorizer name.
func resourceCustomAuthorizerImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <gateway_id>/<name>")
	}
	gatewayId := parts[0]
	name := parts[1]

	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := authorizer.ListOpts{
		GatewayID: gatewayId,
		Name:      name,
	}
	resp, err := authorizer.List(client, opts)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error retrieving OpenTelekomCloud APIGW custom authorizer: %s", err)
	}

	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmt.Errorf("unable to find OpenTelekomCloud APIGW custom authorizer (%s) form server: %s",
			name, err)
	}
	d.SetId(resp[0].ID)

	return []*schema.ResourceData{d}, d.Set("gateway_id", gatewayId)
}

func flattenCustomAuthorizerIdentities(identities []authorizer.Identity) []map[string]interface{} {
	if len(identities) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(identities))
	for i, val := range identities {
		result[i] = map[string]interface{}{
			"name":       val.Name,
			"location":   val.Location,
			"validation": val.Validation,
		}
	}
	return result
}
