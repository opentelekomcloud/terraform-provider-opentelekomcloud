package cfg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/unknwon/com"
)

var maxTimeout = 20 * time.Second

// RoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default http client RoundTripper to allow for logging.
type RoundTripper struct {
	Rt          http.RoundTripper
	OsDebug     bool
	MaxRetries  int
	SignOptions *golangsdk.SignOptions
}

func retryTimeout(count int) time.Duration {
	timeout := time.Duration(math.Pow(2, float64(count))) * time.Second
	if timeout > maxTimeout {
		timeout = maxTimeout
	}
	half := timeout / 2
	return half + time.Duration(rand.Float64()*float64(half))
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

	canReplay := request.Body == nil || request.GetBody != nil

	response, err := lrt.roundTrip(request, canReplay)
	retry := 1
	for response == nil {
		if !canReplay {
			return nil, err
		}
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
		response, err = lrt.roundTrip(request, canReplay)
		retry += 1
	}

	if lrt.OsDebug {
		log.Printf("[DEBUG] OpenTelekomCloud Response Code: %d", response.StatusCode)
		log.Printf("[DEBUG] OpenTelekomCloud Response Headers:\n%s", formatHeaders(response.Header, "\n"))

		response.Body, err = lrt.logResponse(response.Body, response.Header.Get("Content-Type"))
	}

	return response, err
}

func cloneRequest(request *http.Request, canReplay bool) (*http.Request, error) {
	cloned := request.Clone(request.Context())
	cloned.Header = request.Header.Clone()

	if request.Body != nil {
		if canReplay {
			body, err := request.GetBody()
			if err != nil {
				return nil, err
			}
			cloned.Body = body
		} else {
			cloned.Body = request.Body
		}
	}

	return cloned, nil
}

func (lrt *RoundTripper) roundTrip(request *http.Request, canReplay bool) (*http.Response, error) {
	req, err := cloneRequest(request, canReplay)
	if err != nil {
		return nil, err
	}
	if lrt.SignOptions != nil && canReplay {
		golangsdk.ReSign(req, *lrt.SignOptions)
	}

	if lrt.OsDebug {
		log.Printf("[DEBUG] OpenTelekomCloud Request URL: %s %s", req.Method, req.URL)
		log.Printf("[DEBUG] OpenTelekomCloud Request Headers:\n%s", formatHeaders(req.Header, "\n"))

		if req.Body != nil {
			var err error
			req.Body, err = lrt.logRequest(req.Body, req.Header.Get("Content-Type"))
			if err != nil {
				return nil, err
			}
		}
	}

	return lrt.Rt.RoundTrip(req)
}

// logRequest will log the HTTP Request details.
// If the body is JSON, it will attempt to be pretty-formatted.
func (lrt *RoundTripper) logRequest(original io.ReadCloser, contentType string) (io.ReadCloser, error) {
	if strings.HasPrefix(contentType, "application/octet-stream") {
		log.Printf("[DEBUG] OpenTelekomCloud Request Body: not logging binary request body")
		return original, nil
	}

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

	return io.NopCloser(strings.NewReader(bs.String())), nil
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
		return io.NopCloser(strings.NewReader(bs.String())), nil
	}

	var buf bytes.Buffer
	defer func() { _ = original.Close() }()
	if _, err := io.Copy(&buf, original); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] the response is: %s", buf.String())
	log.Printf("[DEBUG] Not logging because OpenTelekomCloud response body isn't JSON")
	return io.NopCloser(strings.NewReader(buf.String())), nil
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
