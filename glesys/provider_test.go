package glesys

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/glesys/glesys-go/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testNamePrefix = "tf-acc-test-"

var testGlesysProvider *schema.Provider
var testGlesysProviders map[string]*schema.Provider

func init() {
	testGlesysProvider = Provider()
	testGlesysProviders = map[string]*schema.Provider{
		"glesys": testGlesysProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProviderURLOverride(t *testing.T) {
	apiURL := "https://dev.example.test"

	rawProvider := Provider()
	raw := map[string]interface{}{
		"userid":       "cl12345",
		"token":        "MYTOKEN",
		"api_endpoint": apiURL,
	}

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}

	meta := rawProvider.Meta()
	if meta == nil {
		t.Fatalf("Expected metadata, got nil")
	}

	client := meta.(*glesys.Client)
	if client.BaseURL.String() != apiURL {
		t.Fatalf("Expected %s, got %s", apiURL, client.BaseURL.String())
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Skip("TF_ACC not set, skipping acceptance tests")
	}

	if v := os.Getenv("GLESYS_USERID"); v == "" {
		t.Fatal("GLESYS_USERID must be set for acceptance tests")
	}
	if v := os.Getenv("GLESYS_TOKEN"); v == "" {
		t.Fatal("GLESYS_TOKEN must be set for acceptance tests")
	}
}

func randomTestName(additionalNames ...string) string {
	prefix := testNamePrefix
	for _, n := range additionalNames {
		prefix += "-" + strings.ReplaceAll(n, " ", "_")
	}
	return randomName(prefix, 10)
}

func randomName(prefix string, length int) string {
	return fmt.Sprintf("%s%s", prefix, acctest.RandString(length))
}

func diagnosticsToString(diags diag.Diagnostics) string {
	diagsAsStrings := make([]string, len(diags))
	for i, diag := range diags {
		diagsAsStrings[i] = diag.Summary
	}

	return strings.Join(diagsAsStrings, "; ")
}
