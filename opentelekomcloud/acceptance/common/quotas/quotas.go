package quotas

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/semaphore"
)

// Quota is a wrapper around a semaphore providing simple control over shared resources with quotas
type Quota struct {
	sem *semaphore.Weighted
	ctx context.Context
}

// NewQuota creates a new Quota with persistent context inside
func NewQuota(count int64) *Quota {
	q := &Quota{
		sem: semaphore.NewWeighted(count),
		ctx: context.Background(),
	}
	return q
}

// NewQuotaWithTimeout creates a new Quota with timing out context
func NewQuotaWithTimeout(count int64, timeout time.Duration) (*Quota, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	quota := &Quota{
		sem: semaphore.NewWeighted(count),
		ctx: ctx,
	}
	return quota, cancelFunc
}

// Acquire decrease count of available resources by 1
func (q *Quota) Acquire() error {
	return q.AcquireMultiple(1)
}

// AcquireMultiple decrease count of available resources by n
func (q *Quota) AcquireMultiple(n int64) error {
	return q.sem.Acquire(q.ctx, n)
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
)
