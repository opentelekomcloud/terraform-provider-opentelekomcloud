package suspendresume

import "github.com/huaweicloud/golangsdk"

// SuspendResult is the response from a Suspend operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type SuspendResult struct {
	golangsdk.ErrResult
}

// UnsuspendResult is the response from an Unsuspend operation. Call
// its ExtractErr method to determine if the request succeeded or failed.
type UnsuspendResult struct {
	golangsdk.ErrResult
}
