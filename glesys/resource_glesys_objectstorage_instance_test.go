package glesys

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccObjectStorageInstance_basic(t *testing.T) {
	name := "glesys_objectstorage_instance.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: `resource "glesys_objectstorage_instance" "test" {
					datacenter = "dc-sto1"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "datacenter", "dc-sto1"),
					resource.TestCheckResourceAttr(name, "description", ""),
					resource.TestCheckResourceAttrSet(name, "secretkey"),
					resource.TestCheckResourceAttrSet(name, "accesskey"),
					resource.TestCheckResourceAttrSet(name, "created"),
				),
			},
			{
				Config: `resource "glesys_objectstorage_instance" "test" {
					datacenter = "dc-sto1"
					description = "test"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", "test"),
				),
			},
		},
	})
}

func TestAccObjectStorageInstance_updateDescription(t *testing.T) {
	name := "glesys_objectstorage_instance.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: `resource "glesys_objectstorage_instance" "test" {
					datacenter = "dc-sto1"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "datacenter", "dc-sto1"),
					resource.TestCheckResourceAttr(name, "description", ""),
					resource.TestCheckResourceAttrSet(name, "secretkey"),
					resource.TestCheckResourceAttrSet(name, "accesskey"),
					resource.TestCheckResourceAttrSet(name, "created"),
				),
			},
			{
				Config: `resource "glesys_objectstorage_instance" "test" {
					datacenter = "dc-sto1"
					description = "tf-test"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "datacenter", "dc-sto1"),
					resource.TestCheckResourceAttr(name, "description", "tf-test"),
					resource.TestCheckResourceAttrSet(name, "secretkey"),
					resource.TestCheckResourceAttrSet(name, "accesskey"),
					resource.TestCheckResourceAttrSet(name, "created"),
				),
			},
		},
	})
}
