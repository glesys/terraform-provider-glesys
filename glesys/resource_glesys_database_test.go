package glesys

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccGlesysDatabase_basic(t *testing.T) {
	resourceName := "glesys_database.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGlesysDatabaseConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testdb"),
					resource.TestCheckResourceAttr(resourceName, "engine", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
				),
			},
		},
	})
}

func testAccGlesysDatabaseConfig() string {
	return `
		resource "glesys_database" "test" {
            datacenterkey = "dc-fbg1"
            name = "tf-test1"
            engine = "mysql"
            engineversion = "8.0"
            plankey = "plan-1core-4gib-25gib"
            allowlist = ["127.0.0.1","127.0.0.2"]
        } `
}
