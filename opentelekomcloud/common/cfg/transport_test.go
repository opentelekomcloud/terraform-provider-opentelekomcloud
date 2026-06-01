package cfg

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
)

type failHandler struct {
	ExpectedFailures int
	ErrorCode        int
	FailCount        int
	OkCode           int
	OkResponse       string

	mut *sync.RWMutex
}

func (f *failHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f.mut == nil {
		f.mut = new(sync.RWMutex)
	}

	if f.OkCode == 0 {
		f.OkCode = 200
	}

	defer func() { _ = r.Body.Close() }()
	if f.FailCount < f.ExpectedFailures {
		f.mut.Lock()
		f.FailCount += 1
		f.mut.Unlock()
		w.WriteHeader(f.ErrorCode)
	} else {
		w.WriteHeader(f.OkCode)
		_, _ = fmt.Fprint(w, f.OkResponse)
	}
}

const tokenOutput = `
{
   "token":{
      "methods":[
         "password"
      ],
      "roles":[],
      "expires_at":"2017-06-03T02:19:49.000000Z",
      "project":{},
      "catalog":[],
      "user": {},
      "issued_at":"2017-06-03T01:19:49.000000Z"
   }
}
`

type recordingTransport struct {
	failCount    int
	maxFailures  int
	sdkDates     []string
	authHeaders  []string
}

func (rt *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.sdkDates = append(rt.sdkDates, req.Header.Get("X-Sdk-Date"))
	rt.authHeaders = append(rt.authHeaders, req.Header.Get("Authorization"))

	if rt.failCount < rt.maxFailures {
		rt.failCount++
		return nil, fmt.Errorf("connection refused")
	}
	return &http.Response{
		StatusCode: 200,
		Body:       http.NoBody,
	}, nil
}

func TestRoundTripperRetryResign(t *testing.T) {
	inner := &recordingTransport{maxFailures: 1}

	signOpts := golangsdk.SignOptions{
		AccessKey: "test-ak",
		SecretKey: "test-sk",
	}

	req, err := http.NewRequest("GET", "https://example.com/v1/resource", nil)
	th.CheckNoErr(t, err)

	// Set stale signature headers so that ReSign is guaranteed to produce different ones
	staleDate := "20200101T000000Z"
	staleAuth := "SDK-HMAC-SHA256 stale-signature"
	req.Header.Set("Host", req.URL.Host)
	req.Header.Set("X-Sdk-Date", staleDate)
	req.Header.Set("Authorization", staleAuth)

	rt := &RoundTripper{
		Rt:          inner,
		MaxRetries:  1,
		SignOptions: &signOpts,
	}

	resp, err := rt.RoundTrip(req)
	th.CheckNoErr(t, err)
	th.AssertEquals(t, 200, resp.StatusCode)

	th.AssertEquals(t, 2, len(inner.sdkDates))
	th.AssertEquals(t, staleDate, inner.sdkDates[0])

	if inner.sdkDates[1] == staleDate {
		t.Error("X-Sdk-Date should be updated on retry, but was still stale")
	}
	if inner.authHeaders[1] == staleAuth {
		t.Error("Authorization should be updated on retry, but was still stale")
	}
}

func TestRoundTripperRetry(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	failHandler := &failHandler{
		ExpectedFailures: 1,
		ErrorCode:        502,
		OkCode:           201,
		OkResponse:       tokenOutput,
	}

	th.Mux.Handle("/", failHandler)

	cfg := &Config{MaxRetries: failHandler.ExpectedFailures}

	_, err := cfg.genClient(golangsdk.AuthOptions{
		IdentityEndpoint: th.Endpoint() + "v3",
		Username:         "user",
		Password:         "qwerty!",
		DomainName:       "DOMAIN001",
	})

	th.CheckNoErr(t, err)
	th.AssertEquals(t, failHandler.ExpectedFailures, failHandler.FailCount)
}
