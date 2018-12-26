package lockunlock

import (
	"github.com/huaweicloud/golangsdk"
)

// LockResult and UnlockResult are the responses from a Lock and Unlock
// operations respectively. Call their ExtractErr methods to determine if the
// requests suceeded or failed.
type LockResult struct {
	golangsdk.ErrResult
}

type UnlockResult struct {
	golangsdk.ErrResult
}
