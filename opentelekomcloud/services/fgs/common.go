package fgs

import "strings"

const (
	errCreationV2Client = "error creating OpenTelekomCloud FuncGraphV2 client: %w"
	fgsClientV2         = "fgs-v2-client"
)

/*
 * Parse urn according from fun_urn.
 * If the separator is not ":" then return to the original value.
 */
func resourceFgsFunctionUrn(urn string) string {
	index := strings.LastIndex(urn, ":")
	if index != -1 {
		urn = urn[0:index]
	}
	return urn
}
