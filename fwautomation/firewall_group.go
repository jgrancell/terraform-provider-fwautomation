package fwautomation

type FirewallGroup struct {
    Groupname string
    Hostname  string
    IPAddress string
}

type FirewallResponse struct {
  Status  string `json:"status"`
  Reason  string `json:"reason,omitempty"`
  Date    string `json:"date"`
  Version string `json:"version"`
}
