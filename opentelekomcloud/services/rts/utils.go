package rts

import (
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/stacks"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
)

func normalizeStackTemplate(templateString interface{}) (string, error) {
	if common.LooksLikeJsonString(templateString) {
		return common.NormalizeJsonString(templateString.(string))
	}

	return common.CheckYamlString(templateString)
}

func flattenStackOutputs(stackOutputs []*stacks.Output) map[string]string {
	outputs := make(map[string]string, len(stackOutputs))
	for _, o := range stackOutputs {
		outputs[*o.OutputKey] = *o.OutputValue
	}
	return outputs
}

// flattenStackParameters is flattening list of
//
//	stack Parameters and only returning existing
//
// parameters to avoid clash with default values
func flattenStackParameters(stackParams map[string]string,
	originalParams map[string]interface{}) map[string]string {
	params := make(map[string]string, len(stackParams))
	for key, value := range stackParams {
		_, isConfigured := originalParams[key]
		if isConfigured {
			params[key] = value
		}
	}
	return params
}
