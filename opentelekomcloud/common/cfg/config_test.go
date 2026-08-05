package cfg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"text/template"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
)

func writeYamlTemplate(tmpl string, filename string, data *Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	cloudTemplate, _ := template.New("yaml").Parse(tmpl)
	return cloudTemplate.Execute(file, data)
}

func checkConfigField(t *testing.T, act *Config, excp *Config, fieldName string) {
	actual := reflect.ValueOf(*act).FieldByName(fieldName).String()
	expected := reflect.ValueOf(*excp).FieldByName(fieldName).String()
	if actual != expected {
		t.Errorf("field %s: expected %s, got %s", fieldName, expected, actual)
	}
}

const fileName = "./clouds.yaml"

func TestReadCloudsYaml(t *testing.T) {
	tmpl := `
clouds:
  useless_cloud:
    auth:
      auth_url: http://localhost/
  {{.Cloud}}:
    auth:
      auth_url: {{.IdentityEndpoint}}
      username: {{.Username}}
      password: {{.Password}}
      project_name: {{.TenantName}}
      domain_name: {{.DomainName}}
    region_name: {{.Region}}
    verify: {{not .Insecure}}
    cert: {{.ClientCertFile}}
    key: {{.ClientKeyFile}}
    cacert: {{.CACertFile}}
`

	referenceConfig := &Config{
		Cloud:            "otc",
		Username:         "demouser",
		Password:         "qwerty!1234",
		Region:           "eu-de",
		TenantName:       "eu-de_sub",
		DomainName:       "OTC1354835",
		IdentityEndpoint: "http://localhost:33666",
		Insecure:         true,
		ClientCertFile:   "cert_file.crt",
		ClientKeyFile:    "key_file.key",
		CACertFile:       "ca.crt",
	}
	_ = os.Setenv("OS_CLOUD", "otc")

	err := writeYamlTemplate(tmpl, fileName, referenceConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(fileName) }()

	c := &Config{Cloud: referenceConfig.Cloud}
	if err := c.Load(); err != nil {
		t.Fatal()
	}

	comparedFields := []string{
		"IdentityEndpoint", "Region",
		"TenantName", "Username", "Password",
		"Insecure", "ClientCertFile",
		"ClientKeyFile", "CACertFile",
	}

	for _, field := range comparedFields {
		t.Run(field, func(tInt *testing.T) {
			checkConfigField(tInt, c, referenceConfig, field)
		})
	}
}

func TestGenClientsRegionPrecedence(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Subject-Token", "test-token")
		_, _ = fmt.Fprint(w, `{
			"token": {
				"project": {
					"id": "test-project-id",
					"name": "eu-de_test-project",
					"domain": {"id": "test-domain-id"}
				},
				"user": {
					"id": "test-user-id",
					"domain": {"id": "test-domain-id"}
				},
				"catalog": []
			}
		}`)
	})

	projectAuth := golangsdk.AuthOptions{
		IdentityEndpoint: fmt.Sprintf("%sv3/", th.Endpoint()),
		TokenID:          "test-token",
		TenantName:       "eu-de_test-project",
	}
	domainAuth := golangsdk.AuthOptions{
		IdentityEndpoint: fmt.Sprintf("%sv3/", th.Endpoint()),
		TokenID:          "test-token",
	}

	t.Run("IAM region is used as fallback", func(t *testing.T) {
		cfg := &Config{}
		th.CheckNoErr(t, cfg.genClients(projectAuth, domainAuth))
		th.AssertEquals(t, "eu-de", cfg.Region)
	})

	t.Run("configured region is preserved", func(t *testing.T) {
		cfg := &Config{Region: "eu-nl"}
		th.CheckNoErr(t, cfg.genClients(projectAuth, domainAuth))
		th.AssertEquals(t, "eu-nl", cfg.Region)
	})
}

func genTemplate(def, attr, option, name string) string {
	return fmt.Sprintf(`
clouds:
  {{.Cloud}}:
    auth:
      auth_url: {{.IdentityEndpoint}}
      %s: {{.%s}}
      %s: {{.%s}}
`, def, attr, option, name)
}

func TestDomain(t *testing.T) {
	projectDefinition := map[string]string{
		"TenantID":   "project_id",
		"TenantName": "project_name",
	}
	synonyms := map[string][]string{
		"DomainName": {"user_domain_name", "domain_name", "project_domain_name"},
		"DomainID":   {"user_domain_id", "domain_id", "project_domain_id", "default_domain"},
	}
	for attr, def := range projectDefinition {
		for name, options := range synonyms {
			for _, option := range options {
				cloudName := fmt.Sprintf("otc-%s", def)
				_ = os.Setenv("OS_CLOUD", cloudName)

				tmpl := genTemplate(def, attr, option, name)
				var referenceConfig = &Config{
					Cloud:            cloudName,
					IdentityEndpoint: "https://localhost:9903/v3",
					TenantID:         "4b04680e-c627-4acb-a972-918cc661bcba",
					TenantName:       "eu-de",
					DomainName:       "OTC12392130",
					DomainID:         "19299b82-9df8-453d-a571-3681f5a4d883",
				}
				t.Run(fmt.Sprintf("%s/%s/%s", attr, name, option), func(tSyn *testing.T) {
					err := writeYamlTemplate(tmpl, fileName, referenceConfig)
					if err != nil {
						tSyn.Fatal(err)
					}
					defer func() { _ = os.Remove(fileName) }()

					config := &Config{
						Cloud:       referenceConfig.Cloud,
						environment: openstack.NewEnv(osPrefix, false),
					}
					if err = config.Load(); err != nil {
						tSyn.Fatal()
					}

					checkConfigField(tSyn, config, referenceConfig, name)
				})
			}
		}
	}

	defer func() { _ = os.Remove(fileName) }()
}

