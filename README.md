# Terraform Provider WX-One

## References

The project setup was done according to the documentation: 

- https://developer.hashicorp.com/terraform/plugin/framework
- https://github.com/hashicorp/terraform-provider-scaffolding-framework
- https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework

## Local Development

- run `go env GOBIN`
- create `.terraformrc` with the following content

```
provider_installation {

  dev_overrides {
      "hashicorp.com/edu/wx-one" = "/Users/christian/go/bin" // this is for a mac setup
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

### Create a new build and run

```bash
go install .  # create new build
cd examples/provider-install-verification # go to terraform example
WX_ONE_HOST=http://localhost:5000 WX_ONE_USERNAME=christian.wolf@wizardtales.com WX_ONE_PASSWORD=xxx TF_LOG_PROVIDER=DEBUG terraform plan # use your username and password and the correct host to run terraform plan 
```

- run `go install .` to create a build
- run `cd examples/provider-install-verification && terraform plan` to check if terraform is able to work with local provider build 

## Add new graphql queries/mutations or update schema

- for schema updates update file `internal/provider/schema.graphql` (it contains `customerSchema.gql` concatenated with `commonSchema.gql`)
  ```
  cat ../api-gateway/lib/graphql/commonSchema.gql > internal/provider/schema.graphql && \
  echo "" >> internal/provider/schema.graphql && \
  cat ../api-gateway/lib/graphql/customerSchema.gql >> internal/provider/schema.graphql
  ```
- to add or update queries and mutations add them to the file `internal/provider/genqlient.graphql`
- to add additional go bindings for types update `internal/provider/genqlient.yaml` (for a complete list of configuration options see https://github.com/Khan/genqlient/blob/main/docs/genqlient.yaml)

after updating run

```
cd internal/provider
go run github.com/Khan/genqlient
```

## Generate documentation

```
cd tools; go generate ./...
```

# Terraform Provider Scaffolding (Terraform Plugin Framework)

_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding). See [Which SDK Should I Use?](https://developer.hashicorp.com/terraform/plugin/framework-benefits) in the Terraform documentation for additional information._

This repository is a *template* for a [Terraform](https://www.terraform.io) provider. It is intended as a starting point for creating Terraform providers, containing:

- A resource and a data source (`internal/provider/`),
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.

These files contain boilerplate code that you will need to edit to create your own Terraform provider. Tutorials for creating Terraform providers can be found on the [HashiCorp Developer](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework) platform. _Terraform Plugin Framework specific guides are titled accordingly._

Please see the [GitHub template repository documentation](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) for how to create a new repository from this template on GitHub.

Once you've written your provider, you'll want to [publish it on the Terraform Registry](https://developer.hashicorp.com/terraform/registry/providers/publishing) so that others can use it.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
