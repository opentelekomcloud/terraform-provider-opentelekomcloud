package opentelekomcloud

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
	"text/template"

	"github.com/huaweicloud/golangsdk"
	th "github.com/huaweicloud/golangsdk/testhelper"
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
		t.Errorf("Field %s: expected %s, got %s", fieldName, expected, actual)
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

	err := writeYamlTemplate(tmpl, fileName, referenceConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(fileName) }()

	c := &Config{Cloud: referenceConfig.Cloud}
	err = readCloudsYaml(c)
	if err != nil {
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
				tmpl := fmt.Sprintf(`
clouds:
  {{.Cloud}}:
    auth:
      auth_url: {{.IdentityEndpoint}}
      %s: {{.%s}}
      %s: {{.%s}}`, def, attr, option, name)
				var referenceConfig = &Config{
					Cloud:            "otc",
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

					c := &Config{Cloud: referenceConfig.Cloud}
					err = readCloudsYaml(c)
					if err != nil {
						tSyn.Fatal()
					}

					checkConfigField(tSyn, c, referenceConfig, name)
				})
				_ = os.Remove(fileName)
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
		defer r.Body.Close()
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error hadling test request")
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
	_, err := genClient(cfg, golangsdk.AuthOptions{
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
