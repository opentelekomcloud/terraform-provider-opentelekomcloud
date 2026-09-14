package swr

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
)

const ClientError = "error creating SWR V2 client: %w"

// alreadyGone reports whether err is a 404, which SWR returns once the resource
// no longer exists. Deleting such a resource is a no-op rather than a failure.
func alreadyGone(err error) bool {
	_, ok := err.(golangsdk.ErrDefault404)
	return ok
}

func organization(d *schema.ResourceData) string {
	return d.Get("organization").(string)
}

func repository(name string) string {
	return strings.ReplaceAll(name, "/", "$")
}
