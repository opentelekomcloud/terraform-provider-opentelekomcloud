package opentelekomcloud

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ValidateStringList(v interface{}, k string, l []string) (ws []string, errors []error) {
	value := v.(string)
	for i := range l {
		if value == l[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, l))
	return
}

// Not currently used.
func ValidateInsensitiveStringList(v interface{}, k string, l []string) (ws []string, errors []error) {
	value := v.(string)
	for i := range l {
		if strings.EqualFold(value, l[i]) {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, l))
	return
}

func ValidateIntRange(v interface{}, k string, l int, h int) (ws []string, errors []error) {
	i, ok := v.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("%q must be an integer", k))
		return
	}
	if i < l || i > h {
		errors = append(errors, fmt.Errorf("%q must be between %d and %d", k, l, h))
		return
	}
	return
}

func validateTrueOnly(v interface{}, k string) (ws []string, errors []error) {
	if b, ok := v.(bool); ok && b {
		return
	}
	if v, ok := v.(string); ok && v == "true" {
		return
	}
	errors = append(errors, fmt.Errorf("%q must be true", k))
	return
}

func validateS3BucketLifecycleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", value))
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as RFC3339 Timestamp Format", value))
	}

	return
}

func validateS3BucketLifecycleExpirationDays(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) <= 0 {
		errors = append(errors, fmt.Errorf(
			"%q must be greater than 0", k))
	}

	return
}

func validateS3BucketLifecycleRuleId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 255 {
		errors = append(errors, fmt.Errorf(
			"%q cannot exceed 255 characters", k))
	}
	return
}

func validateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := normalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

func validateName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 || len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q must contain more than 1 and less than 64 characters", k))
	}

	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf("only alphanumeric characters, hyphens, and underscores allowed in %q", k))
	}

	return
}

func validateStackTemplate(v interface{}, k string) (ws []string, errors []error) {
	if looksLikeJsonString(v) {
		if _, err := normalizeJsonString(v); err != nil {
			errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
		}
	} else {
		if _, err := checkYamlString(v); err != nil {
			errors = append(errors, fmt.Errorf("%q contains an invalid YAML: %s", k, err))
		}
	}
	return
}

func validateIP(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	ipAddr := net.ParseIP(value)

	if ipAddr == nil || value != ipAddr.String() {
		errors = append(errors, fmt.Errorf("%q must contain a valid network IP address, got %q", k, value))
	}

	return
}

func validateCIDR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, got %q", k, value))
	}

	return
}

