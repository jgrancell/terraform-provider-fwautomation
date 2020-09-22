package fwautomation

import (
  "bytes"
  "context"
  "errors"
  "fmt"
  "encoding/json"
  "os"
  "regexp"

  "github.com/hashicorp/go-uuid"
  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
  "golang.org/x/crypto/ssh"
)

func resourceFirewallGroup() *schema.Resource {
  return &schema.Resource{
    CreateContext: resourceFirewallGroupCreate,
    ReadContext:   resourceFirewallGroupRead,
    //UpdateContext: resourceFirewallGroupUpdate,
    DeleteContext: resourceFirewallGroupDelete,
    Schema: map[string]*schema.Schema{
      "group_name": &schema.Schema{
        Type:     schema.TypeString,
        ForceNew: true,
        Required: true,
        ValidateFunc: func(val interface{}, key string) (warns[]string, errs []error) {
          v := val.(string)
          if match, _ := regexp.MatchString("[A-Z_]*", v); ! match {
            errs = append(errs, fmt.Errorf("%q includes invalid characters. May contain [uppercase letters, underscores].", key))
          }
          return
        },
      },
      "hostname": &schema.Schema{
        Type:     schema.TypeString,
        ForceNew: true,
        Required: true,
        ValidateFunc: func(val interface{}, key string) (warns[]string, errs []error) {
          v := val.(string)
          if match, _ := regexp.MatchString("[a-z\\.-]*", v); ! match {
            errs = append(errs, fmt.Errorf("%q must be a fully qualified domain name. May contain [letters, hyphens, periods].", key))
          }
          return
        },
      },
      "ip_address": &schema.Schema{
        Type:     schema.TypeString,
        ForceNew: true,
        Required: true,
        ValidateFunc: func(val interface{}, key string) (warns[]string, errs []error) {
          v := val.(string)
          if match, _ := regexp.MatchString("[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+", v); ! match {
            errs = append(errs, fmt.Errorf("%q includes invalid characters. May contain [uppercase letters, underscores].", key))
          }
          return
        },
      },
    },
    SchemaVersion: 1,
  }
}

func resourceFirewallGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  client := m.(*ssh.Client)
  output, err := runResourceFirewallGroupsTask(client, d, "add")
  debugLogOutput("create status", output.Status)
  debugLogOutput("create reason", output.Reason)
  if err != nil {
    return diag.FromErr(err)
  }

  newUuid, _ := uuid.GenerateUUID()
  d.SetId(newUuid)

  return diags
}

func resourceFirewallGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  return diags
}

func resourceFirewallGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  return resourceFirewallGroupRead(ctx, d, m)
}

func resourceFirewallGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  client := m.(*ssh.Client)
  output, err := runResourceFirewallGroupsTask(client, d, "remove")
  if err != nil {
    return diag.FromErr(err)
  }

  debugLogOutput(d.Id(), output.Status)
  debugLogOutput(d.Id(), output.Reason)
  d.SetId("")

  return diags
}

func runResourceFirewallGroupsTask(c *ssh.Client, d *schema.ResourceData, method string) (FirewallResponse, error) {
  resp := FirewallResponse{}
  session, err := c.NewSession()
  debugLogOutput("task runner debug 1", "passed")
  if err != nil {
    debugLogOutput("task runner debug 1", err.Error())
    return resp, err
  }
  defer session.Close()
  debugLogOutput("task runner debug 2", "passed")
  var b bytes.Buffer
  debugLogOutput("task runner debug 3", "passed")
  session.Stdout = &b
  debugLogOutput("task runner debug 4", "passed")
  cmd := generateCommand(d, method)
  debugLogOutput("task runner debug 5", "passed")
  err2 := session.Run(cmd)
  debugLogOutput("task runner debug 6", "passed")
  if err2 != nil {
    debugLogOutput("task runner debug 6", err.Error())
    return resp, err2
  }

  str := b.String()
  debugLogOutput("creation raw output", str)
  json.Unmarshal([]byte(str), &resp)

  if resp.Status == "failed" {
    return resp, errors.New("Could not add firewall group association. " +resp.Reason)
  }
  return resp, nil
}

func getValue(d *schema.ResourceData, key string, method string) string {
  var val string

  if d.HasChange(key) && method == "remove" {
    val_interface, _ := d.GetChange(key)
    val = val_interface.(string)
  } else {
    val = d.Get(key).(string)
  }

  return val
}

func generateCommand(d *schema.ResourceData, method string) string {
  group_name := getValue(d, "group_name", method)
  hostname := getValue(d, "hostname", method)
  ip_address := getValue(d, "ip_address", method)

  return "manage group group="+group_name+" hostname="+hostname+" ip="+ip_address+" method="+method
}

func debugLogOutput(id string, output string) {
  //Debug log for development
  f, _ := os.OpenFile("./terraform-provider-fwautomation.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  defer f.Close()
  _, err := f.WriteString(id+": "+output+"\n")
  if err != nil {
    panic(err)
  }
  f.Sync()
}
