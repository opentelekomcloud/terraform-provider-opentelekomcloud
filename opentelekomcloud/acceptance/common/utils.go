package common

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

// GenerateRandomDomain is a method Used to generate the domain names and the domain IDs,
// which cannot start with a digit. All elements' length are the same.
func GenerateRandomDomain(count, strLen int) []string {
	if count < 1 || strLen < 1 {
		return nil
	}
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = acctest.RandStringFromCharSet(strLen, "abcdef")
	}
	return result
}

// Reverse is a function that used to reverse the order of the characters in the given string.
func Reverse(s string) string {
	bs := []byte(s)
	for left, right := 0, len(s)-1; left < right; left++ {
		bs[left], bs[right] = bs[right], bs[left]
		right--
	}
	return string(bs)
}
