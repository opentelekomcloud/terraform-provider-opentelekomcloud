package opentelekomcloud

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/huaweicloud/golangsdk/openstack/rts/v1/stacks"
	"gopkg.in/yaml.v2"
)

// Takes list of pointers to strings. Expand to an array
// of raw strings and returns a []interface{}
// to keep compatibility w/ schema.NewSetschema.NewSet
func flattenStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}

func pointersMapToStringList(pointers map[string]*string) map[string]interface{} {
	list := make(map[string]interface{}, len(pointers))
	for i, v := range pointers {
		list[i] = *v
	}
	return list
}

// a convenience wrapper type for the schema.Set map[string]interface{}
// Set operations only alter the underlying map if the value is not nil
type setMap map[string]interface{}

// SetString sets m[key] = *value only if `value != nil`
func (s setMap) SetString(key string, value *string) {
	if value == nil {
		return
	}

	s[key] = *value
}

// SetStringMap sets key to value as a map[string]interface{}, stripping any nil
// values. The value parameter can be a map[string]interface{}, a
// map[string]*string, or a map[string]string.
func (s setMap) SetStringMap(key string, value interface{}) {
	// because these methods are meant to be chained without intermediate
	// checks for nil, we are likely to get interfaces with dynamic types but
	// a nil value.
	if reflect.ValueOf(value).IsNil() {
		return
	}

	m := make(map[string]interface{})

	switch value := value.(type) {
	case map[string]string:
		for k, v := range value {
			m[k] = v
		}
	case map[string]*string:
		for k, v := range value {
			if v == nil {
				continue
			}
			m[k] = *v
		}
	case map[string]interface{}:
		for k, v := range value {
			if v == nil {
				continue
			}

			switch v := v.(type) {
			case string:
				m[k] = v
			case *string:
				if v != nil {
					m[k] = *v
				}
			default:
				panic(fmt.Sprintf("unknown type for SetString: %T", v))
			}
		}
	}

	// catch the case where the interface wasn't nil, but we had no non-nil values
	if len(m) > 0 {
		s[key] = m
	}
}

// Set assigns value to s[key] if value isn't nil
func (s setMap) Set(key string, value interface{}) {
	if reflect.ValueOf(value).IsNil() {
		return
	}

	s[key] = value
}

// Map returns the raw map type for a shorter type conversion
func (s setMap) Map() map[string]interface{} {
	return map[string]interface{}(s)
}

// MapList returns the map[string]interface{} as a single element in a slice to
// match the schema.Set data type used for structs.
func (s setMap) MapList() []map[string]interface{} {
	return []map[string]interface{}{s.Map()}
}

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func normalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// The error is intentionally ignored here to allow empty policies to passthrough validation.
	// This covers any interpolated values
	bytes, _ := json.Marshal(j)

	return string(bytes[:]), nil
}

// Takes a value containing YAML string and passes it through
// the YAML parser. Returns either a parsing
// error or original YAML string.
func checkYamlString(yamlString interface{}) (string, error) {
	var y interface{}

	if yamlString == nil || yamlString.(string) == "" {
		return "", nil
	}

	s := yamlString.(string)

	err := yaml.Unmarshal([]byte(s), &y)
	if err != nil {
		return s, err
	}

	return s, nil
}

func normalizeStackTemplate(templateString interface{}) (string, error) {
	if looksLikeJsonString(templateString) {
		return normalizeJsonString(templateString.(string))
	}

	return checkYamlString(templateString)
}
func flattenStackOutputs(stackOutputs []*stacks.Output) map[string]string {
	outputs := make(map[string]string, len(stackOutputs))
	for _, o := range stackOutputs {
		outputs[*o.OutputKey] = *o.OutputValue
	}
	return outputs
}

// flattenStackParameters is flattening list of
//  stack Parameters and only returning existing
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

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}
