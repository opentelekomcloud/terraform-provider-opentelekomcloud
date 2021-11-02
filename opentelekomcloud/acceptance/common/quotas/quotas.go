package quotas

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/sync/semaphore"
)

const (
	timeoutMsg = "reached timeout waiting for quota `%s` to be acquired"
	tooManyMsg = "can't acquire more resources (%d) than exist (%d) for quota %s"
)

// Quota is a wrapper around a semaphore providing simple control over shared resources with quotas
type Quota struct {
	Name    string
	Size    int64
	Current int64

	sem *semaphore.Weighted
	ctx context.Context
}

// NewQuota creates a new Quota with persistent context inside
func NewQuota(count int64) *Quota {
	q := &Quota{
		sem:     semaphore.NewWeighted(count),
		ctx:     context.Background(),
		Size:    count,
		Current: count,
	}
	return q
}

// NewQuotaWithTimeout creates a new Quota with timing out context
func NewQuotaWithTimeout(count int64, timeout time.Duration) (*Quota, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	quota := &Quota{
		sem:     semaphore.NewWeighted(count),
		ctx:     ctx,
		Size:    count,
		Current: count,
	}
	return quota, cancelFunc
}

// Acquire decrease count of available resources by 1
func (q *Quota) Acquire() error {
	return q.AcquireMultiple(1)
}

// AcquireMultiple decrease count of available resources by n
func (q *Quota) AcquireMultiple(n int64) error {
	if n > q.Size {
		return fmt.Errorf(tooManyMsg, n, q.Size, q.Name)
	}
	if err := q.sem.Acquire(q.ctx, n); err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf(timeoutMsg, q.Name)
		}
	}
	q.Current -= n
	return nil
}

// Release increase count of available resources by 1
func (q *Quota) Release() {
	q.ReleaseMultiple(1)
}

// ReleaseMultiple increase count of available resources by n
func (q *Quota) ReleaseMultiple(n int64) {
	q.sem.Release(n)
	q.Current += n
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
	q := NewQuota(count)
	q.Name = strings.ToLower(envVar)
	return q
}

// Shared quotas

var (
	// Compute

	// Server - shared compute instance quota (number of instances only)
	Server = FromEnv("OS_SERVER_QUOTA", 10)
	CPU    = FromEnv("OS_CPU_QUOTA", 40)
	RAM    = FromEnv("OS_RAM_QUOTA", 160*1024)

	// Networking

	// FloatingIP - shared floating IP quota
	FloatingIP = FromEnv("OS_FLOATING_IP_QUOTA", 3)
	// Router - shared router(VPC) quota
	Router        = FromEnv("OS_ROUTER_QUOTA", 7) // safe value
	Subnet        = FromEnv("OS_SUBNET_QUOTA", 50)
	Network       = FromEnv("OS_NETWORK_QUOTA", 50)
	SecurityGroup = FromEnv("OS_SG_QUOTA", 100)

	// Volumes

	// Volume - quota for block storage volumes
	Volume = FromEnv("OS_VOLUME_QUOTA", 50)
	// VolumeSize - quota for block storage total size, GB
	VolumeSize = FromEnv("OS_VOLUME_SIZE_QUOTA", 12500)

	// LoadBalancer - quota for load balancer instances
	LoadBalancer = FromEnv("OS_LB_QUOTA", 50)
)

// ExpectedQuota is a simple container of quota + count used for `Multiple` operations
type ExpectedQuota struct {
	Q     *Quota
	Count int64
}

// X multiples quota returning new `ExpectedQuota` instance
func (q ExpectedQuota) X(multiplier int64) *ExpectedQuota {
	return &ExpectedQuota{
		Q:     q.Q,
		Count: q.Count * multiplier,
	}
}

type MultipleQuotas []*ExpectedQuota

// X multiples quota returning new `MultipleQuotas` instance
func (q MultipleQuotas) X(multiplier int64) MultipleQuotas {
	newOne := make(MultipleQuotas, len(q))
	for i, q := range q {
		newOne[i] = q.X(multiplier)
	}
	return newOne
}

// AcquireMultipleQuotas tries to acquire all given quotas, reverting on failure
func AcquireMultipleQuotas(e []*ExpectedQuota, interval time.Duration) error {
	// validate if all Count values of ExpectQuota are correct
	var mErr *multierror.Error
	for _, q := range e {
		if q.Count > q.Q.Size {
			mErr = multierror.Append(mErr, fmt.Errorf(tooManyMsg, q.Count, q.Q.Size, q.Q.Name))
		}
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}
	var ok bool
	var err error
	for !ok {
		ok, err = tryAcquireMultiple(e)
		if err != nil {
			return err
		}
		if !ok {
			time.Sleep(interval) // successfully acquired all quotas
		}
	}
	return nil
}

var multipleLock = &sync.Mutex{}

func tryAcquireMultiple(e []*ExpectedQuota) (bool, error) {
	multipleLock.Lock()
	defer multipleLock.Unlock()

	var acquired []*ExpectedQuota
	var ok bool
	defer func() {
		if !ok {
			ReleaseMultipleQuotas(acquired)
		}
	}()

	for _, q := range e {
		if err := q.Q.ctx.Err(); err != nil {
			if _, ok := q.Q.ctx.Deadline(); ok {
				return false, fmt.Errorf(timeoutMsg, q.Q.Name)
			}
			return false, fmt.Errorf("unknown error trying to obtain multiple quotas: %w", err)
		}
		if ok := q.Q.sem.TryAcquire(q.Count); ok {
			q.Q.Current -= q.Count
			acquired = append(acquired, q)
		}
	}
	if len(acquired) == len(e) { // all quotas are acquired
		ok = true
	}
	return ok, nil
}

func ReleaseMultipleQuotas(e []*ExpectedQuota) {
	for _, q := range e {
		q.Q.ReleaseMultiple(q.Count)
	}
}
