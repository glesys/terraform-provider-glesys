package glesys

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPrivateNetwork_basic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-gle-pn")
	name := "glesys_privatenetwork.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGlesysPrivateNetworkBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "ipv6aggregate"),
				),
			},
		},
	})
}

func TestAccPrivateNetworkUpdate(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-pn-edit")

	newName := acctest.RandomWithPrefix("tf-pn-edit-upd")
	name := "glesys_privatenetwork.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGlesysPrivateNetworkBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "ipv6aggregate"),
				),
			},
			{
				Config: testAccGlesysPrivateNetworkUpdatedName(newName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", newName),
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "ipv6aggregate"),
				),
			},
		},
	})
}

func testAccGlesysPrivateNetworkBase(name string) string {
	return fmt.Sprintf(`
		resource "glesys_privatenetwork" "test" {
			name = "%s"
		}`, name)
}

func testAccGlesysPrivateNetworkUpdatedName(name string) string {
	return fmt.Sprintf(`
		resource "glesys_privatenetwork" "test" {
			name = "%s"
		}`, name)
}
