package cfg

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFormatJSONSmartConnectPasswords(t *testing.T) {
	rt := &RoundTripper{}
	for _, raw := range []string{
		`{"source_task":{"redis_password":"source-test-secret","redis_address":"192.0.2.1:6379"},"sink_task":{"redis_password":"sink-test-secret"}}`,
		`{"source_task":null,"sink_task":null}`,
		`{"source_task":{},"sink_task":{"redis_address":"192.0.2.1:6379"}}`,
		`{"task_name":"test"}`,
		`{"auth":{"identity":{"password":{"user":{"password":"iam-test-secret"}}}}}`,
	} {
		formatted := rt.formatJSON([]byte(raw))
		for _, secret := range []string{"source-test-secret", "sink-test-secret", "iam-test-secret"} {
			if strings.Contains(formatted, secret) {
				t.Fatal("password exposed in formatted debug body")
			}
		}
		var before, after map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &before); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(formatted), &after); err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"source_task", "sink_task"} {
			if original, ok := before[name].(map[string]interface{}); ok {
				actual := after[name].(map[string]interface{})
				if _, exists := original["redis_password"]; exists && actual["redis_password"] != "***" {
					t.Fatal("Redis password was not masked")
				}
				if actual["redis_address"] != original["redis_address"] {
					t.Fatal("non-secret task address changed")
				}
			}
		}
	}
}

func TestRoundTripperSmartConnectPasswords(t *testing.T) {
	const payload = `{"source_task":{"redis_password":"source-test-secret"},"sink_task":{"redis_password":"sink-test-secret"},"task_name":"test"}`
	var output bytes.Buffer
	writer := log.Writer()
	log.SetOutput(&output)
	t.Cleanup(func() { log.SetOutput(writer) })
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if string(body) != payload {
			t.Error("logging changed the transmitted request")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, payload)
	}))
	defer server.Close()
	request, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	rt := &RoundTripper{Rt: http.DefaultTransport, OsDebug: true}
	response, err := rt.RoundTrip(request)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != payload {
		t.Error("logging changed the received response")
	}
	for _, secret := range []string{"source-test-secret", "sink-test-secret"} {
		if strings.Contains(output.String(), secret) {
			t.Fatal("Redis password exposed in HTTP debug logs")
		}
	}
	if strings.Count(output.String(), `"redis_password": "***"`) != 4 {
		t.Fatal("source and sink passwords must be masked in both request and response logs")
	}
}
