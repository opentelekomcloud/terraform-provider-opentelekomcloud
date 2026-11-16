package cfg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/topics"
)

func TestSmnV2ClientCloudsDomainReauthentication(t *testing.T) {
	for _, mode := range []string{"password", "token"} {
		for _, selector := range []string{"user_domain_name", "domain_name", "project_domain_name", "user_domain_id", "domain_id"} {
			t.Run(mode+"/"+selector, func(t *testing.T) {
				domainKey, domainValue := "name", "test-domain"
				if strings.HasSuffix(selector, "_id") {
					domainKey, domainValue = "id", "test-domain-id"
				}
				var authRequests, smnRequests atomic.Int32
				var endpoint string
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					if r.URL.Path != "/v3/auth/tokens" {
						if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, "/notifications/topics") ||
							r.Header.Get("X-Auth-Token") != "test-token" {
							t.Errorf("unexpected SMN request: %s %s", r.Method, r.URL.Path)
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						smnRequests.Add(1)
						_, _ = io.WriteString(w, `{"topic_urn":"urn:smn:eu-de:test-project:topic","request_id":"test-request"}`)
						return
					}
					authRequests.Add(1)
					var body struct {
						Auth struct {
							Identity struct {
								Password struct {
									User struct {
										Domain map[string]string `json:"domain"`
									} `json:"user"`
								} `json:"password"`
							} `json:"identity"`
							Scope struct {
								Project struct {
									Name   string            `json:"name"`
									Domain map[string]string `json:"domain"`
								} `json:"project"`
								Domain map[string]string `json:"domain"`
							} `json:"scope"`
						} `json:"auth"`
					}
					if r.Method == http.MethodPost {
						if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
							t.Errorf("decode synthetic IAM request: %v", err)
							w.WriteHeader(http.StatusBadRequest)
							return
						}
					} else if mode != "token" || r.Method != http.MethodGet {
						t.Errorf("unexpected IAM request method: %s", r.Method)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					checkDomain := func(domain map[string]string) {
						if len(domain) != 1 || domain[domainKey] != domainValue {
							t.Errorf("expected exactly one %s domain selector, got %#v", domainKey, domain)
						}
					}
					if mode == "password" {
						checkDomain(body.Auth.Identity.Password.User.Domain)
					}
					projectName := body.Auth.Scope.Project.Name
					if projectName != "" {
						checkDomain(body.Auth.Scope.Project.Domain)
					} else {
						if mode == "password" {
							checkDomain(body.Auth.Scope.Domain)
						}
						projectName = "eu-de"
					}
					w.Header().Set("X-Subject-Token", "test-token")
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusCreated)
					}
					_, _ = fmt.Fprintf(w, `{"token":{
						"project":{"id":%q,"name":%q,"domain":{"id":"test-domain-id"}},
						"user":{"id":"test-user-id","domain":{"id":"test-domain-id"}},
						"catalog":[{"type":"smnv2","name":"smn","endpoints":[{
							"interface":"public","region":"eu-de","url":%q
						}]}]
					}}`, projectName+"-id", projectName, endpoint+"/v2/"+projectName+"-id/")
				}))
				defer server.Close()
				endpoint = server.URL
				config := syntheticSmnCloudConfig(t, endpoint, mode, selector, domainValue)
				if err := config.LoadAndValidate(); err != nil {
					t.Fatalf("initial clouds authentication failed: %v", err)
				}
				if config.DomainID != "test-domain-id" || config.GetDomainID() != "test-domain-id" {
					t.Fatal("authentication-derived account ID must remain available")
				}
				originalClient := config.HwClient
				projectNames := []ProjectName{"eu-de", "eu-de_other", "eu-de"}
				if mode == "token" {
					projectNames[1] = "eu-de"
				}
				for _, projectName := range projectNames {
					client, err := config.SmnV2Client(projectName)
					if err != nil {
						t.Fatalf("SMN reauthentication with %s failed: %v", selector, err)
					}
					if client.ProjectID != string(projectName)+"-id" {
						t.Fatalf("wrong SMN project scope: %s", client.ProjectID)
					}
					if _, err := topics.Create(client, topics.CreateOpts{Name: "test-topic"}); err != nil {
						t.Fatalf("synthetic SMN topic creation failed: %v", err)
					}
				}
				if mode == "password" && selector == "user_domain_name" {
					t.Setenv("OS_DOMAIN_ID", "explicit-conflicting-domain-id")
					if _, err := config.SmnV2Client("eu-de"); err == nil ||
						!strings.Contains(err.Error(), "exactly one of DomainID or DomainName") {
						t.Fatalf("explicit conflicting domain selectors must still be rejected: %v", err)
					}
				}
				if config.HwClient != originalClient || config.TenantName != "eu-de" ||
					config.DomainID != "test-domain-id" {
					t.Fatal("SMN project reconfiguration mutated the shared provider configuration")
				}
				if authRequests.Load() != 8 || smnRequests.Load() != 3 {
					t.Fatalf("unexpected request counts: IAM=%d SMN=%d", authRequests.Load(), smnRequests.Load())
				}
			})
		}
	}
}

