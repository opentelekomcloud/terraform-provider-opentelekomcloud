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
