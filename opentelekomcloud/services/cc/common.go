package cc

import "strconv"

const (
	errCreationV3Client = "error creating OpenTelekomCloud Cloud Connect v3 client: %w"
	ccClientV3          = "cc-v3-client"
)

func stringFilter(v string) []string {
	if v == "" {
		return nil
	}
	return []string{v}
}

func boolFilter(v string) (*bool, error) {
	if v == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
