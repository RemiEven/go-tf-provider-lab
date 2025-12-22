terraform {
  required_providers {
    jsonfile = {
      source = "github.com/remieven/json-file"
    }
  }
}

provider "jsonfile" {
    folder_path = "/workspaces/go-tf-provider-lab/myquotes"
}

resource "jsonfile_quote" "joke1" {
    author = "adibou"
    message = "Coucou me revoilou 0"
}

resource "jsonfile_quote" "joke2" {
    author = "adibou"
    message = "Coucou me revoilou fjgn"
}

# import {
#   to = jsonfile_quote.joke3
#   identity = {
#     id = "0cd0b60e-8250-4a86-b155-c5d3568b47a0"
#   }
# }

output "joke2_id" {
  value = jsonfile_quote.joke2.id
}
