package glesys

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccLoadBalancerTarget_multiBackend covers
// https://github.com/glesys/terraform-provider-glesys/issues/241: Read
// previously cleared the resource's ID (and thus tripped Terraform's
// "produced inconsistent result after apply" error) whenever the
// LoadBalancer had more than one backend, because it aborted as soon as it
// saw the first non-matching backend instead of scanning the whole list.
// With only a single backend, that bug can't reproduce - hence this test
// deliberately creates two.
func TestAccLoadBalancerTarget_multiBackend(t *testing.T) {
	lbName := acctest.RandomWithPrefix("tf-acc-lb")
	tcpBackendName := acctest.RandomWithPrefix("tf-acc-be-tcp")
	httpBackendName := acctest.RandomWithPrefix("tf-acc-be-http")
	targetName := acctest.RandomWithPrefix("tf-acc-target")

	resourceName := "glesys_loadbalancer_target.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGlesysLoadBalancerTargetMultiBackend(lbName, tcpBackendName, httpBackendName, targetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", targetName),
					resource.TestCheckResourceAttr(resourceName, "backend", tcpBackendName),
					resource.TestCheckResourceAttr(resourceName, "targetip", "172.16.0.10"),
					resource.TestCheckResourceAttr(resourceName, "port", "8898"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func testAccGlesysLoadBalancerTargetMultiBackend(lbName, tcpBackendName, httpBackendName, targetName string) string {
	return fmt.Sprintf(`
		resource "glesys_loadbalancer" "test" {
			datacenter = "Falkenberg"
			name       = "%s"
		}

		resource "glesys_loadbalancer_backend" "tcp" {
			loadbalancerid = glesys_loadbalancer.test.id
			name           = "%s"
			mode           = "tcp"
		}

		resource "glesys_loadbalancer_backend" "http" {
			loadbalancerid = glesys_loadbalancer.test.id
			name           = "%s"
			mode           = "http"
		}

		resource "glesys_loadbalancer_target" "test" {
			loadbalancerid = glesys_loadbalancer.test.id
			backend        = glesys_loadbalancer_backend.tcp.id

			name     = "%s"
			port     = 8898
			targetip = "172.16.0.10"
			weight   = 15

			depends_on = [glesys_loadbalancer_backend.http]
		}
	`, lbName, tcpBackendName, httpBackendName, targetName)
}
