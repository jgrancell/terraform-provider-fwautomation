package fwautomation

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/ssh"
)

type ManagementConfig struct {
	Server string
	Domain string
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"management_server": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("FWGROUPS_SERVER", nil),
			},
			"domain": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("FWGROUPS_DOMAIN", nil),
			},
			"authentication_key_path": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("FWGROUPS_AUTH_KEY_PATH", nil),
			},
		},
		ResourcesMap:   map[string]*schema.Resource{
			"fwautomation_fwgroup": resourceFirewallGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Loading the private key for authentications
	key, err := ioutil.ReadFile(d.Get("authentication_key_path").(string))
	if err != nil{
		return nil, diag.FromErr(err)
	}

	// Creating the Signer for the Private Key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// Creating the AuthMethod object with the Signer.
	config := &ssh.ClientConfig{
		User: "automate",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: time.Duration(5)*time.Second,
	}

	c, err := ssh.Dial("tcp", d.Get("management_server").(string), config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
