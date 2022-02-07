package eventstore

import (
  "context"
  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
  "github.com/madedotcom/eventstore-client-go/eventstore"
)

// Provider -
func Provider() *schema.Provider {
  return &schema.Provider{
    Schema: map[string]*schema.Schema{
      "url": &schema.Schema{
        Type:        schema.TypeString,
        Optional:    true,
        DefaultFunc: schema.EnvDefaultFunc("EVENTSTORE_URL", nil),
      },
      "username": &schema.Schema{
        Type:        schema.TypeString,
        Optional:    true,
        DefaultFunc: schema.EnvDefaultFunc("EVENTSTORE_USERNAME", nil),
      },
      "password": &schema.Schema{
        Type:        schema.TypeString,
        Optional:    true,
        Sensitive:   true,
        DefaultFunc: schema.EnvDefaultFunc("EVENTSTORE_PASSWORD", nil),
      },
    },
    ResourcesMap: map[string]*schema.Resource{
      "eventstore_user":         resourceUser(),
      "eventstore_subscription": resourceSubscription(),
    },
    DataSourcesMap: map[string]*schema.Resource{
      "eventstore_user":     dataSourceUser(),
    },
    ConfigureContextFunc: providerConfigure,
  }
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
  url      := d.Get("url").(string)
  username := d.Get("username").(string)
  password := d.Get("password").(string)

  client, err := eventstore.NewClient(url, username, password)

  var diags diag.Diagnostics

  if err != nil {
    diags = append(diags, diag.Diagnostic{
      Severity: diag.Error,
      Summary:  "Unable to create Eventstore client",
      Detail:   err.Error(),
    })
    return nil, diags
  }

  return client, diags
}
