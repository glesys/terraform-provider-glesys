package glesys

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
