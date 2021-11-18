package quotas

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
)

func TestQuota_Release(t *testing.T) {
	q := NewQuota(1)
	th.AssertNoErr(t, q.acquireMultiple(1))
	q.releaseMultiple(1)
	th.AssertNoErr(t, q.acquireMultiple(1))
}

var namelessTimeout = fmt.Sprintf(timeoutMsg, "")

func TestQuota_Timeout(t *testing.T) {
	t.Parallel()
	q, _ := NewQuotaWithTimeout(1, 1*time.Millisecond)
	th.AssertNoErr(t, q.acquireMultiple(1))
	err := q.acquireMultiple(1)
	th.AssertEquals(t, namelessTimeout, err.Error())
}

func TestQuota_AcquireTooMuch(t *testing.T) {
	q := NewQuota(1)
	err := q.acquireMultiple(2)
	th.AssertEquals(t, "can't acquire more resources (2) than exist (1) for quota ", err.Error())
}

func TestFromEnv_Default(t *testing.T) {
	vName := tools.RandomString("OS_", 10)
	vDef := int64(tools.RandomInt(1, 100))
	q := FromEnv(vName, vDef)
	th.AssertEquals(t, q.Size, vDef)
}

func TestFromEnv(t *testing.T) {
	vName := tools.RandomString("OS_", 10)
	vDef := int64(tools.RandomInt(1, 100))
	_ = os.Setenv(vName, strconv.Itoa(int(vDef)))
	q := FromEnv(vName, 0)
	th.AssertEquals(t, q.Size, vDef)
}

func TestFromEnv_InvalidVar(t *testing.T) {
	vName := tools.RandomString("OS_", 10)
	vDef := int64(tools.RandomInt(1, 100))
	_ = os.Setenv(vName, tools.RandomString("var", 3))
	q := FromEnv(vName, vDef)
	th.AssertEquals(t, q.Size, vDef)
}

// Check that deadlock really appears when using simple acquiring
func TestQuota_TimeoutDeadlock(t *testing.T) {
	q1, _ := NewQuotaWithTimeout(1, 2*time.Millisecond)
	q2, _ := NewQuotaWithTimeout(1, 2*time.Millisecond)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		th.AssertNoErr(t, q1.acquireMultiple(1))
		time.Sleep(1 * time.Millisecond)
		err := q2.acquireMultiple(1)
		th.AssertEquals(t, namelessTimeout, err.Error())
	}()
	go func() {
		defer wg.Done()
		th.AssertNoErr(t, q2.acquireMultiple(1))
		time.Sleep(1 * time.Millisecond)
		err := q1.acquireMultiple(1)
		th.AssertEquals(t, namelessTimeout, err.Error())
	}()
	wg.Wait()
}

func TestQuota_Multiple(t *testing.T) {
	q1, _ := NewQuotaWithTimeout(1, 10*time.Minute)
	q2, _ := NewQuotaWithTimeout(2, 10*time.Minute)
	qts := []*ExpectedQuota{{q1, 1}, {q2, 1}}

	th.AssertNoErr(t, acquireMultipleQuotas(qts, 0))
	releaseMultipleQuotas(qts)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Millisecond)
		th.AssertNoErr(t, acquireMultipleQuotas(qts, 0))
		releaseMultipleQuotas(qts)
	}()
	go func() {
		defer wg.Done()
		th.AssertNoErr(t, q2.acquireMultiple(1))
		defer q2.releaseMultiple(1)
		th.AssertNoErr(t, q1.acquireMultiple(1))
		defer q1.releaseMultiple(1)
	}()
	wg.Wait()
}

func TestQuota_MultipleNotExhausting(t *testing.T) {
	q1, _ := NewQuotaWithTimeout(1, 10*time.Millisecond)
	q2, _ := NewQuotaWithTimeout(2, 10*time.Millisecond)
	qts := []*ExpectedQuota{{q1, 1}, {q2, 1}}

	th.AssertNoErr(t, acquireMultipleQuotas(qts, 0))
	th.AssertNoErr(t, q2.acquireMultiple(1))
}

func TestQuota_MultipleTooMany(t *testing.T) {
	q1, _ := NewQuotaWithTimeout(1, 10*time.Millisecond)
	q2, _ := NewQuotaWithTimeout(2, 10*time.Millisecond)
	qts := []*ExpectedQuota{{q1, 10}, {q2, 10}}

	err := acquireMultipleQuotas(qts, 0)
	th.AssertEquals(t,
		"2 errors occurred:\n\t* can't acquire more resources (10) than exist (1) for quota \n\t* can't acquire more resources (10) than exist (2) for quota \n\n",
		err.Error(),
	)
}

func TestQuota_MultipleUnreleased(t *testing.T) {
	q1, _ := NewQuotaWithTimeout(1, 2*time.Millisecond)
	q2, _ := NewQuotaWithTimeout(2, 2*time.Millisecond)
	qts := []*ExpectedQuota{{q1, 1}, {q2, 1}}

	th.AssertNoErr(t, acquireMultipleQuotas(qts, 0))

	err := acquireMultipleQuotas(qts, 0)
	th.AssertEquals(t, namelessTimeout, err.Error())
}

func TestMultipleQuotas(t *testing.T) {
	q1, _ := NewQuotaWithTimeout(3, 2*time.Millisecond)
	q2, _ := NewQuotaWithTimeout(3, 2*time.Millisecond)
	qtsOrig := MultipleQuotas{{q1, 1}, {q2, 1}}

	qts := qtsOrig.X(2)

	th.AssertEquals(t, qts[0].Q, qtsOrig[0].Q)
	th.AssertEquals(t, qts[0].Count, qtsOrig[0].Count*2)

	th.AssertNoErr(t, acquireMultipleQuotas(qts, 0))     // -2
	th.AssertNoErr(t, acquireMultipleQuotas(qtsOrig, 0)) // -1
	err := acquireMultipleQuotas(qtsOrig, 0)
	th.AssertEquals(t, namelessTimeout, err.Error())
}

func TestBookOne(t *testing.T) {
	q, _ := NewQuotaWithTimeout(1, 1*time.Millisecond)
	t.Run("first", func(tt *testing.T) {
		tt.Parallel()
		tt.Logf("started first")
		tt.Logf("quota len: %d", q.Current)
		BookOne(tt, q)
		tt.Logf("quota len: %d", q.Current)
		time.Sleep(20 * time.Millisecond)
	})
	t.Run("second", func(tt *testing.T) {
		tt.Parallel()
		time.Sleep(10 * time.Millisecond)
		err := q.acquireMultiple(1)
		th.AssertEquals(t, namelessTimeout, err.Error())
	})
}
func TestBookMany(t *testing.T) {
	q1, _ := NewQuotaWithTimeout(1, 1*time.Millisecond)
	q2, _ := NewQuotaWithTimeout(1, 1*time.Millisecond)
	qts := MultipleQuotas{{q1, 1}, {q2, 1}}
	t.Run("first", func(tt *testing.T) {
		tt.Parallel()
		BookMany(tt, qts)
		time.Sleep(20 * time.Millisecond)
	})
	t.Run("second", func(tt *testing.T) {
		tt.Parallel()
		time.Sleep(10 * time.Millisecond)
		err := acquireMultipleQuotas(qts, 1*time.Millisecond)
		th.AssertEquals(t, namelessTimeout, err.Error())
	})
}
