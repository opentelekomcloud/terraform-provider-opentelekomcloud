package apigw

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	apis "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/api"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIApiV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIGWApiV2Create,
		ReadContext:   resourceAPIGWApiV2Read,
		UpdateContext: resourceAPIGWApiV2Update,
		DeleteContext: resourceAPIGWApiV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAPIGWApiV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ApiTypePublic), string(ApiTypePrivate),
				}, false),
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(RequestMethodGet), string(RequestMethodPost), string(RequestMethodPut),
					string(RequestMethodDelete), string(RequestMethodHead), string(RequestMethodPatch),
					string(RequestMethodOptions), string(RequestMethodAny),
				}, false),
			},
			"request_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ProtocolTypeHTTP), string(ProtocolTypeHTTPS), string(ProtocolTypeBoth),
				}, false),
			},
			"security_authentication_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ApiAuthTypeNone),
				ValidateFunc: validation.StringInSlice([]string{
					string(ApiAuthTypeNone), string(ApiAuthTypeApp),
					string(ApiAuthTypeIam), string(ApiAuthTypeAuthorizer),
				}, false),
			},
			"security_authentication_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"authorizer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"body_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			},
			"cors": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 255),
				),
			},
			"match_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(MatchModeExact),
				ValidateFunc: validation.StringInSlice([]string{
					string(MatchModePrefix), string(MatchModeExact),
				}, false),
			},
			"success_response": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			},
			"failure_response": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"response_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mock": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"func_graph", "http"},
				Elem:         mockSchemaResource(),
			},
			"func_graph": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     funcGraphSchemaResource(),
			},
			"request_params": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 50,
				Elem:     requestParamsSchemaResource(),
				Set:      resourceBackendParametersHash,
			},
			"backend_params": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 50,
				Elem:     backendParamSchemaResource(),
				Set:      resourceBackendParametersHash,
			},
			"mock_policy": {
				Type:          schema.TypeSet,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"func_graph", "http", "func_graph_policy", "http_policy"},
				Elem:          mockPolicySchemaResource(),
				Description:   "The mock policy backends.",
			},
			"func_graph_policy": {
				Type:          schema.TypeSet,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"mock", "http", "mock_policy", "http_policy"},
				Elem:          functionPolicySchemaResource(),
			},
			"http": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     httpSchemaResource(),
			},
			"http_policy": {
				Type:          schema.TypeSet,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"mock", "func_graph", "mock_policy", "func_graph_policy"},
				Elem:          httpPolicySchemaResource(),
			},
			"registered_at": {
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

func buildApiType(t string) int {
	switch t {
	case string(ApiTypePublic):
		return 1
	default:
		return 2
	}
}

func buildNetworkType(t string) string {
	switch t {
	case string(NetworkTypeV1):
		return "V1"
	default:
		return "V2"
	}
}

func isObjectEnabled(isEnabled bool) *int {
	if isEnabled {
		return pointerto.Int(strBoolEnabled)
	}
	return pointerto.Int(strBoolDisabled)
}

func buildMockStructure(mocks []interface{}) *apis.MockInfo {
	if len(mocks) < 1 {
		return nil
	}

	mockMap := mocks[0].(map[string]interface{})
	return &apis.MockInfo{
		Response:     mockMap["response"].(string),
		AuthorizerID: mockMap["authorizer_id"].(string),
		Description:  mockMap["description"].(string),
		Version:      mockMap["version"].(string),
	}
}

func buildFuncGraphStructure(funcGraphs []interface{}) *apis.FuncInfo {
	if len(funcGraphs) < 1 {
		return nil
	}

	funcMap := funcGraphs[0].(map[string]interface{})
	return &apis.FuncInfo{
		Timeout:        funcMap["timeout"].(int),
		InvocationType: funcMap["invocation_type"].(string),
		FunctionUrn:    funcMap["function_urn"].(string),
		NetworkType:    buildNetworkType(funcMap["network_type"].(string)),
		Version:        funcMap["version"].(string),
		AuthorizerID:   funcMap["authorizer_id"].(string),
	}
}

func buildHttpStructure(webs []interface{}) *apis.BackendApi {
	if len(webs) < 1 {
		return nil
	}

	webMap := webs[0].(map[string]interface{})
	webResp := apis.BackendApi{
		AuthorizerID:    webMap["authorizer_id"].(string),
		UrlDomain:       webMap["url_domain"].(string),
		ReqProtocol:     webMap["request_protocol"].(string),
		Description:     webMap["description"].(string),
		ReqMethod:       webMap["request_method"].(string),
		Version:         webMap["version"].(string),
		ReqUri:          webMap["request_uri"].(string),
		Timeout:         webMap["timeout"].(int),
		EnableClientSSL: pointerto.Bool(webMap["ssl_enable"].(bool)),
		RetryCount:      strconv.Itoa(webMap["retry_count"].(int)),
	}
	if chanId, ok := webMap["vpc_channel_id"]; ok && chanId != "" {
		webResp.VpcChannelStatus = pointerto.Int(strBoolEnabled)
		webResp.VpcChannelInfo = &apis.VpcChannelInfo{
			VpcChannelProxyHost: webMap["vpc_channel_proxy_host"].(string),
			VpcChannelID:        webMap["vpc_channel_id"].(string),
		}
	} else {
		webResp.VpcChannelStatus = pointerto.Int(strBoolDisabled)
		webResp.UrlDomain = webMap["url_domain"].(string)
	}

	return &webResp
}

func buildRequestParameters(requests *schema.Set) []apis.ReqParams {
	if requests.Len() < 1 {
		return nil
	}

	result := make([]apis.ReqParams, requests.Len())
	for i, v := range requests.List() {
		paramMap := v.(map[string]interface{})
		paramType := paramMap["type"].(string)
		param := apis.ReqParams{
			Type:         paramType,
			Name:         paramMap["name"].(string),
			Location:     paramMap["location"].(string),
			DefaultValue: paramMap["default"].(string),
			SampleValue:  paramMap["sample"].(string),
			Required:     isObjectEnabled(paramMap["required"].(bool)),
			ValidEnable:  isObjectEnabled(paramMap["validity_check"].(bool)),
			Description:  paramMap["description"].(string),
			Enumerations: paramMap["enumeration"].(string),
			PassThrough:  isObjectEnabled(paramMap["passthrough"].(bool)),
		}
		switch paramType {
		case string(ParamTypeNumber):
			param.MaxNum = pointerto.Int(paramMap["maximum"].(int))
			param.MinNum = pointerto.Int(paramMap["minimum"].(int))
		case string(ParamTypeString):
			param.MaxSize = pointerto.Int(paramMap["maximum"].(int))
			param.MinSize = pointerto.Int(paramMap["minimum"].(int))
		}
		result[i] = param
	}
	return result
}

func buildBackendParameters(backends *schema.Set) ([]apis.BackendParams, error) {
	result := make([]apis.BackendParams, backends.Len())
	for i, v := range backends.List() {
		pm := v.(map[string]interface{})
		origin := pm["type"].(string)
		if origin == string(ParameterTypeSystem) && pm["system_param_type"].(string) == "" {
			return nil, fmt.Errorf("the 'system_param_type' must set if parameter type is 'SYSTEM'")
		}
		param := apis.BackendParams{
			Origin:   origin,
			Name:     pm["name"].(string),
			Location: pm["location"].(string),
			Value:    buildBackendParameterValue(origin, pm["value"].(string), pm["system_param_type"].(string)),
		}

		if origin != string(ParameterTypeRequest) {
			param.Description = pm["description"].(string)
		}
		result[i] = param
	}

	return result, nil
}

func buildBackendParameterValue(origin, value, paramAuthType string) string {
	internalParams := []string{
		"sourceIp", "stage", "apiId", "appId", "requestId", "serverAddr", "serverName", "handleTime", "providerAppId",
	}
	if origin == "SYSTEM" {
		if paramAuthType == string(SystemParamTypeFrontend) || paramAuthType == string(SystemParamTypeBackend) {
			return fmt.Sprintf("$context.authorizer.%s.%s", paramAuthType, value)
		}
		if common.StrSliceContains(internalParams, value) {
			return fmt.Sprintf("$context.%s", value)
		}
	}
	return value
}

func buildMockPolicy(policies *schema.Set) ([]apis.PolicyMocks, error) {
	if policies.Len() < 1 {
		return nil, nil
	}

	result := make([]apis.PolicyMocks, policies.Len())
	for i, policy := range policies.List() {
		pm := policy.(map[string]interface{})
		params, err := buildBackendParameters(pm["backend_params"].(*schema.Set))
		if err != nil {
			return nil, err
		}
		result[i] = apis.PolicyMocks{
			AuthorizerID:  pm["authorizer_id"].(string),
			Name:          pm["name"].(string),
			Response:      pm["response"].(string),
			EffectMode:    pm["effective_mode"].(string),
			Conditions:    buildPolicyConditions(pm["conditions"].(*schema.Set)),
			BackendParams: params,
		}
	}
	return result, nil
}

func buildPolicyConditions(conditions *schema.Set) []apis.Conditions {
	if conditions.Len() < 1 {
		return nil
	}
	result := make([]apis.Conditions, conditions.Len())
	for i, v := range conditions.List() {
		cm := v.(map[string]interface{})
		condition := apis.Conditions{
			ReqParamName:    cm["param_name"].(string),
			ConditionOrigin: cm["origin"].(string),
			ConditionValue:  cm["value"].(string),
		}
		conType := cm["type"].(string)
		if v, ok := conditionType[conType]; ok {
			condition.ConditionType = v
		}
		result[i] = condition
	}
	return result
}

func buildFuncGraphPolicy(policies *schema.Set) ([]apis.PolicyFunctions, error) {
	if policies.Len() < 1 {
		return nil, nil
	}

	result := make([]apis.PolicyFunctions, policies.Len())
	for i, policy := range policies.List() {
		pm := policy.(map[string]interface{})
		params, err := buildBackendParameters(pm["backend_params"].(*schema.Set))
		if err != nil {
			return nil, err
		}
		result[i] = apis.PolicyFunctions{
			FunctionUrn:    pm["function_urn"].(string),
			InvocationType: pm["invocation_mode"].(string),
			NetworkType:    buildNetworkType(pm["network_type"].(string)),
			Version:        pm["version"].(string),
			Timeout:        pm["timeout"].(int),
			EffectMode:     pm["effective_mode"].(string),
			Name:           pm["name"].(string),
			Conditions:     buildPolicyConditions(pm["conditions"].(*schema.Set)),
			BackendParams:  params,
			AuthorizerID:   pm["authorizer_id"].(string),
		}
	}
	return result, nil
}

func buildHttpPolicy(policies *schema.Set) ([]apis.PolicyHttps, error) {
	if policies.Len() < 1 {
		return nil, nil
	}

	result := make([]apis.PolicyHttps, policies.Len())
	for i, policy := range policies.List() {
		pm := policy.(map[string]interface{})
		params, err := buildBackendParameters(pm["backend_params"].(*schema.Set))
		if err != nil {
			return nil, err
		}
		wp := apis.PolicyHttps{
			UrlDomain:     pm["url_domain"].(string),
			ReqProtocol:   pm["request_protocol"].(string),
			ReqMethod:     pm["request_method"].(string),
			ReqUri:        pm["request_uri"].(string),
			Timeout:       pointerto.Int(pm["timeout"].(int)),
			RetryCount:    strconv.Itoa(pm["retry_count"].(int)),
			EffectMode:    pm["effective_mode"].(string),
			Name:          pm["name"].(string),
			BackendParams: params,
			Conditions:    buildPolicyConditions(pm["conditions"].(*schema.Set)),
			AuthorizerID:  pm["authorizer_id"].(string),
		}
		if chanId, ok := pm["vpc_channel_id"]; ok {
			if chanId != "" {
				wp.VpcChannelStatus = pointerto.Int(strBoolEnabled)
				wp.VpcChannelInfo = &apis.VpcChannelInfo{
					VpcChannelID:        pm["vpc_channel_id"].(string),
					VpcChannelProxyHost: pm["vpc_channel_proxy_host"].(string),
				}
			} else {
				wp.VpcChannelStatus = pointerto.Int(strBoolDisabled)
			}
		}
		result[i] = wp
	}
	return result, nil
}

func buildApiCreateOpts(d *schema.ResourceData) (apis.CreateOpts, error) {
	authType := d.Get("security_authentication_type").(string)
	opts := apis.CreateOpts{
		GatewayID:           d.Get("gateway_id").(string),
		GroupID:             d.Get("group_id").(string),
		Name:                d.Get("name").(string),
		Type:                buildApiType(d.Get("type").(string)),
		Version:             d.Get("version").(string),
		ReqProtocol:         d.Get("request_protocol").(string),
		ReqMethod:           d.Get("request_method").(string),
		ReqUri:              d.Get("request_uri").(string),
		AuthType:            authType,
		Cors:                d.Get("cors").(bool),
		Description:         d.Get("description").(string),
		BodyDescription:     d.Get("body_description").(string),
		ResultNormalSample:  d.Get("success_response").(string),
		ResultFailureSample: d.Get("failure_response").(string),
		AuthorizerID:        d.Get("authorizer_id").(string),
		ResponseID:          d.Get("response_id").(string),
		ReqParams:           buildRequestParameters(d.Get("request_params").(*schema.Set)),
	}
	tagsRaw := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(tagsRaw))
	for i, raw := range tagsRaw {
		tags[i] = raw.(string)
	}
	opts.Tags = tags
	// build match mode
	matchMode := d.Get("match_mode").(string)
	v, ok := matching[matchMode]
	if !ok {
		return opts, fmt.Errorf("invalid match mode: '%s'", matchMode)
	}
	opts.MatchMode = v

	isSimpleAuthEnabled := d.Get("security_authentication_enabled").(bool)
	if authType == string(ApiAuthTypeApp) {
		if isSimpleAuthEnabled {
			opts.AuthOpt = &apis.AuthOpt{
				AppCodeAuthType: string(AppCodeAuthTypeEnable),
			}
		} else {
			opts.AuthOpt = &apis.AuthOpt{
				AppCodeAuthType: string(AppCodeAuthTypeDisable),
			}
		}
	} else if isSimpleAuthEnabled {
		return opts, fmt.Errorf("the security authentication must be 'APP' if simple authentication is true")
	}

	if m, ok := d.GetOk("mock"); ok {
		opts.BackendType = string(BackendTypeMock)
		params, err := buildBackendParameters(d.Get("backend_params").(*schema.Set))
		if err != nil {
			return opts, err
		}
		opts.BackendParams = params
		opts.MockInfo = buildMockStructure(m.([]interface{}))
		policy, err := buildMockPolicy(d.Get("mock_policy").(*schema.Set))
		if err != nil {
			return opts, err
		}
		opts.PolicyMocks = policy
	} else if fg, ok := d.GetOk("func_graph"); ok {
		opts.BackendType = string(BackendTypeFunction)
		params, err := buildBackendParameters(d.Get("backend_params").(*schema.Set))
		if err != nil {
			return opts, err
		}
		opts.BackendParams = params
		opts.FuncInfo = buildFuncGraphStructure(fg.([]interface{}))
		policy, err := buildFuncGraphPolicy(d.Get("func_graph_policy").(*schema.Set))
		if err != nil {
			return opts, err
		}
		opts.PolicyFunctions = policy
	} else {
		opts.BackendType = string(BackendTypeHttp)
		params, err := buildBackendParameters(d.Get("backend_params").(*schema.Set))
		if err != nil {
			return opts, err
		}
		opts.BackendParams = params
		opts.BackendApi = buildHttpStructure(d.Get("http").([]interface{}))
		policy, err := buildHttpPolicy(d.Get("http_policy").(*schema.Set))
		if err != nil {
			return opts, err
		}
		opts.PolicyHttps = policy
	}

	log.Printf("[DEBUG] The API Opts is : %+v", opts)
	return opts, nil
}

