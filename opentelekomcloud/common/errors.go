package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
)

// ConvertExpected403ErrInto404Err is a method used to parsing 403 error and try to convert it to 404 error according
// to the right error code.
// Arguments:
// + err: The error response obtained through HTTP/HTTPS request.
// + errCodeKey: The key name of the error code in the error response body, e.g. 'error_code', 'err_code'.
// + specErrCodes: One or more error codes that you wish to match against the current error, e.g. 'APIGW.0001'.
// Notes: If you missing specErrCodes input, this function will convert all 403 errors into 404 errors.
// How to use it:
// + For the general cases, their error code key is 'error_code', and we should call as follows:
//   - utils.ConvertExpected403ErrInto404Err(err, "error_code")
//   - utils.ConvertExpected403ErrInto404Err(err, "error_code", "DWS.0001")
//   - utils.ConvertExpected403ErrInto404Err(err, "error_code", []string{"DWS.0001", "DLM.3028"}...)
func ConvertExpected403ErrInto404Err(err error, errCodeKey string, specErrCodes ...string) error {
	var err403 golangsdk.ErrDefault403
	if !errors.As(err, &err403) {
		log.Printf("[WARN] Unable to recognize expected error type, want 'golangsdk.ErrDefault403', but got '%s'",
			reflect.TypeOf(err).String())
		return err
	}
	var apiError interface{}
	if jsonErr := json.Unmarshal(err403.Body, &apiError); jsonErr != nil {
		return err
	}

	errCode := PathSearch(errCodeKey, apiError, nil)
	if errCode == nil {
		// 4xx means the client parsing was failed.
		return golangsdk.ErrDefault400{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("Unable to find the error code from the error body using given error code key (%s), the error is: %#v",
					errCodeKey, apiError)),
			},
		}
	}

	if len(specErrCodes) < 1 {
		log.Printf("[INFO] Identified 403 error parsed it as 404 error (without the error code control)")
		return golangsdk.ErrDefault404{}
	}
	if StrSliceContains(specErrCodes, fmt.Sprint(errCode)) {
		log.Printf("[INFO] Identified 403 error with code '%v' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	log.Printf("[WARN] Unable to recognize expected error code (%v), want %v", errCode, specErrCodes)
	return err
}

// ConvertExpected400ErrInto404Err works like ConvertExpected403ErrInto404Err for services that answer a request
// for a missing object with 400 and a dedicated error code, e.g. CFW returns 400 'CFW.00200005' for a deleted
// address group or service group.
func ConvertExpected400ErrInto404Err(err error, errCodeKey string, specErrCodes ...string) error {
	var err400 golangsdk.ErrDefault400
	if !errors.As(err, &err400) {
		return err
	}
	var apiError interface{}
	if jsonErr := json.Unmarshal(err400.Body, &apiError); jsonErr != nil {
		return err
	}

	errCode := PathSearch(errCodeKey, apiError, nil)
	if errCode == nil {
		return err
	}
	if len(specErrCodes) < 1 || StrSliceContains(specErrCodes, fmt.Sprint(errCode)) {
		log.Printf("[INFO] Identified 400 error with code '%v' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	return err
}
