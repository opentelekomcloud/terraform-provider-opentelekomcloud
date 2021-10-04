package quotas

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/sync/semaphore"
)

const (
	timeoutMsg = "reached timeout waiting for quota to be acquired"
	tooManyMsg = "can't acquire more resources (%d) than exist (%d)"
)

// Quota is a wrapper around a semaphore providing simple control over shared resources with quotas
type Quota struct {
	sem  *semaphore.Weighted
	ctx  context.Context
	size int64
}

// NewQuota creates a new Quota with persistent context inside
func NewQuota(count int64) *Quota {
	q := &Quota{
		sem:  semaphore.NewWeighted(count),
		ctx:  context.Background(),
		size: count,
	}
	return q
}

// NewQuotaWithTimeout creates a new Quota with timing out context
func NewQuotaWithTimeout(count int64, timeout time.Duration) (*Quota, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	quota := &Quota{
		sem:  semaphore.NewWeighted(count),
		ctx:  ctx,
		size: count,
	}
	return quota, cancelFunc
}

// Acquire decrease count of available resources by 1
func (q *Quota) Acquire() error {
	return q.AcquireMultiple(1)
}

// AcquireMultiple decrease count of available resources by n
func (q *Quota) AcquireMultiple(n int64) error {
	if n > q.size {
		return fmt.Errorf(tooManyMsg, n, q.size)
	}
	if err := q.sem.Acquire(q.ctx, n); err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf(timeoutMsg)
		}
	}
	return nil
}

// Release increase count of available resources by 1
func (q *Quota) Release() {
	q.ReleaseMultiple(1)
}

// ReleaseMultiple increase count of available resources by n
func (q *Quota) ReleaseMultiple(n int64) {
	q.sem.Release(n)
}

// FromEnv creates quota instance with limit set to env var or default value if the variable
// is not set or not an integer.
func FromEnv(envVar string, def int64) *Quota {
	count := def
	if eq := os.Getenv(envVar); eq != "" {
		if n, err := strconv.Atoi(eq); err == nil {
			count = int64(n)
		} else {
			log.Printf("failed to read env var %s, using default value %d: %s", envVar, def, err.Error())
		}
	}
	return NewQuota(count)
}

// Shared quotas

var (
	// Compute

	// Server - shared compute instance quota (number of instances only)
	Server = FromEnv("OS_SERVER_QUOTA", 10)
	CPU    = FromEnv("OS_CPU_QUOTA", 40)
	RAM    = FromEnv("OS_RAM_QUOTA", 160)

	// Networking

	// FloatingIP - shared floating IP quota
	FloatingIP = FromEnv("OS_FLOATING_IP_QUOTA", 3)
	// Router - shared router(VPC) quota
	Router = FromEnv("OS_ROUTER_QUOTA", 10)

	// Volumes

	// Volume - quota for block storage volumes
	Volume = FromEnv("OS_VOLUME_QUOTA", 50)
	// VolumeSize - quota for block storage total size, GB
	VolumeSize = FromEnv("OS_VOLUME_SIZE_QUOTA", 12500)
)

// ExpectedQuota is a simple container of quota + count used for `Multiple` operations
type ExpectedQuota struct {
	Q     *Quota
	Count int64
}

// AcquireMultipleQuotas tries to acquire all given quotas, reverting on failure
func AcquireMultipleQuotas(e []*ExpectedQuota, interval time.Duration) error {
	var acquired []*ExpectedQuota
	repeated := false
	// validate if all Count values of ExpectQuota are correct
	var mErr *multierror.Error
	for _, q := range e {
		if q.Count > q.Q.size {
			mErr = multierror.Append(mErr, fmt.Errorf(tooManyMsg, q.Count, q.Q.size))
		}
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}
	for len(acquired) != len(e) {
		ReleaseMultipleQuotas(acquired)
		acquired = nil
		if repeated {
			time.Sleep(interval)
		}
		for _, q := range e {
			if err := q.Q.ctx.Err(); err != nil {
				if _, ok := q.Q.ctx.Deadline(); ok {
					return fmt.Errorf(timeoutMsg)
				}
			}
			ok := q.Q.sem.TryAcquire(q.Count)
			if ok {
				acquired = append(acquired, q)
			}
		}
		repeated = true
	}
	return nil
}

func ReleaseMultipleQuotas(e []*ExpectedQuota) {
	for _, q := range e {
		q.Q.ReleaseMultiple(q.Count)
	}
}
