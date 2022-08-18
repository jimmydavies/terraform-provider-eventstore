package eventstore

import (
  "context"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
  "github.com/madedotcom/eventstore-client-go/eventstore"
)

func resourceUser() *schema.Resource {
  return &schema.Resource{
    CreateContext: resourceUserCreate,
    ReadContext:   resourceUserRead,
    UpdateContext: resourceUserUpdate,
    DeleteContext: resourceUserDelete,
    Importer: &schema.ResourceImporter{
       StateContext: schema.ImportStatePassthroughContext,
    },
    Schema: map[string]*schema.Schema{
      "username": &schema.Schema{
        Type:     schema.TypeString,
        ForceNew: true,
        Required: true,
      },
      "password": &schema.Schema{
        Type:      schema.TypeString,
        Sensitive: true,
        Required:  true,
      },
      "fullname": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
        Default:  "",
      },
      "groups": {
        Type:     schema.TypeList,
        Optional: true,
        Elem:     &schema.Schema{Type: schema.TypeString},
        DefaultFunc:  func() (interface{}, error) {
          return []interface{}{}, nil
        },
      },
      "disabled": &schema.Schema{
        Type:     schema.TypeBool,
        Optional: true,
        Default:  false,
      },
    },
  }
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*eventstore.Client)
  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  var fullName string
  if d.Get("fullname") == nil {
    fullName = d.Get("username").(string)
  } else {
    fullName = d.Get("fullname").(string)
  }

  groups := []string{}

  for _, group := range d.Get("groups").([]interface{}) {
    groups = append(groups, group.(string))
  }

  disabled := d.Get("disabled").(bool)

  user, err := client.CreateUser(d.Get("username").(string), d.Get("password").(string), fullName, groups)

  d.SetId(user.UserName)
  d.Set("password",  d.Get("password").(string))

  if err != nil {
    return diag.FromErr(err)
  }

  if disabled == true {
    client.DisableUser(d.Get("username").(string))
  }

  return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*eventstore.Client)

  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  user, err := client.GetUser(d.Id())

  if err != nil {
    return diag.FromErr(err)
  }

  if user != nil {
    d.Set("username", user.UserName)
    d.Set("fullname", user.FullName)
    d.Set("groups", user.Groups)
    d.Set("disabled", user.Disabled)
  } else {
    d.SetId("")
  }

  return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*eventstore.Client)

  var err error

  if d.HasChange("fullname") || d.HasChange("groups") {

    var groups []string

    for _, group := range d.Get("groups").([]interface{}) {
      groups = append(groups, group.(string))
    }

    _, err := client.UpdateUser(d.Get("username").(string), d.Get("fullname").(string), groups)

    if err != nil {
      return diag.FromErr(err)
    }
  }

  if d.HasChange("password") {
    if !client.SetUserPassword(d.Get("username").(string), d.Get("password").(string)) {
      return diag.Errorf("Failed to update password")
    }
  }

  if d.HasChange("disabled") {
    if d.Get("disabled").(bool) {
      _, err = client.DisableUser(d.Get("username").(string))
    } else {
      _, err = client.EnableUser(d.Get("username").(string))
    }

    if err != nil {
      return diag.Errorf("Failed to Disable/Enable User")
    }
  }
    
  return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*eventstore.Client)

  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  success := client.DeleteUser(d.Get("username").(string))

  if !success {
    return diag.Errorf("Failed to delete user")
  }

  d.SetId("")

  return diags
}