func resourceAPIGWApiV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts, err := buildApiCreateOpts(d)
	if err != nil {
		return diag.Errorf("unable to build the OpenTelekomCloud APIGW API create opts: %s", err)
	}
	api, err := apis.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud APIGW API: %s", err)
	}
	d.SetId(api.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAPIGWApiV2Read(clientCtx, d, meta)
}

func resourceAPIGWApiV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	resp, err := apis.Get(client, d.Get("gateway_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud APIGW API")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("group_id", resp.GroupID),
		d.Set("name", resp.Name),
		d.Set("authorizer_id", resp.AuthorizerID),
		d.Set("request_protocol", resp.ReqProtocol),
		d.Set("request_method", resp.ReqMethod),
		d.Set("request_uri", resp.ReqUri),
		d.Set("security_authentication_type", resp.AuthType),
		d.Set("cors", resp.Cors),
		d.Set("description", resp.Description),
		d.Set("body_description", resp.BodyDescription),
		d.Set("success_response", resp.ResultNormalSample),
		d.Set("failure_response", resp.ResultFailureSample),
		d.Set("response_id", resp.ResponseID),
		d.Set("type", analyseApiType(resp.Type)),
		d.Set("request_params", flattenApiRequestParams(resp.ReqParams)),
		d.Set("backend_params", flattenBackendParameters(resp.BackendParams)),
		d.Set("match_mode", analyseApiMatchMode(resp.MatchMode)),
		d.Set("security_authentication_enabled", analyseAppSecurityAuth(resp.AuthOpt)),
		d.Set("mock", flattenMockStructure(resp.MockInfo)),
		d.Set("mock_policy", flattenMockPolicy(resp.PolicyMocks)),
		d.Set("func_graph", flattenFuncGraphStructure(resp.FuncInfo)),
		d.Set("func_graph_policy", flattenFuncGraphPolicy(resp.PolicyFunctions)),
		d.Set("http", flattenHttpStructure(resp.BackendApi, d.Get("http.0.ssl_enable").(bool))),
		d.Set("http_policy", flattenHttpPolicy(resp.PolicyHttps)),
		d.Set("registered_at", resp.RegisterTime),
		d.Set("updated_at", resp.UpdateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving  OpenTelekomCloud APIGW API fields: %s", err)
	}
	return nil
}

