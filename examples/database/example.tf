resource "glesys_database" "mydb" {
  datacenterkey = "dc-fbg1"
  name = "tf-test1"
  engine = "mysql"
  engineversion = "8.0"
  plankey = "plan-1core-4gib-25gib"
  allowlist = ["127.0.0.1","127.0.0.2"]
}