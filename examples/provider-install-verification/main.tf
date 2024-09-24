terraform {
  required_providers {
    wxone = {
      source = "hashicorp.com/edu/wx-one"
    }
  }
}

provider "wxone" {
}

data "wxone_project" "default" {}

resource "wxone_key" "test" {
  name = "test"
  project_wide = false
}

output "test_key" {
  value = wxone_key.test
}