func resourceAPIGWApiV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts, err := buildApiCreateOpts(d)
	if err != nil {
		return diag.Errorf("unable to build the OpenTelekomCloud APIGW API updateOpts: %s", err)
	}
	_, err = apis.Update(client, d.Id(), opts)
	if err != nil {
		return diag.Errorf("error updating OpenTelekomCloud APIGW API (%s): %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAPIGWApiV2Read(clientCtx, d, meta)
}

func resourceAPIGWApiV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	if err = apis.Delete(client, d.Get("gateway_id").(string), d.Id()); err != nil {
		return diag.Errorf("unable to delete the OpenTelekomCloud APIGW API (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceAPIGWApiV2ImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf(errCreationV2Client, err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}
	name := parts[1]
	gatewayId := parts[0]
	apiId, err := GetApiIdByName(client, gatewayId, name)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(apiId)
	err = d.Set("gateway_id", gatewayId)
	if err != nil {
		return nil, fmt.Errorf("error setting OpenTelekomCloud APIGW gateway_id attribute: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

// GetApiIdByName is a method to get a specifies API ID from a APIGW instance by name.
func GetApiIdByName(client *golangsdk.ServiceClient, gatewayId, name string) (string, error) {
	opts := apis.ListOpts{
		GatewayID: gatewayId,
		Name:      name,
	}
	apiList, err := apis.List(client, opts)
	if err != nil {
		return "", fmt.Errorf("error retrieving OpenTelekomCloud APIGW APIs: %s", err)
	}
	if len(apiList) < 1 {
		return "", fmt.Errorf("unable to find the OpenTelekomCloud APIGW API (%s) from cloud: %s", name, err)
	}
	return apiList[0].ID, nil
}

func analyseBackendParameterValue(origin, value string) (paramType, paramValue string) {
	log.Printf("[ERROR] The value of the backend parameter is: %s", value)
	if origin == string(ParameterTypeSystem) {
		regex := regexp.MustCompile(`\$context\.authorizer\.(frontend|backend)\.([\w-]+)`)
		result := regex.FindStringSubmatch(value)
		if len(result) == 3 {
			paramType = result[1]
			paramValue = result[2]
			return
		}

		regex = regexp.MustCompile(`\$context\.([\w-]+)`)
		result = regex.FindStringSubmatch(value)
		if len(result) == 2 {
			paramType = string(SystemParamTypeInternal)
			paramValue = result[1]
			return
		}
		log.Printf("[ERROR] The system parameter format is invalid, want '$context.xxx' (internal parameter), "+
			"'$context.authorizer.frontend.xxx' or '$context.authorizer.frontend.xxx', but '%s'.", value)
		return
	}
	paramValue = value
	return
}

func flattenBackendParameters(backendParams []apis.BackendParams) []map[string]interface{} {
	if len(backendParams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(backendParams))
	for i, v := range backendParams {
		origin := v.Origin
		paramAuthType, paramValue := analyseBackendParameterValue(v.Origin, v.Value)
		param := map[string]interface{}{
			"type":     origin,
			"name":     v.Name,
			"location": v.Location,
			"value":    paramValue,
		}
		if paramAuthType != "" {
			param["system_param_type"] = paramAuthType
		}
		if origin != string(ParameterTypeRequest) {
			param["description"] = v.Description
		}
		result[i] = param
	}
	return result
}

func analyseConditionType(conType string) string {
	for k, v := range conditionType {
		if v == conType {
			return k
		}
	}
	return ""
}

func analyseApiType(t int) string {
	apiType := map[int]string{
		1: "Public",
		2: "Private",
	}
	if v, ok := apiType[t]; ok {
		return v
	}
	return ""
}

func analyseNetworkType(t string) string {
	networkType := map[string]string{
		"NON-VPC": "V1",
		"VPC":     "V2",
	}
	if v, ok := networkType[t]; ok {
		return v
	}
	return ""
}

func analyseApiMatchMode(mode string) string {
	for k, v := range matching {
		if v == mode {
			return k
		}
	}
	return ""
}

func analyseAppSecurityAuth(opt *apis.AuthOpt) bool {
	return opt.AppCodeAuthType == "HEADER"
}

func parseObjectEnabled(objStatus *int) bool {
	if objStatus == pointerto.Int(strBoolEnabled) {
		return true
	}
	if objStatus != pointerto.Int(strBoolDisabled) {
		log.Printf("[DEBUG] unexpected object value, want '1'(yes) or '2'(no), but got '%d'", objStatus)
	}
	return false
}

func flattenApiRequestParams(reqParams []apis.ReqParams) []map[string]interface{} {
	if len(reqParams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(reqParams))
	for i, v := range reqParams {
		param := map[string]interface{}{
			"name":        v.Name,
			"location":    v.Location,
			"type":        v.Type,
			"required":    parseObjectEnabled(v.Required),
			"passthrough": parseObjectEnabled(v.PassThrough),
			"enumeration": v.Enumerations,
			"sample":      v.SampleValue,
			"default":     v.DefaultValue,
			"description": v.Description,
		}
		switch v.Type {
		case string(ParamTypeNumber):
			param["maximum"] = v.MaxNum
			param["minimum"] = v.MinNum
		case string(ParamTypeString):
			param["maximum"] = v.MaxSize
			param["minimum"] = v.MinSize
		}
		result[i] = param
	}
	return result
}

func flattenMockStructure(mockResp *apis.MockInfo) []map[string]interface{} {
	if mockResp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"response":      mockResp.Response,
			"authorizer_id": mockResp.AuthorizerID,
			"version":       mockResp.Version,
			"description":   mockResp.Description,
		},
	}
}

func flattenFuncGraphStructure(funcResp *apis.FuncInfo) []map[string]interface{} {
	if funcResp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"function_urn":    funcResp.FunctionUrn,
			"timeout":         funcResp.Timeout,
			"invocation_type": funcResp.InvocationType,
			"version":         funcResp.Version,
			"authorizer_id":   funcResp.AuthorizerID,
			"description":     funcResp.Description,
			"network_type":    analyseNetworkType(funcResp.NetworkType),
		},
	}
}

func flattenHttpStructure(webResp *apis.BackendApi, sslEnabled bool) []map[string]interface{} {
	if webResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"request_uri":      webResp.ReqUri,
		"request_method":   webResp.ReqMethod,
		"request_protocol": webResp.ReqProtocol,
		"timeout":          webResp.Timeout,
		"ssl_enable":       sslEnabled,
		"authorizer_id":    webResp.AuthorizerID,
		"retry_count":      common.StringToInt(&webResp.RetryCount),
	}
	if webResp.VpcChannelInfo.VpcChannelID != "" {
		result["vpc_channel_id"] = webResp.VpcChannelInfo.VpcChannelID
		result["vpc_channel_proxy_host"] = webResp.VpcChannelInfo.VpcChannelProxyHost
	} else {
		result["url_domain"] = webResp.UrlDomain
	}

	return []map[string]interface{}{
		result,
	}
}

