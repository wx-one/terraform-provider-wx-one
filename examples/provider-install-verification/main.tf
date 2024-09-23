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