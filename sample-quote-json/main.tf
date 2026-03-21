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
    message = "Auriez-vous projeté de mettre des dinosaures dans votre parc à dinosaures ?"
}

resource "jsonfile_quote" "joke2" {
    author = "adibou"
    message = "- Dieu crée les dinosaures. Dieu détruit les dinosaures. Dieu crée l'homme. L'homme détruit Dieu. L'homme crée les dinosaures. - Les dinosaures mangent l'homme. Et la femme hérite de la Terre."
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
