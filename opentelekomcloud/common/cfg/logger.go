package cfg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/unknwon/com"
)

var maxTimeout = 10 * time.Minute

// RoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default http client RoundTripper to allow for logging.
type RoundTripper struct {
	Rt         http.RoundTripper
	OsDebug    bool
	MaxRetries int
}

func retryTimeout(count int) time.Duration {
	seconds := math.Pow(2, float64(count))
	timeout := time.Duration(seconds) * time.Second
	if timeout > maxTimeout { // won't wait more than maxTimeout
		timeout = maxTimeout
	}
	return timeout
}

// RoundTrip performs a round-trip HTTP request and logs relevant information about it.
func (lrt *RoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	defer func() {
		if request.Body != nil {
			_ = request.Body.Close()
		}
	}()

	// for future reference, this is how to access the Transport struct:
	// tlsconfig := lrt.Rt.(*http.Transport).TLSClientConfig

	var err error

	if lrt.OsDebug {
		log.Printf("[DEBUG] OpenTelekomCloud Request URL: %s %s", request.Method, request.URL)
		log.Printf("[DEBUG] OpenTelekomCloud Request Headers:\n%s", formatHeaders(request.Header, "\n"))

		if request.Body != nil {
			request.Body, err = lrt.logRequest(request.Body, request.Header.Get("Content-Type"))
			if err != nil {
				return nil, err
			}
		}
	}

	response, err := lrt.Rt.RoundTrip(request)
	// Retrying connection
	retry := 1
	for response == nil {
		if retry > lrt.MaxRetries {
			if lrt.OsDebug {
				log.Printf("[DEBUG] OpenTelecomCloud connection error, retries exhausted. Aborting")
			}
			err = fmt.Errorf("OpenTelecomCloud connection error, retries exhausted. Aborting. Last error was: %s", err)
			return nil, err
		}

		if lrt.OsDebug {
			log.Printf("[DEBUG] OpenTelecomCloud connection error, retry number %d: %s", retry, err)
		}
		time.Sleep(retryTimeout(retry))
		response, err = lrt.Rt.RoundTrip(request)
		retry += 1
	}

	if lrt.OsDebug {
		log.Printf("[DEBUG] OpenTelekomCloud Response Code: %d", response.StatusCode)
		log.Printf("[DEBUG] OpenTelekomCloud Response Headers:\n%s", formatHeaders(response.Header, "\n"))

		response.Body, err = lrt.logResponse(response.Body, response.Header.Get("Content-Type"))
	}

	return response, err
}

// logRequest will log the HTTP Request details.
// If the body is JSON, it will attempt to be pretty-formatted.
func (lrt *RoundTripper) logRequest(original io.ReadCloser, contentType string) (io.ReadCloser, error) {
	defer func() { _ = original.Close() }()

	var bs bytes.Buffer
	_, err := io.Copy(&bs, original)
	if err != nil {
		return nil, err
	}

	// Handle request contentType
	if strings.HasPrefix(contentType, "application/json") {
		debugInfo := lrt.formatJSON(bs.Bytes())
		log.Printf("[DEBUG] OpenTelekomCloud Request Body: %s", debugInfo)
	} else {
		log.Printf("[DEBUG] OpenTelekomCloud Request Body: %s", bs.String())
	}

	return ioutil.NopCloser(strings.NewReader(bs.String())), nil
}

// logResponse will log the HTTP Response details.
// If the body is JSON, it will attempt to be pretty-formatted.
func (lrt *RoundTripper) logResponse(original io.ReadCloser, contentType string) (io.ReadCloser, error) {
	if strings.HasPrefix(contentType, "application/json") {
		var bs bytes.Buffer
		defer func() { _ = original.Close() }()
		_, err := io.Copy(&bs, original)
		if err != nil {
			return nil, err
		}
		debugInfo := lrt.formatJSON(bs.Bytes())
		if debugInfo != "" {
			log.Printf("[DEBUG] OpenTelekomCloud Response Body: %s", debugInfo)
		}
		return ioutil.NopCloser(strings.NewReader(bs.String())), nil
	}

	var buf bytes.Buffer
	defer func() { _ = original.Close() }()
	if _, err := io.Copy(&buf, original); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] the response is: %s", buf.String())
	log.Printf("[DEBUG] Not logging because OpenTelekomCloud response body isn't JSON")
	return ioutil.NopCloser(strings.NewReader(buf.String())), nil
}

// formatJSON will try to pretty-format a JSON body.
// It will also mask known fields which contain sensitive information.
func (lrt *RoundTripper) formatJSON(raw []byte) string {
	var data map[string]interface{}

	err := json.Unmarshal(raw, &data)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse OpenTelekomCloud JSON: %s", err)
		return string(raw)
	}

	// Mask known password fields
	if v, ok := data["auth"].(map[string]interface{}); ok {
		if v, ok := v["identity"].(map[string]interface{}); ok {
			if v, ok := v["password"].(map[string]interface{}); ok {
				if v, ok := v["user"].(map[string]interface{}); ok {
					v["password"] = "***"
				}
			}
		}
	}

	// Ignore the catalog
	if v, ok := data["token"].(map[string]interface{}); ok {
		if _, ok := v["catalog"]; ok {
			return ""
		}
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("[DEBUG] Unable to re-marshal OpenTelekomCloud JSON: %s", err)
		return string(raw)
	}

	return string(pretty)
}

// formatHeaders processes a headers object plus a deliminator, returning a string
func formatHeaders(headers http.Header, separator string) string {
	redactedHeaders := redactHeaders(headers)
	sort.Strings(redactedHeaders)

	return strings.Join(redactedHeaders, separator)
}

// List of headers that need to be redacted
var headersToRedact = []string{
	"x-auth-token",
	"x-auth-key",
	"x-service-token",
	"x-storage-token",
	"x-account-meta-temp-url-key",
	"x-account-meta-temp-url-key-2",
	"x-container-meta-temp-url-key",
	"x-container-meta-temp-url-key-2",
	"set-cookie",
	"x-subject-token",
}

// redactHeaders processes a headers object, returning a redacted list
func redactHeaders(headers http.Header) (processedHeaders []string) {
	for name, header := range headers {
		for _, v := range header {
			if com.IsSliceContainsStr(headersToRedact, name) {
				processedHeaders = append(processedHeaders, fmt.Sprintf("%v: %v", name, "***"))
			} else {
				processedHeaders = append(processedHeaders, fmt.Sprintf("%v: %v", name, v))
			}
		}
	}
	return
}
