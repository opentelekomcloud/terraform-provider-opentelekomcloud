package common

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
)

func TestCheckDeletedRecognizesWrappedNotFound(t *testing.T) {
	tests := map[string]struct {
		err       error
		wantError bool
		wantID    string
	}{
		"direct not found": {
			err:    golangsdk.ErrDefault404{},
			wantID: "",
		},
		"wrapped not found": {
			err:    fmt.Errorf("error retrieving resource: %w", golangsdk.ErrDefault404{}),
			wantID: "",
		},
		"other error": {
			err:       errors.New("request failed"),
			wantError: true,
			wantID:    "resource-id",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, nil)
			d.SetId("resource-id")

			err := CheckDeleted(d, tc.err, "error reading resource")
			if (err != nil) != tc.wantError {
				t.Fatalf("CheckDeleted() error = %v, wantError %t", err, tc.wantError)
			}
			if d.Id() != tc.wantID {
				t.Fatalf("resource ID = %q, want %q", d.Id(), tc.wantID)
			}
		})
	}
}

func TestCheckDeletedDiagRecognizesWrappedNotFound(t *testing.T) {
	tests := map[string]struct {
		err       error
		wantError bool
		wantID    string
	}{
		"direct not found": {
			err:    golangsdk.ErrDefault404{},
			wantID: "",
		},
		"wrapped not found": {
			err:    fmt.Errorf("error retrieving resource: %w", golangsdk.ErrDefault404{}),
			wantID: "",
		},
		"other error": {
			err:       errors.New("request failed"),
			wantError: true,
			wantID:    "resource-id",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, nil)
			d.SetId("resource-id")

			diags := CheckDeletedDiag(d, tc.err, "error reading resource")
			if diags.HasError() != tc.wantError {
				t.Fatalf("CheckDeletedDiag() diagnostics = %v, wantError %t", diags, tc.wantError)
			}
			if d.Id() != tc.wantID {
				t.Fatalf("resource ID = %q, want %q", d.Id(), tc.wantID)
			}
		})
	}
}
