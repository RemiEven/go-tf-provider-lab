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

resource "jsonfile_quote" "joke0" {
    author = "adibou"
    message = "Coucou me revoilou 0"
}

resource "jsonfile_quote" "joke2" {
    author = "adibou"
    message = "Coucou me revoilou fjgn"
}

output "ijijij" {
  value = jsonfile_quote.joke2.id
}

# output "blblblblblbl" {
#   value = jsonfile_quote.joke2.id
# }