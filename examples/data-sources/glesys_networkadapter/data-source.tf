# glesys_networkadapter datasource
data "glesys_networkadapter" "nic1" {
  id = "bfcb4eac-b831-4124-9e5a-629d6516d205"
}

output "nic_networkid" {
  value = data.glesys_networkadapter.nic1.networkid
}
