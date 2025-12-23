terraform registry
Registry protocol is open, no "official" implementation but eg. Artifactory supports it


## to check

terraform registry manifest file
lock file

---

```sh
git clone https://github.com/hashicorp/terraform-provider-scaffolding-framework
mv terraform-provider-scaffolding-framework terraform-provider-json-file
cd terraform-provider-json-file/
go mod edit -module terraform-provider-json-file
go mod tidy
```

/// pour générer la doc : go generate ./... (folder examples)


attention quand on fait le plan avec un import (terraform plan -generate-config-out=generated.tf)
-> ca ajoute dans le block generé un field `provider =` qui fait vriller le validator de versions lockées par TF jusqu'à ce qu'on le retire 
