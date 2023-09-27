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
				Config: glesysEmailAliasSkeleton(newDomain, "goto = \""+rName+"@"+newDomain+"\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "emailalias", "alice@"+newDomain),
					resource.TestCheckResourceAttr(name, "goto", rName+"@"+newDomain),
				),
			},
			{
				Config: glesysEmailAliasSkeleton(newDomain, "goto = \"kurt@"+newDomain+"\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "goto", "kurt@"+newDomain),
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
func glesysEmailAliasUpdate(domain string, s string) string {
	return fmt.Sprintf(
		`resource "glesys_dnsdomain" "test" {
			name = "%s"
		}

		resource "glesys_emailalias" "test" {
			emailalias = "alice@${glesys_dnsdomain.test.name}"
			%s
		}`, domain, s)
}
