package common

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var SuccessHTTPCodes = []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}

func IsEmptyValue(v reflect.Value) (bool, error) {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0, nil
	case reflect.Bool:
		return !v.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0, nil
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0, nil
	case reflect.Interface, reflect.Ptr:
		return v.IsNil(), nil
	case reflect.Invalid:
		return true, nil
	}
	return false, fmt.Errorf("isEmptyValue:: unknown type")
}

func ReplaceVars(d *schema.ResourceData, linkTmpl string, kv map[string]string) (string, error) {
	re := regexp.MustCompile("{([[:word:]]+)}")

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if kv != nil {
			if v, ok := kv[m]; ok {
				return v
			}
		}
		if m == "project" {
			return "replace_holder"
		}
		if d != nil {
			if m == "id" {
				return d.Id()
			}
			v, ok := d.GetOk(m)
			if ok {
				v1, _ := convertToStr(v)
				return v1
			}
		}
		return ""
	}

	s := re.ReplaceAllStringFunc(linkTmpl, replaceFunc)
	return strings.Replace(s, "replace_holder/", "", 1), nil
}

func ReplaceVarsForTest(rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{([[:word:]]+)}")

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return "replace_holder"
		}
		if rs != nil {
			if m == "id" {
				return rs.Primary.ID
			}
			v, ok := rs.Primary.Attributes[m]
			if ok {
				return v
			}
		}
		return ""
	}

	s := re.ReplaceAllStringFunc(linkTmpl, replaceFunc)
	return strings.Replace(s, "replace_holder/", "", 1), nil
}

func NavigateValue(d interface{}, index []string, arrayIndex map[string]int) (interface{}, error) {
	for n, i := range index {
		if d == nil {
			return nil, nil
		}
		if d1, ok := d.(map[string]interface{}); ok {
			d, ok = d1[i]
			if !ok {
				msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
				return nil, fmt.Errorf("%s: '%s' may not exist", msg, i)
			}
		} else {
			msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
			return nil, fmt.Errorf("%s: Can not convert (%s) to map", msg, reflect.TypeOf(d))
		}

		if arrayIndex != nil {
			if j, ok := arrayIndex[strings.Join(index[:n+1], ".")]; ok {
				if d == nil {
					return nil, nil
				}
				if d2, ok := d.([]interface{}); ok {
					if len(d2) == 0 {
						return nil, nil
					}
					if j >= len(d2) {
						msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
						return nil, fmt.Errorf("%s: The index is out of array", msg)
					}

					d = d2[j]
				} else {
					msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
					return nil, fmt.Errorf("%s: Can not convert (%s) to array, index=%s.%v", msg, reflect.TypeOf(d), i, j)
				}
			}
		}
	}

	return d, nil
}

func convertToStr(v interface{}) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	} else if i, ok := v.(int); ok {
		return strconv.Itoa(i), nil
	} else if b, ok := v.(bool); ok {
		return strconv.FormatBool(b), nil
	}

	return "", fmt.Errorf("can't convert to string")
}