func testRequestRetry(t *testing.T, count int) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	retryCount := count

	var info = struct {
		retries int
		mut     *sync.RWMutex
	}{
		0,
		new(sync.RWMutex),
	}

	th.Mux.HandleFunc("/route/", func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()
		_, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("error hadling test request")
		}
		if info.retries < retryCount {
			info.mut.RLock()
			info.retries += 1
			info.mut.RUnlock()
			panic(err) // simulate EOF
		}
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, `%v`, info.retries)
	})

	cfg := &Config{MaxRetries: retryCount}
	_, err := cfg.genClient(golangsdk.AuthOptions{
		IdentityEndpoint: fmt.Sprintf("%s/route", th.Endpoint()),
	})
	_, ok := err.(golangsdk.ErrDefault500)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, retryCount, info.retries)
}

func TestRequestRetry(t *testing.T) {
	t.Run("TestRequestMultipleRetries", func(t *testing.T) { testRequestRetry(t, 2) })
	t.Run("TestRequestSingleRetry", func(t *testing.T) { testRequestRetry(t, 1) })
	t.Run("TestRequestZeroRetry", func(t *testing.T) { testRequestRetry(t, 0) })
}

type retryCaptureRoundTripper struct {
	calls          int
	xSdkDates      []string
	authorizations []string
	bodies         []string
}

func (rt *retryCaptureRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.calls += 1
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		rt.bodies = append(rt.bodies, string(body))
	}
	rt.xSdkDates = append(rt.xSdkDates, req.Header.Get("X-Sdk-Date"))
	rt.authorizations = append(rt.authorizations, req.Header.Get("Authorization"))

	if rt.calls == 1 {
		return nil, fmt.Errorf("temporary connection error")
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("{}")),
		Request:    req,
	}, nil
}

func TestRoundTripperReSignsAndRestoresBodyOnRetry(t *testing.T) {
	signOptions := golangsdk.SignOptions{
		AccessKey: "test-ak",
		SecretKey: "test-sk",
	}
	req, err := http.NewRequest(http.MethodPost, "https://example.com/v1/resource", strings.NewReader(`{"name":"test"}`))
	th.CheckNoErr(t, err)
	req.Header.Set("Content-Type", "application/json")
	golangsdk.Sign(req, signOptions)

	capture := &retryCaptureRoundTripper{}
	rt := &RoundTripper{
		Rt:          capture,
		MaxRetries:  1,
		SignOptions: &signOptions,
	}

	resp, err := rt.RoundTrip(req)
	th.CheckNoErr(t, err)
	th.AssertEquals(t, http.StatusOK, resp.StatusCode)
	th.AssertEquals(t, 2, capture.calls)
	th.AssertDeepEquals(t, []string{`{"name":"test"}`, `{"name":"test"}`}, capture.bodies)

	if capture.xSdkDates[0] == "" || capture.xSdkDates[0] == capture.xSdkDates[1] {
		t.Fatalf("expected retry to use a fresh X-Sdk-Date, got %q then %q", capture.xSdkDates[0], capture.xSdkDates[1])
	}
	if capture.authorizations[0] == "" || capture.authorizations[0] == capture.authorizations[1] {
		t.Fatalf("expected retry to use a fresh Authorization header")
	}
}

func TestRoundTripperDoesNotRetryNonReplayableBody(t *testing.T) {
	signOptions := golangsdk.SignOptions{
		AccessKey: "test-ak",
		SecretKey: "test-sk",
	}
	req, err := http.NewRequest(http.MethodPut, "https://example.com/v2/images/id/file", io.NopCloser(strings.NewReader("image-bytes")))
	th.CheckNoErr(t, err)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Sdk-Date", "20200101T000000Z")
	req.Header.Set("Authorization", "SDK-HMAC-SHA256 stale-signature")

	capture := &retryCaptureRoundTripper{}
	rt := &RoundTripper{
		Rt:          capture,
		OsDebug:     true,
		MaxRetries:  2,
		SignOptions: &signOptions,
	}

	_, err = rt.RoundTrip(req)
	if err == nil {
		t.Fatal("expected first transport error")
	}
	th.AssertEquals(t, 1, capture.calls)
	th.AssertDeepEquals(t, []string{"image-bytes"}, capture.bodies)
	th.AssertDeepEquals(t, []string{"20200101T000000Z"}, capture.xSdkDates)
}

// TestGenClientAKSKReSignOptions checks that SignOptions is set after AK/SK
// authentication, so retried requests are re-signed instead of reusing a stale
// signature (#3392).
func TestGenClientAKSKReSignOptions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	// AK/SK authentication lists the service catalog; an empty catalog is enough.
	th.Mux.HandleFunc("/v3/auth/catalog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"catalog": []}`)
	})

	cfg := &Config{MaxRetries: 1}
	// Ending the endpoint with /v3/ lets ChooseVersion skip version negotiation.
	client, err := cfg.genClient(golangsdk.AKSKAuthOptions{
		IdentityEndpoint: fmt.Sprintf("%sv3/", th.Endpoint()),
		AccessKey:        "test-ak",
		SecretKey:        "test-sk",
		ProjectId:        "test-project-id",
	})
	th.CheckNoErr(t, err)

	rt, ok := client.HTTPClient.Transport.(*RoundTripper)
	th.AssertEquals(t, true, ok)

	if rt.SignOptions == nil {
		t.Fatal("expected RoundTripper.SignOptions to be set for AK/SK auth, got nil")
	}
	th.AssertEquals(t, "test-ak", rt.SignOptions.AccessKey)
	th.AssertEquals(t, "test-sk", rt.SignOptions.SecretKey)
}
