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
		_, _ = fmt.Fprintf(w, f.OkResponse)
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
