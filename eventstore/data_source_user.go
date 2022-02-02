package eventstore

import (
  "context"
  "strconv"
  "time"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
  "github.com/madedotcom/eventstore-client-go/eventstore"
)

func dataSourceUser() *schema.Resource {
  return &schema.Resource{
    ReadContext: dataSourceUserRead,
    Schema: map[string]*schema.Schema{
      "username": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
      },
      "fullname": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "groups": {
        Type:     schema.TypeList,
        Computed: true,
        Elem:     &schema.Schema{Type: schema.TypeString},
      },
      "disabled": &schema.Schema{
        Type:     schema.TypeBool,
        Computed: true,
      },
    },
  }
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*eventstore.Client)

  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  user, err := client.GetUser(d.Get("username").(string))
  if err != nil {
    return diag.FromErr(err)
  }

  if err := d.Set("fullname", user.FullName); err != nil {
    return diag.FromErr(err)
  }

  if err := d.Set("groups", user.Groups); err != nil {
    return diag.FromErr(err)
  }

  if err := d.Set("disabled", user.Disabled); err != nil {
    return diag.FromErr(err)
  }

  // always run
  d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

  return diags
}
