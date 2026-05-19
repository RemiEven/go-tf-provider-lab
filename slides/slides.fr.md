# Provider Terraform
<!-- .slide: class="page-title" -->



## Plan
<!-- .slide: class="toc" -->

- Système de plugins de Terraform
- Définition du provider
- Tester son provider en local
- Définition d'une ressource
- Tests automatisés
- Documentation & publication



### Système de plugins de Terraform (1/2)

![](resources/terraform-plugin-overview.png)

- Core : définit une interface commune, et découvre...
- Plugins : binaires éxecutables, écrits en Go
	- principalement des providers (il y a aussi des provisioners)
	- chaque plugin est spécialisé
	- le core les lance et communique avec eux par gRPC (protocole v6 pour TF 1.0)

Quand `terraform plan` ou `terraform apply` est lancé, le core orchestre la communication avec les plugins adéquats. Ceux-ci ont juste à traiter des requêtes unitaires, comme "créé *cette* ressource", ou "quel est l'état de *celle-ci*".

Notes:
Le protocole par gRPC évolue (v6 pour TF 1.0), et il ne faut pas essayer d'en dépendre (aucune garantie de stabilité).
https://developer.hashicorp.com/terraform/plugin/how-terraform-works



### Système de plugins de Terraform (2/2)

Quand `terraform init` est lancé, Terraform :

- lit les fichiers de configuration dans le répertoire courant pour déterminer les pluings nécessaires
- cherche les plugins installés à plusieurs endroits
- parfois, télécharge des plugins
- décide quelle version des plugins utiliser
- écrit un *lock file*

Terraform lira ensuite ce lock file pour s'assurer que la même version des plugins sera utilisée jusqu'à ce que `terraform init` soit lancé de nouveau.



### Définition du provider

Pour créer un provider, il est recommandé d'utiliser le **plugin framework** (go module `github.com/hashicorp/terraform-plugin-framework`). C'est le successeur du **plugin SDKv2**.

