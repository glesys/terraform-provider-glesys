package glesys

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGlesysDNSDomain_Basic(t *testing.T) {
	domainName := randomTestName() + ".com"

	dataName := "data.glesys_dnsdomain.exampledata"
	resName := "glesys_dnsdomain.example"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: glesysResourceDNSDomainSkeleton(domainName) + glesysDataSourceDNSDomainSkeleton(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "name", domainName),
					resource.TestCheckResourceAttrPair(resName, "ttl", dataName, "ttl"),
					resource.TestCheckResourceAttrPair(resName, "expire", dataName, "expire"),
					resource.TestCheckResourceAttrPair(resName, "retry", dataName, "retry"),
					resource.TestCheckResourceAttrPair(resName, "refresh", dataName, "refresh"),
					resource.TestCheckResourceAttrPair(resName, "minimum", dataName, "minimum"),
				),
			},
		},
	})
}

func glesysResourceDNSDomainSkeleton(domain string) string {
	return fmt.Sprintf(
		`resource "glesys_dnsdomain" "example" {
	    name       = "%s"
        ttl = 3600
     }
	 `, domain)
}

func glesysDataSourceDNSDomainSkeleton() string {
	return `data "glesys_dnsdomain" "exampledata" {
            name  = glesys_dnsdomain.example.name
         }`
}
