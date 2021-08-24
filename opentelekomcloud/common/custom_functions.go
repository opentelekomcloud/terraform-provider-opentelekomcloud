package common

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// FirstOneSet select first nonempty string or returns error if both are empty
func FirstOneSet(res map[string]interface{}, k1, k2 string) (interface{}, error) {
	v1 := res[k1]
	v2 := res[k2]
	if v1 == "" && v2 == "" {
		return nil, fmt.Errorf("none of %s and %s are set", k1, k2)
	}
	if v1 != "" {
		return v1, nil
	}
	return v2, nil
}

func InstallScriptHashSum(script string) string {
	// Check whether the preinstall/postinstall is not Base64 encoded.
	// Always calculate hash of base64 decoded value since we
	// check against double-encoding when setting it
	v, base64DecodeError := base64.StdEncoding.DecodeString(script)
	if base64DecodeError != nil {
		v = []byte(script)
	}

	hash := sha1.Sum(v)
	return hex.EncodeToString(hash[:])
}

func InstallScriptEncode(script string) string {
	if _, err := base64.StdEncoding.DecodeString(script); err != nil {
		return base64.StdEncoding.EncodeToString([]byte(script))
	}
	return script
}

func GetHashOrEmpty(v interface{}) string {
	switch v := v.(type) {
	case string:
		return InstallScriptHashSum(v)
	default:
		return ""
	}
}
