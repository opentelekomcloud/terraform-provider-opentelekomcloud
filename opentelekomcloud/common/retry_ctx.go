package common

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// RetryFunc is the function retried until it succeeds.
// The first return parameter is the result of the retry func.
// The second return parameter indicates whether a retry is required.
// The last return parameter is the error of the func.
type RetryFunc func() (res interface{}, retry bool, err error)

type RetryContextWithWaitForStateParam struct {
	Ctx context.Context
	// The func that need to be retried
	RetryFunc RetryFunc
	// The wait func when the retry which returned by the retry func is true
	WaitFunc resource.StateRefreshFunc
	// The target of the wait func
	WaitTarget []string
	// The pending of the wait func
	WaitPending []string
	// The timeout of the retry func and wait func
	Timeout time.Duration
	// The delay timeout of the retry func and wait func
	DelayTimeout time.Duration
	// The poll interval of the retry func and wait func
	PollInterval time.Duration
}

// RetryContextWithWaitForState The RetryFunc will be called first
// if the error of the return is nil, the retry will be ended and the res of the return will be returned
// if the retry of the return is true, the RetryFunc will be retried, and the WaitFunc will be called if it is not nil
// if the retry of the return is false, the retry will be ended and the error of the retry func will be returned
func RetryContextWithWaitForState(param *RetryContextWithWaitForStateParam) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"retryable"},
		Target:       []string{"success"},
		Timeout:      param.Timeout,
		Delay:        param.DelayTimeout,
		PollInterval: param.PollInterval,
		Refresh: func() (interface{}, string, error) {
			res, retry, err := param.RetryFunc()
			if err == nil {
				if res != nil {
					return res, "success", nil
				}
				// If we didn't find the resource, convert it to "", otherwise,
				// it will report an error in WaitForStateContext.
				return "", "success", nil
			}

			if !retry {
				return nil, "quit", err
			}

			if param.WaitFunc != nil {
				stateConf := &resource.StateChangeConf{
					Target:       param.WaitTarget,
					Pending:      param.WaitPending,
					Refresh:      param.WaitFunc,
					Timeout:      param.Timeout,
					Delay:        param.DelayTimeout,
					PollInterval: param.PollInterval,
				}
				if _, err := stateConf.WaitForStateContext(param.Ctx); err != nil {
					return nil, "quit", err
				}
			}
			return "", "retryable", nil
		},
	}

	return stateConf.WaitForStateContext(param.Ctx)
}
