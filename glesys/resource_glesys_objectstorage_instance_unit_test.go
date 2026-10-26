package glesys

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestTrackedObjectStorageCredential(t *testing.T) {
	original := glesys.ObjectStorageCredential{CredentialID: "c-1", AccessKey: "ORIGINAL", SecretKey: "original-secret"}
	added := glesys.ObjectStorageCredential{CredentialID: "c-2", AccessKey: "ADDED", Description: "ObjectStoreUI"}

	tests := []struct {
		name      string
		creds     []glesys.ObjectStorageCredential
		accessKey string
		want      string
		wantOK    bool
	}{
		{"no credentials", nil, "ORIGINAL", "", false},
		{"no access key in state uses the first", []glesys.ObjectStorageCredential{added, original}, "", "ADDED", true},
		{"tracked credential listed second is kept", []glesys.ObjectStorageCredential{added, original}, "ORIGINAL", "ORIGINAL", true},
		{"tracked credential listed first is kept", []glesys.ObjectStorageCredential{original, added}, "ORIGINAL", "ORIGINAL", true},
		{"tracked credential deleted uses the first", []glesys.ObjectStorageCredential{added}, "ORIGINAL", "ADDED", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := trackedObjectStorageCredential(tt.creds, tt.accessKey)
			if ok != tt.wantOK || got.AccessKey != tt.want {
				t.Errorf("got (%q, %v), want (%q, %v)", got.AccessKey, ok, tt.want, tt.wantOK)
			}
		})
	}
}

// instanceDetailsClient returns a client whose objectstorage/instancedetails
// call answers with the given credentials, without secret keys.
func instanceDetailsClient(t *testing.T, credentials string) *glesys.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/objectstorage/instancedetails" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"response": {"instance": {"id": "os-ab123", "datacenter": "dc-sto1",
			"description": "", "created": "2026-10-08T08:11:53+00:00", "credentials": %s}}}`, credentials)
	}))
	t.Cleanup(srv.Close)

	client := glesys.NewClient("CL12345", "token", "test")
	if err := client.SetBaseURL(srv.URL); err != nil {
		t.Fatal(err)
	}
	return client
}

func TestObjectStorageInstanceRead_credentials(t *testing.T) {
	// Newer credential first, as in observed responses.
	const both = `[
		{"id": "c-2", "accesskey": "ADDED", "description": "ObjectStoreUI", "created": "2026-10-20T08:46:56+00:00"},
		{"id": "c-1", "accesskey": "ORIGINAL", "created": "2026-10-08T08:11:53+00:00"}]`
	const addedOnly = `[
		{"id": "c-2", "accesskey": "ADDED", "description": "ObjectStoreUI", "created": "2026-10-20T08:46:56+00:00"}]`

	tests := []struct {
		name          string
		credentials   string
		stateAccess   string
		stateSecret   string
		wantAccessKey string
		wantSecretKey string
	}{
		{"refresh keeps the tracked credential", both, "ORIGINAL", "original-secret", "ORIGINAL", "original-secret"},
		{"empty state (create, import) uses the first", both, "", "", "ADDED", ""},
		{"tracked credential deleted uses the first", addedOnly, "ORIGINAL", "original-secret", "ADDED", "original-secret"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := instanceDetailsClient(t, tt.credentials)
			d := schema.TestResourceDataRaw(t, resourceGlesysObjectStorageInstance().Schema, map[string]interface{}{
				"datacenter": "dc-sto1",
			})
			d.SetId("os-ab123")
			if tt.stateAccess != "" {
				if err := d.Set("accesskey", tt.stateAccess); err != nil {
					t.Fatal(err)
				}
			}
			if tt.stateSecret != "" {
				if err := d.Set("secretkey", tt.stateSecret); err != nil {
					t.Fatal(err)
				}
			}

			if diags := resourceGlesysObjectStorageInstanceRead(context.Background(), d, client); diags.HasError() {
				t.Fatalf("read: %v", diags)
			}
			if got := d.Get("accesskey").(string); got != tt.wantAccessKey {
				t.Errorf("accesskey = %q, want %q", got, tt.wantAccessKey)
			}
			if got := d.Get("secretkey").(string); got != tt.wantSecretKey {
				t.Errorf("secretkey = %q, want %q", got, tt.wantSecretKey)
			}
		})
	}
}
