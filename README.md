# Terraform Provider for Oneuptime

OpenAPI specification for OneUptime. This document describes the API endpoints, request and response formats, and other details necessary for developers to interact with the OneUptime API.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0 **or** [OpenTofu](https://opentofu.org/docs/intro/install/) >= 1.6
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
```sh
git clone https://github.com/oneuptime/terraform-provider-oneuptime
cd terraform-provider-oneuptime
```

2. Build the provider using the Go `install` command:
```sh
go build
```

## Using the Provider

```terraform
terraform {
  required_providers {
    oneuptime = {
      source = "oneuptime/oneuptime"
      version = "12.0.34"
    }
  }
}

provider "oneuptime" {
  oneuptime_url = "https://oneuptime.com" # or your self-hosted instance URL
  api_key       = var.oneuptime_api_key
}
```

The source address carries no registry hostname on purpose: OpenTofu resolves it against `registry.opentofu.org` and Terraform against `registry.terraform.io`, and the provider is published to both. Drive it with `tofu` or `terraform` — see [docs/guides/opentofu.md](./docs/guides/opentofu.md).

## Reusable Modules

[`modules/`](./modules) holds hand-written, engine-agnostic HCL modules:

- [`monitoring-and-incident-response`](./modules/monitoring-and-incident-response) — HTTP monitors for a service, an on-call rotation paged when they fail, and a status page listing them.

```terraform
module "storefront" {
  source = "github.com/OneUptime/terraform-provider-oneuptime//modules/monitoring-and-incident-response?ref=v12.0.34"

  service_name = "storefront"
  monitors     = { homepage = { url = "https://example.com" } }
}
```

Pin `ref` to a published tag — this repository is regenerated on every release.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go build`. This will build the provider and put the provider binary in the current directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
make testacc
```

## Local Installation

To install the provider locally for testing:

```sh
make install
```

This will build and install the provider to your local Terraform plugins directory.

## Testing

To run unit tests:

```sh
go test ./...
```

To run acceptance tests:

```sh
TF_ACC=1 go test ./... -v -timeout 120m
```

## Documentation

Documentation is generated using [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs). Run the following command to generate documentation:

```sh
go generate
```

## Contributing

1. This is a read-only repository. The source code is generated from the OneUptime OpenAPI specification. You can check the main repository at [OneUptime](https://github.com/oneuptime/oneuptime). Please fork the main repository and make changes there.
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
