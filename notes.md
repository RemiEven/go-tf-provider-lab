terraform core : binary, provides common interface, and discovers...
terraform plugins : executable binaries written in Go, the core communicates with them over gRPC (protocol v6 for tf 1.0).
Main type of plugins: providers (now there are also provisioners)

to create a provider, we can use the "plugin framework" (recommended, go module github.com/hashicorp/terraform-plugin-framework), or the older "plugin SDKv2"

doc https://developer.hashicorp.com/terraform/plugin
godoc https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework
tutorials https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework (148 minutes advertised)
template repository https://github.com/hashicorp/terraform-provider-scaffolding-framework


When terraform init is run, Terraform reads configuration files in the working directory to determine which plugins are necessary, searches for installed plugins in several locations, sometimes downloads additional plugins, decides which plugin versions to use, and writes a lock file to ensure Terraform will use the same plugin versions in this directory until terraform init runs again.

handling errors: diagnostics (package diag) https://developer.hashicorp.com/terraform/plugin/framework/diagnostics

https://developer.hashicorp.com/terraform/plugin/how-terraform-works

terraform registry
Registry protocol is open, no "official" implementation but eg. Artifactory supports it

to use a locally built provider: https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#prepare-terraform-for-local-provider-install

## to check

terraform registry manifest file
lock file
