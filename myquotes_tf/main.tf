terraform {
  required_providers {
    citation2000 = {
      source = "github.com/remieven/citation2000"
    }
  }
}

provider "citation2000" {
    folder_path = "/workspaces/go-tf-provider-lab/myquotes"
}

resource "citation2000_quote" "malcolm" {
    author = "Ian Malcolm"
    message = "Auriez-vous projeté de mettre des dinosaures dans votre parc à dinosaures ?"
}

resource "citation2000_quote" "sattler" {
    author = "Ellie Sattler"
    message = "- Dieu crée les dinosaures. Dieu détruit les dinosaures. Dieu crée l'homme. L'homme détruit Dieu. L'homme crée les dinosaures. - Les dinosaures mangent l'homme. Et la femme hérite de la Terre."
}

output "sattler_quote_id" {
  value = citation2000_quote.sattler.id
}
