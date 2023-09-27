package glesys

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEmailAlias_basic(t *testing.T) {
	time.Sleep(2 * time.Second)
	newDomain := "tfemail-" + acctest.RandString(6) + ".com"
	rName := acctest.RandString(6)
	name := "glesys_emailalias.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: glesysEmailAliasSkeleton(newDomain, fmt.Sprintf("goto = \"%s@%s\"", rName, newDomain)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "emailalias", fmt.Sprintf("alice@%s", newDomain)),
					resource.TestCheckResourceAttr(name, "goto", fmt.Sprintf("%s@%s", rName, newDomain)),
				),
			},
			{
				Config: glesysEmailAliasSkeleton(newDomain, fmt.Sprintf("goto = \"kurt@%s\"", newDomain)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "goto", fmt.Sprintf("kurt@%s", newDomain)),
				),
			},
		},
	})
}

func glesysEmailAliasSkeleton(domain string, s string) string {
	return fmt.Sprintf(
		`resource "glesys_dnsdomain" "test" {
			name = "%s"
		}

		resource "glesys_emailalias" "test" {
			emailalias = "alice@${glesys_dnsdomain.test.name}"
			%s
		}`, domain, s)
}
