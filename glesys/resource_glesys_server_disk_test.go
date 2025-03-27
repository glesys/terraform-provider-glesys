package glesys

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServerDiskVMware_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-vmw-disk")
	sName := acctest.RandomWithPrefix("tf-srv-vmw")

	name := "glesys_server_disk.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGlesysServerBaseVMware(sName) + testAccGlesysServerDiskVMware(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "size", "20"),
					resource.TestCheckResourceAttr(name, "type", "silver"),
				),
			},
		},
	})
}

func testAccGlesysServerDiskVMware(name string) string {
	return fmt.Sprintf(`
		resource "glesys_server_disk" "test" {
			serverid = glesys_server.test.id
			name     = "%s"
			size     = 20
			type     = "silver"
		} `, name)
}
