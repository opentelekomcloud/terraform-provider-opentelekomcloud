package pauseunpause

import "github.com/huaweicloud/golangsdk"

// PauseResult is the response from a Pause operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type PauseResult struct {
	golangsdk.ErrResult
}

// UnpauseResult is the response from an Unpause operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type UnpauseResult struct {
	golangsdk.ErrResult
}
