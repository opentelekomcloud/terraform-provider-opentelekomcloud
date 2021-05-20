package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"

	ver "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func jsonBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

// ConvertStructToMap converts an instance of struct to a map object, and
// changes each key of fields to the value of 'nameMap' if the key in it
// or to its corresponding lowercase.
func ConvertStructToMap(obj interface{}, nameMap map[string]string) (map[string]interface{}, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("Error converting struct to map, marshal failed:%v", err)
	}

	m, err := regexp.Compile(`"[a-z0-9A-Z_]+":`)
	if err != nil {
		return nil, fmt.Errorf("Error converting struct to map, compile regular express failed")
	}
	nb := m.ReplaceAllFunc(
		b,
		func(src []byte) []byte {
			k := fmt.Sprintf("%s", src[1:len(src)-2])
			v, ok := nameMap[k]
			if !ok {
				v = strings.ToLower(k)
			}
			return []byte(fmt.Sprintf("\"%s\":", v))
		},
	)
	log.Printf("[DEBUG]ConvertStructToMap:: before change b =%s", b)
	log.Printf("[DEBUG]ConvertStructToMap:: after change nb=%s", nb)

	p := make(map[string]interface{})
	err = json.Unmarshal(nb, &p)
	if err != nil {
		return nil, fmt.Errorf("Error converting struct to map, unmarshal failed:%v", err)
	}
	log.Printf("[DEBUG]ConvertStructToMap:: map= %#v\n", p)
	return p, nil
}

func LooksLikeJsonString(s interface{}) bool {
	return regexp.MustCompile(`^\s*{`).MatchString(s.(string))
}

func Base64IfNot(src string) string {
	_, err := base64.StdEncoding.DecodeString(src)
	if err == nil {
		return src
	}
	return base64.StdEncoding.EncodeToString([]byte(src))
}

type versionSlice []*ver.Version

func (v versionSlice) Len() int {
	return len(v)
}

func (v versionSlice) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v versionSlice) Less(i, j int) bool {
	return v[i].LessThan(v[j])
}

func (v versionSlice) ToStringSlice() []string {
	res := make([]string, len(v))
	for i, version := range v {
		res[i] = version.Original()
	}
	return res
}

func sortAsStringSlice(src []string) []string {
	res := make([]string, len(src))
	copy(res, src)
	sort.Sort(sort.Reverse(sort.StringSlice(res)))
	return res
}

// SortVersions sorts versions from newer to older.
// If non-version-like string will be found in the slice,
// slice will be sorted as string slice in reversed order (z-a)
func SortVersions(src []string) []string {
	verSlice := make(versionSlice, len(src))
	for i, v := range src {
		val, err := ver.NewVersion(v)
		if err != nil {
			return sortAsStringSlice(src) // in case it's not version-like
		}
		verSlice[i] = val
	}
	sort.Sort(sort.Reverse(verSlice))
	return verSlice.ToStringSlice()
}

// BuildRequest takes an opts struct and builds a request body for
// Gophercloud to execute
func BuildRequest(opts interface{}, parent string) (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	b = AddValueSpecs(b)
	return map[string]interface{}{parent: b}, nil
}

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, msg string) error {
	_, ok := err.(golangsdk.ErrDefault404)
	if ok {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("%s: %s", msg, err)
}

// AddValueSpecs expands the 'value_specs' object and removes 'value_specs'
// from the request body.
func AddValueSpecs(body map[string]interface{}) map[string]interface{} {
	if body["value_specs"] != nil {
		for k, v := range body["value_specs"].(map[string]interface{}) {
			body[k] = v
		}
		delete(body, "value_specs")
	}

	return body
}

// MapValueSpecs converts ResourceData into a map
func MapValueSpecs(d cfg.SchemaOrDiff) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("value_specs").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

// MapResourceProp converts ResourceData property into a map
func MapResourceProp(d *schema.ResourceData, prop string) map[string]interface{} {
	m := make(map[string]interface{})
	for key, val := range d.Get(prop).(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func CheckForRetryableError(err error) *resource.RetryError {
	switch err.(type) {
	case golangsdk.ErrDefault409, golangsdk.ErrDefault500, golangsdk.ErrDefault503:
		return resource.RetryableError(err)
	default:
		return resource.NonRetryableError(err)
	}
}

func IsResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(golangsdk.ErrDefault404)
	return ok
}

func ExpandToStringSlice(v []interface{}) []string {
	s := make([]string, len(v))
	for i, val := range v {
		if strVal, ok := val.(string); ok {
			s[i] = strVal
		}
	}

	return s
}

// StrSliceContains checks if a given string is contained in a slice
// When anybody asks why Go needs generics, here you go.
func StrSliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func GetAllAvailableZones(d *schema.ResourceData) []string {
	rawZones := d.Get("available_zones").([]interface{})
	zones := make([]string, len(rawZones))
	for i, raw := range rawZones {
		zones[i] = raw.(string)
	}
	log.Printf("[DEBUG] getAvailableZones: %#v", zones)

	return zones
}

func StringInSlice(str string, slice []string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func BuildComponentID(parts ...string) string {
	return strings.Join(parts, "/")
}
