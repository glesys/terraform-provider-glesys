package glesys

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/glesys/glesys-go/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIP_basic(t *testing.T) {
	name := "glesys_ip.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGlesysProviders,
		CheckDestroy: testAccIPResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: `resource "glesys_ip" "test" {
				    datacenter = "Stockholm"
				    platform   = "KVM"
				    version    = 4

				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "datacenter", "Stockholm"),
					resource.TestCheckResourceAttr(name, "platform", "KVM"),
					resource.TestCheckResourceAttr(name, "version", "4"),
					resource.TestCheckResourceAttrSet(name, "address"),
					resource.TestCheckResourceAttrSet(name, "broadcast"),
					resource.TestCheckResourceAttrSet(name, "gateway"),
					resource.TestCheckResourceAttrSet(name, "cost.0.amount"),
					resource.TestCheckResourceAttrSet(name, "cost.0.currency"),
					resource.TestCheckResourceAttrSet(name, "cost.0.time_period"),
					resource.TestCheckResourceAttrSet(name, "locked_to_account"),
					resource.TestCheckResourceAttrSet(name, "name_servers.0"),
					resource.TestCheckResourceAttrSet(name, "netmask"),
					resource.TestCheckResourceAttrSet(name, "platforms.0"),
					resource.TestCheckResourceAttrSet(name, "ptr"),
					resource.TestCheckResourceAttrSet(name, "reserved"),
				),
			},
		},
	})
}

func TestAccIP_updatePTR(t *testing.T) {
	name := "glesys_ip.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGlesysProviders,
		CheckDestroy: testAccIPResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: `resource "glesys_ip" "test" {
				    datacenter = "Stockholm"
				    platform   = "KVM"
				    version    = 4

				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "datacenter", "Stockholm"),
					resource.TestCheckResourceAttr(name, "platform", "KVM"),
					resource.TestCheckResourceAttr(name, "version", "4"),
					resource.TestCheckResourceAttrSet(name, "address"),
					resource.TestCheckResourceAttrSet(name, "ptr"),
					testAccIPDefaultPTRValue(name),
				),
			},
			{
				Config: `resource "glesys_ip" "test" {
				    datacenter = "Stockholm"
				    platform   = "KVM"
				    version    = 4
				    ptr        = "test.ptr."
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "datacenter", "Stockholm"),
					resource.TestCheckResourceAttr(name, "platform", "KVM"),
					resource.TestCheckResourceAttr(name, "version", "4"),
					resource.TestCheckResourceAttr(name, "ptr", "test.ptr."),
					resource.TestCheckResourceAttrSet(name, "address"),
				),
			},
		},
	})
}

func testAccIPResourceDestroy(s *terraform.State) error {
	client := testGlesysProvider.Meta().(*glesys.Client)

	for _, resource := range s.RootModule().Resources {
		if resource.Type != "glesys_ip" {
			continue
		}

		ip, err := client.IPs.Details(context.Background(), resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Unexpected error checking for released ip %s", err)
		}

		if ip.Reserved == "yes" {
			return fmt.Errorf("IP %s is still reserved after destroy", resource.Primary.ID)
		}
	}

	return nil
}

func testAccIPDefaultPTRValue(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", resourceName)
		}

		address := resource.Primary.Attributes["address"]
		ptr := resource.Primary.Attributes["ptr"]

		expectedPTR := fmt.Sprint(strings.ReplaceAll(address, ".", "-"), "-static.glesys.net.")
		if ptr != expectedPTR {
			return fmt.Errorf("Expected ptr value %s, got %s", expectedPTR, ptr)
		}

		return nil
	}
}