func flattenPolicyConditions(conditions []apis.Conditions) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(conditions))
	for i, v := range conditions {
		result[i] = map[string]interface{}{
			"origin":     v.ConditionOrigin,
			"param_name": v.ReqParamName,
			"type":       analyseConditionType(v.ConditionType),
			"value":      v.ConditionValue,
		}
	}
	return result
}

func flattenMockPolicy(policies []apis.PolicyMocks) []map[string]interface{} {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		result[i] = map[string]interface{}{
			"name":           policy.Name,
			"response":       policy.Response,
			"effective_mode": policy.EffectMode,
			"authorizer_id":  policy.AuthorizerID,
			"backend_params": flattenBackendParameters(policy.BackendParams),
			"conditions":     flattenPolicyConditions(policy.Conditions),
		}
	}

	return result
}

func flattenFuncGraphPolicy(policies []apis.PolicyFunctions) []map[string]interface{} {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		result[i] = map[string]interface{}{
			"name":            policy.Name,
			"function_urn":    policy.FunctionUrn,
			"version":         policy.Version,
			"invocation_mode": policy.InvocationType,
			"effective_mode":  policy.EffectMode,
			"timeout":         policy.Timeout,
			"authorizer_id":   policy.AuthorizerID,
			"backend_params":  flattenBackendParameters(policy.BackendParams),
			"conditions":      flattenPolicyConditions(policy.Conditions),
		}
	}

	return result
}

func flattenHttpPolicy(policies []apis.PolicyHttps) []map[string]interface{} {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		retryCount := policy.RetryCount
		wp := map[string]interface{}{
			"url_domain":       policy.UrlDomain,
			"name":             policy.Name,
			"request_protocol": policy.ReqProtocol,
			"request_method":   policy.ReqMethod,
			"effective_mode":   policy.EffectMode,
			"request_uri":      policy.ReqUri,
			"timeout":          policy.Timeout,
			"retry_count":      common.StringToInt(&retryCount),
			"authorizer_id":    policy.AuthorizerID,
			"backend_params":   flattenBackendParameters(policy.BackendParams),
			"conditions":       flattenPolicyConditions(policy.Conditions),
		}
		if policy.VpcChannelInfo.VpcChannelID != "" {
			wp["vpc_channel_id"] = policy.VpcChannelInfo.VpcChannelID
			wp["vpc_channel_proxy_host"] = policy.VpcChannelInfo.VpcChannelProxyHost
		} else {
			wp["url_domain"] = policy.UrlDomain
		}

		result[i] = wp
	}

	return result
}
