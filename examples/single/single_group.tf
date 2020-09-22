provider "fwautomation" {
  domain                  = "example.com"
  management_server       = "manager.example.com:22"
  authentication_key_path = "/home/example/.ssh/id_rsa"
}

resource "fwautomation_fwgroup" "singleton" {
  group_name = "TEST-001"
  hostname   = "new-server.example.com"
  ip_address = "1.1.1.1"
}
