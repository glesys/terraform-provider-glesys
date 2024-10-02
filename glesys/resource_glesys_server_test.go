package glesys

import (
	"fmt"
	"testing"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_getTemplate(t *testing.T) {
	srv := &glesys.ServerDetails{}
	for _, tt := range []struct {
		name           string
		tfTemplate     string
		readTemplate   string
		readTemplateID string
		readTags       []string
		want           string
	}{
		{
			name:         "KVM_instance",
			tfTemplate:   "ubuntu-20-04",
			readTemplate: "Ubuntu 20.04 LTS (Focal Fossa)",
			readTags:     []string{"ubuntu", "ubuntu-lts", "ubuntu-20-04"},
			want:         "ubuntu-20-04",
		},
		{
			name:           "KVM_instance_UUID_Template",
			tfTemplate:     "fc5d38f7-4c9d-4920-a3a0-3252f71fe2c5",
			readTemplate:   "Ubuntu 20.04 LTS (Focal Fossa)",
			readTemplateID: "fc5d38f7-4c9d-4920-a3a0-3252f71fe2c5",
			readTags:       []string{"ubuntu", "ubuntu-lts", "ubuntu-20-04"},
			want:           "fc5d38f7-4c9d-4920-a3a0-3252f71fe2c5",
		},
		{
			name:         "VMware_instance",
			tfTemplate:   "Ubuntu 20.04 LTS 64-bit",
			readTemplate: "Ubuntu 20.04 LTS 64-bit",
			readTags:     []string{},
			want:         "Ubuntu 20.04 LTS 64-bit",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			srv.Template = tt.readTemplate
			srv.InitialTemplate.Name = tt.readTemplate
			srv.InitialTemplate.CurrentTags = tt.readTags
			srv.InitialTemplate.ID = tt.readTemplateID
			if got := getTemplate(tt.tfTemplate, srv); got != tt.want {
				t.Errorf("got: %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccServerVMware_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-srv-vmw")

	name := "glesys_server.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGlesysServerBase_VMware(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "hostname", rName),
					resource.TestCheckResourceAttr(name, "datacenter", "Falkenberg"),
					resource.TestCheckResourceAttr(name, "platform", "VMware"),
					resource.TestCheckResourceAttrSet(name, "ipv4_address"),
					resource.TestCheckResourceAttrSet(name, "ipv6_address"),
				),
			},
		},
	})
}

func testAccGlesysServerBase_VMware(name string) string {
	return fmt.Sprintf(`
		resource "glesys_server" "test" {
			hostname   = "%s"
			datacenter = "Falkenberg"
			platform   = "VMware"
			bandwidth  = 100
			cpu        = 1
			memory     = 1024
			storage    = 10
			template   = "Debian 12 64-bit"

			user {
		          username   = "acctestuser"
		          publickeys = ["ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINOCh8br7CwZDMGmINyJgBip943QXgkf7XdXrDMJf5Dl acctestuser@example-host"]
			  password   = "hunter123!"
			}
		} `, name)
}
