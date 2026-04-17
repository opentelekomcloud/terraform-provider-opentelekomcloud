package common

import (
	"fmt"
	"log"
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

// suppressEntireMapDiff handles changes to the entire map
func suppressEntireMapDiff(paramKey string, d *schema.ResourceData) bool {
	originMapCategory := fmt.Sprintf("%s_origin", paramKey)
	log.Printf("[DEBUG][EntireMapDiff] Handling entire map change for '%s', origin map category: '%s'",
		paramKey, originMapCategory)

	originMapVal := d.Get(originMapCategory)
	if originMapVal == nil {
		log.Printf("[DEBUG][EntireMapDiff] Origin map '%s' is nil, suppressing diff for entire map '%s'",
			originMapCategory, paramKey)
		return true
	}

	originMap := TryMapValueAnalysis(originMapVal)
	if len(originMap) == 0 {
		log.Printf("[DEBUG][EntireMapDiff] Origin map '%s' is empty, suppressing diff for entire map '%s'",
			originMapCategory, paramKey)
		return true
	}

	// For entire map changes, we need to check if all keys in the new value exist in origin
	// This is a simplified approach - in practice, you might want more sophisticated comparison
	log.Printf("[DEBUG][EntireMapDiff] Entire map '%s' change detected, checking against origin", paramKey)
	return false // For now, report the change and let individual key suppression handle it
}

// keyExists checks if a key exists in the map
func keyExists(m map[string]interface{}, key string) bool {
	_, exists := m[key]
	return exists
}

// determineSuppression determines whether to suppress the diff based on key existence
func determineSuppression(existsInCurrent, existsInOrigin, isOriginEmpty bool, keyName string) bool {
	if existsInCurrent {
		// The key exists in current configuration
		if isOriginEmpty {
			// Origin is empty or nil, report the change (locally added)
			log.Printf("[DEBUG] The key '%s' found in current config but origin is empty", keyName)
			return false
		}

		if existsInOrigin {
			// The key exists in both current config and origin, report this change
			log.Printf("[DEBUG] The key '%s' found in both current config and origin", keyName)
			return false
		}

		// The key exists in current config but not in origin (locally added)
		log.Printf("[DEBUG] The key '%s' found in current config but not in origin", keyName)
		return false
	}

	// The key does not exist in current configuration
	if isOriginEmpty {
		// Origin is empty or nil, suppress the diff (remotely added)
		log.Printf("[DEBUG] The key '%s' not found in current config and origin is empty, suppressing diff",
			keyName)
		return true
	}

	if existsInOrigin {
		// The key exists in origin but not in current config (locally removed)
		log.Printf("[DEBUG] The key '%s' found in origin but not in current config (locally removed)",
			keyName)
		return false
	}

	// The key does not exist in either current config or origin (remotely added)
	log.Printf("[DEBUG] The key '%s' not found in either current config or origin (remotely added), suppressing diff",
		keyName)
	return true
}

// suppressSingleKeyDiff handles changes to a single key
func suppressSingleKeyDiff(paramKey string, d *schema.ResourceData) bool {
	categories := strings.Split(paramKey, ".")
	mapCategory := strings.Join(categories[:len(categories)-1], ".")
	originMapCategory := fmt.Sprintf("%s_origin", mapCategory)
	keyName := categories[len(categories)-1]

	log.Printf("[DEBUG][SingleKeyDiff] Processing key '%s', mapCategory='%s', originMapCategory='%s', keyName='%s'",
		paramKey, mapCategory, originMapCategory, keyName)

	// Get origin map (last local configuration)
	originMapVal := d.Get(originMapCategory)
	originMap := TryMapValueAnalysis(originMapVal)
	log.Printf("[DEBUG][SingleKeyDiff] Origin map '%s' content: %+v", originMapCategory, originMap)

	// Get current configuration map
	currentMapVal := GetNestedObjectFromRawConfig(d.GetRawConfig(), mapCategory)
	if currentMapVal == nil {
		log.Printf("[DEBUG][SingleKeyDiff] Current map '%s' is nil, suppressing diff for key '%s'",
			mapCategory, keyName)
		return true
	}

	currentMap := TryMapValueAnalysis(currentMapVal)
	log.Printf("[DEBUG][SingleKeyDiff] Current map '%s' content: %+v", mapCategory, currentMap)

	// Check if the key exists in current configuration
	existsInCurrent := keyExists(currentMap, keyName)
	existsInOrigin := keyExists(originMap, keyName)
	isOriginEmpty := originMapVal == nil || len(originMap) == 0

	// Determine suppression based on key existence
	return determineSuppression(existsInCurrent, existsInOrigin, isOriginEmpty, keyName)
}

func SuppressMapDiffs() schema.SchemaDiffSuppressFunc {
	return func(paramKey, o, n string, d *schema.ResourceData) bool {
		log.Printf("[DEBUG][SuppressMapDiffs] Called with paramKey='%s', old='%s', new='%s'", paramKey, o, n)

		// Ignore length change judgment, because this method will judge each changed key one by one
		if strings.HasSuffix(paramKey, ".%") {
			log.Printf("[DEBUG][SuppressMapDiffs] Ignoring length change for '%s'", paramKey)
			return true
		}

		// Handle the case where the entire map is being changed
		if !strings.Contains(paramKey, ".") {
			return suppressEntireMapDiff(paramKey, d)
		}

		// Handle single key changes
		return suppressSingleKeyDiff(paramKey, d)
	}
}