La façon la plus simple est de se baser sur le repository "template" fourni par Hashicorp : [https://github.com/hashicorp/terraform-provider-scaffolding-framework](https://github.com/hashicorp/terraform-provider-scaffolding-framework)

Notes:
Il existe un adapteur pour migrer les providers du SDKv2 vers le plugin framework.



### Définition du provider : interface Provider

```go
type Provider interface {
    // Metadata returns the name and version of the provider
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

    // Schema returns what should be in the provider block in HCL files
	Schema(context.Context, SchemaRequest, *SchemaResponse)

	// Configure initializes API client (if any)
	Configure(context.Context, ConfigureRequest, *ConfigureResponse)

	// DataSources lists the datasources of the provider
	DataSources(context.Context) []func() datasource.DataSource

	// Resources lists the resources of the provider.
	Resources(context.Context) []func() resource.Resource
}
```

Notes:
Les 3 premières méthodes ont le même type de signature que des handlers HTTP
-> On développe bien un serveur (ici, en gRPC plutôt qu'en HTTP, mais cela ne change rien)



### Définition du provider

```go
// This type implements the framework's Provider interface.
// In Go, we only need to implements the necessary methods, there's no "implements" keyword.
type Citation2000Provider struct {
	version string
}
```



### Définition du provider : schéma (1/2)

Les *schémas* permettent d'indiquer à Terraform le contenu attendu dans un block HCL.

```go
schema.Schema{
    Attributes: map[string]schema.Attribute{
        "folder_path": schema.StringAttribute{
            MarkdownDescription: "Folder containing the files",
            Required:            true,
        },
    },
}
```

pour

```terraform
provider "citation2000" {
    folder_path = "/workspaces/go-tf-provider-lab/myquotes"
}
```



### Définition du provider : schéma (2/2)

La méthode `Schema` est celle appelée par le framework quand Terraform Core a besoin de cette information.

```go
func (p *Citation2000Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"folder_path": schema.StringAttribute{
				MarkdownDescription: "Path of the folder containing the json files",
				Required:            true,
			},
		},
	}
}
```



### Définition du provider : model

Le schema est utile à Terraform Core, mais n'est pas un objet Go dont on pourrait utiliser les propriétés.
Pour cela, on peut définir un type modèle :

```go
/* Used by Terraform to parse the HCL block of the provider. */
/* Matches the provider's schema */
type JsonFileProviderModel struct {
	FolderPath types.String `tfsdk:"folder_path"`
}
```

Notes:
Basé sur des tags, comme par exemple le parsing JSON.



### Configuration du provider

```go
/* Called at the beginning of the provider's lifecycle. */
func (p *Citation2000Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Citation2000ProviderModel.       // empty model
    diagnostics := req.Config.Get(ctx, &data) // hydrate it with the request's data
	resp.Diagnostics.Append(diagnostics...)   // add any potential problems to the response
	if resp.Diagnostics.HasError() {          // if something requires more than a warning, stop here
		return
	}

	resp.ResourceData = data.FolderPath.ValueString() // the folder path value will be passed to resources
}
```

Notes:
Pour ce provider on initialise pas de client d'API : à la place on remplit ResourceData, qui sera automatiquement passé aux ressources.



### Définition du provider : gestion des erreurs

```go
diagnostics := req.Config.Get(ctx, &data)
resp.Diagnostics.Append(diagnostics...)
if resp.Diagnostics.HasError() {
    return
}
```

Plutôt que de se baser sur le type `error` de la lib standard, on utilise des `diag.Diagnostic` du framework.
Même principe (__errors as values__) mais permet d'accumuler plusieurs erreurs/warnings.

On peut créer ses propres diagnostics :

```go
diag.NewErrorDiagnostic(
    "failed to create quote",
    "failed to create quote: "+err.Error(),
)
```

Notes:
https://developer.hashicorp.com/terraform/plugin/framework/diagnostics



### Définition du provider : listing des ressources et datasources gérées

```go
func (p *Citation2000Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewQuoteResource, // we'll see how to create this later
	}
}

func (p *Citation2000Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
```



### Définition du provider : fichier main

```go
var version string = "dev"

func main() {
	opts := providerserver.ServeOpts{
		Address: "github.com/remieven/citation2000",
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
```



### Tester son provider en local (1/2)

- Vérifier la valeur de la variable d'environnement `$GOBIN` (ex: `/go/bin`)
- Depuis le dossier du provider, `go install .`
    - À relancer à chaque fois qu'on modifie le code
- Dans `~/.terraformrc` :

```terraform
provider_installation {

  dev_overrides {
      "github.com/remieven/citation2000" = "/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Notes:
https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#prepare-terraform-for-local-provider-install



### Tester son provider en local (2/2)

On peut ensuite utiliser notre provider dans un projet :

```terraform
terraform {
  required_providers {
    jsonfile = {
      source = "github.com/remieven/citation2000"
    }
  }
}

provider "jsonfile" {
    # configuration spécifique au provider
    folder_path = "/workspaces/go-tf-provider-lab/myquotes"
}
```

Avec un provider en local, il n'est ni nécessaire, ni recommandé d'éxecuter `terraform init`.



### Définition d'une ressource : interface Resource

Comme pour le provider, on implémente une interface :

```go
type Resource interface {
	// Metadata returns the full name of the resource
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

	// Schema returns the schema for this resource
	Schema(context.Context, SchemaRequest, *SchemaResponse)

    /* CRUD operations */
    Create(context.Context, CreateRequest, *CreateResponse)
	Read(context.Context, ReadRequest, *ReadResponse)
	Update(context.Context, UpdateRequest, *UpdateResponse)
	Delete(context.Context, DeleteRequest, *DeleteResponse)
}
```

Il y a des interfaces complémentaires qu'on peut également satisfaire pour implémenter des fonctionnalités plus avancées (ex: `ResourceWithImportState`, `ResourceWithModifyPlan`).



### Définition d'une ressource

```go
// This type implements the framework's Resource interface.
type QuoteResource struct {
	folderPath string
}

func (r *QuoteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quote"
}
```



### Définition d'une ressource : schéma

```go
func (r *QuoteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Quote",
		Attributes: map[string]schema.Attribute{
			"message": schema.StringAttribute{
				MarkdownDescription: "Message of the quote",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the quote",
				Computed:            true,
				PlanModifiers: []planmodifier.String{ // this is how we tell Terraform this field's lifecyle isn't regular
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
```



### Définition d'une ressource : model

```go
type QuoteResourceModel struct {
	Message types.String `tfsdk:"message"`
	Author  types.String `tfsdk:"author"`
	ID      types.String `tfsdk:"id"`
}
```



### Définition d'une ressource : configuration (optionel)

```go
func (r *QuoteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	folderPath, ok := req.ProviderData.(string) // we retrieve the provider's configuration...
	if !ok {                                    // and check that it matches the type we expect
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected %T, got: %T. Please report this issue to the provider developers.", folderPath, req.ProviderData),
		)
		return
	}

	r.folderPath = folderPath
}
```



### Définition d'une ressource : méthode Create

```go
func (r *QuoteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    /* Retrieve input data from request */
	var data QuoteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

    /* Actually perform the resource creation */
	q := quote.Quote{ data.Message.ValueString() }
	id, err := quote.CreateQuoteFile(r.folderPath, q)
	if err != nil {
        diagnostic := diag.NewErrorDiagnostic(
            "failed to create quote",
            "failed to create quote: "+err.Error(),
        )
		resp.Diagnostics.Append(diagnostic)
		return
	}

    /* Complete input data with new fields and set it in the response */
	data.ID = types.StringValue(id)
	tflog.Trace(ctx, "created quote "+id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```



### Définition d'une ressource : import (optionel)

Implémentation de l'interface `resource.ResourceWithImportState`.
Si notre ressource implémente `resource.ResourceWithIdentity`, on peut utiliser un helper :

```go
// ImportState implements resource.ResourceWithImportState.
func (r *QuoteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("id"), req, resp)
}
```



### Logs

Les plugins n'écrivent pas leurs logs eux-mêmes : c'est Terraform Core qui s'en charge.

Package `tflog` du framework : `tflog.Trace(ctx, "created quote "+id)`.

Par défaut, Terraform n'affiche aucun log des providers.
Il faut les activer avec par exemple `TF_LOG=TRACE`, `TF_LOG=ERROR`...

Notes:
Pour des cas plus complexes il est également possible de lancer son provider en debug avec delve, mais c'est difficile et ça a des effets de bord (car le cycle de vie du processus du plugin n'est plus géré par Terraform core).



### Tests automatisés

Basé sur le système de tests unitaires de Go, mais avec une exécution contrôlée par un ensemble de fonctions utilitaires définies par le framework.

Les ressources sont vraiment créées : plus proche de tests d'intégration/d'acceptance que de tests unitaires.
Par défaut, pour les lancer, il faut passer la variable d'environnement `TF_ACC=true` à `go test ./...`.

On donne une configuration terraform puis on fait des assertions sur le state.
Les ressources créées sont automatiquement supprimées à la fin.



### Documentation

- Bien remplir dans chaque schéma les champs `MarkdownDescription`
- Ajouter des exemples de configuration dans le dossier `examples`

```
examples/
    provider/
        provider.tf
    resources/
        jsonfile_quote/
            resourcee.tf
```

Puis utiliser `hashicorp/terraform-plugin-docs/cmd` (souvent avec `go:generate`) pour générer les fichiers `.md` de la documentation.

Si on s'est basé sur le repository de scaffolding : `go generate ./...`.



### Publication

Voir https://developer.hashicorp.com/terraform/registry/providers/publishing .

Publier un provider dépend beaucoup de l'endroit où on le publie.
La registry Terraform publique propose une CD basée sur des webhooks.



### Provider Terraform : pour aller plus loin

- [Documentation](https://developer.hashicorp.com/terraform/plugin)
- [Godoc du framework](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework)
- [Tutoriel officiel](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework) (~2h30)