func validateVBSPolicyName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if strings.HasPrefix(strings.ToLower(value), "default") {
		errors = append(errors, fmt.Errorf(
			"%q cannot start with default: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

func validateVBSPolicyFrequency(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 14 {
		errors = append(errors, fmt.Errorf(
			"%q should be in the range of 1-14: %d", k, value))
	}
	return
}

func validateVBSPolicyStatus(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "ON" && value != "OFF" {
		errors = append(errors, fmt.Errorf(
			"%q should be either ON or OFF: %q", k, value))
	}
	return
}

func validateVBSPolicyRetentionNum(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 2 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be less than 2: %d", k, value))
	}
	return
}

func validateVBSPolicyRetainBackup(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Y" && value != "N" {
		errors = append(errors, fmt.Errorf(
			"%q should be either N or Y: %q", k, value))
	}
	return
}

func validateVBSTagKey(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 36 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 36 characters: %q", k, value))
	}
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

func validateVBSTagValue(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 43 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 43 characters: %q", k, value))
	}
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

func validateVBSBackupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if strings.HasPrefix(strings.ToLower(value), "autobk") {
		errors = append(errors, fmt.Errorf(
			"%q cannot start with autobk: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

func validateVBSBackupDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}
	pattern := `^[^<>]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

func validateAntiDdosTrafficPosID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 9 {
		errors = append(errors, fmt.Errorf(
			"%q should be in the range of 1-9: %d", k, value))
	}
	return
}

func validateAntiDdosHttpRequestPosID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 15 {
		errors = append(errors, fmt.Errorf(
			"%q should be in the range of 1-15: %d", k, value))
	}
	return
}

func validateAntiDdosCleaningAccessPosID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 8 {
		errors = append(errors, fmt.Errorf(
			"%q should be in the range of 1-8: %d", k, value))
	}
	return
}

func validateAntiDdosAppTypeID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 1 {
		errors = append(errors, fmt.Errorf(
			"%q should be 0 or 1: %d", k, value))
	}
	return
}

func validateECSTagValue(v interface{}, k string) (ws []string, errors []error) {
	tagmap := v.(map[string]interface{})
	vv := regexp.MustCompile(`^[0-9a-zA-Z-_]+$`)
	for k, v := range tagmap {
		value := v.(string)
		if !vv.MatchString(value) {
			errors = append(errors, fmt.Errorf("Tag value must be string only contains digits, letters, underscores(_) and hyphens(-), but got %s=%s", k, value))
			break
		}
	}
	return
}

func validateK8sTagsMap(v interface{}, k string) (ws []string, errors []error) {
	values := v.(map[string]interface{})
	pattern := regexp.MustCompile(`^[.\-_A-Za-z0-9]+$`)

	for key, value := range values {
		valueString := value.(string)
		if len(key) < 1 {
			errors = append(errors, fmt.Errorf("key %q cannot be shorter than 1 characters: %q", k, key))
		}

		if len(valueString) < 1 {
			errors = append(errors, fmt.Errorf("value %q cannot be shorter than 1 characters: %q", k, value))
		}

		if len(key) > 63 {
			errors = append(errors, fmt.Errorf("key %q cannot be longer than 63 characters: %q", k, key))
		}

		if len(valueString) > 63 {
			errors = append(errors, fmt.Errorf("value %q cannot be longer than 63 characters: %q", k, value))
		}

		if !pattern.MatchString(key) {
			errors = append(errors, fmt.Errorf("key %q doesn't comply with restrictions (%q): %q", k, pattern, key))
		}

		if !pattern.MatchString(valueString) {
			errors = append(errors, fmt.Errorf("value %q doesn't comply with restrictions (%q): %q", k, pattern, valueString))
		}
	}

	return
}

func validateDDSStartTime(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	re := regexp.MustCompile(`(?P<hh>\d{2}):(?P<mm>\d{2})-(?P<HH>\d{2}):(?P<MM>\d{2})`)
	timeRange := re.FindStringSubmatch(value)

	paramsMap := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(timeRange) {
			paramsMap[name] = timeRange[i]
		}
	}

	startHour, err := strconv.Atoi(paramsMap["hh"])
	if err != nil {
		errors = append(errors, fmt.Errorf("%q must be int convertable: %s", paramsMap["hh"], err))
	}
	endHour, err := strconv.Atoi(paramsMap["HH"])
	if err != nil {
		errors = append(errors, fmt.Errorf("%q must be int convertable: %s", paramsMap["HH"], err))
	}
	startMinutes, err := strconv.Atoi(paramsMap["mm"])
	if err != nil {
		errors = append(errors, fmt.Errorf("%q must be int convertable: %s", paramsMap["mm"], err))
	}
	endMinutes, err := strconv.Atoi(paramsMap["MM"])
	if err != nil {
		errors = append(errors, fmt.Errorf("%q must be int convertable: %s", paramsMap["MM"], err))
	}
	if len(errors) != 0 {
		return
	}
	if startHour+1 != endHour {
		errors = append(errors, fmt.Errorf("the `HH` value must be 1 greater than the `hh` value: %s", v))
	}
	if startMinutes != endMinutes {
		errors = append(errors, fmt.Errorf("the values from `mm` and `MM` must be the same: %s", v))
	}
	if startMinutes%15 != 0 {
		errors = append(errors, fmt.Errorf("the values from `mm` and `MM` must be set to any of the 00, 15, 30, or 45: %s", v))
	}

	return
}

func resourceASGroupValidateListenerId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	listenerIds := strings.Split(value, ",")
	if len(listenerIds) <= 3 {
		return
	}
	errors = append(errors, fmt.Errorf("%q supports binding up to 3 ELB listeners which are separated by a comma", k))
	return
}
