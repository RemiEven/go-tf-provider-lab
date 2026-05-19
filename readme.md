# Go/Terraform provider lab

This repository contains examples and resources for anyone trying to learn about how to create a simple Terraform provider, using Hashicorp's Terraform plugin-framework.
The devcontainer configuration offers a ready-to-use environment, with both the Go toolchain and Terraform already installed.

The example is build around an application (Citation2000) that manages quotes from movies, TV shows, books, ...

### Description of each folder

##### myquotes

This is where the quotes are stored, as json files.

##### citation2000_webui

A simple HTTP server to display a page with all quotes in the `myquotes` folder.
To run it:

```bash
cd ./quote-server
go run main.go
```

Then open http://localhost:5678/ (or the appropriate host if using codespaces).

##### myquotes_tf

An example of Terraform project that uses the provider defined in the `terraform-provider-json-file` folder.

##### slides

Contains learning resources. Slides use [Zenika's sensei framework](https://github.com/Zenika/sensei).

##### terraform-provider-citation2000

Contains the actual custom provider. Based on [the Terraform Provider Scaffolding repository from Hashicorp](https://github.com/hashicorp/terraform-provider-scaffolding-framework).

For an actual provider terraform provider that's supposed to be published on e.g. the Terraform registry, this folder would be the root of your repository.
See https://github.com/mgappa/terraform-provider-ncz-json-file .

### Testing the provider locally

When using devcontainer, `$GOBIN` is `/go/bin`: this is where our provider's executable binary will be after running `go install .` from `/workspaces/go-tf-provider-lab/terraform-provider-citation2000`.
Still assuming you're using devcontainer, the home folder is `/home/vscode`.
This is where we need to add a `terraformrc` file with the following content:

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

This is how Terraform will know it should use the provider we have locally.

NB:

- don't forget to `go install .` each time you make a modification to your provider's code.
- do not run `terraform init` in the `myquotes_tf/` folder: it's neither necessary nor recommended when working with a "local" provider.
