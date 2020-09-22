provider "fwautomation" {
  domain                  = "example.com"
  management_server       = "manager.example.com:22"
  authentication_key_path = "/home/example/.ssh/id_rsa"
}

locals {
  firewall_groups = [
    "TEST-1",
    "TEST-2",
  ]
}

resource "fwautomation_fwgroup" "multiples" {
  count = length(local.firewall_groups)
  group_name = local.firewall_groups[count.index]
  hostname   = "my-server.example.com"
  ip_address = "1.1.1.1"
}
