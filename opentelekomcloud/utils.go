package opentelekomcloud

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

// convertStructToMap converts an instance of struct to a map object, and
// changes each key of fileds to the value of 'nameMap' if the key in it
// or to its corresponding lowercase.
func convertStructToMap(obj interface{}, nameMap map[string]string) (map[string]interface{}, error) {
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
	log.Printf("[DEBUG]convertStructToMap:: before change b =%s", b)
	log.Printf("[DEBUG]convertStructToMap:: after change nb=%s", nb)

	p := make(map[string]interface{})
	err = json.Unmarshal(nb, &p)
	if err != nil {
		return nil, fmt.Errorf("Error converting struct to map, unmarshal failed:%v", err)
	}
	log.Printf("[DEBUG]convertStructToMap:: map= %#v\n", p)
	return p, nil
}

func looksLikeJsonString(s interface{}) bool {
	return regexp.MustCompile(`^\s*{`).MatchString(s.(string))
}

func base64IfNot(src string) string {
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
		res[i] = version.String()
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
