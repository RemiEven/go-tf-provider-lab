# Go/Terraform provider lab

This repository contains examples and resources for anyone trying to learn about how to create a simple Terraform provider, using Hashicorp's Terraform plugin-framework.
The devcontainer configuration offers a ready-to-use environment, with both the Go toolchain and Terraform already installed.

The example is build around an application that manages quotes from movies, TV shows, books, ...

### Description of each folder

##### myquotes

This is where the quotes are stored, as json files.

##### quote-server

A simple HTTP server to display a page with all quotes in the `myquotes` folder.
To run it:

```bash
cd ./quote-server
go run main.go
```

Then open http://localhost:5678/ (or the appropriate host if using codespaces).

##### sample-quote-json

An example of Terraform project that uses the provider defined in the `terraform-provider-json-file` folder.

##### slides

Contains learning resources. Slides use [Zenika's sensei framework](https://github.com/Zenika/sensei).

##### terraform-provider-json-file

Contains the actual custom provider. Based on [the Terraform Provider Scaffolding repository from Hashicorp](https://github.com/hashicorp/terraform-provider-scaffolding-framework).

For an actual provider terraform provider that's supposed to be published on e.g. the Terraform registry, this folder would be the root of your repository.
See https://github.com/mgappa/terraform-provider-ncz-json-file .
