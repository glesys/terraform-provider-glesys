resource "glesys_objectstorage_instance" "objectstore_example_fbg" {
  count = 1
  datacenter = "dc-sto1"
  description = "tf-test-fbg-${count.index}-2"
}

resource "glesys_objectstorage_credential" "objectstore_example_fbg_cred" {
  count = 1
  instanceid = glesys_objectstorage_instance.objectstore_example_fbg[0].id
  description = "tf-test-fbg-cred-${count.index}"
}
