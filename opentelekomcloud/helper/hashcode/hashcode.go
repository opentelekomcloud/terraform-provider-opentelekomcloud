package hashcode

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"
)

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// Strings hashes a list of strings to a unique hashcode.
func Strings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", String(buf.String()))
}

func DecodeHashAndHexEncode(v interface{}) string {
	switch v := v.(type) {
	case string:
		return installScriptHashSum(v)
	default:
		return ""
	}
}

func installScriptHashSum(script string) string {
	// Check whether the script is not Base64 encoded.
	// Always calculate hash of base64 decoded value since we
	// check against double-encoding when setting it
	v, base64DecodeError := base64.StdEncoding.DecodeString(script)
	if base64DecodeError != nil {
		v = []byte(script)
	}

	hash := sha1.Sum(v)
	return hex.EncodeToString(hash[:])
}

// TryBase64EncodeString will encode a string with base64.
// If the string is already base64 encoded, returns it directly.
func TryBase64EncodeString(str string) string {
	if _, err := base64.StdEncoding.DecodeString(str); err != nil {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	return str
}
