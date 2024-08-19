package common

import (
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awspolicy "github.com/jen20/awspolicyequivalence"
)

func SuppressEquivalentAwsPolicyDiffs(_, old, new string, _ *schema.ResourceData) bool {
	equivalent, err := awspolicy.PoliciesAreEquivalent(old, new)
	if err != nil {
		return false
	}

	return equivalent
}

// SuppressDiffAll suppress all changes?
func SuppressDiffAll(_, _, _ string, _ *schema.ResourceData) bool {
	return true
}

// SuppressMinDisk suppress changes if we get a computed min_disk_gb if value is unspecified (default 0)
func SuppressMinDisk(_, old, new string, _ *schema.ResourceData) bool {
	return new == "0" || old == new
}

// SuppressExternalGateway suppress changes if we don't specify an external gateway, but one is specified for us
func SuppressExternalGateway(_, old, new string, _ *schema.ResourceData) bool {
	return new == "" || old == new
}

// SuppressComputedFixedWhenFloatingIp suppress changes if we get a fixed ip when not expecting one,
// if we have a floating ip (generates fixed ip).
func SuppressComputedFixedWhenFloatingIp(_, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("floating_ip"); ok && v != "" {
		return new == "" || old == new
	}
	return false
}

func SuppressRdsNameDiffs(_, old, new string, _ *schema.ResourceData) bool {
	if strings.HasPrefix(old, new) && strings.HasSuffix(old, "_node0") {
		return true
	}
	return false
}

func SuppressStringSepratedByCommaDiffs(_, old, new string, _ *schema.ResourceData) bool {
	if len(old) != len(new) {
		return false
	}
	oldArray := strings.Split(old, ",")
	newArray := strings.Split(new, ",")
	sort.Strings(oldArray)
	sort.Strings(newArray)

	return reflect.DeepEqual(oldArray, newArray)
}

func SuppressSmartVersionDiff(_, old, new string, _ *schema.ResourceData) bool {
	compiledVer := regexp.MustCompile(`v(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:-(\w+))?$`)
	oldArray := compiledVer.FindStringSubmatch(old)
	newArray := compiledVer.FindStringSubmatch(new)
	if oldArray == nil {
		return false
	}
	for i := 1; i < len(newArray); i++ {
		if oldArray[i] == "" || newArray[i] == "" {
			return true
		}
		if oldArray[i] != newArray[i] {
			return false
		}
	}
	return false
}

func SuppressCaseInsensitive(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

func SuppressEqualZoneNames(_, old, new string, _ *schema.ResourceData) bool {
	oldShort := strings.TrimSuffix(old, ".")
	newShort := strings.TrimSuffix(new, ".")
	return oldShort == newShort
}

func SuppressStrippedNewLines(_, old, new string, _ *schema.ResourceData) bool {
	newline := "\n"
	return strings.Trim(old, newline) == strings.Trim(new, newline)
}

func SuppressEmptyStringSHA(k, old, new string, d *schema.ResourceData) bool {
	// Sometimes the API responds with the equivalent, empty SHA1 sum
	// echo -n "" | shasum
	if (old == "da39a3ee5e6b4b0d3255bfef95601890afd80709" && new == "") ||
		(old == "" && new == "da39a3ee5e6b4b0d3255bfef95601890afd80709") {
		return true
	}
	return false
}
