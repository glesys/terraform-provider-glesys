package glesys

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEmailAccount_basic(t *testing.T) {
	time.Sleep(10 * time.Second)
	newDomain := "tfemail-" + acctest.RandString(6) + ".com"
	name := "glesys_emailaccount.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: glesysEmailAccountSkeleton(newDomain, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "emailaccount", "alice@"+newDomain),
					resource.TestCheckResourceAttr(name, "antispamlevel", "3"),
					resource.TestCheckResourceAttr(name, "antivirus", "yes"),
					resource.TestCheckResourceAttr(name, "autorespond", "no"),
					resource.TestCheckResourceAttr(name, "quotaingib", "1"),
				),
			},
			{
				Config: glesysEmailAccountSkeleton(newDomain, "quotaingib = 2\nrejectspam = \"yes\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "quotaingib", "2"),
					resource.TestCheckResourceAttr(name, "rejectspam", "yes"),
				),
			},
		},
	})
}

func glesysEmailAccountSkeleton(domain string, s string) string {
	return fmt.Sprintf(
		`resource "glesys_dnsdomain" "test" {
			name = "%s"
		}

		resource "glesys_emailaccount" "test" {
			emailaccount = "alice@${glesys_dnsdomain.test.name}"
			%s
		}`, domain, s)
}
