package glesys

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPrivateNetworkSegment_basic(t *testing.T) {
	pName := acctest.RandomWithPrefix("tf-pn")
	rName := acctest.RandomWithPrefix("tf-pn-seg")

	name := "glesys_privatenetwork_segment.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGlesysPrivateNetworkBase(pName) + testAccGlesysPrivateNetworkSegmentBase(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "datacenter", "dc-fbg1"),
					resource.TestCheckResourceAttr(name, "platform", "kvm"),
					resource.TestCheckResourceAttr(name, "ipv4subnet", "192.168.2.0/24"),
					resource.TestCheckResourceAttrSet(name, "ipv6subnet"),
				),
			},
		},
	})
}

func testAccGlesysPrivateNetworkSegmentBase(name string) string {
	return fmt.Sprintf(`
		resource "glesys_privatenetwork_segment" "test" {
			privatenetworkid = glesys_privatenetwork.test.id
			name = "%s"
			datacenter = "dc-fbg1"
			platform = "kvm"
			ipv4subnet = "192.168.2.0/24"
		} `, name)
}
