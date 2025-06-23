package ces

import "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"

const (
	errCreationClient = "error creating CESv1 client: %w"
	cesClientV1       = "ces-v1-client"
)

func InterfaceToInt64(i interface{}) int64 {
	v, ok := i.(int)
	if !ok {
		panic(fmterr.Errorf("InterfaceToInt64: value is not of type int: %#v", i))
	}
	return int64(v)
}
