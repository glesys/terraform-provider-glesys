package glesys

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testGlesysProvider *schema.Provider
var testGlesysProviders map[string]terraform.ResourceProvider

func init() {
	testGlesysProvider = Provider().(*schema.Provider)
	testGlesysProviders = map[string]terraform.ResourceProvider{
		"glesys": testGlesysProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GLESYS_USERID"); v == "" {
		t.Fatal("GLESYS_USERID must be set for acceptance tests")
	}
	if v := os.Getenv("GLESYS_TOKEN"); v == "" {
		t.Fatal("GLESYS_TOKEN must be set for acceptance tests")
	}
}
