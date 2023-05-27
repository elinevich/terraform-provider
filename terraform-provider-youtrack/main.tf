terraform {
  required_providers {
    youtrack = {
      source = "terraform.local/local/youtrack"
      version = "1.0.0"
    }
  }
}

provider "youtrack" {
  api_version = "v1"
  base_url = var.youtrack_base_url
  token = var.youtrack_token
}

data "name_users" "user_elinevich" {
  provider = youtrack
  login = "anna.e"
}

output "user" {
  value = data.name_users.user_elinevich
}
