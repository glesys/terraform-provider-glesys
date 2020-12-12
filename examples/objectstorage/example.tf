resource "glesys_objectstorage_instance" "objectstore_example_sto" {
  count = 1
  datacenter = "dc-sto1"
  description = "tf-test-sto-${count.index}-2"
}

resource "glesys_objectstorage_credential" "objectstore_example_sto_cred" {
  count = 1
  instanceid = glesys_objectstorage_instance.objectstore_example_sto[0].id
  description = "tf-test-sto-cred-${count.index}"
}
