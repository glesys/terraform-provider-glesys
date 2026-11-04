package glesys

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorage_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-bsv")
	sName := acctest.RandomWithPrefix("tf-bsv-vm")

	rNameRename := acctest.RandomWithPrefix("tf-bsv-new")
	name := "glesys_blockstorage.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testGlesysAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGlesysServerBaseKVM(sName) + testAccGlesysBlockStorageVolume(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "size", "20"),
					resource.TestCheckResourceAttr(name, "tier", "capacity"),
				),
			},
			{
				// test rename
				Config: testAccGlesysServerBaseKVM(sName) + testAccGlesysBlockStorageVolume(rNameRename),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rNameRename),
					resource.TestCheckResourceAttr(name, "size", "20"),
					resource.TestCheckResourceAttr(name, "tier", "capacity"),
				),
			},
			{
				// test resize of blockstorage volume
				Config: testAccGlesysServerBaseKVM(sName) + testAccGlesysBlockStorageVolumeResize(rNameRename),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rNameRename),
					resource.TestCheckResourceAttr(name, "size", "23"),
					resource.TestCheckResourceAttr(name, "tier", "capacity"),
				),
			},
		},
	})
}

func testAccGlesysBlockStorageVolume(volName string) string {
	return fmt.Sprintf(`
		resource "glesys_blockstorage" "test" {
			serverid   = glesys_server.test.id
		  datacenter = "dc-fbg1"
			name       = "%s"
			size       = 20
			tier       = "capacity"
		} `, volName)
}

func testAccGlesysBlockStorageVolumeResize(volName string) string {
	return fmt.Sprintf(`
		resource "glesys_blockstorage" "test" {
			serverid   = glesys_server.test.id
		  datacenter = "dc-fbg1"
			name       = "%s"
			size       = 23
			tier       = "capacity"
		} `, volName)
}