func syntheticSmnCloudConfig(t *testing.T, endpoint, mode, selector, domainValue string) *Config {
	t.Helper()
	dir := t.TempDir()
	auth := "      username: test-user\n      password: test-password\n"
	if mode == "token" {
		auth = "      token: test-token\n"
	}
	content := fmt.Sprintf("clouds:\n  smn-auth-test:\n    auth:\n      auth_url: %s/v3/\n      project_name: eu-de\n      %s: %s\n%s    region_name: eu-de\n",
		endpoint, selector, domainValue, auth)
	path := filepath.Join(dir, "clouds.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	empty := filepath.Join(dir, "empty.yaml")
	if err := os.WriteFile(empty, []byte("clouds: {}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	const prefix = "OS_"
	for _, key := range []string{
		"AUTH_URL", "USERNAME", "USER_ID", "PASSWORD", "TOKEN", "AUTH_TOKEN",
		"PROJECT_ID", "PROJECT_NAME", "TENANT_ID", "TENANT_NAME",
		"DOMAIN_ID", "DOMAIN_NAME", "USER_DOMAIN_ID", "USER_DOMAIN_NAME",
		"PROJECT_DOMAIN_ID", "PROJECT_DOMAIN_NAME", "DEFAULT_DOMAIN",
		"REGION_NAME", "REGION_ID", "CACERT", "CERT", "KEY",
	} {
		t.Setenv(prefix+key, "")
	}
	t.Setenv(prefix+"CLIENT_CONFIG_FILE", path)
	t.Setenv(prefix+"CLIENT_SECURE_FILE", empty)
	t.Setenv(prefix+"CLIENT_VENDOR_FILE", empty)
	return &Config{Cloud: "smn-auth-test", environment: openstack.NewEnv(prefix)}
}

func TestReconfigProjectNamePreservesOtherAuthDomains(t *testing.T) {
	for _, mode := range []string{"token", "aksk"} {
		t.Run(mode, func(t *testing.T) {
			config := syntheticSmnCloudConfig(t, "https://iam.example.com", "password", "user_domain_name", "test-domain")
			config.Username, config.Password = "test-user", "test-password"
			config.DomainID, config.DomainName = "configured-domain-id", "test-domain"
			config.Swauth = true
			if mode == "token" {
				config.Token = "test-token"
			} else {
				config.AccessKey, config.SecretKey = "test-ak", "test-sk"
			}
			actual, err := reconfigProjectName(*config, "eu-de")
			if err != nil {
				t.Fatal(err)
			}
			if actual.DomainID != "configured-domain-id" || actual.DomainName != "test-domain" {
				t.Fatalf("%s domain metadata changed during project reconfiguration", mode)
			}
		})
	}
}
