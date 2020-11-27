package glesys

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccObjectStorageCredential_basic(t *testing.T) {
	name := "glesys_objectstorage_credential.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: glesysObjectStorageSkeleton(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", ""),
					resource.TestCheckResourceAttrSet(name, "secretkey"),
					resource.TestCheckResourceAttrSet(name, "accesskey"),
					resource.TestCheckResourceAttrSet(name, "created"),
					resource.TestCheckResourceAttrSet(name, "instanceid"),
				),
			},
			{
				Config: glesysObjectStorageSkeleton("description = \"tf-test\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", "tf-test"),
				),
			},
		},
	})
}

func glesysObjectStorageSkeleton(s string) string {
	return fmt.Sprintf(
		`resource "glesys_objectstorage_instance" "test" {
			datacenter = "dc-sto1"
		}

		resource "glesys_objectstorage_credential" "test" {
			instanceid = glesys_objectstorage_instance.test.id
			%s
		}`, s)
}
