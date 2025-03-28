<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# LederWorks IaC Tools Terraform Provider

The iactools Terraform Provider brings custom functions which is not possible to implement with pure terraform logic.

When using the iactools provider we recommend using the latest version of Terraform Core ([the latest version can be found here](https://developer.hashicorp.com/terraform/install)).

* [Terraform Website](https://www.terraform.io)
* [IaC Tools Provider Documentation](https://registry.terraform.io/providers/lederworks/iactools/latest/docs)
* [IaC Tools Provider Usage Examples](https://github.com/lederworks/terraform-provider-iactools/tree/main/examples)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.23

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

```hcl
# 1. Specify the version of the IaC Tools Provider to use
terraform {
  required_providers {
    iactools = {
      source = "lederworks/iactools"
      version = "0.2.0"
    }
  }
}

# 2. Configure the IaC Tools Provider
provider "iactools" {
 # requires no configuration
}

# 3. Run the inverse_cidr function
output "inverse_cidr" {
  value = provider::iactools::inverse_cidr("192.168.0.0/16", "192.168.1.0/24")
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
