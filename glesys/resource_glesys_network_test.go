package glesys

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGlesysNetwork_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-vlan")

	name := "glesys_network.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGlesysNetwork(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", rName),
					resource.TestCheckResourceAttr(name, "datacenter", "Falkenberg"),
					resource.TestCheckResourceAttr(name, "public", "no"),
				),
			},
		},
	})
}

func testAccGlesysNetwork(description string) string {
	return fmt.Sprintf(`
		resource "glesys_network" "test" {
			description = "%s"
			datacenter  = "Falkenberg"
		} `, description)
}
