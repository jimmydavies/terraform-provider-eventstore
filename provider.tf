terraform {
  required_providers {
    eventstore = {
      version = "0.0.1"
      source = "made.com/made/eventstore"
    }
  }
}

provider "eventstore" {
  url = "http://eventstore.service.test.consul:2113"
}

data "eventstore_user" "admin" {
  username = "admin"
}

output "admin_fullname" {
  value = data.eventstore_user.admin.fullname
}
