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

		CustomizeDiff: customAuthorizerFunctionVersionDiff,

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
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressCustomAuthorizerFunctionUrnDiff,
			},
			"function_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"function_alias_uri": {
				Type:     schema.TypeString,
				Optional: true,
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

	// The API stores the function URN without version and the version in a separate field,
	// so the version suffix (if any) is stripped from the URN and sent as `authorizer_version`.
	functionUrn, urnVersion := parseFunctionUrn(d.Get("function_urn").(string))
	configuredVersion, isConfigured := getConfiguredFunctionVersion(d)
	functionVersion := resolveFunctionVersion(configuredVersion, isConfigured, urnVersion,
		d.Get("function_version").(string))

	return authorizer.CreateOpts{
		GatewayID:      d.Get("gateway_id").(string),
		Name:           d.Get("name").(string),
		Type:           d.Get("type").(string),
		AuthorizerType: "FUNC", // The custom authorizer only support 'FUNC'.
		FunctionUrn:    functionUrn,
		Version:        functionVersion,
		AliasUrn:       d.Get("function_alias_uri").(string),
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
		d.Set("function_urn", flattenCustomAuthorizerFunctionUrn(d.Get("function_urn").(string), resp)),
		d.Set("function_version", resp.Version),
		d.Set("function_alias_uri", resp.AliasUrn),
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
		Limit:     500,
	}
	resp, err := authorizer.List(client, opts)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error retrieving OpenTelekomCloud APIGW custom authorizer: %s", err)
	}

	// The name filter may not be applied by the server, so the result is filtered by the exact name.
	for _, item := range resp {
		if item.Name == name {
			d.SetId(item.ID)
			return []*schema.ResourceData{d}, d.Set("gateway_id", gatewayId)
		}
	}

	return []*schema.ResourceData{d}, fmt.Errorf("unable to find OpenTelekomCloud APIGW custom authorizer (%s) from server", name)
}

// The function URN format is: urn:fss:<region>:<project_id>:function:<package>:<function_name>[:<version>]
var functionUrnWithVersionRegex = regexp.MustCompile(`^(urn:[^:]+:[^:]+:[^:]+:function:[^:]+:[^:]+):([^:]+)$`)

// parseFunctionUrn splits the function URN into the URN without version and the version (or alias).
// If the URN does not contain a version, the original value and an empty version are returned.
func parseFunctionUrn(urn string) (string, string) {
	if matches := functionUrnWithVersionRegex.FindStringSubmatch(urn); len(matches) == 3 {
		return matches[1], matches[2]
	}
	return urn, ""
}

// getConfiguredFunctionVersion returns the function version from the configuration and whether it is configured.
func getConfiguredFunctionVersion(d *schema.ResourceData) (string, bool) {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() || !rawConfig.Type().IsObjectType() ||
		!rawConfig.Type().HasAttribute("function_version") {
		return "", false
	}
	v := rawConfig.GetAttr("function_version")
	if v.IsNull() || !v.IsKnown() {
		return "", false
	}
	return v.AsString(), true
}

// resolveFunctionVersion determines the function version which will be sent to the API:
//  1. the explicitly configured `function_version`;
//  2. the version suffix of `function_urn` (backward compatibility);
//  3. the current `function_version` value (e.g. read from the API after import).
func resolveFunctionVersion(configuredVersion string, isConfigured bool, urnVersion, currentVersion string) string {
	if isConfigured && configuredVersion != "" {
		return configuredVersion
	}
	if urnVersion != "" {
		return urnVersion
	}
	return currentVersion
}

// suppressCustomAuthorizerFunctionUrnDiff suppresses the difference between the function URN with version suffix and
// the URN without version, when the version suffix is the same as the `function_version` value.
// The API always returns the function URN without version and the version in a separate field.
func suppressCustomAuthorizerFunctionUrnDiff(_, o, n string, d *schema.ResourceData) bool {
	if o == "" || n == "" || o == n {
		return false
	}
	oldUrn, oldVersion := parseFunctionUrn(o)
	newUrn, newVersion := parseFunctionUrn(n)
	if oldUrn != newUrn || (oldVersion != "" && newVersion != "") {
		return false
	}
	version := oldVersion + newVersion // only one of them is not empty
	return version == d.Get("function_version").(string)
}

// customAuthorizerFunctionVersionDiff updates the planned `function_version` value when it is omitted and
// the version is specified as the suffix of the `function_urn`.
func customAuthorizerFunctionVersionDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if !d.HasChange("function_urn") {
		return nil
	}
	if rawConfig := d.GetRawConfig(); !rawConfig.IsNull() && rawConfig.IsKnown() &&
		rawConfig.Type().IsObjectType() && rawConfig.Type().HasAttribute("function_version") &&
		!rawConfig.GetAttr("function_version").IsNull() {
		return nil
	}
	if _, version := parseFunctionUrn(d.Get("function_urn").(string)); version != "" {
		return d.SetNew("function_version", version)
	}
	return nil
}

// flattenCustomAuthorizerFunctionUrn keeps the function URN with the version suffix in the state, if it is equivalent to
// the URN and version returned by the API, to avoid differences for the configurations with the version suffix.
func flattenCustomAuthorizerFunctionUrn(currentUrn string, resp *authorizer.AuthorizerResp) string {
	if urn, version := parseFunctionUrn(currentUrn); version != "" && urn == resp.FunctionUrn && version == resp.Version {
		return currentUrn
	}
	return resp.FunctionUrn
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
