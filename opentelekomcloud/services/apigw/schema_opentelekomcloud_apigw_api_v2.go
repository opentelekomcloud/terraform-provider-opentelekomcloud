package apigw

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func backendParamSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ParameterTypeRequest),
					string(ParameterTypeConstant),
					string(ParameterTypeSystem),
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z][\w-.]*$`),
						"Only letters, digits, hyphens (-), underscores (_) and periods (.) are allowed, and must "+
							"start with a letter."),
					validation.StringLenBetween(1, 32),
				),
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ParamLocationPath),
					string(ParamLocationQuery),
					string(ParamLocationHeader),
				}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[^<>]*$"),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 255),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[^<>]*$"),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 255),
				),
			},
			"system_param_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(SystemParamTypeInternal),
					string(SystemParamTypeFrontend),
					string(SystemParamTypeBackend),
				}, false),
			},
		},
	}
}

func policyConditionSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"param_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ConditionTypeEqual),
				ValidateFunc: validation.StringInSlice([]string{
					string(ConditionTypeEqual),
					string(ConditionTypeEnumerated),
					string(ConditionTypeMatching),
				}, false),
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ConditionSourceParam),
				ValidateFunc: validation.StringInSlice([]string{
					string(ConditionSourceParam),
					string(ConditionSourceSource),
				}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func requestParamsSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z][\w-.]*$`),
						"Only letters, digits, hyphens (-), underscores (_) and periods (.) are allowed, "+
							"and must start with a letter."),
					validation.StringLenBetween(1, 32),
				),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  ParamTypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(ParamTypeString), string(ParamTypeNumber),
				}, false),
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  ParamLocationPath,
				ValidateFunc: validation.StringInSlice([]string{
					string(ParamLocationPath), string(ParamLocationQuery), string(ParamLocationHeader),
				}, false),
			},
			"default": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 255),
				),
			},
			"sample": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 255),
				),
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"validity_check": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"enumeration": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maximum": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"minimum": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"passthrough": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func mockSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"response": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 2048),
			},
			"authorizer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func funcGraphSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"function_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5000,
				ValidateFunc: validation.IntBetween(1, 600000),
			},
			"invocation_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  InvocationTypeSync,
				ValidateFunc: validation.StringInSlice([]string{
					string(InvocationTypeAsync), string(InvocationTypeSync),
				}, false),
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  NetworkTypeV1,
				ValidateFunc: validation.StringInSlice([]string{
					string(NetworkTypeV2), string(NetworkTypeV1),
				}, false),
			},
			"authorizer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func mockPolicySchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"effective_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(EffectiveModeAny),
				ValidateFunc: validation.StringInSlice([]string{
					string(EffectiveModeAll),
					string(EffectiveModeAny),
				}, false),
			},
			"response": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(8, 2048),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z]\w*$`),
						"Only letters, digits and underscores (_) are allowed, and start with a letter."),
					validation.StringLenBetween(3, 64),
				),
			},
			"backend_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     backendParamSchemaResource(),
				Set:      resourceParametersHash,
			},
			"conditions": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 5,
				Elem:     policyConditionSchemaResource(),
			},
			"authorizer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func functionPolicySchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"function_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"invocation_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  InvocationTypeSync,
				ValidateFunc: validation.StringInSlice([]string{
					string(InvocationTypeAsync),
					string(InvocationTypeSync),
				}, false),
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  NetworkTypeV1,
				ValidateFunc: validation.StringInSlice([]string{
					string(NetworkTypeV1),
					string(NetworkTypeV2),
				}, false),
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5000,
				ValidateFunc: validation.IntBetween(1, 600000),
			},
			"effective_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(EffectiveModeAny),
				ValidateFunc: validation.StringInSlice([]string{
					string(EffectiveModeAll),
					string(EffectiveModeAny),
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z]\w*$`),
						"Only letters, digits and underscores (_) are allowed, and must start with a letter."),
					validation.StringLenBetween(3, 64),
				),
			},
			"backend_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     backendParamSchemaResource(),
				Set:      resourceParametersHash,
			},
			"conditions": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 5,
				Elem:     policyConditionSchemaResource(),
			},
			"authorizer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func httpSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"authorizer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ProtocolTypeHTTPS),
				ValidateFunc: validation.StringInSlice([]string{
					string(ProtocolTypeHTTP),
					string(ProtocolTypeHTTPS),
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(RequestMethodGet),
					string(RequestMethodPost),
					string(RequestMethodPut),
					string(RequestMethodDelete),
					string(RequestMethodHead),
					string(RequestMethodPatch),
					string(RequestMethodOptions),
					string(RequestMethodAny),
				}, false),
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5000,
				ValidateFunc: validation.IntBetween(1, 600000),
			},
			"ssl_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"retry_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"vpc_channel_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_channel_proxy_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func httpPolicySchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ProtocolTypeHTTP),
					string(ProtocolTypeHTTPS),
				}, false),
			},
			"request_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(RequestMethodGet),
					string(RequestMethodPost),
					string(RequestMethodPut),
					string(RequestMethodDelete),
					string(RequestMethodHead),
					string(RequestMethodPatch),
					string(RequestMethodOptions),
					string(RequestMethodAny),
				}, false),
			},
			"request_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5000,
				ValidateFunc: validation.IntBetween(1, 600000),
			},
			"retry_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"effective_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(EffectiveModeAny),
				ValidateFunc: validation.StringInSlice([]string{
					string(EffectiveModeAll),
					string(EffectiveModeAny),
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z]\w*$`),
						"Only letters, digits and underscores (_) are allowed, and must start with a letter."),
					validation.StringLenBetween(3, 64),
				),
			},
			"backend_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     backendParamSchemaResource(),
				Set:      resourceParametersHash,
			},
			"conditions": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 5,
				Elem:     policyConditionSchemaResource(),
			},
			"authorizer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_channel_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_channel_proxy_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceParametersHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["type"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
	}
	if m["name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	}

	return hashcode.String(buf.String())
}
